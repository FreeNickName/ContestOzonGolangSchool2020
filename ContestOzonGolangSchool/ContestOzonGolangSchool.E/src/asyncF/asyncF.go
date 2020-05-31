package asyncF

type PkgName struct {}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	go func() {
		for i := 0; i < n; i++ {
			sum := 0
			select {
				case res, ok := <-in1:
					// println(i, ok, res)
					if !ok {
						panic("in1 is closed")
					}
					sum += f(<-in2) + f(res)
					// println(i, ok, res, sum)
				case res, ok := <-in2:
					// println(i, ok, res)
					if !ok {
						panic("in2 is closed")
					}
					sum += f(<-in1) + f(res)
			}
			out <- sum
		}
	}()
}
