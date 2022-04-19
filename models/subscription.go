package models

import (
	"fmt"
	"net/http"
)

type SubscriptionType struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}
type SubscriptionTypeList struct {
	Subscriptions []SubscriptionType `json:"subscriptions"`
}

func (i *SubscriptionType) Bind(r *http.Request) error {
	if i.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	return nil
}
func (*SubscriptionTypeList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (*SubscriptionType) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type SubscriptionAction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
}
type SubscriptionActionList struct {
	Actions []SubscriptionAction `json:"actions"`
}

func (i *SubscriptionAction) Bind(r *http.Request) error {
	return nil
}
func (*SubscriptionActionList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (*SubscriptionAction) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type LogResponse struct {
	Success bool   `json:"success"`
	Details string `json:"details"`
	Count   int    `json:"count"`
	Limit   int    `json:"limit"`
}

type Healthcheck struct {
	Up      bool   `json:"up"`
	Details string `json:"details"`
}

func (*Healthcheck) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*LogResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type SubscriptionUserAction struct {
	GUID                  string `json:"GUID"`
	SubscriptionAccountId string `json:"subscriptionAccountId"`
	ActionType            string `json:"actionType"`
	UsageAmount           int    `json:"usageAmount"`
	Product               string `json:"product"`
}
