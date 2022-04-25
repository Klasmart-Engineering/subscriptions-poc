package db

import (
	"database/sql"
	"fmt"
	"subscriptions.demo/models"
	"time"
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

func (db Database) IsSubscriptionActive(subscriptionAccountId int) (bool, error) {
	var activated string
	if err := db.Conn.QueryRow(`
			SELECT active 
			FROM subscription_account WHERE account_holder_id = $1`, subscriptionAccountId).Scan(&activated); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("unable to check if subscription is active: %s", err)
		}
	}

	return activated == "true", nil
}

func (db Database) UpdateLastProcessed(subscription *models.SubscriptionEvaluation) {

	sqlStatement := `
UPDATE subscription_account
 SET last_processed = NOW()
WHERE id = $1;`

	db.Conn.Exec(sqlStatement, subscription.ID)

}

func (db Database) UpdateSubscriptionStatus(subscriptionId int, active bool) {

	sqlStatement := `
UPDATE subscription_account
 SET active = $1
WHERE id = $1;`

	db.Conn.Exec(sqlStatement, &active, &subscriptionId)

}

func (db Database) UsageOfSubscription(subscriptionEvaluation models.SubscriptionEvaluation) (int, error) {

	var subscriptionUsage int

	var countInteractionsSql = `
			SELECT COUNT(1) AS subscription_usage 
			FROM subscription_account_user_log 
			WHERE subscription_account_id = $1 AND product = $2 and interaction_at > $3`

	if err := db.Conn.QueryRow(countInteractionsSql,
		subscriptionEvaluation.ID, subscriptionEvaluation.Product, subscriptionEvaluation.LastProcessedTime).Scan(&subscriptionUsage); err != nil {
		if err == sql.ErrNoRows {
			return subscriptionUsage, fmt.Errorf("unknown usage on subscription: %s", subscriptionEvaluation.ID)
		}
	}
	return subscriptionUsage, nil
}

func (db Database) SubscriptionsToProcess() (*models.SubscriptionEvaluations, error) {

	list := &models.SubscriptionEvaluations{}
	rows, err := db.Conn.Query(`
		SELECT subAccount.id, subProduct.product, subProduct.threshold, subProduct.name, subAccount.last_processed
		FROM subscription_account subAccount
		JOIN subscription_account_product subProduct
		  ON subAccount.id = subProduct.subscription_id
		WHERE subAccount.last_processed IS NULL OR (now() < last_processed + ((SELECT run_frequency_minutes from subscription_account)||' minutes')::interval)`)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var subscriptionEvaluation models.SubscriptionEvaluation
		var lastProcessed sql.NullString
		err := rows.Scan(&subscriptionEvaluation.ID, &subscriptionEvaluation.Product, &subscriptionEvaluation.Threshold, &subscriptionEvaluation.Name, &lastProcessed)

		if lastProcessed.Valid {
			subscriptionEvaluation.LastProcessedTime = lastProcessed.String
		} else {
			subscriptionEvaluation.LastProcessedTime = ""
		}

		if err != nil {
			return list, err
		}

		list.SubscriptionEvaluations = append(list.SubscriptionEvaluations, subscriptionEvaluation)
	}
	return list, nil
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
			INSERT INTO subscription_account_user_log (GUID, subscription_account_id, action_type, usage, product, interaction_at)
			VALUES ($1, $2, $3, $4, $5, NOW())`)
	if es != nil {
		panic(es.Error())
	}

	_, er := stmt.Exec(userAction.GUID, userAction.SubscriptionAccountId, userAction.ActionType, userAction.UsageAmount, userAction.Product)
	if er != nil {
		panic(er.Error())
	}

}

func (db Database) CountInteractionsForSubscription(userAction models.SubscriptionUserAction) (int, error) {

	var lastProcessedTime time.Time
	if err := db.Conn.QueryRow(`
			SELECT last_processed
			FROM subscription_account
			WHERE account_holder_id = $1 `,
		userAction.SubscriptionAccountId, userAction.Product).Scan(&lastProcessedTime); err != nil {
		if err == sql.ErrNoRows {
			panic(err)
		}
	}

	var countInteractionsSql = `
			SELECT COUNT(1) AS user_interactions 
			FROM subscription_account_user_log 
			WHERE subscription_account_id = $1 AND product = $2 `
	var countUserInteractions int
	if !lastProcessedTime.IsZero() {
		countInteractionsSql = countInteractionsSql + "AND interaction_at > $3"
		if err := db.Conn.QueryRow(countInteractionsSql,
			userAction.SubscriptionAccountId, userAction.Product, lastProcessedTime).Scan(&countUserInteractions); err != nil {
			if err == sql.ErrNoRows {
				return countUserInteractions, fmt.Errorf("unknown count on user: %s", userAction.GUID)
			}
		}
	} else {
		if err := db.Conn.QueryRow(countInteractionsSql,
			userAction.SubscriptionAccountId, userAction.Product).Scan(&countUserInteractions); err != nil {
			if err == sql.ErrNoRows {
				return countUserInteractions, fmt.Errorf("unknown count on user: %s", userAction.GUID)
			}
		}
	}
	return countUserInteractions, nil
}

func (db Database) GetThresholdForSubscriptionProduct(userAction models.SubscriptionUserAction) (int, error) {

	var subscriptionThreshold int
	if err := db.Conn.QueryRow(`
			SELECT sap.threshold 
			FROM subscription_account_product sap 
			JOIN subscription_account sa
			  ON sap.subscription_id = sa.id
			WHERE sa.account_holder_id = $1 AND sap.product = $2 `,
		userAction.SubscriptionAccountId, userAction.Product).Scan(&subscriptionThreshold); err != nil {
		if err == sql.ErrNoRows {
			return subscriptionThreshold, fmt.Errorf("unknown threshold on user: %s", userAction.GUID)
		}
	}
	return subscriptionThreshold, nil
}
