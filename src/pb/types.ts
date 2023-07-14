/**
 * This file was @generated using pocketbase-typegen
 */

export enum Collections {
	Links = 'links',
	Users = 'users'
}

// Alias types for improved usability
export type IsoDateString = string;
export type RecordIdString = string;
export type HTMLString = string;

// System fields
export type BaseSystemFields<T = never> = {
	id: RecordIdString;
	created: IsoDateString;
	updated: IsoDateString;
	collectionId: string;
	collectionName: Collections;
	expand?: T;
};

export type AuthSystemFields<T = never> = {
	email: string;
	emailVisibility: boolean;
	username: string;
	verified: boolean;
} & BaseSystemFields<T>;

// Record types for each collection

export enum LinksTagsOptions {
	'thb-ext' = 'thb-ext',
	'thb-index' = 'thb-index',
	'obsolete' = 'obsolete'
}
export type LinksRecord = {
	slug: string;
	title?: string;
	url: string;
	tags?: LinksTagsOptions[];
	view?: number;
	vercelView?: number;
	author: RecordIdString;
	enabled?: boolean;
};

export enum UsersGroupsOptions {
	'tohomoeuser' = 'tohomoeuser',
	'tohomoeop' = 'tohomoeop'
}
export type UsersRecord = {
	name?: string;
	avatar?: string;
	groups?: UsersGroupsOptions[];
	avatarUrl?: string;
};

// Response types include system fields and match responses from the PocketBase API
export type LinksResponse<Texpand = unknown> = Required<LinksRecord> & BaseSystemFields<Texpand>;
export type UsersResponse<Texpand = unknown> = Required<UsersRecord> & AuthSystemFields<Texpand>;

// Types containing all Records and Responses, useful for creating typing helper functions

export type CollectionRecords = {
	links: LinksRecord;
	users: UsersRecord;
};

export type CollectionResponses = {
	links: LinksResponse;
	users: UsersResponse;
};
