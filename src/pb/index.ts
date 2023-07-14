import type { LinksResponse, UsersResponse } from './types';

export * from './const';
export * from './types';

export type LinksResponseExpaned = LinksResponse<{ author: UsersResponse }>;
