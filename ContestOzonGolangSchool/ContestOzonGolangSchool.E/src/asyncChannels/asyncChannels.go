package asyncChannels

type PkgName struct {}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	go ProcessMerge(f, in1, in2, out, n)
}

func ProcessMerge(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	f1 := make(chan int, 10)
	f2 := make(chan int, 10)

	go UseFChannel(f, in1, f1, n)
	go UseFChannel(f, in2, f2, n)
	go SumChannels(f1, f2, out, n)
}

func SumChannels(in1 <-chan int, in2 <-chan int, out chan<- int, max int) {
	i := 0
	for ; i < max; {
		sum := 0
		ok := false
		select {
			case sum, ok = <-in1:
				if ok {
					sum += <-in2
				} else {
					println("in1 is closed")
				}
			case sum, ok = <-in2:
				if ok {
					sum += <-in1
				} else {
					println("in2 is closed")
				}
		}
		if ok {
			out <- sum
			i++
		}
	}
	// defer println("SumChannels is done", i)
}

func UseFChannel(f func(int) int, in <-chan int, out chan<- int, max int) {
	defer close(out)

	i := 0
	for ; i < max; i++ {
		select {
			case res, ok := <-in:
				if !ok {
					panic("in closed")
					// return
				} 
				out <- f(res)
		}
	}
	// defer println("UseFChannel is done", i)
}
