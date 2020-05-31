package sync2WG

import (
	"sync"
)

type PkgName struct {}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	// println()
	go Sum2ChannelsN(f, in1, in2, out, n)
}

func SumChannels(f func(int) int, in1 int, in2 int, out chan<- int, key int, wg *sync.WaitGroup, pwg *sync.WaitGroup, max int) {
	res := f(in1) + f(in2)
    if key > 0 {
		pwg.Wait()
	}
	out <- res
	defer wg.Done()
	// println(key, ":", res)
	if key == max - 1 {
		close(out)
	}
}

func Sum2ChannelsN(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	prevCalcWg := new(sync.WaitGroup)
	i := 0
	for  {
		currCalcWg := new(sync.WaitGroup)
		currCalcWg.Add(1)

		x1 := <-in1
		x2 := <-in2
		go SumChannels(f, x1, x2, out, i, currCalcWg, prevCalcWg, n)
		prevCalcWg = currCalcWg
		i++
		if i == n {
			return
		}
	}
}
