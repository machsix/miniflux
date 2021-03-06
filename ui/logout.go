// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/cookie"
	"miniflux.app/http/request"
	"miniflux.app/http/response/html"
	"miniflux.app/http/route"
	"miniflux.app/logger"
	"miniflux.app/ui/session"
)

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	sess := session.New(h.store, request.SessionID(r))
	user, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	sess.SetLanguage(user.Language)
	sess.SetTheme(user.Theme)
	sess.SetView(user.View)

	if err := h.store.RemoveUserSessionByToken(user.ID, request.UserSessionToken(r)); err != nil {
		logger.Error("[UI:Logout] %v", err)
	}

	http.SetCookie(w, cookie.Expired(
		cookie.CookieUserSessionID,
		h.cfg.IsHTTPS,
		h.cfg.BasePath(),
	))

	html.Redirect(w, r, route.Path(h.router, "login"))
}
