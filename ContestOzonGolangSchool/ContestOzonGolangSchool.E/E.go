package main

import (
	"fmt"
	"time"
	"sync"
)

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	var wg sync.WaitGroup 
	wg.Add(1)
	go Merge2Channels22(f, in1, in2, out, n,&wg)
	wg.Wait()
}

func waitChan(in <-chan int, f func(params ...int) int, wg *sync.WaitGroup) {
loopW:
	for {
		select {
		case res := <-in:
			f(res)
			wg.Done()
			fmt.Println("r", res)
			// println("w(", 2, "): ", res)
		default:
			time.Sleep(1)
			break loopW
		}
	}
}

func getsum() func(params ...int) int {
	sum := 0
	return func(params ...int) int {
		for _, num := range params {
			sum += num
		}
		return sum
	}
}

func mergeChan(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, key int) {
	var wg sync.WaitGroup
	wg.Add(2)
	fmt.Println("s", "run")
	// sumChan := make(chan int, 2)
	sumka := getsum()
	// keyL := key
	
	go waitChan(in1, sumka, &wg)
	go waitChan(in2, sumka, &wg)
	fmt.Println("s", "wait")
	wg.Wait()
	out <- sumka()

	fmt.Println("s", sumka())

// 	time.Sleep(1)

// 	cnt := 1
// 	sum := 0
// loopM:
// 	for {
// 		select {
// 		case res := <-sumChan:
// 			println("s(", keyL, "): ", sum, "+", res)
// 			sum += f(res)
// 			if cnt > 0 {

// 				println("s(", keyL, "): goto loop", cnt)
// 				cnt -= 1
// 				break loopM
// 			} else {

// 				close(sumChan)
// 				out <- sum
// 			}
// 		default:
// 			time.Sleep(2)
// 			println("s(", keyL, "): def")
// 			break loopM
// 		}

// 	}

}

func Merge2Channels22(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		mergeChan(f, in1, in2, out, i)
	}
	wg.Done()
}

func Merge2Channels2(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {

	for i := 0; i < n; i++ {
		println("run", i)

		var start = time.Now()
		var sum = 0
	loop:

		for {
			
			select {
			case res2 := <-in2:
				sum += f(res2)
				println(time.Since(start), "(2): ", res2)
			loop2:

				for {
					//var sum2 = sum
					select {

					case res := <-in1:
						sum += f(res)
						println(time.Since(start), "(1): ", res)
						out <- sum //f(x1) + f(x2)
						sum = 0
						println(time.Since(start), "(sum): ", sum)
					default:
						println(time.Since(start), "(sum): df")
						time.Sleep(6)
						break loop2
					}
				}

			default:
				println(time.Since(start), "(sum): df2")
				time.Sleep(6)
				break loop
			}
		}

		//x1 := <-in1
		//x2 := <-in2

	}
}

func f(x int) int {
	return x + 0
}

func pull(in1 chan<- int, x int, s int) {
	//time.Sleep(10)
	in1 <- x
}

func main() {
	time.Sleep(1)
	in1 := make(chan int, 10)
	in2 := make(chan int, 10)
	out := make(chan int, 10)
	pull(in1, 4, 5)  //in1 <- 4
	pull(in2, 5, 3)  //in2 <- 5
	pull(in1, 44, 5) //in1 <- 4
	pull(in2, 55, 3) //in2 <- 5

	Merge2Channels(f, in1, in2, out, 10)

	time.Sleep(200000)
	fmt.Println("done")
	// fmt.Println(<-out)
	// fmt.Println(<-out)

}
