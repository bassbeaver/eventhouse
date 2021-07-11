package main

import (
	"context"
	"fmt"
	"github.com/bojand/ghz/runner"
	"sync"
)

type Entity struct {
	EntityId      string
	EntityType    string
	EventsCounter *MutexCounter
}

// --------------------------

type RunAct struct {
	Report *runner.Report
	Err    error
}

// --------------------------

type ConcurrentSlice interface {
	Append(newEntry interface{}) ConcurrentSlice
	Get(i int) (interface{}, bool)
	Len() int
}

// --------------------------

type SelfCleaningConcurrentSlice struct {
	slice              []interface{}
	writeChan          chan interface{}
	totalEntitiesCount int
	maxEntities        int
}

func (s *SelfCleaningConcurrentSlice) Append(newEntry interface{}) ConcurrentSlice {
	s.writeChan <- newEntry

	return s
}

func (s *SelfCleaningConcurrentSlice) Get(i int) (interface{}, bool) {
	if s.Len() <= i {
		return nil, false
	}

	return s.slice[i], true
}

func (s *SelfCleaningConcurrentSlice) Len() int {
	return len(s.slice)
}

func (s *SelfCleaningConcurrentSlice) Total() int {
	return s.totalEntitiesCount
}

func (s *SelfCleaningConcurrentSlice) startAppendListener() {
	go func() {
		for newEntry := range s.writeChan {
			if s.Len() >= s.maxEntities {
				from := int(s.Len() / 2)
				s.slice = s.slice[from:]
				fmt.Println("SelfCleaningConcurrentSlice truncated")
			}

			s.slice = append(s.slice, newEntry)
			s.totalEntitiesCount++
		}
	}()
}

func NewSelfCleaningConcurrentSlice(maxElements int) *SelfCleaningConcurrentSlice {
	s := &SelfCleaningConcurrentSlice{
		slice:       make([]interface{}, 0),
		writeChan:   make(chan interface{}),
		maxEntities: maxElements,
	}

	s.startAppendListener()

	return s
}

// --------------------------

type MutexCounter struct {
	counter  int
	mutexObj sync.RWMutex
}

func (mc *MutexCounter) Current() int {
	mc.mutexObj.RLock()
	cnt := mc.counter
	mc.mutexObj.RUnlock()

	return cnt
}

func (mc *MutexCounter) Increment() int {
	mc.mutexObj.RLock()
	mc.counter++
	cnt := mc.counter
	mc.mutexObj.RUnlock()

	return cnt
}

func NewMutexCounter() *MutexCounter {
	return &MutexCounter{mutexObj: sync.RWMutex{}}
}

// --------------------------

type ClientAuthCredentials struct{}

func (t ClientAuthCredentials) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Basic " + apiClientCredentials[0]}, nil
}

func (ClientAuthCredentials) RequireTransportSecurity() bool {
	return false
}
