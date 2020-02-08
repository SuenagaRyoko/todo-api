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

type UpdateTaskRequest struct {
	UserID string `json:"userID"`
	TaskID string `json:"taskID"`
	TaskName string `json:"taskName"`
	Deadline string `json:"deadline"`
	Status string `json:"status"`
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

	updateTaskRequest := &UpdateTaskRequest{}

	updateTaskRequest.UserID = request.PathParameters["userID"]
	updateTaskRequest.TaskID = request.PathParameters["taskID"]

	if err := json.Unmarshal(jsonBytes, updateTaskRequest); err != nil {
		fmt.Println("[Request parse Error]", err)
	}

	fmt.Println("[updateTaskRequest]", updateTaskRequest)

	sess, sessErr := session.NewSession()
	if sessErr != nil {
		fmt.Println("[Session Error]", sessErr)
	}
	svc := dynamodb.New(sess)

	updateParams := &dynamodb.UpdateItemInput{
		TableName: aws.String("todo"),
		ExpressionAttributeNames: map[string]*string{
			"#Status":  aws.String("status"),
			"#Deadline":  aws.String("deadline"),
			"#TaskName":  aws.String("taskName"),
    },
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(updateTaskRequest.Status),
			},
			":deadline": {
				S: aws.String(updateTaskRequest.Deadline),
			},
			":taskName": {
				S: aws.String(updateTaskRequest.TaskName),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"userID": {
					S: aws.String(updateTaskRequest.UserID),
			},
			"taskID": {
					S: aws.String(updateTaskRequest.TaskID),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET #Status = :status, #Deadline = :deadline, #TaskName = :taskName"),
	}

	updateItem, updateErr := svc.UpdateItem(updateParams)
	if updateErr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, updateErr
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintln(updateItem),
	}, nil
}

func main() {
	lambda.Start(handler)
}
