<script lang="ts">
	import {
		Badge,
		Button,
		CheckboxGroup,
		Drawer,
		Input,
		Pagination,
		Portal,
		Select,
		Table
	} from 'stwui';
	import type { PageData } from './$types';
	import { goto } from '$app/navigation';
	import { encodeSearchParams } from 'stwui/utils';
	import type { TableColumn } from 'stwui/types';
	import {
		pb,
		type LinksResponse,
		type UsersResponse,
		type LinksResponseExpaned
	} from '../../../pb';
	import { user } from '../../../stores/auth';
	import type { ListResult } from 'pocketbase';

	export let data: PageData;

	let perPage: number;
	let baseUrl: string;
	let orderBy: string;
	let order: 'asc' | 'desc';
	let currentPage: number;
	let slug: string | undefined;
	let list: ListResult<
		LinksResponse<{
			author: UsersResponse;
		}>
	>;
	let isNew = true;
	let item: Partial<LinksResponseExpaned> | undefined;
	let tags: {
		name: string;
		color: string;
	}[] = [];

	$: {
		perPage = data.perPage;
		baseUrl = data.baseUrl;
		orderBy = data.orderBy;
		order = data.order;
		currentPage = data.currentPage;
		slug = data.slug;
		tags = data.tags;
		list = data.list;
		isNew = data.isNew;
		item = data.item;
	}

	const columns: TableColumn[] = [
		{
			column: 'slug',
			label: '短语',
			placement: 'left',
			class: 'w-[20%]'
		},
		{
			column: 'title',
			label: '标题',
			placement: 'left',
			class: 'w-[30%]'
		},
		{
			column: 'tags',
			label: '标签',
			placement: 'left',
			class: 'w-[30%]'
		},
		{
			column: 'view',
			label: '浏览',
			placement: 'right',
			class: 'w-[10%]'
		},
		{
			column: 'enabled',
			label: '启用',
			placement: 'left',
			class: 'w-[10%]'
		}
	];

	function onPreviousClick() {
		const newPage = Math.max(currentPage - 1, 1);
		if (newPage !== currentPage) {
			goto(
				`${baseUrl}` +
					encodeSearchParams({
						orderBy,
						order,
						page: `${newPage}`
					})
			);
		}
	}
	function onNextClick() {
		const newPage = currentPage + 1;
		if (newPage !== currentPage) {
			goto(
				`${baseUrl}` +
					encodeSearchParams({
						orderBy,
						order,
						page: `${newPage}`
					})
			);
		}
	}
	function onPageClick(page: number) {
		const newPage = page;
		if (newPage !== currentPage) {
			goto(
				`${baseUrl}` +
					encodeSearchParams({
						orderBy,
						order,
						page: `${newPage}`
					})
			);
		}
	}
	function onColumnHeaderClick(column: string) {
		goto(
			`${baseUrl}` +
				encodeSearchParams({
					orderBy: column,
					order: column === orderBy && order === 'asc' ? 'desc' : 'asc',
					page: `${currentPage}`
				})
		);
	}

	async function upsertOne(data: Partial<LinksResponseExpaned>): Promise<LinksResponseExpaned> {
		if (data.id) {
			return await pb.collection(`links`).update(
				data.id,
				{
					slug: data.slug,
					title: data.title,
					url: data.url,
					tags: data.tags,
					author: data.author,
					enabled: data.enabled
				},
				{ expand: `author` }
			);
		} else {
			return await pb.collection(`links`).create(data, { expand: `author` });
		}
	}

	async function onSubmit(event: SubmitEvent) {
		if (item) {
			const data: Partial<LinksResponseExpaned> = { id: item.id };
			const formData = new FormData(event.target as HTMLFormElement);
			formData.forEach((value, key) => {
				if (value == null || value === '') return;

				if (key === 'tags' && typeof value === 'string') {
					value = JSON.parse(value).map((t: { value: string }) => t.value);
				}
				Reflect.set(data, key, value);
			});
			data.enabled = !!data.enabled;
			if (!data.author) data.author = $user?.id;

			item = await upsertOne(data);
			goto(`/link/${item.id}${location.search}`, { invalidateAll: true });
		}
	}

	function onClickPush() {
		return pb.send('/api/push', { method: 'POST' });
	}
	function onClickView() {
		return pb.send('/api/view', { method: 'POST' });
	}
