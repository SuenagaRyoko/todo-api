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

type PostTaskRequest struct {
	UserID string `json:"userID"`
	TaskID string `json:"taskID"`
	TaskName string `json:"taskName"`
	Deadline string `json:"deadline"`
	Status string `json:"status"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body
	jsonBytes := ([]byte)(body)
	task := &PostTaskRequest{}

	task.UserID = request.PathParameters["userID"]

	if err := json.Unmarshal(jsonBytes, task); err != nil {
		fmt.Println("[Parse Error]", err)
	}

	fmt.Println("[task]", task)

	sess, sessErr := session.NewSession()
	if sessErr != nil {
		fmt.Println("[Session Error]", sessErr)
	}
	svc := dynamodb.New(sess)

	putParams, putParamsErr := dynamodbattribute.MarshalMap(task)
	if putParamsErr != nil {
		fmt.Println("[putParams Error]", putParamsErr)
	}

	_, paramErr := svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("todo"),
		Item:      putParams,
	})
	if paramErr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, paramErr
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintln("success"),
	}, nil
}

func main() {
	lambda.Start(handler)
}
