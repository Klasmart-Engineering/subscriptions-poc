package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"subscriptions.demo/database"
)

var dbInstance db.Database

func NewHandler(db db.Database) http.Handler {
	router := chi.NewRouter()
	dbInstance = db
	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	router.Route("/subscription-types", subscriptionTypes)
	router.Route("/subscription-actions", subscriptionActions)
	router.Route("/log-action", logAction)
	router.Route("/healthcheck", healthcheck)
	router.Route("/liveness", liveness)
	router.Route("/evaluate-subscriptions", evaluateSubscriptions)
	return router
}
func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(405)
	render.Render(w, r, ErrMethodNotAllowed)
}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(400)
	render.Render(w, r, ErrNotFound)
}
