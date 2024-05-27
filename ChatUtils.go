package main

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func GetChatFromCaseId(caseID string) (Chat, error) {
	filt := expression.Name("_id").Equal(expression.Value(caseID))

	//log case id
	log.Printf("Case ID: %s", caseID)

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return Chat{}, err
	}

	result, err := dynamo.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 &ChatsTable,
	})

	if err != nil {
		return Chat{}, err
	}

	if len(result.Items) == 0 {
		return Chat{}, nil
	}

	if len(result.Items) > 1 {
		return Chat{}, err
	}

	ult := result.Items[0]

	//create new chat return object
	var messages []Message
	for _, m := range ult["messages"].L {
		log.Printf("Message: %v", m)
		message := Message{
			Text:      *m.M["text"].S,
			Sender:    *m.M["sender"].S,
			Timestamp: *m.M["timestamp"].S,
		}
		messages = append(messages, message)
	}

	var selectedDocs []string
	for _, d := range ult["selected_docs"].L {
		selectedDocs = append(selectedDocs, *d.S)
	}

	myChat := Chat{
		ID:           *ult["_id"].S,
		Messages:     messages,
		SelectedDocs: selectedDocs,
		UserID:       *ult["user_id"].S,
	}

	return myChat, err

}
