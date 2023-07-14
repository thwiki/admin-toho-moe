<script lang="ts">
	import '../app.postcss';
	import { page } from '$app/stores';
	import Header from '../lib/Header.svelte';
	import { Layout } from 'stwui';
</script>

<Layout class="h-full w-full">
	<Header path={$page.url.pathname} />
	<Layout.Content class="h-[calc(100%-64px)] bg-background">
		<Layout.Content.Body class="relative h-full w-full overflow-x-hidden overflow-y-auto p-4">
			<slot />
		</Layout.Content.Body>
	</Layout.Content>
</Layout>
