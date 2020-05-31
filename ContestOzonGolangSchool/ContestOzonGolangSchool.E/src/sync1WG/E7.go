package sync1WG

import (
	"sync"
)

type PkgName struct {}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	// println()
	go Sum2ChannelsN(f, in1, in2, out, n)
}

func SumChannels(f func(int) int, in1 int, in2 int, out chan<- int, key int, wg *sync.WaitGroup, max int) {
	defer wg.Done()
	out <- f(in1) + f(in2)
	// println(key, ":", res)
}

func Sum2ChannelsN(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	i := 0
	for  {
		currCalcWg := new(sync.WaitGroup)
		currCalcWg.Add(1)

		x1, ok1 := <-in1
		x2, ok2 := <-in2
		if !ok1 || !ok2 {
			return
		}
		go SumChannels(f, x1, x2, out, i, currCalcWg, n)
        currCalcWg.Wait()
		i++
		if i == n {
			return
		}
	}
}
