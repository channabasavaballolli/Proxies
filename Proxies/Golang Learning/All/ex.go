package main

// import (
// 	"fmt"
// )

// func main() {
// 	ch := make(chan int)
// 	fmt.Println("Now the msg will be recieved")
// 	go func() {
// 		for i := 0; i <= 10; i++ {
// 			ch <- i
// 		}
// 		close(ch)
// 	}()

// 	for msg := range ch {
// 		fmt.Println("Recieved the message", msg)
// 	}
// }
