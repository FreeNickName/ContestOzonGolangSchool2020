package sync2WG

import (
	"sync"
)

type PkgName struct {}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	go ProcessMerge(f, in1, in2, out, n)
}

func SumChannels(f func(int) int, in1 int, in2 int, out chan<- int, idx int, cwg *sync.WaitGroup, pwg *sync.WaitGroup, max int) {
	defer cwg.Done()
	res := f(in1) + f(in2)
    if idx > 0 {
		pwg.Wait()
	}
	out <- res
	// println(idx, ":", res)
	// if idx == max {
	// 	close(out)
	// }
}

func ProcessMerge(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	var prevWg *sync.WaitGroup
	i := 0
	for ; i < n; i++ {
		currWg := new(sync.WaitGroup)
		currWg.Add(1)

		x1 := <-in1
		x2 := <-in2
		go SumChannels(f, x1, x2, out, i, currWg, prevWg, n)
		prevWg = currWg
	}
}
