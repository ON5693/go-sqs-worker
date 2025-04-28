package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	QueueUrl = "http://localhost:4566/000000000000/my-queue"
	TopicArn = "arn:aws:sns:us-east-1:000000000000:my-topic"
)

func main() {
	go initSQS()

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:4566"),
		Region:   aws.String("us-east-1"),
	}))

	svc := sns.New(sess)

	publishParams := &sns.PublishInput{
		Message:  aws.String("Hello, World!"),
		TopicArn: aws.String(TopicArn),
	}

	_, err := svc.Publish(publishParams)
	if err != nil {
		fmt.Println("Error publishing message:", err)
		return
	}

	fmt.Println("Message published successfully")
}

func initSQS() {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:4566"),
		Region:   aws.String("us-east-1"),
	}))

	svc := sqs.New(sess)

	receiverParams := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(QueueUrl),
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(20),
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signalCh:
			fmt.Println("Received signal to terminate")
			return
		default:
			result, err := svc.ReceiveMessage(receiverParams)
			if err != nil {
				fmt.Println("Error receiving message:", err)
				time.Sleep(1 * time.Second)
				continue
			}

			for _, message := range result.Messages {
				fmt.Printf("Received message: %s \n", *message.Body)

				deleteParams := &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(QueueUrl),
					ReceiptHandle: message.ReceiptHandle,
				}

				_, err := svc.DeleteMessage(deleteParams)
				if err != nil {
					fmt.Println("Error deleting message:", err)
					time.Sleep(1 * time.Second)
					continue
				}
			}
		}
	}
}
