package db

import (
	"database/sql"
	"fmt"
	"subscriptions.demo/models"
)

func (db Database) Healthcheck() (bool, error) {
	var up int
	if err := db.Conn.QueryRow(`
			SELECT 1 AS up 
			FROM subscription_account`).Scan(&up); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("unable to get connection to the database: %s", err)
		}
	}

	return up == 1, nil
}

func (db Database) GetAllSubscriptions() (*models.SubscriptionTypeList, error) {
	list := &models.SubscriptionTypeList{}
	rows, err := db.Conn.Query("SELECT * FROM subscription_type ORDER BY id DESC")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var subscription models.SubscriptionType
		err := rows.Scan(&subscription.ID, &subscription.Name, &subscription.UpdatedAt, &subscription.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Subscriptions = append(list.Subscriptions, subscription)
	}
	return list, nil
}

func (db Database) GetAllSubscriptionActions() (*models.SubscriptionActionList, error) {
	list := &models.SubscriptionActionList{}
	rows, err := db.Conn.Query("SELECT * FROM subscription_action")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var action models.SubscriptionAction
		err := rows.Scan(&action.Name, &action.Description, &action.Unit)
		if err != nil {
			return list, err
		}
		list.Actions = append(list.Actions, action)
	}
	return list, nil
}

func (db Database) LogUserAction(userAction models.SubscriptionUserAction) {
	stmt, es := db.Conn.Prepare(`
			INSERT INTO subscription_account_user_log (GUID, subscription_account_id, action_type, interaction_at)
			VALUES ($1, $2, $3, NOW())`)
	if es != nil {
		panic(es.Error())
	}

	_, er := stmt.Exec(userAction.GUID, userAction.SubscriptionAccountId, userAction.ActionType)
	if er != nil {
		panic(er.Error())
	}

}

func (db Database) CountUserInteractionsForSubscription(userAction models.SubscriptionUserAction) (int, error) {

	var countUserInteractions int
	if err := db.Conn.QueryRow(`
			SELECT COUNT(1) AS user_interactions 
			FROM subscription_account_user_log 
			WHERE GUID = $1`,
		userAction.GUID).Scan(&countUserInteractions); err != nil {
		if err == sql.ErrNoRows {
			return countUserInteractions, fmt.Errorf("unknown count on user: %s", userAction.GUID)
		}
	}
	return countUserInteractions, nil

}

func (db Database) GetThresholdForSubscription(userAction models.SubscriptionUserAction) (int, error) {

	var subscriptionThreshold int
	if err := db.Conn.QueryRow(`
			SELECT threshold 
			FROM subscription_account 
			WHERE account_holder_id = $1`,
		userAction.SubscriptionAccountId).Scan(&subscriptionThreshold); err != nil {
		if err == sql.ErrNoRows {
			return subscriptionThreshold, fmt.Errorf("unknown threshold on user: %s", userAction.GUID)
		}
	}
	return subscriptionThreshold, nil
}
