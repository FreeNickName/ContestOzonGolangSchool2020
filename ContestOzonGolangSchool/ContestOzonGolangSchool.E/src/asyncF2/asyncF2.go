package asyncF2

import "sync"

type PkgName struct {}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	// println(n, len(in1), len(in2), len(out))
	go func() {
		defer close(out)
		// defer println(n, len(in1), len(in2), len(out))
		var wg sync.WaitGroup
		var x1, x2 int
		var ok1, ok2 bool
		for i := 0; i < n; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				var res int
				res, ok1 = <-in1
				if ok1 {
					x1 = f(res)
				}
			}()
			go func() {
				defer wg.Done()
				var res int
				res, ok2 = <-in2
				if ok2 {
					x2 = f(res)
				}
			}()
			wg.Wait()
			if !(ok1 && ok2) {
				println(ok1, ok2)
				return
			}
			out <- x1 + x2
		}
	}()
}
