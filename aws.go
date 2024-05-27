package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	dynamo *dynamodb.DynamoDB
	s3Client *s3.S3
)

// connectDynamo returns a dynamoDB client
func InitDynamoDBTClient() (db *dynamodb.DynamoDB) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
}

// connectS3 returns a S3 client
func InitS3Client() (s3Client *s3.S3) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return s3.New(sess)
}
