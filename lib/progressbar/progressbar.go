package progressbar

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type bar struct {
	curPos    int
	taskCount int
	workers   int
}

const line = "                                                            "

func NewBar(taskCount int, workers int) (chan<- struct{}, func()) {
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
	return ch, b.done
}

func (b *bar) done() {
	b.draw()
	fmt.Println(" ✅")
}

func (b *bar) consume(c <-chan struct{}) {
	for range c {
		b.curPos++
	}
}

func (b *bar) clearLine() {
	fmt.Print("\r" + line + "\r")
}

func (b *bar) draw() {
	b.clearLine()
	fmt.Printf("%v/%v", b.curPos, b.taskCount)
}

func maxWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	sttyOut, err := cmd.Output()
	if err != nil {
		return 80
	}
	width, err := strconv.Atoi(strings.Split(string(sttyOut), " ")[0])
	if err != nil {
		return 80
	}
	return width
}
