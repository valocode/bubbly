package adapter

import (
	"fmt"
	"testing"
)

// This function will look to see that the current subscription model runs efficiently
// and doesn't run into any immediate design issues
func TestSubscriptions(t *testing.T) {
	release := &Release{
		ID: "1234",
	}

	actionType := &ActionType{
		ID:   "1234",
		Type: "Send Email",
	}
	sub := &Subscription{
		Release:   *release,
		EventType: NewCVE,
	}
	action := &Action{
		Subscription: *sub,
		ActionType:   *actionType,
		Notification: func(cve CVE) {
			fmt.Printf("A new CVE with ID: %s has been found\n", cve.ID)
		},
	}
	// Create a list of actions and subscriptions
	actionList := make(map[string]Action)
	actionList[sub.ID] = *action
	subs := Subscriptions{*sub}

	// A new CVE is found and added to a Release
	cve, eventType := newCVE(release.ID)

	// loop through the subscriptions and see if any conditions are met
	// in a DB environment, it will be easier to just query based on EventType, rather than loop through *every* sub
	for _, s := range subs {
		if s.EventType == eventType {
			// mimic fetching from DB based on subscription ID and send the notification
			actionList[s.ID].Notification(*cve)
		}
	}
}

// Creates a new CVE and returns the eventType NewCVE to simulate a trigger
func newCVE(releaseID string) (*CVE, EventType) {
	return &CVE{
		ID: releaseID,
	}, NewCVE
}

type Release struct {
	ID string
}

// in this model, subscriptions are tied to a Release and an eventtype, this way, if a release has a subscription to
// the eventType NewCVE, we can easily see which releases need to be informed
type Subscription struct {
	ID string
	Release
	EventType
}

// an action is something that happens after an event. This is something that the Subscription will do. As of now, I added
// the Notification function as a way to show something that can be done, because presumably, someone will want to be
// notified of an action.
type Action struct {
	Subscription
	ActionType
	Notification notification
}

// An actionType will have an ID and a Type
type ActionType struct {
	ID   string
	Type string
}

// just a dummy CVE object
type CVE struct {
	ID string
}

type notification func(cve CVE)
type Subscriptions []Subscription

type EventType string

// define event types
const (
	NewCVE EventType = "New CVE"
)
