package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

// ContextKey is used to store the auth context value in the request context
type ContextKey string

var contextKey ContextKey = "AUTH_CONTEXT"

// Session is used to store authentication data as a context value
type Session struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
}

// Config defines the configuration used for the authentication middleware and
// the token verifier
type Config struct {
	ProviderURL  string
	ClientID     string
	ClientSecret string
	RedirectURL  string

	Scopes  []string
	Enabled bool
}

// Provider stores the runtime configuration and verifier for the authorizer
// and the middleware
type Provider struct {
	oidc       *oidc.Provider
	oidcConfig *oauth2.Config
	Verifier   *oidc.IDTokenVerifier
	Disabled   bool
}

// NewProvider creates a new OIDC provider instance
func NewProvider(ctx context.Context, conf *Config) (*Provider, error) {
	provider, err := oidc.NewProvider(ctx, conf.ProviderURL)
	if err != nil {
		return nil, err
	}

	cf := &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		RedirectURL:  conf.RedirectURL,
		Scopes:       conf.Scopes,
		Endpoint:     provider.Endpoint(),
	}

	return &Provider{
		oidc:       provider,
		oidcConfig: cf,
		Verifier:   provider.Verifier(&oidc.Config{ClientID: cf.ClientID}),
		Disabled:   !conf.Enabled,
	}, nil
}

// AuthorizeHandler validates the authentication code from the client and
// returns a access token as a response
func (p *Provider) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	oauth2Token, err := p.oidcConfig.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("access_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusUnauthorized)
		return
	}
	idToken, err := p.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Write(data)
}

// EchoAuthorizeHandler returns a wrapped http handler for the echo router
func (p *Provider) EchoAuthorizeHandler() echo.HandlerFunc {
	return echo.WrapHandler(http.HandlerFunc(p.AuthorizeHandler))
}

// Middleware reads and verifies the bearer tokens and injects the extracted
// data to the request context
func (p *Provider) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if p.Disabled {
			session := &Session{
				UserID: "ANONYMOUS",
				Email:  "noemail",
			}
			ctx = context.WithValue(ctx, contextKey, session)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
			return
		}

		rawIDToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		idToken, err := p.Verifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID Token", http.StatusUnauthorized)
			return
		}

		session := &Session{}

		if err := idToken.Claims(&session); err != nil {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, contextKey, session)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

// EchoMiddleware returns an instance of the Middleware wrapped for Echo
func (p *Provider) EchoMiddleware() echo.MiddlewareFunc {
	return echo.WrapMiddleware(p.Middleware)
}

// GetSession is a helper function to return the session struct stored in
// the request context
func GetSession(ctx context.Context) *Session {
	session, ok := ctx.Value(contextKey).(*Session)
	if !ok {
		return nil
	}
	return session
}
