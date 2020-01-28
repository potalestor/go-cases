package pubsub

import "fmt"

type Subscriber interface {
	EventChan() chan<- string
	Key() string
}

type Sub struct {
	key       string
	eventChan chan string
}

func NewSub(key string) Subscriber {
	s := &Sub{key, make(chan string)}
	go s.execute()
	return s
}
func (s *Sub) Key() string {
	return s.key
}
func (s *Sub) EventChan() chan<- string {
	return s.eventChan
}

func (s *Sub) execute() {
	for s := range s.eventChan {
		fmt.Println(s)
	}
	fmt.Println(s.key, "sub closed")
}
