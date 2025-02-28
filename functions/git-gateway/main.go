package main

import (
	"log"
	"time"

	"github.com/netlify/git-gateway/cmd"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	c := make(chan error, 1)

	go func() {
		c <- cmd.RootCommand().Execute()
	}()

	select {
		case err:= <-c:
			log.Fatal(err)
			return &events.APIGatewayProxyResponse{
				StatusCode:        400,
				Body:              "Something didn't work",
			}, nil
		case <-time.After(2 * time.Second):
			return &events.APIGatewayProxyResponse{
				StatusCode:        200,
				Body:              "Ran Git Gateway for 2s, then quit.",
			}, nil
	}
}

func main() {
	lambda.Start(handler)
}
