import type { LayoutLoad } from './$types';
import { pb } from '../../pb';

async function loadTags() {
	const colors = [
		'bg-red-500',
		'bg-orange-500',
		'bg-amber-500',
		'bg-yellow-500',
		'bg-lime-500',
		'bg-green-500',
		'bg-emerald-500',
		'bg-teal-500',
		'bg-cyan-500',
		'bg-sky-500',
		'bg-blue-500',
		'bg-indigo-500',
		'bg-violet-500',
		'bg-purple-500',
		'bg-fuchsia-500',
		'bg-pink-500',
		'bg-rose-500'
	];

	const options: { maxSelect: number; values: string[] } = await pb.send('/api/tags', {
		method: 'GET'
	});

	return options.values.map((value) => {
		let sum = 0;
		for (let i = 0; i < value.length; ++i) {
			sum += value.charCodeAt(i);
		}
		return {
			name: value,
			color: colors[sum % colors.length]
		};
	});
}

export const load = (async () => {
	const tags = await loadTags();
	return { tags };
}) satisfies LayoutLoad;
