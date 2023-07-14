import type { PageLoad } from './$types';
import { logout } from '../../stores/auth';

export const load = (() => {
	logout();
}) satisfies PageLoad;
