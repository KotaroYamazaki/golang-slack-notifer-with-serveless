package main

import (
	"github.com/KotaroYamazaki/golang-slack-notifer-with-serveless/app"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(app.Run)
}
