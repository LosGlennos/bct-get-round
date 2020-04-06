package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Handler struct {
	DynamoClient *dynamodb.DynamoDB
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
		DynamoClient: getDynamoClient(),
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
	params := &dynamodb.ScanInput{
		TableName: aws.String("BeerCartingTour"),
	}
	result, err := handler.DynamoClient.Scan(params)
	if err != nil {
		panic(err)
	}

	var rounds []Round
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &rounds)
	if err != nil {
		panic(err)
	}

	return rounds
}

func (handler *Handler) getRound(round string) Round {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("BeerCartingTour"),
		IndexName: aws.String("Round"),
		AttributesToGet: []*string{aws.String("PlayerName"), aws.String("Points")},
		KeyConditions: map[string]*dynamodb.Condition{
			"Round": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(round),
					},
				},
			},
		},
	}

	var response, err = handler.DynamoClient.Query(queryInput)
	if err != nil {
		panic(err)
	}

	var roundResult Round
	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &roundResult)
	if err != nil {
		panic(err)
	}

	return roundResult
}

func main() {
	lambda.Start(HandleRequest)
}
