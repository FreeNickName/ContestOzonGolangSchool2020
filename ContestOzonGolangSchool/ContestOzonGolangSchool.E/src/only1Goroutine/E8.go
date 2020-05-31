package only1Goroutine

type PkgName struct {}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	go Sum2ChannelsN(f, in1, in2, out, n)
}

func Sum2ChannelsN(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	i := 0
	for  {
		i++
		if i == n {
			return
		}

		sum := 0
		ok := false
		select {
			case sum, ok = <-in1:
				sum = f(<-in2) + f(sum)
			case sum, ok = <-in2:
				sum = f(<-in1) + f(sum)
		}
		if !ok {
			return
		}
		out <- sum
	}
}
