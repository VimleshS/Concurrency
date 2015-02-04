package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	msg  string
	wait chan bool
}

func main() {
	c := fanIn(boring("Ann"), boring("Joe"))
	for i := 0; i < 10; i++ {
		msg1 := <-c
		fmt.Println(msg1.msg)
		msg2 := <-c
		fmt.Println(msg2.msg)

		//its value can be anything.. just a trigger point
		msg1.wait <- true
		msg2.wait <- true
	}
	fmt.Println("Exiting......")
}

//https://talks.golang.org/2012/concurrency.slide#28
func fanIn(input1, input2 <-chan Message) <-chan Message {
	c := make(chan Message)
	//go func() {
	//	for {
	//		c <- <-input1
	//	}
	//}()
	//go func() {
	//	for {
	//		c <- <-input2
	//	}
	//}()

	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			case <-time.After(1 * time.Second):
				fmt.Println("You are too slow")
				return
			}
		}
	}()
	return c
}

//Generator: function that returns a channel
//Making use of a closure..
func boring(msg string) <-chan Message {
	c := make(chan Message)

	//single instance of waitForIt shared across all goroutines.. closures
	waitForIt := make(chan bool)
	go func() {
		for i := 0; ; i++ {
			c <- Message{fmt.Sprintf("%s %d", msg, i), waitForIt}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			<-waitForIt
		}
	}()
	return c // Return the channel to the caller.
}
