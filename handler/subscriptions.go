package handler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	redpanda "subscriptions.demo/messaging"
	models "subscriptions.demo/models"
	"time"
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
	subscriptions, err := dbInstance.SubscriptionsToProcess()

	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	// Look at subscription account on last processed and next run
	// if last processed is NULL OR last_processed time + run frequency > NOW()
	// get the subscription ID, product, threshold, name, last_processed_Time for each

	for _, subscriptionToEvaluate := range subscriptions.SubscriptionEvaluations {
		usageAmount, err := dbInstance.UsageOfSubscription(subscriptionToEvaluate)

		if err != nil {
			panic(err)
		}

		now := time.Now()
		var prod []models.EvaluatedSubscriptionProduct
		prod = append(prod, models.EvaluatedSubscriptionProduct{Name: subscriptionToEvaluate.Product, Type: subscriptionToEvaluate.Name, UsageAmount: usageAmount})

		var evaluatedSubscription = models.EvaluatedSubscription{SubscriptionId: subscriptionToEvaluate.ID, Products: prod, DateFromEpoch: subscriptionToEvaluate.LastProcessedTime, DateToEpoch: strconv.FormatInt(now.Unix(), 10)}
		marshalledEvalutedSubscription, err := json.Marshal(evaluatedSubscription)
		redpanda.Publish(string(marshalledEvalutedSubscription))
		dbInstance.UpdateLastProcessed(&subscriptionToEvaluate)
	}

	// look in log for each since last_processed_time ands tally that up
	// also use NOW on that query. Return that and also set that on the JSON response
	// send message to Kafka topic and then update last_processed.
	// do this as part of a transaction
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
	interactions, err := dbInstance.CountInteractionsForSubscription(subscriptionUserAction)
	if err != nil {
		panic(err)
	}

	var response models.LogResponse
	threshold, er := dbInstance.GetThresholdForSubscriptionProduct(subscriptionUserAction)
	if er != nil {
		panic(er)
	}

	active, err := dbInstance.IsSubscriptionActive(subscriptionUserAction.SubscriptionAccountId)

	if err != nil {
		panic(err)
	}

	if !active {
		var inactiveResponse = models.LogResponse{Success: false, Details: "BLOCKED. Subscription not active", Count: interactions, Limit: threshold}
		if err := render.Render(w, r, &inactiveResponse); err != nil {
			render.Render(w, r, ErrorRenderer(err))
		}
	}

	if interactions > threshold {
		response = models.LogResponse{Success: false, Details: "BLOCKED", Count: interactions, Limit: threshold}
	} else {
		response = models.LogResponse{Success: true, Details: "WITHIN LIMITS", Count: interactions, Limit: threshold}
	}

	if err := render.Render(w, r, &response); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}
