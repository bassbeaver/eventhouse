package main

import (
	"encoding/json"
	"fmt"
	"github.com/bassbeaver/eventhouse/api/compiled/event"
	"math/rand"
	"time"
)

const (
	subscriptionEntityType        = "Subscription"
	subscriptionEventTypeCreated  = "Created"
	subscriptionEventTypeRenewed  = "Renewed"
	subscriptionEventTypeCanceled = "Canceled"
	subscriptionStatusActive      = "active"
	subscriptionStatusCanceled    = "canceled"
)

type plan struct {
	level    string
	amount   float32
	duration string
}

type eventPayloadTransaction struct {
	Id     string  `json:"id"`
	Amount float32 `json:"amount"`
}

type eventPayloadPlan struct {
	Level    string `json:"level"`
	Duration string `json:"duration"`
}

type SubscriptionSequence struct {
	id                int
	transactionsCount int
	plan              plan
	status            string
	events            []*event.PushRequest
}

func (s *SubscriptionSequence) GetEntityId() string {
	return fmt.Sprintf("%s%d", subscriptionEntityType, s.id)
}

func (s *SubscriptionSequence) GetLastEvent() *event.PushRequest {
	return s.events[len(s.events)-1]
}

func (s *SubscriptionSequence) IsActive() bool {
	return subscriptionStatusActive == s.status
}

func (s *SubscriptionSequence) IsCanceled() bool {
	return subscriptionStatusCanceled == s.status
}

func (s *SubscriptionSequence) IsNew() bool {
	return "" == s.status
}

func (s *SubscriptionSequence) GenerateNextEvent() bool {
	// If subscription new initialized - create it
	if nil == s.events {
		s.status = subscriptionStatusActive
		s.events = []*event.PushRequest{
			newCreatedPushRequest(s),
		}
		s.transactionsCount++

		return true
	}

	if subscriptionStatusActive != s.status {
		return false
	}

	closeDice := rand.Intn(100)
	if closeDice >= 85 { // 15% chance to close subscription
		s.events = append(s.events, newCanceledPushRequest(s))
		s.status = subscriptionStatusCanceled
	} else { // Renew subscription
		s.events = append(s.events, newRenewalPushRequest(s))
		s.transactionsCount++
	}

	return true
}

func CreateNewSubscriptionSequence(id int) *SubscriptionSequence {
	return &SubscriptionSequence{
		id:   id,
		plan: plansPool[rand.Intn(len(plansPool))],
	}
}

func newCreatedPushRequest(s *SubscriptionSequence) *event.PushRequest {
	payload := struct {
		Amount      float32                 `json:"amount"`
		Plan        eventPayloadPlan        `json:"plan"`
		Transaction eventPayloadTransaction `json:"transaction"`
	}{
		Amount: s.plan.amount,
		Plan: eventPayloadPlan{
			Level:    s.plan.level,
			Duration: s.plan.duration,
		},
		Transaction: eventPayloadTransaction{
			Id:     fmt.Sprintf("%s_transaction%d", s.GetEntityId(), s.transactionsCount),
			Amount: s.plan.amount,
		},
	}

	serializedPayload, jsonError := json.Marshal(payload)
	if nil != jsonError {
		panic("failed to serialize to JSON: " + jsonError.Error())
	}

	result := &event.PushRequest{
		IdempotencyKey: fmt.Sprintf("%s:%d", s.GetEntityId(), time.Now().UnixNano()),
		EntityType:     subscriptionEntityType,
		EntityId:       s.GetEntityId(),
		EventType:      subscriptionEventTypeCreated,
		Payload:        string(serializedPayload),
	}

	return result
}

func newRenewalPushRequest(s *SubscriptionSequence) *event.PushRequest {
	payload := struct {
		Transaction eventPayloadTransaction
	}{
		Transaction: eventPayloadTransaction{
			Id:     fmt.Sprintf("%s_transaction%d", s.GetEntityId(), s.transactionsCount),
			Amount: s.plan.amount,
		},
	}

	serializedPayload, jsonError := json.Marshal(payload)
	if nil != jsonError {
		panic("failed to serialize to JSON: " + jsonError.Error())
	}

	result := &event.PushRequest{
		IdempotencyKey: fmt.Sprintf("%s:%d", s.GetEntityId(), time.Now().UnixNano()),
		EntityType:     subscriptionEntityType,
		EntityId:       s.GetEntityId(),
		EventType:      subscriptionEventTypeRenewed,
		Payload:        string(serializedPayload),
	}

	return result
}

func newCanceledPushRequest(s *SubscriptionSequence) *event.PushRequest {
	result := &event.PushRequest{
		IdempotencyKey: fmt.Sprintf("%s:%d", s.GetEntityId(), time.Now().UnixNano()),
		EntityType:     subscriptionEntityType,
		EntityId:       s.GetEntityId(),
		EventType:      subscriptionEventTypeCanceled,
	}

	return result
}

var plansPool = []plan{
	{level: "basic", amount: 10, duration: "1m"},
	{level: "basic", amount: 100, duration: "1y"},
	{level: "plus", amount: 15.5, duration: "1m"},
	{level: "plus", amount: 155, duration: "1y"},
	{level: "premium", amount: 20, duration: "1m"},
	{level: "premium", amount: 200, duration: "1y"},
}
