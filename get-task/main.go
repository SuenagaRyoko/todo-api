package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type GetTaskRequest struct {
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

	getTaskRequest := &GetTaskRequest{}

	if err := json.Unmarshal(jsonBytes, getTaskRequest); err != nil {
		fmt.Println("[Request Parse Error]", err)
	}

	getTaskRequest.UserID = request.PathParameters["userID"]
	getTaskRequest.TaskID = request.PathParameters["taskID"]

	fmt.Println("[getTaskRequest]", getTaskRequest)

	sess, sessErr := session.NewSession()
	if sessErr != nil {
		fmt.Println("[Session Error]", sessErr)
	}
	svc := dynamodb.New(sess)

	getParams := &dynamodb.GetItemInput{
		TableName: aws.String("todo"),
		Key: map[string]*dynamodb.AttributeValue{
			"userID": {
				S: aws.String(getTaskRequest.UserID),
			},
			"taskID": {
				S: aws.String(getTaskRequest.TaskID),
			},
		},
	}

	fmt.Println("[getParams]", getParams)

	getItem, getErr := svc.GetItem(getParams)
	if getErr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, getErr
	}

	tasks := &Tasks{}
	if err := dynamodbattribute.UnmarshalMap(getItem.Item, tasks); err != nil {
		fmt.Println("[Unmarshal Error]", err)
	}

	fmt.Println("[getItem]", getItem)

	taskJsonBytes, _ := json.Marshal(tasks)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(taskJsonBytes),
	}, nil
}

func main() {
	lambda.Start(handler)
}
