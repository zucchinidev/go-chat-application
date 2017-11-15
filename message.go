package main

import "time"

//The message type will encapsulate the message string itself, but we have also added
//the Name and When fields that respectively hold the user's name and a timestamp of
//when the message was sent.
type message struct {
	Name    string
	Message string
	When    time.Time
}
