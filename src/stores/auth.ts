import { goto } from '$app/navigation';
import { error } from '@sveltejs/kit';
import { pb, providerName, redirectUrl, type UsersResponse } from '../pb';
import { writable } from 'svelte/store';

export const isValid = writable(pb.authStore.isValid);
export const user = writable(pb.authStore.model);

pb.authStore.onChange((_, model) => {
	user.set(model);
	isValid.set(pb.authStore.isValid);
});

export const guard = async () => {
	if (!pb.authStore.isValid) {
		const authMethods = await pb.collection('users').listAuthMethods();

		for (const provider of authMethods.authProviders) {
			if (provider.name === providerName) {
				const authUrl = new URL(provider.authUrl);
				authUrl.searchParams.delete('scope');
				authUrl.searchParams.set('redirect_uri', redirectUrl);

				localStorage.setItem('back', location.pathname);
				localStorage.setItem('provider', JSON.stringify(provider));

				location.href = authUrl.href;
				return;
			}
		}
	}

	if (
		!Array.isArray(pb.authStore.model?.groups) ||
		(!pb.authStore.model?.groups.includes('tohomoeop') &&
			!pb.authStore.model?.groups.includes('tohomoeuser'))
	) {
		throw error(403, '无效用户组');
	}
};

export const callback = async () => {
	const params = new URL(location.href).searchParams;
	const back = localStorage.getItem('back') ?? '/';
	const provider = JSON.parse(localStorage.getItem('provider') ?? '{}') ?? {};

	// compare the redirect's state param and the stored provider's one
	if (provider.state !== params.get('state')) {
		throw error(403, "State parameters don't match.");
	}

	// authenticate
	await pb.collection('users').authWithOAuth2Code<UsersResponse>(
		provider.name,
		params.get('code') ?? '',
		provider.codeVerifier,
		redirectUrl,
		// pass optional user create data
		{
			emailVisibility: false
		}
	);

	localStorage.removeItem('back');
	localStorage.removeItem('provider');

	goto(back);
};

export const logout = async () => {
	pb.authStore.clear();
};
