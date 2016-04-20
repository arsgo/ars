package mqservice

type IMQPublisher interface {
	Publish(string) error
}
type IMQConsumer interface {
	Consume(func(string))
}

type IMQService interface {
	NewPublisher() IMQPublisher
	NewConsumer() IMQConsumer
}
