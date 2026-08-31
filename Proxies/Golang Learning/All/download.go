package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func main() {
// 	var wg sync.WaitGroup
// 	urls := []string{
// 		"https://google.com",
// 		"https://golang.org",
// 		"https://github.com",
// 	}
// 	for _, url := range urls {
// 		wg.Add(1)
// 		go download(url, &wg)
// 	}
// 	wg.Wait()
// 	fmt.Println("All downloads completed")
// }

// func download(url string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	fmt.Printf("Downloading from %s \n", url)
// 	time.Sleep(800 * time.Millisecond)
// 	fmt.Printf("finished downloading from %s \n", url)

// }
