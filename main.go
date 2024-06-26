package main

import (
	"fmt"
)
type Actor struct {
	Ch chan Message
	State interface{}
	Organizer *MessageOrganizer
	Name string

}
// Message is a struct that holds the body of the message
type Message struct {
	body string 
	detail string
	From string
}

// MessageOrganizer is a struct that holds the actors and also takes care of the message passing
type MessageOrganizer struct {
	Actors map[string]*Actor
}

func NewActor(name string, state interface{}, organizer *MessageOrganizer) *Actor{
	// make a new channel of 16 messages
	actor := &Actor{
		Name: name,
		State: state,
		Organizer: organizer,
	}
	return actor

}
func (actor *Actor) StartCustomerActor() {
	ch := make(chan Message, 16)
	actor.Ch = ch
	go actor.CustomerActorLoop(ch)
}

func (actor *Actor) StartWaiterActor() {
	ch := make(chan Message, 16)
	actor.Ch = ch
	go actor.WaiterActorLoop(ch)
}
func (actor *Actor) StartPieCaseActor() {
	ch := make(chan Message, 16)
	actor.Ch = ch
	go actor.PieCaseActorLoop(ch)
}
func (actor *Actor) StartIceCreamActor() {
	ch := make(chan Message, 16)
	actor.Ch = ch
	go actor.IceCreamActorLoop(ch)
}


func (actor *Actor) CustomerActorLoop(ch <-chan Message) {
	for{
		msg := <- ch
		switch msg.body{
		case "order":
			if msg.detail == "" {
				fmt.Println("Invalid order, please specify pie or icecream")
			}
			if msg.detail == "pie"  || msg.detail == "both"{
				actor.Organizer.SendMessage(actor.Name, "WaiterActor", Message{body: "order", detail: "pie", From: actor.Name})
		    }
			if msg.detail == "icecream" || msg.detail == "both"{
				actor.Organizer.SendMessage(actor.Name, "WaiterActor", Message{body: "order", detail: "icecream", From: actor.Name})
			}
		case "put on table":
			if msg.detail == "pie"{
				fmt.Println("Pie is on the table for customer: ", actor.Name)
			}
			if msg.detail == "icecream"{
				fmt.Println("Icecream is on the table for customer: ", actor.Name)
			}
		case "no pie left", "no icecream left":
			fmt.Println(actor.Name, msg.body)
		default:
			fmt.Println("Invalid message")

	}
}
}
func (actor *Actor) WaiterActorLoop(ch <-chan Message) {
	for{
		msg := <- ch
		switch msg.body{
		case "order":
            if msg.detail == "pie"{
				actor.Organizer.SendMessage(actor.Name, "PieCaseActor", Message{body: "get slice", From: msg.From})
			}
			if msg.detail == "icecream"{
				// for a quick and dirty solution, we will just copy the msg.from to the icecream actor. Since we're doing this so late we can't optimize the codde to make it cleaner
				actor.Organizer.SendMessage(actor.Name, "IceCreamActor", Message{body: "get scoop", From: msg.From})
			}
		case "Add To Order":
			if msg.detail == "pie"{
				fmt.Println("Adding pie to order to customer: ", msg.From)
			}
			if msg.detail == "icecream"{
				fmt.Println("Adding icecream to order to customer: ", msg.From)
			}
		case "no pie left":
			fmt.Println("No pie left")
	}
}
}
func (actor *Actor) PieCaseActorLoop(ch <-chan Message) {
	for{
		msg := <- ch
		switch msg.body{
		case "get slice":
			if len(actor.State.([]string)) > 0 {
				actor.Organizer.SendMessage(actor.Name, msg.From, Message{body: "put on table", detail: "pie"})
				actor.State = actor.State.([]string)[1:]
				actor.Organizer.SendMessage(actor.Name, "WaiterActor", Message{body: "Add To Order", detail: "pie", From: msg.From})
			} else {
				actor.Organizer.SendMessage(actor.Name, msg.From, Message{body: "no pie left"})
			}
		default:
			fmt.Println("Invalid message")
	}
}
}

func (actor *Actor) IceCreamActorLoop(ch <-chan Message) {
	for{
		msg := <- ch
		switch msg.body{
		case "get scoop":
			if len(actor.State.([]string)) > 0 {
				actor.Organizer.SendMessage(actor.Name, msg.From, Message{body: "put on table", detail: "icecream"})
				actor.State = actor.State.([]string)[1:]
				actor.Organizer.SendMessage(actor.Name, "WaiterActor", Message{body: "Add To Order", detail: "icecream", From: msg.From})
			} else {
				actor.Organizer.SendMessage(actor.Name, msg.From, Message{body: "no icecream left"})
			}
		default:
			fmt.Println("Invalid message")
	}
}
}
// Create a message organizer which registers the actors and also takes care of the message passing
func NewMessageOrganizer() *MessageOrganizer{
	return &MessageOrganizer{
		Actors: make(map[string]*Actor),
	}
}


// Register the actors by passing in a pointer to the actor
func (mo *MessageOrganizer) RegisterActor(actor *Actor) {
	mo.Actors[actor.Name] = actor
}

// SendMessage takes in the sender, receiver and the message to be sent, then routes it to the receiver
func (mo *MessageOrganizer) SendMessage(sender string, receiver string, msg Message) {
	// make sure the sender and receiver are registered
	go func() {
	if _, ok := mo.Actors[sender]; !ok {
		fmt.Println("Sender not registered", sender)
		return
	}
	if _, ok := mo.Actors[receiver]; !ok {
		fmt.Println("Receiver not registered: ", receiver)
		return
	}
	mo.Actors[receiver].Ch <- msg
	}()
}



func main() {
	// start pie case actor

    messageOrganizer := NewMessageOrganizer() 
	CustomerActor := NewActor("CustomerActor",nil,messageOrganizer) 
	CustomerActor.StartCustomerActor()
	messageOrganizer.RegisterActor(CustomerActor)

	CustomerActor2 := NewActor("CustomerActor2",nil,messageOrganizer) 
	CustomerActor2.StartCustomerActor()
	messageOrganizer.RegisterActor(CustomerActor2)

	// Waiteractor
	WaiterActor := NewActor("WaiterActor",nil,messageOrganizer)
	WaiterActor.StartWaiterActor()
	messageOrganizer.RegisterActor(WaiterActor)

	// PieCaseActor
	PieCaseActor := NewActor("PieCaseActor", []string{"apple", "cherry", "blueberry"},messageOrganizer)
	PieCaseActor.StartPieCaseActor()
	messageOrganizer.RegisterActor(PieCaseActor)

	// IceCreamActor
    IceCreamActor := NewActor("IceCreamActor", []string{"vanilla", "chocolate", "strawberry"},messageOrganizer)
	IceCreamActor.StartIceCreamActor()
	messageOrganizer.RegisterActor(IceCreamActor)


	CustomerActor.Ch <- Message{body: "order", detail: "pie"}
	CustomerActor2.Ch <- Message{body: "order", detail: "icecream"}

	CustomerActor2.Ch <- Message{body: "order", detail: "both"}
	CustomerActor.Ch <- Message{body: "order", detail: "both"}
	CustomerActor2.Ch <- Message{body: "order", detail: "both"}
	CustomerActor.Ch <- Message{body: "order", detail: "both"}

}
