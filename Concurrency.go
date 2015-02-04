package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	c := boring("Vimlesh")
	for i := 0; i <= 5; i++ {
		fmt.Println(<-c)
	}
	fmt.Println("Exiting......")
}

//Generator: function that returns a channel
//Making use of a closure..
func boring(msg string) <-chan string { // Returns receive-only channel of strings.
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i) // Expression to be sent can be any suitable value.
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // Return the channel to the caller.
}
