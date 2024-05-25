package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func createUser(user User) error {
	//print user
	log.Printf("User: %+v", user)

	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(generateRandomString(16)),
			},
			"email": {
				S: aws.String(user.Email),
			},
			"cases": {
				L: []*dynamodb.AttributeValue{},
			},
			"first_name": {
				S: aws.String(user.FirstName),
			},
			"last_name": {
				S: aws.String(user.LastName),
			},
			"organization": {
				S: aws.String(user.Organization),
			},
			"password": {
				S: aws.String(user.Password),
			},
			"profile_picture": {
				S: aws.String(user.ProfilePicture),
			},
		},
		TableName: &UsersTable,
	})

	return err
}

func getUserFromEmail(email string) (User, error) {

	log.Printf("Email: %s", email)
	filt := expression.Name("email").Equal(expression.Value(email))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(UsersTable),
	}

	result, err := dynamo.Scan(params)

	if err != nil {
		return User{}, err
	}

	///if results has length 0, return empty user
	if len(result.Items) == 0 {
		return User{}, nil
	}

	//if results has length > 1, return error
	if len(result.Items) > 1 {
		return User{}, err
	}

	new_result := result.Items[0]

	user := User{
		ID:             *new_result["_id"].S,
		Email:          *new_result["email"].S,
		Cases:          []string{},
		FirstName:      *new_result["first_name"].S,
		LastName:       *new_result["last_name"].S,
		Organization:   *new_result["organization"].S,
		Password:       *new_result["password"].S,
		ProfilePicture: *new_result["profile_picture"].S,
	}

	log.Printf("User: %+v", user)

	return user, nil
}
