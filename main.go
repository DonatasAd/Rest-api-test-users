package main

import (
	"net/http"
	"os"
	"rest-api-test-users/pkg/handlers"
	"rest-api-test-users/pkg/user"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func main() {
	// create AWS session
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return
	}
	// Create DynamoDB client
	user.DynaClient = dynamodb.New(awsSession)
	//
	lambda.Start(router)
}

// router calls function based on http method
func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case http.MethodGet:
		return handlers.GetUsersHandler(req)
	case http.MethodPost:
		return handlers.PostUserHandler(req)
	case http.MethodPut:
		return handlers.PutUserHandler(req)
	case http.MethodDelete:
		return handlers.DeleteUserHandler(req)
	default:
		return handlers.UnhandledMethod()
	}
}
