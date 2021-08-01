package event_reader_reducer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	apiEvent "github.com/bassbeaver/eventhouse/api/compiled/event"
	"strconv"
	"time"
)

const (
	EventReducerServiceAlias      = "EventReducerService"
	entityTypeSubscription        = "Subscription"
	eventTypeSubscriptionCreated  = "Created"
	eventTypeSubscriptionRenewed  = "Renewed"
	eventTypeSubscriptionCanceled = "Canceled"
	subscriptionStatusActive      = "active"
	subscriptionStatusCanceled    = "canceled"
	planDurationMonth             = "m"
	planDurationYear              = "y"
)

var unknownEntityTypeError = errors.New("unknown entity type")

type EventReducer struct {
	reducers map[string]map[string]func(eventObj *apiEvent.Event)
	db       *sql.DB
}

func (er *EventReducer) Reduce(eventObj *apiEvent.Event) {
	reducer, reducerIsOk := er.reducers[eventObj.GetEntityType()][eventObj.GetEventType()]
	if !reducerIsOk {
		fmt.Printf(
			"No reducer for EntityType:%s EventType:%s EventId:%s\n",
			eventObj.GetEntityType(),
			eventObj.GetEventType(),
			eventObj.GetEventId(),
		)

		return
	}

	reducer(eventObj)
}

func (er *EventReducer) getEntityLastEventId(entityType string, entityId string) (string, error) {
	entitiesTables := map[string]string{
		entityTypeSubscription: "subscription",
	}

	entityTable, entityTableIsOk := entitiesTables[entityType]
	if !entityTableIsOk {
		return "", unknownEntityTypeError
	}

	var lastEventId string
	queryError := er.db.
		QueryRow("SELECT last_event_id FROM "+entityTable+" WHERE id = ?", entityId).
		Scan(&lastEventId)

	// If no rows were returned - interpret that like "empty last event id"
	if sql.ErrNoRows == queryError {
		queryError = nil
	}

	return lastEventId, queryError
}

func (er *EventReducer) subscriptionCreatedReducer(eventObj *apiEvent.Event) {
	lastEventId, lastEventIdError := er.getEntityLastEventId(eventObj.EntityType, eventObj.EntityId)
	if nil != lastEventIdError {
		fmt.Printf("Failed to reduce event %s, failed to get last event id. Error: %s. \n", eventObj.EventId, lastEventIdError.Error())

		return
	}

	if "" != lastEventId {
		fmt.Printf("Failed to reduce event %s. Subscription %s is already in DB. \n", eventObj.EventId, eventObj.EntityId)

		return
	}

	payload := struct {
		Amount float32 `json:"amount"`
		Plan   struct {
			Level    string `json:"level"`
			Duration string `json:"duration"`
		} `json:"plan"`
		Transaction struct {
			Id     string  `json:"id"`
			Amount float32 `json:"amount"`
		} `json:"transaction"`
	}{}

	if jsonError := json.Unmarshal([]byte(eventObj.GetPayload()), &payload); nil != jsonError {
		fmt.Printf("Failed to reduce event %s. Failed to unmarshall payload, error: %s. \n", eventObj.EventId, jsonError.Error())

		return
	}

	expired, expiredError := addPlanDuration(eventObj.GetRecorded().AsTime(), payload.Plan.Duration)
	if nil != expiredError {
		fmt.Printf("Failed to reduce event %s. Failed to calculate expiration, error: %s. \n", eventObj.EventId, expiredError)

		return
	}

	tx, txError := er.db.BeginTx(context.Background(), nil)
	if nil != txError {
		fmt.Printf("Failed to reduce event %s. Failed to start new DB transaction, error: %s. \n", eventObj.EventId, txError.Error())

		return
	}

	_, dbError := tx.Exec(
		"INSERT INTO subscription(id, plan_level, plan_duration, created, expired, price, status, last_event_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		eventObj.GetEntityId(),
		payload.Plan.Level,
		payload.Plan.Duration,
		eventObj.GetRecorded().AsTime().Format("2006-01-02 15:04:05.999999"),
		expired.Format("2006-01-02 15:04:05.999999"),
		payload.Transaction.Amount,
		subscriptionStatusActive,
		eventObj.GetEventId(),
	)
	if nil != dbError {
		fmt.Printf("Failed to reduce event %s. Failed to save new Subscription to DB, error: %s. \n", eventObj.EventId, dbError.Error())
		transactionRollbackAndLog(tx)

		return
	}

	_, dbError = tx.Exec(
		"INSERT INTO transaction(id, subscription_id, created, amount) VALUES (?, ?, ?, ?)",
		payload.Transaction.Id,
		eventObj.GetEntityId(),
		eventObj.GetRecorded().AsTime().Format("2006-01-02 15:04:05.999999"),
		payload.Transaction.Amount,
	)
	if nil != dbError {
		fmt.Printf("Failed to reduce event %s. Failed to save new Transaction to DB, error: %s. \n", eventObj.EventId, dbError.Error())
		transactionRollbackAndLog(tx)

		return
	}

	transactionCommitAndLog(tx)
}

