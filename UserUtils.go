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

func getUserFromId(id string) (User, error) {

	log.Printf("ID: %s", id)
	filt := expression.Name("_id").Equal(expression.Value(id))

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

	//if results has length 0, return empty user
	if len(result.Items) == 0 {
		return User{}, nil
	}

	//if results has length > 1, return error
	if len(result.Items) > 1 {
		return User{}, err
	}

	new_result := result.Items[0]
	log.Printf("New Result: %+v", new_result)
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

func deleteUserFromId(id string, email string) (user User, err error) {
	//check in database if user exists
	
	user, err = getUserFromId(id)
	if err != nil {
		return user, err
	}
	if user.ID == "" {
		fakeUser := User{}
		return fakeUser, err
	}
	
	_, err = dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(id),
			},
		},
		TableName: &UsersTable,
	})
	log.Printf("Deleted user with ID: %s and email: %s", id, email)
	return user, err
}

func updateUser(user User) ( User,  error) {
	//get user by id
	var return_user User

	var err error

	log.Printf("FromUpdateUser Before User: %+v", user)
	return_user, err = getUserFromId(user.ID)
	log.Printf("FromUpdateUser AFter User: %+v", return_user)
	if err != nil {
		return return_user, err
	}
	if return_user.ID == "" {
		return return_user, err
	}

	_, err = dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#A": aws.String("first_name"),
			"#B": aws.String("last_name"),
			"#C": aws.String("organization"),
			"#D": aws.String("profile_picture"),
			"#E": aws.String("email"),
			"#F":  aws.String("password"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":first_name": {
				S: aws.String(user.FirstName),
			},
			":last_name": {
				S: aws.String(user.LastName),
			},
			":organization": {
				S: aws.String(user.Organization),
			},
			":profile_picture": {
				S: aws.String(user.ProfilePicture),
			},
			":email": {
				S: aws.String(user.Email),
			},
			":password": {
				S: aws.String(user.Password),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(user.ID),
			},
		},
		TableName:        &UsersTable,
		UpdateExpression: aws.String("SET #A = :first_name, #B = :last_name, #C = :organization, #D = :profile_picture, #E = :email, #F = :password"),
	})

	return return_user, err
}
	