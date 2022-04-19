package handler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log"
	"net/http"
	models "subscriptions.demo/models"
)

func healthcheck(router chi.Router) {
	router.Get("/", dbHealthcheck)
}

func liveness(router chi.Router) {
	router.Get("/", applicationLiveness)
}

func subscriptionTypes(router chi.Router) {
	router.Get("/", getAllSubscriptionTypes)
}

func subscriptionActions(router chi.Router) {
	router.Get("/", getAllSubscriptionActions)
}

func logAction(router chi.Router) {
	router.Post("/", logUserAction)
}

func evaluateSubscriptions(router chi.Router) {
	router.Post("/", evaluateSubscriptionsUsage)
}

func evaluateSubscriptionsUsage(w http.ResponseWriter, r *http.Request) {
	log.Println("EVALUATING USAGE")
}

func dbHealthcheck(w http.ResponseWriter, r *http.Request) {
	up, err := dbInstance.Healthcheck()

	if err != nil || !up {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}

	health := models.Healthcheck{Up: true, Details: "Successfully connected to the database"}
	if err := render.Render(w, r, &health); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func applicationLiveness(w http.ResponseWriter, r *http.Request) {

	health := models.Healthcheck{Up: true, Details: "Application up"}
	if err := render.Render(w, r, &health); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func getAllSubscriptionTypes(w http.ResponseWriter, r *http.Request) {
	subscriptionTypes, err := dbInstance.GetAllSubscriptions()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, subscriptionTypes); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func getAllSubscriptionActions(w http.ResponseWriter, r *http.Request) {
	subscriptionActions, err := dbInstance.GetAllSubscriptionActions()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, subscriptionActions); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func logUserAction(w http.ResponseWriter, r *http.Request) {
	var subscriptionUserAction models.SubscriptionUserAction
	json.NewDecoder(r.Body).Decode(&subscriptionUserAction)

	dbInstance.LogUserAction(subscriptionUserAction)
	userInteractions, err := dbInstance.CountUserInteractionsForSubscription(subscriptionUserAction)
	if err != nil {
		panic(err)
	}

	var response models.LogResponse
	threshold, er := dbInstance.GetThresholdForSubscription(subscriptionUserAction)
	if er != nil {
		panic(er)
	}

	if userInteractions > threshold {
		response = models.LogResponse{Success: false, Details: "BLOCKED", Count: userInteractions, Limit: threshold}
	} else {
		response = models.LogResponse{Success: true, Details: "WITHIN LIMITS", Count: userInteractions, Limit: threshold}
	}

	if err := render.Render(w, r, &response); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}
