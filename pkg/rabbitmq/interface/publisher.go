package rmqinterface

// Publisher provides the interface for the functionality of Publisher (message broker)
//
//	FYI: to use this interface, we need to implement it somewhere in our code. :)
type Publisher interface {
	PublishLogToLogManager(logType string, msg []byte) error    // FYI: Log Manager will handle how it will be sent into another channel
	PublishLogToElasticSearch(logType string, msg []byte) error // publishes to the elasticsearch database
	PublishLogToMessageBroker(logType string, msg []byte) error // publishes to the message broker (e.g. RabbitMQ)
}
