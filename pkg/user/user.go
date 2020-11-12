package user

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// DynaClient Dynamodb client
var DynaClient dynamodbiface.DynamoDBAPI

// tableName DynamoDB table name
const tableName = "users"

const (
	errorFailedToUnmarshalRecord = "failed to unmarshal record"
	errorFailedToFetchRecord     = "failed to fetch record"
	errorInvalidUserData         = "invalid user data"
	errorCouldNotMarshalItem     = "could not marshal item"
	errorCouldNotDeleteItem      = "could not delete item"
	errorCouldNotDynamoPutItem   = "could not dynamo put item error"
	errorUserAlreadyExists       = "User already exists"
	errorUserDoesNotExists       = "User does not exist"
)

// User struct object
type User struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	LastName string `json:"LastName"`
	Age      int    `json:"Age"`
}

// GetUserByID from DynamoDB
func GetUserByID(id string) (*User, error) {
	// Prepare the input for the query.
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	}

	// Retrieve the item from DynamoDB.
	result, err := DynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(errorFailedToFetchRecord)
	}
	// if user not found return error
	if result.Item == nil {
		return nil, errors.New(errorUserDoesNotExists)
	}

	// The result.Item object returned has the underlying type
	// map[string]*AttributeValue. We can use the UnmarshalMap helper
	// to parse this straight into the fields of a struct. Note:
	// UnmarshalListOfMaps also exists if you are working with multiple
	// items.
	user := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsers find all users from DynamoDB
func GetUsers() (*[]User, error) {
	// Prepare the input for the query.
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	// Retrieve the items from DynamoDB.
	result, err := DynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(errorFailedToFetchRecord)
	}
	// parse result into User struct
	items := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	return items, nil
}

// CreateUser adds new user to DynamoDB
func CreateUser(req events.APIGatewayProxyRequest) (*User, error) {
	var user User
	err := json.Unmarshal([]byte(req.Body), &user)
	if err != nil {
		return nil, errors.New(errorInvalidUserData)
	}
	// Save user

	attributeValues, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(errorCouldNotMarshalItem)
	}
	// Prepare the input for the query.
	input := &dynamodb.PutItemInput{
		Item:                attributeValues,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_not_exists(ID)"),
	}

	_, err = DynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(errorCouldNotDynamoPutItem)
	}
	return &user, nil
}

// PutUser replaces user with new user
func PutUser(req events.APIGatewayProxyRequest) (*User, error) {
	var user User
	err := json.Unmarshal([]byte(req.Body), &user)
	if err != nil {
		return nil, errors.New(errorInvalidUserData)
	}
	// Save user

	attributeValues, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(errorCouldNotMarshalItem)
	}
	// Prepare the input for the query.
	input := &dynamodb.PutItemInput{
		Item:                attributeValues,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_exists(ID)"),
	}

	_, err = DynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(errorCouldNotDynamoPutItem)
	}
	return &user, nil
}

// DeleteUser replaces user with new user
func DeleteUser(id string) error {
	// Prepare the input for the query.
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_exists(ID)"),
	}

	_, err := DynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(errorCouldNotDynamoPutItem)
	}
	return nil
}
