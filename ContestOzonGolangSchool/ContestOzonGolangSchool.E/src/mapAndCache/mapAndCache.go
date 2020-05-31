package mapAndCache

import (
	"sync"
)

type PkgName struct {}

type syncMap struct{
	sync.RWMutex
	m map[int]int
}

func CreateMap() *syncMap {
	return &syncMap{m: make(map[int]int)}
}

func (b *syncMap) Load(key int) (int, bool) {
    b.Lock()
    defer b.Unlock()
    val, ok := b.m[key]
    return val, ok
}

func (b *syncMap) Store(key int, value int) {
    b.Lock()
    defer b.Unlock()
    b.m[key] = value
}

func (b *syncMap) Delete(key int) {
    b.Lock()
    defer b.Unlock()
    delete(b.m, key)
}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	go ProcessMerge(f, in1, in2, out, n)
}

func ProcessMerge(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	f1 := make(chan int, 2)
	f2 := make(chan int, 2)
	b1 := CreateMap()
	b2 := CreateMap()
	s1 := make(chan bool, 2)
	s2 := make(chan bool, 2)
	cache := CreateMap()

	go ChannelToBuff(f, in1, n, b1, s1, cache)
	go ChannelToBuff(f, in2, n, b2, s2, cache)

	go ReadBuffByOrder(b1, f1, n, s1)
	go ReadBuffByOrder(b2, f2, n, s2)

	go SumChannels(f1, f2, out, n)
}

func SumChannels(in1 <-chan int, in2 <-chan int, out chan<- int, max int) {
	// defer println("SumChannels is done")

	for i := 0; i < max; i++ {
		sum := 0
		ok := false
		select {
			case sum, ok = <-in1:
				if ok {
					sum += <-in2
				}
			case sum, ok = <-in2:
				if ok {
					sum += <-in1
				}
		}
		if !ok {
			println("push to out is done")
			// close(out)
			return
		}
		out <- sum
	}
}

func ChannelToBuff(f func(int) int, in <-chan int, max int, b *syncMap, done chan<- bool, cache *syncMap) {
	// defer println("ChannelToBuff is done")
	for i := 0; i < max; i++ {
		res, ok := <-in
		if !ok {
			println("in closed")
			break
		}
		go UseFToBuff(f, i, res, b, done, cache)
	}
}

// После вычисления значения пишем в кэш, чтобы не считать повторно аналогичные значения
func UseFToBuff(f func(int) int, key int, val int, b *syncMap, done chan<- bool, cache *syncMap) {
	res, oks := cache.Load(val)
	if !oks {
		res = f(val)
		cache.Store(val, res)
	}	
	b.Store(key, res)
	// println("f:", key, "=", val)
	done <- true
	
}

// Запись мапа в канал по порядку ключей в asc, signal указывает наполнение мапа 
func ReadBuffByOrder(b *syncMap, out chan<- int, max int, signal <-chan bool) {
	// defer close(out)
	// defer println("ReadBuffByOrder is done")
	for i := 0; i < max; {
		select {
			case _, oks := <-signal:
				if (!oks) {
					println("signal closed")
					return
				}
				for {
					v, ok := b.Load(i)
					// println("b", i, "=", v)
					if !ok {
						break
					}	
					out <- v
					b.Delete(i)
					i++
				}
		}
	}
}
