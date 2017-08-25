package progressbar

import (
	"fmt"
	"time"
)

type bar struct {
	curPos    int
	taskCount int
	workers   int
}

const line = "                                                            "

func NewBar(taskCount int, workers int) chan<- struct{} {
	var ch = make(chan struct{})
	b := &bar{
		taskCount: taskCount,
		workers:   workers,
	}
	go b.consume(ch)
	go func() {
		for range time.NewTicker(100 * time.Millisecond).C {
			b.draw()
		}
	}()
	return ch
}

func (b *bar) consume(c <-chan struct{}) {
	for range c {
		b.curPos++
		// b.draw()
	}
}

func (b *bar) clearLine() {
	fmt.Print("\r" + line + "\r")
}

func (b *bar) draw() {
	b.clearLine()
	fmt.Printf("%v/%v", b.curPos, b.taskCount)
}
