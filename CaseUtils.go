package main

import (
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
		"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

)



func CreateCase(myCase Case) (Case, error) {

	log.Printf("Case: %+v", myCase)

	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(generateRandomString(16)),
			},
			"case_title": {
				S: aws.String(myCase.CaseTitle),
			},
			"attorney_first_name": {
				S: aws.String(myCase.AttorneyFirstName),
			},
			"attorney_last_name": {
				S: aws.String(myCase.AttorneyLastName),
			},
			"bucket_name": {
				S: aws.String(myCase.BucketName),
			},
			"case_info": {
				S: aws.String(myCase.CaseInfo),
			},
			"case_type": {
				S: aws.String(myCase.CaseType),
			},
			"city": {
				S: aws.String(myCase.City),
			},
			"date": {
				S: aws.String(myCase.Date),

			},
			"judge_name": {
				S: aws.String(myCase.JudgeName),
			},
			"number_files": {
				N: aws.String(strconv.Itoa(myCase.NumberFiles)),
			},
			"seed_doc": {
				S: aws.String(myCase.SeedDoc),
			},
			"seed_text": {
				S: aws.String(myCase.SeedText),
			},
			"state": {
				S: aws.String(myCase.State),
			},
		},
		TableName: &CasesTable,
	})




	return myCase, err
}


func GetCaseFromId(caseID string) (Case, error) {
	log.Printf("Case ID: %s", caseID)
	filt := expression.Name("_id").Equal(expression.Value(caseID))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(CasesTable),
	}
	result, err := dynamo.Scan(params)
	if err != nil {
		return Case{}, err
	}
	if *result.Count == 0 {
		return Case{}, nil
	}
	var myCase Case
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &myCase)
	if err != nil {
		return Case{}, err
	}
	return myCase, nil
}
