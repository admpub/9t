package main

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/ActiveState/tail"
	"github.com/mattn/go-colorable"
)

var (
	seekInfoOnStart = &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}
	colorableOutput = colorable.NewColorableStdout()
)

type tailer struct {
	*tail.Tail
	colorCode int
	maxWidth  int
}

func newTailer(filename string, colorCode int, maxWidth int) (*tailer, error) {
	t, err := tail.TailFile(filename, tail.Config{
		Follow:   true,
		Location: seekInfoOnStart,
		Logger:   tail.DiscardingLogger,
	})

	if err != nil {
		return nil, err
	}

	return &tailer{
		Tail:      t,
		colorCode: colorCode,
		maxWidth:  maxWidth,
	}, nil
}

func (t tailer) do(wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range t.Lines {
		fmt.Fprintf(colorableOutput, "\x1b[%dm%*s\x1b[0m: %s\n", t.colorCode, t.maxWidth, t.name(), line.Text)
	}
}

func (t tailer) name() string {
	return path.Base(t.Filename)
}
