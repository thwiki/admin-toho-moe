<script lang="ts">
	import { Breadcrumbs, Layout, Button } from 'stwui';

	export let path: string;

	interface Crumb {
		label: string;
		href: string;
	}

	let crumbs: Crumb[] = [];

	$: {
		crumbs = (path ?? ``)
			.split(`/`)
			.filter((segment) => segment !== ``)
			.map((segment, index, arr) => ({
				label: segment,
				href: `/` + arr.slice(0, index + 1).join(`/`)
			}));
	}
</script>

<Layout.Header class="py-2 px-4">
	<Breadcrumbs>
		<Breadcrumbs.Crumb href={'/'}>
			<Breadcrumbs.Crumb.Label slot="label">ToHoMoe</Breadcrumbs.Crumb.Label>
		</Breadcrumbs.Crumb>
		{#each crumbs as crumb}
			<Breadcrumbs.Crumb href={crumb.href}>
				<Breadcrumbs.Crumb.Label slot="label">{crumb.label}</Breadcrumbs.Crumb.Label>
			</Breadcrumbs.Crumb>
		{/each}
	</Breadcrumbs>
	<Layout.Header.Extra slot="extra">
		<Button class="ml-auto" href="/logout">登出</Button>
	</Layout.Header.Extra>
</Layout.Header>
