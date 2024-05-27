package main

import (
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func CreateCase(myCase Case) (Case, error) {

	log.Printf("Case: %+v", myCase)

	case_id := generateRandomString(16)
	myCase.ID = case_id

	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(case_id),
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
			"user_id": {
				S: aws.String(myCase.UserID),
			},
		},
		TableName: &CasesTable,
	})

	return myCase, err
}

func GetCaseFromId(caseID string) (Case, error) {
	filt := expression.Name("_id").Equal(expression.Value(caseID))

	//log case id
	log.Printf("Case ID: %s", caseID)

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return Case{}, err
	}

	result, err := dynamo.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 &CasesTable,
	})

	if err != nil {
		return Case{}, err
	}

	if len(result.Items) == 0 {
		return Case{}, nil
	}

	if len(result.Items) > 1 {
		return Case{}, err
	}

	ult := result.Items[0]

	numberFiles, err := strconv.Atoi(*ult["number_files"].N)

	//log number of files

	if err != nil {
		return Case{}, err
	}

	myCase := Case{
		ID:                *ult["_id"].S,
		CaseTitle:         *ult["case_title"].S,
		AttorneyFirstName: *ult["attorney_first_name"].S,
		AttorneyLastName:  *ult["attorney_last_name"].S,
		BucketName:        *ult["bucket_name"].S,
		CaseInfo:          *ult["case_info"].S,
		CaseType:          *ult["case_type"].S,
		City:              *ult["city"].S,
		Date:              *ult["date"].S,
		JudgeName:         *ult["judge_name"].S,
		NumberFiles:       numberFiles,
		SeedDoc:           *ult["seed_doc"].S,
		SeedText:          *ult["seed_text"].S,
		State:             *ult["state"].S,
		UserID:            *ult["user_id"].S,
	}

	//log case
	log.Printf("Case: %+v", myCase)

	return myCase, nil
}

func GetCasesByUserId(user_id string) ([]Case, error) {
	filt := expression.Name("user_id").Equal(expression.Value(user_id))
	log.Printf("User ID: %s", user_id)
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return []Case{}, err
	}

	result, err := dynamo.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 &CasesTable,
	})

	if err != nil {
		return []Case{}, err
	}

	if len(result.Items) == 0 {
		return []Case{}, nil
	}

	var cases []Case

	for _, i := range result.Items {
		numberFiles, err := strconv.Atoi(*i["number_files"].N)
		if err != nil {
			return []Case{}, err
		}

		myCase := Case{
			ID:                *i["_id"].S,
			CaseTitle:         *i["case_title"].S,
			AttorneyFirstName: *i["attorney_first_name"].S,
			AttorneyLastName:  *i["attorney_last_name"].S,
			BucketName:        *i["bucket_name"].S,
			CaseInfo:          *i["case_info"].S,
			CaseType:          *i["case_type"].S,
			City:              *i["city"].S,
			Date:              *i["date"].S,
			JudgeName:         *i["judge_name"].S,
			NumberFiles:       numberFiles,
			SeedDoc:           *i["seed_doc"].S,
			SeedText:          *i["seed_text"].S,
			State:             *i["state"].S,
			UserID:            *i["user_id"].S,
		}

		cases = append(cases, myCase)
	}

	log.Printf("Cases: %+v", cases)

	return cases, nil

}

func DeleteCaseById(caseID string) (Case, error) {
	myCase, err := GetCaseFromId(caseID)
	if err != nil {
		return Case{}, err
	}

	if myCase.ID == "" {
		return Case{}, nil
	}

	_, err = dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(caseID),
			},
		},
		TableName: &CasesTable,
	})

	return myCase, err
}

func DeleteCasesByUser(user_id string) ([]Case, error) {
	cases, err := GetCasesByUserId(user_id)
	if err != nil {
		return []Case{}, err
	}

	if len(cases) == 0 {
		return []Case{}, nil
	}

	for _, c := range cases {
		_, err = dynamo.DeleteItem(&dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"_id": {
					S: aws.String(c.ID),
				},
			},
			TableName: &CasesTable,
		})
		if err != nil {
			return []Case{}, err
		}
	}

	return cases, nil
}
