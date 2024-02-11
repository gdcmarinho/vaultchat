package main

import (
	"fmt"
	"time"
)

type User struct {
	UUID			string
	PublicKey		string
}

type Session struct {
	UUID			string
	Participants	[]string
}

type Message struct {
	UUID			string
	SessionUUID		string
	SenderUUID		string
	RecipientUUID	string
	Content			string
	SentAt			time.Time
	Signature		string
}

func main() {

}