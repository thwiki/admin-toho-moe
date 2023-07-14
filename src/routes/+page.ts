import type { PageLoad } from './$types';
import { guard } from '../stores/auth';

export const load = (() => {
	guard();
}) satisfies PageLoad;
