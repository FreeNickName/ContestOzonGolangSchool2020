package map2out

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
	b1 := CreateMap()
	b2 := CreateMap()
	s := make(chan bool, 5)

	go ChannelToBuff(f, in1, n, b1, s)
	go ChannelToBuff(f, in2, n, b2, s)
	go ReadBuffByOrder(b1, b2, out, n, s)
}

func ChannelToBuff(f func(int) int, in <-chan int, max int, b *syncMap, done chan<- bool) {
	for i := 0; i < max; i++ {
		res, ok := <-in
		if !ok {
			panic("in is closed")
			// break
		}
		go UseFToBuff(f, i, res, b, done)
	}
}

func UseFToBuff(f func(int) int, key int, val int, b *syncMap, done chan<- bool) {
	b.Store(key, f(val))
	// println("f:", key, "=", val)
	done <- true
}

// Запись мапа в канал по порядку ключей в asc, signal указывает наполнение мапа 
func ReadBuffByOrder(b1 *syncMap, b2 *syncMap, out chan<- int, max int, signal <-chan bool) {
	// defer close(out)
	for i := 0; i < max; {
		_, oks := <-signal
		if (!oks) {
			panic("signal closed")
			// return
		}
		for {
			if i == max {
				break
			}
			v1, ok1 := b1.Load(i)
			// println("b1", i, "=", v1)
			if !ok1 {
				break
			}
			v2, ok2 := b2.Load(i)
			// println("b2", i, "=", v2)
			if !ok2 {
				break
			}
			out <- v1 + v2
			b1.Delete(i)
			b2.Delete(i)
			i++
		}
	}
}
