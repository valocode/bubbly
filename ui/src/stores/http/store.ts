import { bubblyAPI } from './env';
import axios from 'axios';
import type { AxiosResponse, AxiosError } from 'axios';
import { writable } from 'svelte/store';
import type { User } from '$stores/auth/store';

interface HTTPError {
	name: string;
	message: string;
	// Whether a response was received
	response: boolean;
	// The status code of the response
	responseStatusCode?: number;
	// Whether a request was made
	request: boolean;
}

interface HTTPStore<Data = any> {
	data?: Data;
	error?: HTTPError;
	// statusCode: number;
	fetching: boolean;
}

// the httpStore takes a user from the auth store ($authStore.user)
// and performs some HTTP request to bubbly's API endpoint on behalf of the user
export function httpStore<Data = any>(user: User) {
	const store = writable<HTTPStore<Data>>({ data: undefined, error: undefined, fetching: false });

	async function request(method: string, path: string, params: object = null, data: object = null) {

		// Clear the store as we are about to make a new request
		store.update((value) => {
			value.fetching = true;
			value.error = undefined;
			// value.statusCode = undefined
			value.data = undefined;
			return value;
		});

		// TODO: figure this out
		const baseURL = bubblyAPI + "/api/v1";

		let headers = {}
		if (user !== null) {
			headers = {
				'Content-type': 'application/json',
				Authorization: user.token
			};
		}

		return new Promise<Data>((resolve, reject) => {
			axios
				.request<Data>({
					method: method === 'POST' ? 'post' : 'get',
					url: baseURL + path,
					headers,
					params,
					data
				})
				.then((resp: AxiosResponse) => {
					if (resp.data === null) {
						// Then something did not go well...
						store.update((value) => {
							value.fetching = false;
							value.error = new Error('received null') as HTTPError;
							return value;
						});
					}
					store.update((value) => {
						value.fetching = false;
						value.data = resp.data;
						// value.statusCode = resp.status
						return value;
					});
					resolve(resp.data);
				})
				.catch((reason: AxiosError) => {
					store.update((value) => {
						value.fetching = false;
						value.error = axiosErrorToHTTPError(reason);
						return value;
					});
					reject(reason);
				});
		});
	}
	return {
		...store,
		// request: (method: string, path: string) => request(method, path),
		get: (path: string, params: object = null, data: object = null) =>
			request('GET', path, params, data),
		post: (path: string, params: object, data: object = null) => request('POST', path, params, data)
	};
}

function axiosErrorToHTTPError(error: AxiosError): HTTPError {
	error.request;
	console.log('AxiosError Code: ', error.code);
	return {
		name: error.name,
		message: error.message,
		response: error.response !== undefined,
		responseStatusCode: error.response ? error.response.status : undefined,
		request: error.request !== undefined
	};
}