func (er *EventReducer) subscriptionRenewedReducer(eventObj *apiEvent.Event) {
	if err := er.checkLastEventId(eventObj); nil != err {
		fmt.Printf("Failed to reduce event %s, Error: %s. \n", eventObj.EventId, err.Error())

		return
	}

	payload := struct {
		Transaction struct {
			Id     string  `json:"id"`
			Amount float32 `json:"amount"`
		}
	}{}

	if jsonError := json.Unmarshal([]byte(eventObj.GetPayload()), &payload); nil != jsonError {
		fmt.Printf("Failed to reduce event %s. Failed to unmarshall payload, error: %s. \n", eventObj.EventId, jsonError.Error())

		return
	}

	tx, txError := er.db.BeginTx(context.Background(), nil)
	if nil != txError {
		fmt.Printf("Failed to reduce event %s. Failed to start new DB transaction, error: %s. \n", eventObj.EventId, txError.Error())

		return
	}

	var currentExpiredValue, planDuration string
	queryError := tx.QueryRow(
		"SELECT expired, plan_duration FROM subscription WHERE id = ? LIMIT 1",
		eventObj.GetEntityId(),
	).Scan(&currentExpiredValue, &planDuration)
	if nil != queryError {
		fmt.Printf("Failed to reduce event %s. Failed to get current expiration, error: %s. \n", eventObj.EventId, queryError)
		transactionRollbackAndLog(tx)

		return
	}

	currentExpired, parseError := time.Parse("2006-01-02 15:04:05.999999", currentExpiredValue)
	if nil != parseError {
		fmt.Printf("Failed to reduce event %s. Failed to parse current expiration: %s, error: %s. \n", currentExpiredValue, eventObj.EventId, parseError)
		transactionRollbackAndLog(tx)

		return
	}

	newExpired, expiredError := addPlanDuration(currentExpired, planDuration)
	if nil != expiredError {
		fmt.Printf("Failed to reduce event %s. Failed to calculate expiration, error: %s. \n", eventObj.EventId, expiredError)

		return
	}

	_, dbError := tx.Exec(
		"UPDATE subscription SET expired = ?, last_event_id = ? WHERE id = ?",
		newExpired.Format("2006-01-02 15:04:05.999999"),
		eventObj.GetEventId(),
		eventObj.GetEntityId(),
	)
	if nil != dbError {
		fmt.Printf("Failed to reduce event %s. Failed to update Subscription in DB, error: %s. \n", eventObj.EventId, dbError.Error())
		transactionRollbackAndLog(tx)

		return
	}

	_, dbError = tx.Exec(
		"INSERT INTO transaction(id, subscription_id, created, amount) VALUES (?, ?, ?, ?)",
		payload.Transaction.Id,
		eventObj.GetEntityId(),
		eventObj.GetRecorded().AsTime().Format("2006-01-02 15:04:05.999999"),
		payload.Transaction.Amount,
	)
	if nil != dbError {
		fmt.Printf("Failed to reduce event %s. Failed to save new Transaction to DB, error: %s. \n", eventObj.EventId, dbError.Error())
		transactionRollbackAndLog(tx)

		return
	}

	transactionCommitAndLog(tx)
}

