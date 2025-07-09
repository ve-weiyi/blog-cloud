package queue

import (
	"container/list"
)

// SubscriptionList manages a list of subscriptions.
type SubscriptionList struct {
	subs *list.List
}

// NewSubscriptionList creates a new subscription list.
func NewSubscriptionList() *SubscriptionList {
	return &SubscriptionList{subs: list.New()}
}

// Add adds a subscription to the list.
func (sl *SubscriptionList) Add(sub Subscription) {
	sl.subs.PushBack(sub)
}

// Remove removes a subscription from the list.
func (sl *SubscriptionList) Remove(sub Subscription) {
	for e := sl.subs.Front(); e != nil; e = e.Next() {
		if sub == e.Value.(Subscription) {
			sl.subs.Remove(e)
			return
		}
	}
}

// Get retrieves and removes the first subscription from the list.
func (sl *SubscriptionList) Get() Subscription {
	e := sl.subs.Front()
	if e == nil {
		return nil
	}
	sl.subs.Remove(e)
	return e.Value.(Subscription)
}
