package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"encoding/json"

	"net/http"

	"bitbucket.org/RocksauceStudios/standup-lambda/modules/slack"
	"bitbucket.org/RocksauceStudios/standup-lambda/modules/standup"
	"github.com/altairsix/eventsource"
	"github.com/altairsix/eventsource/dynamodbstore"
	"github.com/apex/go-apex"
	"github.com/apex/log"
	jlog "github.com/apex/log/handlers/json"
	"github.com/gorilla/schema"
	"github.com/segmentio/ksuid"
)

type message struct {
	Body string `json:"body"`
}

type response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

type Event struct {
	ChannelID   string `schema:"channel_id"`
	ChannelName string `schema:"channel_name"`
	Command     string `schema:"command"`
	ResponseURL string `schema:"response_url"`
	TeamDomain  string `schema:"team_domain"`
	TeamID      string `schema:"team_id"`
	Text        string `schema:"text"`
	Token       string `schema:"token"`
	UserID      string `schema:"user_id"`
	UserName    string `schema:"user_name"`
}

func main() {
	log.SetHandler(jlog.Default)
	decoder := schema.NewDecoder()

	apex.HandleFunc(func(msg json.RawMessage, lambdaContext *apex.Context) (interface{}, error) {
		var m message

		if err := json.Unmarshal(msg, &m); err != nil {
			return nil, fmt.Errorf("unable to decode request (%v): %v", string(msg), err)
		}

		var event Event
		values, err := url.ParseQuery(m.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to parse form data (%v): %v", m.Body, err)
		}

		if err := decoder.Decode(&event, values); err != nil {
			return nil, fmt.Errorf("unable to decode form data (%+v): %v", values, err)
		}

		// verify event.Token to ensure request is coming from Slack
		if event.Token != os.Getenv("SLACK_TOKEN") {
			return nil, fmt.Errorf("invalid message token: %#v", event)
		}

		// create dynamodb store for events
		store, err := dynamodbstore.New(
			os.Getenv("AWS_DYNAMODB_TABLE_STATUS_EVENTS"),
			dynamodbstore.WithRegion(os.Getenv("AWS_REGION")),
		)
		if err != nil {
			return nil, fmt.Errorf("error creating store: %v", err)
		}

		// create eventsource repo
		repo := eventsource.New(&standup.Status{},
			eventsource.WithStore(store),
			eventsource.WithSerializer(eventsource.NewJSONSerializer(
				standup.StatusSubmitted{},
			)),
		)

		ctx := context.Background()
		id := ksuid.New()

		err = repo.Dispatch(ctx, &standup.SubmitStatus{
			CommandModel: eventsource.CommandModel{ID: id.String()},
			TeamID:       event.TeamID,
			UserID:       event.UserID,
			Text:         event.Text,
		})
		if err != nil {
			return nil, fmt.Errorf("error dispatching command: %v", err)
		}

		aggregate, err := repo.Load(ctx, id.String())
		if err != nil {
			return nil, fmt.Errorf("error loading aggregate: %v", err)
		}

		// Send slack message to `SLACK_STANDUP_CHANNEL`
		if err := slack.PostStatusToSlack(
			os.Getenv("SLACK_AUTHENTICATION_TOKEN"),
			os.Getenv("SLACK_STANDUP_CHANNEL"),
			aggregate.(*standup.Status),
		); err != nil {
			return nil, fmt.Errorf("error posting to Slack: %v", err)
		}

		var body string = ""
		if event.ChannelID != os.Getenv("SLACK_STANDUP_CHANNEL") {
			body = `{"response_type": "ephemeral", "text": "status submitted in #standup channel"}`
		}

		return response{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: body,
		}, nil
	})
}