</script>

<article>
	<div class="flex mb-4">
		<Button type="primary" href={`/link/_${location.search}`}>创建</Button>
		<Button class="ml-auto" type="primary" on:click={onClickView}>更新浏览</Button>
		<Button class="ml-2" type="danger" on:click={onClickPush}>推送更新</Button>
	</div>

	<Portal>
		{#if slug}
			<Drawer
				handleClose={() => {
					goto(`/link${location.search}`);
				}}
			>
				<Drawer.Header slot="header">
					<span class="text-lg"
						>{#if isNew}创建短链接{:else}编辑短链接： {item?.slug ?? ``}{/if}</span
					>
				</Drawer.Header>
				<Drawer.Content slot="content">
					{#if item}
						{@const itemTag = (item.tags ?? []).map((tag) => ({ value: tag }))}
						<form class="flex flex-col space-y-2" on:submit|preventDefault={onSubmit}>
							<Input name="slug" value={item.slug}>
								<Input.Label slot="label">短语</Input.Label>
							</Input>
							<Input name="title" value={item.title}>
								<Input.Label slot="label">标题</Input.Label>
							</Input>
							<Input name="url" value={item.url}>
								<Input.Label slot="label">网址</Input.Label>
							</Input>

							<Select name="tags" placeholder="标签" value={itemTag} multiple>
								<Select.Label slot="label">标签</Select.Label>
								<Select.Options slot="options">
									{#each tags as tag}
										<Select.Options.Option
											class={tag.color}
											option={{ value: tag.name, label: tag.name }}
										/>
									{/each}
								</Select.Options>
							</Select>

							{#if item.expand?.author}
								<Input name="" value={item.expand?.author.username} readonly>
									<Input.Label slot="label">作者</Input.Label>
								</Input>
							{/if}

							<CheckboxGroup>
								<CheckboxGroup.Checkbox name="enabled" value="enabled" checked={item.enabled}>
									<CheckboxGroup.Checkbox.Label slot="label">启用</CheckboxGroup.Checkbox.Label>
								</CheckboxGroup.Checkbox>
							</CheckboxGroup>

							<Button type="primary" htmlType="submit">保存</Button>
						</form>
					{:else}
						短链接不存在
					{/if}
				</Drawer.Content>
			</Drawer>
		{/if}
	</Portal>

	<Table class="rounded-md overflow-hidden h-full" {columns}>
		<Table.Header slot="header" {order} {orderBy} {onColumnHeaderClick} />
		<Table.Body slot="body">
			{#each list.items as link}
				<Table.Body.Row
					id={link.id}
					on:click={() => {
						goto(`/link/${link.id}${location.search}`);
					}}
				>
					<Table.Body.Row.Cell column={0}>{link.slug}</Table.Body.Row.Cell>
					<Table.Body.Row.Cell column={1}>{link.title}</Table.Body.Row.Cell>
					<Table.Body.Row.Cell column={2}>
						<div class="flex space-x-1">
							{#each link.tags as name}
								{@const tag = tags.find((t) => t.name === name)}
								{#if tag}
									<Badge class={tag.color}>{tag.name}</Badge>
								{/if}
							{/each}
						</div>
					</Table.Body.Row.Cell>
					<Table.Body.Row.Cell column={3}>{link.view + link.vercelView}</Table.Body.Row.Cell>
					<Table.Body.Row.Cell column={4}>{link.enabled}</Table.Body.Row.Cell>
				</Table.Body.Row>
			{/each}
		</Table.Body>
		<Table.Footer slot="footer" class="shadow-none">
			<Pagination
				start={(list.page - 1) * perPage + 1}
				end={(list.page - 1) * perPage + list.items.length}
				total={list.totalItems}
				pageSize={perPage}
				{currentPage}
				{onPreviousClick}
				{onNextClick}
				{onPageClick}
			/>
		</Table.Footer>
	</Table>
</article>
