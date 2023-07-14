import type { PageLoad } from './$types';
import { callback } from '../../stores/auth';

export const load = (() => {
	callback();
}) satisfies PageLoad;
