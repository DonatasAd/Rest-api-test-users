package handlers

import (
	"net/http"
	"rest-api-test-users/pkg/user"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
)

const errorMethodNotAllowed = "method Not allowed"

// ErrorBody struct object
type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

// GetUsersHandler gets user by ID
func GetUsersHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["ID"]
	// if id exist find user by ID
	if len(id) > 0 {
		user, err := user.GetUserByID(id)
		if err != nil {
			errorBody := ErrorBody{aws.String(err.Error())}
			return apiResponse(http.StatusBadRequest, errorBody)
		}
		return apiResponse(http.StatusOK, user)
	}
	// if id do not exist get list of users
	users, err := user.GetUsers()
	if err != nil {
		errorBody := ErrorBody{aws.String(err.Error())}
		return apiResponse(http.StatusBadRequest, errorBody)
	}
	return apiResponse(http.StatusOK, users)
}

// PostUserHandler creates new user in DynamoDB
func PostUserHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, err := user.CreateUser(req)
	if err != nil {
		errorBody := ErrorBody{aws.String(err.Error())}
		return apiResponse(http.StatusBadRequest, errorBody)
	}
	return apiResponse(http.StatusOK, user)
}

// PutUserHandler replace with new user DynamoDB
func PutUserHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, err := user.PutUser(req)
	if err != nil {
		errorBody := ErrorBody{aws.String(err.Error())}
		return apiResponse(http.StatusBadRequest, errorBody)
	}
	return apiResponse(http.StatusOK, user)
}

// DeleteUserHandler deletes user from DynamoDB
func DeleteUserHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["ID"]
	err := user.DeleteUser(id)
	if err != nil {
		errorBody := ErrorBody{aws.String(err.Error())}
		return apiResponse(http.StatusBadRequest, errorBody)
	}
	return apiResponse(http.StatusOK, nil)
}

// UnhandledMethod handles methods that are not implimented
func UnhandledMethod() (events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, errorMethodNotAllowed)
}
