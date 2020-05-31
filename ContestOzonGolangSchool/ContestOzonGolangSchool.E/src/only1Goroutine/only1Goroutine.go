package only1Goroutine

type PkgName struct {}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	go Sum2Channels(f, in1, in2, out, n)
}

func Sum2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	for i := 0; i < n; i++  {
		sum := 0
		select {
			case res, ok := <-in1:
				// println(i, ok, res)
				if !ok {
					return
				}
				sum += f(<-in2) + f(res)
				// println(i, ok, res, sum)
			case res, ok := <-in2:
				// println(i, ok, res)
				if !ok {
					return
				}
				sum += f(<-in1) + f(res)
		}
		out <- sum
	}
}
