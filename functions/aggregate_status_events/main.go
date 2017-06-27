package main

import (
	"context"
	"fmt"
	"os"

	"bitbucket.org/RocksauceStudios/standup-lambda/modules/standup"
	"github.com/altairsix/eventsource"
	"github.com/altairsix/eventsource/dynamodbstore"
	"github.com/apex/go-apex"
	"github.com/apex/go-apex/dynamo"
	"github.com/apex/log"
	jlog "github.com/apex/log/handlers/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	log.SetHandler(jlog.Default)

	dynamo.HandleFunc(func(event *dynamo.Event, bg *apex.Context) (err error) {

		// create dynamodb store for events
		store, err := dynamodbstore.New(
			os.Getenv("AWS_DYNAMODB_TABLE_STATUS_EVENTS"),
			dynamodbstore.WithRegion(os.Getenv("AWS_REGION")),
		)
		if err != nil {
			return fmt.Errorf("error creating store: %v", err)
		}

		// create eventsource repo
		repo := eventsource.New(&standup.Status{},
			eventsource.WithStore(store),
			eventsource.WithSerializer(eventsource.NewJSONSerializer(
				standup.StatusSubmitted{},
			)),
		)

		for _, record := range event.Records {
			ctx := context.Background()

			switch record.EventName {
			case "INSERT", "MODIFY":
				key := record.Dynamodb.Keys["key"]
				if err := aggregate(ctx, repo, *key.S); err != nil {
					return err
				}
			case "REMOVE":
				key := record.Dynamodb.Keys["key"]
				if err := remove(ctx, *key.S); err != nil {
					return err
				}
			}

		}

		return
	})
}

func aggregate(ctx context.Context, repo *eventsource.Repository, id string) error {
	aggregate, err := repo.Load(ctx, id)
	if err != nil {
		log.WithField("id", id).WithError(err).Info("ignoring load error")
		return nil
	}

	// save aggregate to a dynamodb table
	status := aggregate.(*standup.Status)
	item, err := dynamodbattribute.MarshalMap(status)
	if err != nil {
		log.WithError(err).Info("error converting aggregate to dynamodb attribute map")
		return err
	}

	svc := dynamodb.New(session.New())
	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item:                   item,
		TableName:              aws.String(os.Getenv("AWS_DYNAMODB_TABLE_STATUS_AGGREGATES")),
		ReturnConsumedCapacity: aws.String(dynamodb.ReturnConsumedCapacityNone),
	})
	if err != nil {
		log.WithError(err).Info("error with dynamodb PutItem")
		return err
	}

	return nil
}

func remove(_ context.Context, id string) error {
	svc := dynamodb.New(session.New())

	if _, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		TableName: aws.String(os.Getenv("AWS_DYNAMODB_TABLE_STATUS_AGGREGATES")),
	}); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				// item not in table. ignore...
			default:
				return fmt.Errorf("error with dyhnamodb DeleteItem (%v): %v", aerr.Code(), err)
			}
		} else {
			return fmt.Errorf("error with dyhnamodb DeleteItem: %v", err)
		}
	}

	return nil
}
