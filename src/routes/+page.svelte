<script lang="ts">
	import { Button } from 'stwui';
	import { isValid, user } from '../stores/auth';
</script>

<article>
	{#if $isValid && $user}
		<div>欢迎，{$user?.username ?? ''}</div>
		<div class="my-4">
			<Button type="primary" href="/link">查看短链接</Button>
		</div>
	{:else}
		<div>欢迎，请通过 THBWiki 登录</div>
	{/if}
</article>
