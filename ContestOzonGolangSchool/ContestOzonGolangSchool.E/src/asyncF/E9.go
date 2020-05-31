package asyncF

type PkgName struct {}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	f1 := make(chan int, 10)
	f2 := make(chan int, 10)

	go UseFChannels(f, in1, in2, f1, f2, n)

	go SumChannels(f1, f2, out)
}

func SumChannels(in1 <-chan int, in2 <-chan int, out chan<- int) {
	for  {
		sum := 0
		ok := false
		select {
			case sum, ok = <-in1:
				sum += <-in2
			case sum, ok = <-in2:
				sum += <-in1
		}
		if !ok {
			// println("push to out is done")
			// close(out)
			return
		}
		out <- sum
	}
}

func UseFChannels(f func(int) int, in1 <-chan int, in2 <-chan int, out1 chan<- int, out2 chan<- int, n int) {
	i := 0
	for  {
		if i == n {
			break
		}
		i++
		select {
			case res, ok := <-in1:
				if !ok {
					println("in1 closed")
					break
				} 
				out1 <- f(res)
				out2 <- f(<-in2)
			case res, ok := <-in2:
				if !ok {
					println("in2 closed")
					break
				} 
				out2 <- f(res)
				out1 <- f(<-in1)
		}
	}
	close(out1)
	close(out2)
}
