module github.com/mtanzim/event-driven-bookstore/bookstore-server

go 1.15

require (
	github.com/aws/aws-sdk-go v1.36.18 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.3.0
	github.com/klauspost/compress v1.11.4 // indirect
	github.com/mtanzim/event-driven-bookstore/common-server v0.0.0-20201229182559-fbcb61992188
	github.com/rs/cors v1.7.0
	go.mongodb.org/mongo-driver v1.4.4
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	gopkg.in/confluentinc/confluent-kafka-go.v1 v1.5.2
)
