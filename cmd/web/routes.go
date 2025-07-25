package main

import (
	"io/fs"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"snippetbox.gregor-pifko/ui"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Static resources file server
	staticFS, _ := fs.Sub(ui.Files, "static")
	router.ServeFiles("/static/*filepath", http.FS(staticFS))

	// Middleware for dynamic (session-enabled) routes
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Middleware for authorized (logged in) routes
	protected := dynamic.Append(app.requireAuthentication)

	// Ping route for status-checking or uptime monitoring of server
	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// Unprotected routes
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected (authenticated-only) routes
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// Common middleware (used by all routes)
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
