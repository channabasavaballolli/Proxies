package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func runWorker(id int, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	fmt.Printf("Worker %d:Starting task...\n", id)
// 	time.Sleep(500 * time.Millisecond)
// 	fmt.Printf("Worker %d:Task  completed\n", id)

// }

// func main() {
// 	var wg sync.WaitGroup //pointer to waitgroup and go sets the internal waitgroup counter to 0

// 	for i := 1; i <= 3; i++ {
// 		wg.Add(1)            //to inform waitgroup 1 task is started (increases counter by 1)
// 		go runWorker(i, &wg) //starts worker in the background
// 	}
// 	fmt.Println("Main: Waiting for all workers to finish...")
// 	wg.Wait()
// 	fmt.Println("Main: All workers have finished")
// }

// // package main

// // import (
// // 	"fmt"
// // )

// // func main() {
// // 	ch := make(chan string)
// // 	fmt.Println("Now the msg will be recieved")
// // 	go func() {
// // 		ch <- "Greetings from the background goroutine!"
// // 	}()

// // 	msg := <-ch
// // 	fmt.Println("Recieved the message", msg)
// // }

// // package main

// // import (
// // 	"fmt"
// // 	"time"
// // )

// // func printNumbers() {
// // 	for i := 1; i <= 5; i++ {
// // 		fmt.Println("Number :", i)
// // 		time.Sleep(100 * time.Millisecond)
// // 	}
// // }

// // func main() {
// // 	go printNumbers()

// // 	fmt.Println("Hello from main")
// // 	time.Sleep(1 * time.Second)
// // 	fmt.Println("main finished")
// // }
