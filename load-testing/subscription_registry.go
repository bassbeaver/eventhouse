package main

import (
	"math/rand"
	"sync"
	"time"
)

type SubscriptionRegistry struct {
	freeSequences []int // list of indexes of free elements from sequences slice
	sequences     []*SubscriptionSequence
	mtx           sync.Mutex
}

func (sr *SubscriptionRegistry) GetSequence() *SubscriptionSequence {
	createNewSequence := func() *SubscriptionSequence {
		newSequenceIndex := len(sr.sequences)
		newSequence := CreateNewSubscriptionSequence(newSequenceIndex)
		sr.sequences = append(sr.sequences, newSequence)

		return newSequence
	}

	selectFreeSequence := func() *SubscriptionSequence {
		selectedFreeIndex := rand.Intn(len(sr.freeSequences))
		selectedSequenceIndex := sr.freeSequences[selectedFreeIndex]
		sr.freeSequences = append(sr.freeSequences[:selectedFreeIndex], sr.freeSequences[selectedFreeIndex+1:]...)

		return sr.sequences[selectedSequenceIndex]
	}

	var selectedSequence *SubscriptionSequence

	sr.mtx.Lock()

	// If no free sequences - create one
	if 0 == len(sr.freeSequences) {
		selectedSequence = createNewSequence()
	} else {
		// If free sequences exists - 20% chance to create new sequence
		if 4 == rand.Intn(5) {
			selectedSequence = createNewSequence()
		} else {
			selectedSequence = selectFreeSequence()
		}
	}

	sr.mtx.Unlock()

	return selectedSequence
}

func (sr *SubscriptionRegistry) FreeSequence(s *SubscriptionSequence) {
	// If sequence has active subscription - wait some time to emulate
	// real business time between events and return it to free pool
	if s.IsActive() {
		timerObj := time.NewTimer(2 * time.Second)
		go func() {
			<-timerObj.C

			sr.mtx.Lock()
			sr.freeSequences = append(sr.freeSequences, s.id)
			sr.mtx.Unlock()
		}()
	}
}

func NewSubscriptionRegistry() *SubscriptionRegistry {
	return &SubscriptionRegistry{
		freeSequences: make([]int, 0),
		sequences:     make([]*SubscriptionSequence, 0),
	}
}
