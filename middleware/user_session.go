// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"

	"github.com/miniflux/miniflux/http/cookie"
	"github.com/miniflux/miniflux/http/route"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"

	"github.com/gorilla/mux"
)

// UserSession handles the user session middleware.
func (m *Middleware) UserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := m.getSessionFromCookie(r)

		if session == nil {
			logger.Debug("[Middleware:UserSession] Session not found")
			if m.isPublicRoute(r) {
				next.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, route.Path(m.router, "login"), http.StatusFound)
			}
		} else {
			logger.Debug("[Middleware:UserSession] %s", session)
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDContextKey, session.UserID)
			ctx = context.WithValue(ctx, IsAuthenticatedContextKey, true)
			ctx = context.WithValue(ctx, UserSessionTokenContextKey, session.Token)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (m *Middleware) isPublicRoute(r *http.Request) bool {
	route := mux.CurrentRoute(r)
	switch route.GetName() {
	case "login",
		"checkLogin",
		"stylesheet",
		"javascript",
		"oauth2Redirect",
		"oauth2Callback",
		"appIcon",
		"favicon",
		"webManifest":
		return true
	default:
		return false
	}
}

func (m *Middleware) getSessionFromCookie(r *http.Request) *model.UserSession {
	sessionCookie, err := r.Cookie(cookie.CookieUserSessionID)
	if err == http.ErrNoCookie {
		return nil
	}

	session, err := m.store.UserSessionByToken(sessionCookie.Value)
	if err != nil {
		logger.Error("[Middleware:UserSession] %v", err)
		return nil
	}

	return session
}