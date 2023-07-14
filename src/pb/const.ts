import { env } from '$env/dynamic/public';
import PocketBase from 'pocketbase';

export const domain = env.PUBLIC_BASE_URL;
export const pb = new PocketBase(domain);
export const providerName = 'oidc';
export const redirectUrl = `${domain}/redirect.html`;
