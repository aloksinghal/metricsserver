package publish

type MockPublisher struct {}


func (m MockPublisher) Publish (messages [][]byte, key string) error {
	return nil
}