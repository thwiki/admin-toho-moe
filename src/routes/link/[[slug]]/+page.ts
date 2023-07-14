import type { PageLoad } from './$types';
import { guard } from '../../../stores/auth';
import { pb, type LinksResponseExpaned } from '../../../pb';

function getOne(id: string) {
	return pb.collection(`links`).getOne<LinksResponseExpaned>(id, { expand: `author` });
}

export const load = (async ({ url, params }) => {
	guard();

	const perPage = 20;
	const baseUrl = url.pathname;
	const orderBy = url.searchParams.get('orderBy') || 'created';
	const order: 'asc' | 'desc' = url.searchParams.get('order') === 'asc' ? 'asc' : 'desc';
	const currentPage = parseInt(url.searchParams.get('page') ?? '1', 10) || 1;
	const slug = params.slug;

	const list = await pb.collection('links').getList<LinksResponseExpaned>(currentPage, perPage, {
		sort: `${order === 'asc' ? '' : '-'}${orderBy}`
	});

	let isNew = false;
	let item: Partial<LinksResponseExpaned> | undefined = undefined;

	if (slug == null || slug === '_') {
		isNew = true;
		item = {
			slug: ``,
			title: ``,
			url: ``,
			tags: [],
			enabled: true
		};
	} else {
		isNew = false;
		try {
			item = await getOne(slug);
		} catch (_: unknown) {
			item = undefined;
		}
	}

	return {
		perPage,
		baseUrl,
		orderBy,
		order,
		currentPage,
		slug,
		list,
		isNew,
		item
	};
}) satisfies PageLoad;
