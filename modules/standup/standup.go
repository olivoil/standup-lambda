package standup

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/RocksauceStudios/standup-parser"
	"github.com/altairsix/eventsource"
)

// SubmitStatus is a command to submit a status
type SubmitStatus struct {
	eventsource.CommandModel
	TeamID string
	UserID string
	Text   string
}

// Apply converts a command into a series of events
func (item *Status) Apply(ctx context.Context, cmd eventsource.Command) ([]eventsource.Event, error) {
	switch v := cmd.(type) {
	case *SubmitStatus:
		statusSubmitted := &StatusSubmitted{
			Model: eventsource.Model{
				ID:      cmd.AggregateID(),
				Version: item.Version + 1,
				At:      time.Now(),
			},
			TeamID: v.TeamID,
			UserID: v.UserID,
			Text:   v.Text,
		}
		return []eventsource.Event{statusSubmitted}, nil

	default:
		return nil, fmt.Errorf("unhandled command, %v", v)
	}
}

// StatusSubmitted defines a status submitted event
type StatusSubmitted struct {
	eventsource.Model
	TeamID string
	UserID string
	Text   string
}

// Aggregate
type Status struct {
	ID          string
	Version     int
	TeamID      string
	UserID      string
	Statement   *parser.Statement
	Text        string
	SubmittedAt time.Time
}

// On applies an event on an Aggregate
func (item *Status) On(event eventsource.Event) error {
	switch v := event.(type) {
	case *StatusSubmitted:
		item.Version = v.Model.Version
		item.ID = v.Model.ID

		p := parser.New(strings.NewReader(v.Text))
		statement, err := p.Parse()
		if err != nil {
			return err
		}
		item.Statement = statement

		item.TeamID = v.TeamID
		item.UserID = v.UserID
		item.SubmittedAt = time.Now()

	default:
		return fmt.Errorf("unable to handle event, %v", v)
	}

	return nil
}
