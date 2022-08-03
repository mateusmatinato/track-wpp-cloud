package repository

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"wpp-cloud/internal/domain"
)

type Repository interface {
	SaveTrackingInfo(info domain.TrackingInfo) error
	GetTrackingInfo(code string) (domain.TrackingInfo, error)
	ScanTrackingInfo() ([]domain.TrackingInfo, error)
}

type repository struct {
	*dynamodb.DynamoDB
}

func (r repository) ScanTrackingInfo() ([]domain.TrackingInfo, error) {

	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := r.DynamoDB.Scan(params)
	if err != nil {
		log.Printf("[ERROR] Query API call failed: %s", err)
		return nil, err
	}

	if len(result.Items) == 0 {
		return []domain.TrackingInfo{}, nil
	}

	var items []domain.TrackingInfo

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		fmt.Println("[ERROR] Error unmarshalling tracking info from dynamodb")
		return []domain.TrackingInfo{}, errors.New("[ERROR] Error unmarshalling tracking info from dynamodb")
	}

	return items, nil
}

func (r repository) GetTrackingInfo(code string) (domain.TrackingInfo, error) {

	input := &dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			"code": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(code),
					},
				},
			},
		},
		TableName: aws.String(tableName),
	}

	result, readErr := r.DynamoDB.Query(input)
	if readErr != nil {
		fmt.Println("[ERROR] Error searching tracking info on dynamodb: " + readErr.Error())
		return domain.TrackingInfo{}, readErr
	}

	if len(result.Items) == 0 {
		return domain.TrackingInfo{}, nil
	}

	var items []domain.TrackingInfo

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		fmt.Println("[ERROR] Error unmarshalling tracking info from dynamodb")
		return domain.TrackingInfo{}, errors.New("[ERROR] Error unmarshalling tracking info from dynamodb")
	}

	return items[0], nil
}

const tableName = "tracking"

func (r repository) SaveTrackingInfo(info domain.TrackingInfo) error {
	trackingInfo, marshalErr := dynamodbattribute.MarshalMap(info)

	if marshalErr != nil {
		fmt.Println("[ERROR] Failed marshalling tracking info to dynamodb")
		return marshalErr
	}

	input := &dynamodb.PutItemInput{
		Item:      trackingInfo,
		TableName: aws.String(tableName),
	}

	_, writeErr := r.DynamoDB.PutItem(input)

	if writeErr != nil {
		fmt.Println("Failed to write to dynamo")
		return writeErr
	}

	return nil
}

func NewRepository() Repository {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		},
	))

	db := dynamodb.New(sess)

	return &repository{DynamoDB: db}
}
