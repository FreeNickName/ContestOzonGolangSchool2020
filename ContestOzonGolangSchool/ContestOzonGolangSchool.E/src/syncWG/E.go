package syncWG

import (
	"sync"
)

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	// println()
	go Sum2ChannelsN(f, in1, in2, out, n)
}

func SumChannels(f func(int) int, in1 int, in2 int, out chan<- int, key int, wg *sync.WaitGroup, pwg *sync.WaitGroup) {
	res := f(in1) + f(in2)
    if key > 0 {
		pwg.Wait()
	}
	out <- res
	defer wg.Done()
	// println(key, ":", res)
}

func Sum2ChannelsN(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	prevChanWg := new(sync.WaitGroup)
	prevCalcWg := new(sync.WaitGroup)
	i := 0
	for  {
		currChanWg := new(sync.WaitGroup)
		currChanWg.Add(1)
		currCalcWg := new(sync.WaitGroup)
		currCalcWg.Add(1)

		go func(key int, wg *sync.WaitGroup, pwg *sync.WaitGroup, cwg *sync.WaitGroup, pcwg *sync.WaitGroup) {
			if key > 0 {
				pwg.Wait()
			}
			x1 := <-in1
			x2 := <-in2
			wg.Done()
			go SumChannels(f, x1, x2, out, key, cwg, pcwg)
		}(i, currChanWg, prevChanWg, currCalcWg, prevCalcWg)
		prevCalcWg = currCalcWg
		prevChanWg = currChanWg
		i++
	}
}
