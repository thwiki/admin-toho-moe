package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/list"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := pocketbase.New()

	publicDir := "./public"
	indexFallback := true

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(publicDir), indexFallback))
		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/tags",
			Handler: func(c echo.Context) error {
				authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
				if authRecord == nil || !isUser(authRecord) {
					return apis.NewForbiddenError("Only auth records can access this endpoint", nil)
				}

				collection, err := app.Dao().FindCollectionByNameOrId("links")
				if err != nil {
					return apis.NewNotFoundError("Collection not found", nil)
				}

				field := collection.Schema.GetFieldByName("tags")
				return c.JSON(http.StatusOK, field.Options)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/push",
			Handler: func(c echo.Context) error {
				authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
				if authRecord == nil || !isOp(authRecord) {
					return apis.NewForbiddenError("Only auth records can access this endpoint", nil)
				}

				err := pushUpdate()
				if err != nil {
					return err
				}

				return c.JSON(http.StatusOK, map[string]string{
					"message": "Pushed",
				})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/view",
			Handler: func(c echo.Context) error {
				authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
				if authRecord == nil || !isOp(authRecord) {
					return apis.NewForbiddenError("Only auth records can access this endpoint", nil)
				}

				stats, err := getView()
				if err != nil {
					return apis.NewApiError(http.StatusServiceUnavailable, "Backend error", nil)
				}

				for _, stat := range stats.Data {
					slug := strings.Trim(stat.Key, "/ ")
					record, err := app.Dao().FindFirstRecordByData("links", "slug", slug)
					if err != nil {
						continue
					}
					record.Set("vercelView", stat.Total)
					app.Dao().SaveRecord(record)
				}

				return c.JSON(http.StatusOK, map[string]string{
					"message": "Updated",
				})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})
		return nil
	})

	app.OnRecordBeforeAuthWithOAuth2Request().Add(func(e *core.RecordAuthWithOAuth2Event) error {
		if username, ok := e.OAuth2User.RawUser["username"]; ok {
			if username, ok := username.(string); ok {
				e.OAuth2User.Username = username
			}
		}
		if avatar, ok := e.OAuth2User.RawUser["avatar"]; ok {
			if avatar, ok := avatar.(string); ok {
				e.OAuth2User.AvatarUrl = avatar
			}
		}

		return nil
	})

	app.OnRecordAuthRequest().Add(func(e *core.RecordAuthEvent) error {
		meta, ok := e.Meta.(*auth.AuthUser)
		if !ok {
			return nil
		}

		e.Record.SetUsername(meta.Username)
		e.Record.Set("avatarUrl", meta.AvatarUrl)

		if groups, ok := meta.RawUser["groups"]; ok {
			filteredGroups := filterGroups(groups)
			e.Record.Set("groups", filteredGroups)
		}

		if err := app.Dao().SaveRecord(e.Record); err != nil {
			return err
		}

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func filterGroups(val any) []string {
	groups := list.ToUniqueStringSlice(val)
	filtered := []string{}

	for _, group := range groups {
		if group == "tohomoeuser" || group == "tohomoeop" {
			filtered = append(filtered, group)
		}
	}
	return filtered
}

func isUser(authRecord *models.Record) bool {
	groups := authRecord.GetStringSlice("groups")
	for _, group := range groups {
		if group == "tohomoeuser" || group == "tohomoeop" {
			return true
		}
	}
	return false
}

func isOp(authRecord *models.Record) bool {
	groups := authRecord.GetStringSlice("groups")
	for _, group := range groups {
		if group == "tohomoeop" {
			return true
		}
	}
	return false
}

func pushUpdate() (err error) {
	hook := os.Getenv("PUSH_HOOK")
	client := &http.Client{}

	req, err := http.NewRequest("POST", hook, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	return nil
}

func getView() (stats VercelInsightsStats, err error) {
	fromTime := time.Now().Add(time.Hour * -6 * 30 * 24).Format("2006-01-02T15:04:05-0700")
	toTime := time.Now().Format("2006-01-02T15:04:05-0700")
	teamId := os.Getenv("VERCEL_TEAM_ID")
	projectId := os.Getenv("VERCEL_PROJECT_ID")
	filter := "{}"
	environment := "production"
	limit := "150"

	endpoint, err := url.Parse("https://vercel.com/api/web/insights/stats/path")
	if err != nil {
		return
	}

	query := endpoint.Query()
	query.Add("from", fromTime)
	query.Add("to", toTime)
	query.Add("teamId", teamId)
	query.Add("projectId", projectId)
	query.Add("filter", filter)
	query.Add("environment", environment)
	query.Add("limit", limit)
	endpoint.RawQuery = query.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("VERCEL_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &stats)
	if err != nil {
		return
	}

	return stats, nil
}

type VercelInsightsStats struct {
	Data []VercelInsightsStat `json:"data"`
}

type VercelInsightsStat struct {
	Key     string `json:"key"`
	Total   int64  `json:"total"`
	Devices int64  `json:"devices"`
}
