package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Handler struct {
	dynamoClient     *dynamodb.DynamoDB
}

type Round struct {
	Round   string   `json:"round"`
	Results []Result `json:"results"`
}

type Result struct {
	PlayerName string `json:"playerName"`
	Points     int    `json:"points"`
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	handler := Handler{
		dynamoClient: getDynamoClient(),
	}

	queriedRound := request.QueryStringParameters["queriedRound"]

	if queriedRound == "" {
		rounds := handler.getAllRounds()
		returnValue, err := json.Marshal(rounds)
		if err != nil {
			panic(err)
		}
		return events.APIGatewayProxyResponse{Body: string(returnValue), StatusCode: 200}, nil
	} else {
		round := handler.getRound(queriedRound)
		returnValue, err := json.Marshal(round)
		if err != nil {
			panic(err)
		}

		return events.APIGatewayProxyResponse{Body: string(returnValue), StatusCode: 200}, nil
	}
}

func getDynamoClient() *dynamodb.DynamoDB {
	endpointCfg := aws.NewConfig().
		WithRegion("eu-west-1").
		WithCredentialsChainVerboseErrors(true)

	s := session.Must(session.NewSession())
	dynamoClient := dynamodb.New(s, endpointCfg)
	return dynamoClient
}

func (handler *Handler) getAllRounds() []Round {

}

func (handler *Handler) getRound(round string) Round {

}

func main() {
	lambda.Start(HandleRequest)
}
