Run stack local 
```bash 
docker pull localstack/localstack

docker run -it -d -p 4566:4566 localstack/localstack start
```

Initialise queue
```bash	
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name my-queue
```

Send message
```bash
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://localhost:4566/000000000000/my-queue --message-body "Hello, World!"
```

Create Topic
```bash
aws --endpoint-url=http://localhost:4566 sns create-topic --name my-topic
```

Subscribe to topic
```bash
aws --endpoint-url=http://localhost:4566 sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:my-topic --protocol sqs --notification-endpoint arn:aws:sqs:us-east-1:000000000000:my-queue
```