package main

import (
	"os"
	"sync"
	"strconv"
)

type queue struct {
	list []string
}

var mu sync.Mutex

func main() {
	ch := make(chan string)
	signal := make(chan bool)
	wg := sync.WaitGroup{}
	qu := queue{list: []string{}}

	go qu.AutoDeque(signal, ch, &wg)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go qu.InQueue(signal, "test Element" + strconv.Itoa(i))
	}
	go saveLogs(ch, &wg)

	wg.Wait()
	close(ch)
	close(signal)
}

func (w *queue) AutoDeque(signal chan bool, chValue chan string, wg *sync.WaitGroup) {
	for {
		if x, ok := <-signal; x && ok {
			chValue <- w.list[0]
			w.list = w.list[1:]
		} else if !ok {
			break
		}
	}
}

func (w *queue) InQueue(signal chan bool, value string) {
	mu.Lock()
	w.list = append(w.list, value)
	signal <- true
	mu.Unlock()
}


func saveLogs(ch chan string, wg *sync.WaitGroup){
	for {
		if tmp, ok := <-ch; ok {
			f, _ := os.OpenFile("./logs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			defer f.Close()
			f.WriteString(tmp + "\n")
			wg.Done()
		} else {
			break
		}
	}
}