func (er *EventReducer) subscriptionCanceledReducer(eventObj *apiEvent.Event) {
	if err := er.checkLastEventId(eventObj); nil != err {
		fmt.Printf("Failed to reduce event %s, Error: %s. \n", eventObj.EventId, err.Error())

		return
	}

	_, dbError := er.db.Exec(
		"UPDATE subscription SET status = ?, last_event_id = ? WHERE id = ?",
		subscriptionStatusCanceled,
		eventObj.GetEventId(),
		eventObj.GetEntityId(),
	)
	if nil != dbError {
		fmt.Printf("Failed to reduce event %s. Failed to cancel Subscription in DB, error: %s. \n", eventObj.EventId, dbError.Error())

		return
	}
}

func (er *EventReducer) checkLastEventId(eventObj *apiEvent.Event) error {
	lastEventId, lastEventIdError := er.getEntityLastEventId(eventObj.EntityType, eventObj.EntityId)
	if nil != lastEventIdError {
		return errors.New(
			fmt.Sprintf("Failed to get last event id. Error: %s.", lastEventIdError.Error()),
		)
	}

	if eventObj.GetPreviousEventId() != lastEventId {
		return errors.New(
			fmt.Sprintf(
				"Last event id in Projection does not match previous event id from Event. "+
					"Projection.LastEventId: %s, Event.PreviousEventId: %s.",
				lastEventId,
				eventObj.GetPreviousEventId(),
			),
		)
	}

	return nil
}

func transactionRollbackAndLog(tx *sql.Tx) {
	if err := tx.Rollback(); nil != err {
		fmt.Printf("Error during transaction rollback: %s. \n", err.Error())
	}
}

func transactionCommitAndLog(tx *sql.Tx) {
	if err := tx.Commit(); nil != err {
		fmt.Printf("Error during transaction commit: %s. \n", err.Error())
	}
}

func validatePlanDuration(planDuration string) error {
	if 2 != len(planDuration) ||
		(string(planDuration[1]) != planDurationMonth && string(planDuration[1]) != planDurationYear) {
		return errors.New("invalid plan duration: " + planDuration)
	}

	if _, durationValueError := strconv.Atoi(string(planDuration[0])); nil != durationValueError {
		return errors.New(fmt.Sprintf("invalid plan duration %s, error: %s", planDuration, durationValueError))
	}

	return nil
}

func addPlanDuration(date time.Time, planDuration string) (*time.Time, error) {
	if invalidDurationErr := validatePlanDuration(planDuration); nil != invalidDurationErr {
		return nil, invalidDurationErr
	}

	var newDate time.Time

	durationValue, _ := strconv.Atoi(string(planDuration[0]))

	switch string(planDuration[1]) {
	case planDurationYear:
		newDate = date.AddDate(durationValue, 0, 0)
	case planDurationMonth:
		newDate = date.AddDate(0, durationValue, 0)
	default:
		return &newDate, errors.New("unknown plan duration period: " + planDuration)
	}

	return &newDate, nil
}

func NewEventReducer(db *sql.DB) *EventReducer {
	r := &EventReducer{
		db: db,
	}
	r.reducers = map[string]map[string]func(eventObj *apiEvent.Event){
		entityTypeSubscription: {
			eventTypeSubscriptionCreated:  r.subscriptionCreatedReducer,
			eventTypeSubscriptionRenewed:  r.subscriptionRenewedReducer,
			eventTypeSubscriptionCanceled: r.subscriptionCanceledReducer,
		},
	}

	return r
}
