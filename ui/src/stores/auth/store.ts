// import {
// 	AuthProvider,
// 	EmailAuthProvider,
// 	GithubAuthProvider,
// 	GoogleAuthProvider,
// 	IdTokenResult,
// 	OAuthProvider,
// 	onAuthStateChanged,
// 	signInWithEmailAndPassword,
// 	signInWithRedirect
// } from 'firebase/auth';

import { readable, Readable, Subscriber, Unsubscriber } from 'svelte/store';
import { browser } from '$app/env';
// import { loadAuth } from './loadAuth';

export enum AuthState {
	Authenticating = 1,
	SigningIn = 2,
	SignedOut = 3,
	SignedIn = 4
}

export enum Provider {
	UsernamePassword = 1,
	Google = 2,
	GitHub = 3,
	Microsoft = 4
}

interface SignInValues {
	provider: Provider;
	email?: string;
	password?: string;
}

export interface AuthStatus {
	user: User | null;
	known: boolean;
	state: AuthState;
}

// A type representing an authenticated user
export interface User {
	id: string;
	name: string;
	email: string;
	photo_url: string;
	token: string;
}

// createUserFromClaims creates a User from given claims and id token
// TODO: adopt Firebase's User and UserInfo types
// https://firebase.google.com/docs/reference/js/v9/auth.userinfo
const createUserFromClaims = (claims, token): User => ({
	id: claims.user_id,
	name: claims.name,
	email: claims.email,
	photo_url: claims.picture,
	token: token
});

// authStore is a readable Svelte Store used to store the client-only
// auth state of a user (SSR unsupported)
export const authStore: Readable<AuthStatus> = readable<AuthStatus>(
	// default store state
	{
		user: null,
		known: false, // known false lets us assert that the store is in its default state client-side
		state: AuthState.SignedOut
	},
	// function called on subscription.
	function start(set: Subscriber<AuthStatus>): Unsubscriber {
		if (browser) {
			// listen to auth changes _only_ on client
			listen(set);
		} else {
			// SSR disabled, so no auth on server
			set({ user: null, known: true, state: null });
		}

		// function called on unsubscription
		// we set the auth state to represent a signed out user
		return function stop() {
			set({
				user: null,
				known: true,
				state: AuthState.SignedOut
			});
		};
	}
);

// listen listens for changes to the auth state and updates
// the store (via set) accordingly
async function listen(set: Subscriber<AuthStatus>) {
	// const auth = loadAuth();
	const auth = null;

	// onAuthStateChanged(
	// 	auth,
	// 	async (authUser) => {
	// 		// if the auth state has changed with a valid user in the callback
	// 		// grab the claims and update store state for that user
	// 		if (authUser) {
	// 			const token: IdTokenResult = await authUser.getIdTokenResult();
	// 			const user: User = createUserFromClaims(token.claims, token.token);
	// 			set({
	// 				user: user,
	// 				known: true,
	// 				state: AuthState.SignedIn
	// 			});
	// 		}
	// 		// if the auth state has changed with no user in the callback,
	// 		// the user must have signed out
	// 		else {
	// 			set({
	// 				user: null,
	// 				known: true,
	// 				state: AuthState.SignedOut
	// 			});
	// 		}
	// 	},
	// 	(err) => console.error(err.message)
	// );
}

// providerFor returns a Firebase AuthProvider for the given Provider
// or throws on an unsupported Provider
// function providerFor(provider: Provider): AuthProvider {
// 	switch (provider) {
// 		case Provider.Google:
// 			return new GoogleAuthProvider();
// 		case Provider.GitHub:
// 			return new GithubAuthProvider();
// 		case Provider.Microsoft:
// 			return new OAuthProvider('microsoft.com');
// 		case Provider.UsernamePassword:
// 			return new EmailAuthProvider();

// 		default:
// 			throw 'unknown provider ' + provider;
// 	}
// }

// signInWith signs in a user with the provided SignInValues instance
export async function signInWith(values: SignInValues) {
	// const auth = loadAuth();
	// const provider = providerFor(values.provider);

	// if (values.provider === Provider.UsernamePassword) {
	// 	await signInWithEmailAndPassword(auth, values.email, values.password);
	// } else {
	// 	await signInWithRedirect(auth, provider);
	// }
}

export async function signOut() {
	// const auth = loadAuth();
	// await auth.signOut();
}
