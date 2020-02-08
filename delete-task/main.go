package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DeleteTaskRequest struct {
	UserID string `json:"userID"`
	TaskID string `json:"taskID"`
}

type Tasks struct {
	UserID string `json:"userID"`
	TaskID string `json:"taskID"`
	TaskName string `json:"taskName"`
	Deadline string `json:"deadline"`
	Status string `json:"status"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body
	jsonBytes := ([]byte)(body)

	deleteTaskRequest := &DeleteTaskRequest{}

	if err := json.Unmarshal(jsonBytes, deleteTaskRequest); err != nil {
		fmt.Println("[Request Parse Error]", err)
	}

	deleteTaskRequest.UserID = request.PathParameters["userID"]
	deleteTaskRequest.TaskID = request.PathParameters["taskID"]

	fmt.Println("[deleteTaskRequest]", deleteTaskRequest)

	sess, sessErr := session.NewSession()
	if sessErr != nil {
		fmt.Println("[Session Error]", sessErr)
	}
	svc := dynamodb.New(sess)

	deleteParams := &dynamodb.DeleteItemInput{
		TableName: aws.String("todo"),
		Key: map[string]*dynamodb.AttributeValue{
			"userID": {
				S: aws.String(deleteTaskRequest.UserID),
			},
			"taskID": {
				S: aws.String(deleteTaskRequest.TaskID),
			},
		},
	}

	fmt.Println("[deleteParams]", deleteParams)

	if _, deleteErr := svc.DeleteItem(deleteParams); deleteErr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, deleteErr
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintln("success"),
	}, nil
}

func main() {
	lambda.Start(handler)
}
