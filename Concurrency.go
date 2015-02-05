package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web   = fakeSearch("Web")
	Image = fakeSearch("Image")
	Video = fakeSearch("Video")
)

type Result string
type Search func(query string) Result

func Google(query string) (results []Result) {
	/* Synchronus call to get search result */
	//results = append(results, Web(query))
	//results = append(results, Image(query))
	//results = append(results, Video(query))
	//return

	/*  Go routines  */

	c := make(chan Result)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Image(query) }()

	timeout := time.After(80 * time.Millisecond)

	for i := 0; i < 3; i++ {
		//result := <-c
		//results = append(results, result)
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

//Replicating the servers..
func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}

type Message struct {
	msg  string
	wait chan bool
}

func f(left, right chan int) {
	fmt.Printf("%p   %p", left, right)
	fmt.Println()
	left <- 1 + <-right
}

func main() {

	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)

	//Getting the first output..
	//rand.Seed(time.Now().UnixNano())
	//start := time.Now()
	//result := First("golang",
	//	fakeSearch("replica 1"),
	//	fakeSearch("replica 2"))
	//elapsed := time.Since(start)
	//fmt.Println(result)
	//fmt.Println(elapsed)

	//const n = 10
	//leftmost := make(chan int)
	//right := leftmost
	//left := leftmost
	//for i := 0; i < n; i++ {
	//	right = make(chan int)
	//	go f(left, right)
	//	left = right
	//}

	//go func(c chan int) {
	//	c <- 1
	//}(right)

	//fmt.Println(<-leftmost)
	//fmt.Printf("%p   %p   %p", left, right, leftmost)
}

//https://talks.golang.org/2012/concurrency.slide#28
//Also Know as Multiplexer
func fanIn(input1, input2 <-chan Message) <-chan Message {
	c := make(chan Message)
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
