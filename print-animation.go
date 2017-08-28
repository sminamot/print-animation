package printAnimation

import (
	"fmt"
	"time"

	termbox "github.com/nsf/termbox-go"
)

const (
	defaultDuration     = 50 * time.Millisecond
	defaultLastDuration = 500 * time.Millisecond
)

type PrintAnimation struct {
	printDuration       time.Duration
	lastDuration        time.Duration
	stringSlice         stringSlice
	terminateKey        termbox.Key
	isTerminateByAnyKey bool
}

type stringSlice []string

func New() *PrintAnimation {
	return &PrintAnimation{
		printDuration:       defaultDuration,
		lastDuration:        defaultLastDuration,
		stringSlice:         []string{},
		terminateKey:        termbox.KeyEsc,
		isTerminateByAnyKey: false,
	}
}

func (p *PrintAnimation) SetDuration(d time.Duration) {
	p.printDuration = d
}

func (p *PrintAnimation) SetLastDuration(d time.Duration) {
	p.lastDuration = d
}

func (p *PrintAnimation) AddString(s string) {
	p.stringSlice = append(p.stringSlice, s)
}

func (p *PrintAnimation) SetStrings(s []string) {
	p.stringSlice = s
}

func (p *PrintAnimation) Print() {
	for _, v := range p.stringSlice {
		fmt.Println(v)
	}
}

func (p *PrintAnimation) PrintAnimation() {
	max := p.maxLength()
	t := time.NewTicker(p.printDuration)
	defer t.Stop()

	ch := make(chan struct{})
	i := 0
	// 最後まで表示したら終了
	go func() {
		for i = 1; i <= max; i++ {
			select {
			case <-t.C:
				draw(p.outputStringSlice(i))
			}
		}
		time.Sleep(p.lastDuration)
		ch <- struct{}{}
	}()
	// Escキーが押されたら終了
	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					ch <- struct{}{}
				}
			}
		}
	}()
	time.Sleep(p.lastDuration)
	for {
		select {
		case <-ch:
			p.stringSlice = p.outputStringSlice(i)
			return
		}
	}
}

func (p *PrintAnimation) maxLength() (max int) {
	for _, s := range p.stringSlice {
		if max < len(s) {
			max = len(s)
		}
	}
	return
}

func (p *PrintAnimation) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i, v := range p.stringSlice {
		drawLine(0, i, v)
	}

	termbox.Flush()
}

func (p *PrintAnimation) outputStringSlice(target int) stringSlice {
	v := make(stringSlice, len(p.stringSlice))
	for i, vv := range p.stringSlice {
		v[i] = vv
		if len(vv) > target {
			v[i] = vv[:target]
		}
	}
	return v
}

func draw(ss stringSlice) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i, v := range ss {
		drawLine(0, i, v)
	}

	termbox.Flush()
}

func drawLine(x, y int, str string) {
	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault
	runes := []rune(str)

	for i := 0; i < len(runes); i += 1 {
		termbox.SetCell(x+i, y, runes[i], color, backgroundColor)
	}
}
