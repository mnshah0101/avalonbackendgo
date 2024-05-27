package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

	var messages []Message
	for _, m := range ult["messages"].L {
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

func AddMessageToChat(caseID string, message string, sender string, date string) error {
	chat, err := GetChatFromCaseId(caseID)
	if err != nil {
		return fmt.Errorf("failed to get chat, %v", err)
	}

	if chat.ID == "" {
		return fmt.Errorf("chat not found")
	}

	newMessage := Message{
		Text:      message,
		Sender:    sender,
		Timestamp: date,
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(ChatsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {S: &caseID}},
		UpdateExpression: aws.String("SET messages = list_append(messages, :newMessage)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":newMessage": {
				L: []*dynamodb.AttributeValue{
					{
						M: map[string]*dynamodb.AttributeValue{
							"text":      {S: &newMessage.Text},
							"sender":    {S: &newMessage.Sender},
							"timestamp": {S: &newMessage.Timestamp},
						},
					},
				},
			},
		},
	}

	_, err = dynamo.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update item, %v", err)
	}

	fmt.Println("Added message to chat")

	return nil
}

func CreateChat(caseID string, userID string) error {
	chat := Chat{
		ID:           caseID,
		Messages:     []Message{},
		SelectedDocs: []string{},
		UserID:       userID,
	}

	//add one message to chat
	newMessage := Message{
		Text:      "Welcome to the chat",
		Sender:    "System",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	chat.Messages = append(chat.Messages, newMessage)



	av, err := dynamodbattribute.MarshalMap(chat)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(ChatsTable),
	}

	_, err = dynamo.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
