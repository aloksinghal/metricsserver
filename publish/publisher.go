package publish

type Publisher interface {
	Publish(messages [][]byte, key string) error
}