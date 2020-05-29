package main

import (
	"sync"
	"time"
)

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	go Sum2ChannelsN(f, in1, in2, out, n)
}

func ReadChannel(in <-chan int, f func(v int), wg *sync.WaitGroup) {
// loopW:
// 	for {
// 		select {
// 		case
		 res := <-in
		// :
			f(res)
			wg.Done()			
	// 		return
	// 	default:
	// 		time.Sleep(1)
	// 		break loopW
	// 	}
	// }
}

func GetAdder() func(params ...int) int {
	sum := int(0)
	return func(params ...int) (ret int) {
		for _, num := range params {
			sum += num
		}
		ret = sum
		return
	}
}

func SumChannels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int) {
	var wg sync.WaitGroup
	wg.Add(2)

	adder := GetAdder()
	proccessor := func(v int) {
		adder(f(v))
	}
	
	go ReadChannel(in1, proccessor, &wg)
	go ReadChannel(in2, proccessor, &wg)

	wg.Wait()
	out <- adder()
}

func Sum2ChannelsN(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	for i := 0; i < n; i++ {
		SumChannels(f, in1, in2, out)
	}
}

func f(x int) int {
	time.Sleep(10000)
	sum := 0
	for i := 0; i < 10000000; i++ {
		sum = i
	}
	sum -= 1000000

	return x + 2
}

func pull(in1 chan<- int, x int, s int) {
	in1 <- x
}

func makeEvenGenerator() func() uint {
    i := uint(0)
    return func() (ret uint) {
        ret = i
        i += 2
        return
    }
}
func main() {
    // nextEven := makeEvenGenerator()
    // println(nextEven()) // 0
	// println(nextEven()) // 2
	
	// nextEven = makeEvenGenerator()
    // println(nextEven()) // 0
    // println(nextEven()) // 2

	in1 := make(chan int, 10)
	in2 := make(chan int, 10)
	out := make(chan int, 10)
	pull(in1, 4, 5)  //in1 <- 4
	pull(in2, 5, 3)  //in2 <- 5

	pull(in1, 13, 5) //in1 <- 4
	pull(in2, 7, 3) //in2 <- 5

	pull(in1, 1234243, 5) //in1 <- 4
	pull(in2, 724243, 3) //in2 <- 5

	pull(in1, 1342433, 5) //in1 <- 4
	pull(in2, 723232, 3) //in2 <- 5

	pull(in1, 123233, 5) //in1 <- 4
	pull(in2, 75656, 3) //in2 <- 5

	pull(in1, 13000, 5) //in1 <- 4
	pull(in2, 65657, 3) //in2 <- 5

	pull(in1, 1300076, 5) //in1 <- 4
	pull(in2, 6565712, 3) //in2 <- 5

	pull(in1, 1303400, 5) //in1 <- 4
	pull(in2, 6563457, 3) //in2 <- 5

	Merge2Channels(f, in1, in2, out, 10)

	time.Sleep(1)
	// var wg sync.WaitGroup
	// wg.Add(10)
	println("M:wait")
// loop:	
	for {
		select {
		case res := <-out:
			println("M:", res)
			time.Sleep(2000)
			// wg.Done()
			// break loop
		default:
			time.Sleep(200)
			println("M:...")
			// break loop
		}
	}
//	wg.Wait()
}
