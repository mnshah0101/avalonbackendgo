package main

import (
	"log"
	"strconv"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"bytes"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func GetDocumentsByCaseId(caseID string) ([]Document, error) {
	filt := expression.Name("case").Equal(expression.Value(caseID))
	log.Printf("Case ID: %s", caseID)
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return []Document{}, err
	}

	result, err := dynamo.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 &DocumentsTable,
	})

	if err != nil {
		return []Document{}, err
	}

	if len(result.Items) == 0 {
		return []Document{}, nil
	}

	var documents []Document

	for _, i := range result.Items {
		relevancy, err := strconv.ParseFloat(*i["relevancy"].N, 64)
		stored, err := strconv.ParseBool(strconv.FormatBool(*i["stored"].BOOL))
		if err != nil {
			return []Document{}, err
		}

		doc := Document{
			ID:        *i["_id"].S,
			FileNames: *i["file_name"].S,
			CaseID:    *i["case"].S,
			Date:      *i["date"].S,
			FileURL:   *i["file_url"].S,
			Relevancy: relevancy,
			Stored:    stored,
		}

		documents = append(documents, doc)
	}

	return documents, nil

}

func GetDocumentById(documentID string) (Document, error) {
	filt := expression.Name("_id").Equal(expression.Value(documentID))
	log.Printf("Document ID: %s", documentID)
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return Document{}, err
	}

	result, err := dynamo.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 &DocumentsTable,
	})

	if err != nil {
		return Document{}, err
	}

	if len(result.Items) == 0 {
		return Document{}, nil
	}

	if len(result.Items) > 1 {
		return Document{}, err
	}

	relevancy, err := strconv.ParseFloat(*result.Items[0]["relevancy"].N, 64)
	stored, err := strconv.ParseBool(strconv.FormatBool(*result.Items[0]["stored"].BOOL))
	if err != nil {
		return Document{}, err
	}

	doc := Document{
		ID:        *result.Items[0]["_id"].S,
		FileNames: *result.Items[0]["file_name"].S,
		CaseID:    *result.Items[0]["case"].S,
		Date:      *result.Items[0]["date"].S,
		FileURL:   *result.Items[0]["file_url"].S,
		Relevancy: relevancy,
		Stored:    stored,
	}

	return doc, nil
}

func DeleteDocumentById(documentID string) error {
	_, err := dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: &documentID,
			},
		},
		TableName: &DocumentsTable,
	})

	return err
}

func DeleteDocumentsByCaseId(caseID string) ([]Document, error) {
	documents, err := GetDocumentsByCaseId(caseID)
	if err != nil {
		return nil, err
	}

	for _, doc := range documents {
		err = DeleteDocumentById(doc.ID)
		if err != nil {
			return nil, err
		}
	}

	return documents, nil
}




func uploadFileToS3(fileName string, fileContent []byte) error {
	input := &s3.PutObjectInput{
		Bucket:  &Bucket,
		Key:    aws.String(fileName),
		Body:   aws.ReadSeekCloser(bytes.NewReader(fileContent)),
	}

	_, err := s3Client.PutObject(input)
	return err
}
