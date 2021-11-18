import { dev } from '$app/env';

// If dev mode, hard code the localhost URL. Else use empty string and paths with
// default to "/"
export const bubblyAPI = dev ? "http://localhost:8111" : ""
