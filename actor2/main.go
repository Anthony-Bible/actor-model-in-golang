package main

import (
	"fmt"
	"sync"
)

type User1 struct {
	Ch      chan Message
	handler func(Message)
	wg      *sync.WaitGroup
}
type Message struct {
	body   string
	sender string
}

func newUser1(handler func(Message)) *User1 {
	channel := make(chan Message)
	return &User1{
		Ch:      channel,
		handler: handler,
		wg:      new(sync.WaitGroup),
	}
}
func (u *User1) Start() {
	go func() {
		for {
			msg := <-u.Ch
			u.handler(msg)
			u.wg.Done()
		}
	}()
}

func handler1(msg Message) {
	fmt.Println(msg.sender, msg.body)
}
func (u *User1) Stop() {
	close(u.Ch)
}
func (u *User1) Send(msg Message) {
	u.wg.Add(1)
	u.Ch <- msg

}

func main() {
	user1 := newUser1(handler1)
	user1.Start()
	user1.Send(Message{body: "Hello", sender: "user1"})
	user1.Send(Message{body: "World", sender: "user1"})
	user1.wg.Wait()
}
