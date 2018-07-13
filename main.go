package main

import (
  "fmt"
  "sync"
  "time"
)


// just playing with channels

var wg sync.WaitGroup
var ch chan int = make(chan int, 30)

func consume(i int) {
  defer wg.Done()

  for node := range ch {

    fmt.Println(i, node)
  }
}


func main() {

  wg.Add(8)
  go func() {
    defer wg.Done()

    for i := 0; i < 100; i++ {
      ch <- i
      if 0 == i % 50 {
        time.Sleep(time.Second * 10)
      }
    }
    close(ch)
  }()



  for i := 0; i < 7; i++ {
    go consume(i)
  }

  wg.Wait()

  fmt.Println("Done")
}

//func main() {
//  ch := make(chan int, 2)
//
//  go func() {
//    for i := 0; i < 10; i++ {
//      ch <- i
//    }
//  }()
//
//  for i := 0; i < 11; i++ {
//    fmt.Println(<- ch)
//  }
//}