package hclog

import (
	"bytes"
	"fmt"
	"io"
)

type writer struct {
	b     bytes.Buffer
	w     io.Writer
	color ColorOption
}

func newWriter(w io.Writer, color ColorOption) *writer {
	return &writer{w: w, color: color}
}

func (w *writer) Flush(level Level) (err error) {
	var unwritten = w.b.Bytes()

	fmt.Printf("||||||| Flush: w.b.Bytes() = %+v\n", string(w.b.Bytes()))

	if w.color != ColorOff {
		color := _levelToColor[level]
		unwritten = []byte(color.Sprintf("%s", unwritten))
	}

	if lw, ok := w.w.(LevelWriter); ok {
		_, err = lw.LevelWrite(level, unwritten)
	} else {
		_, err = w.w.Write(unwritten)
	}
	w.b.Reset()
	return err
}

func (w *writer) Write(p []byte) (int, error) {
	fmt.Printf("||||||| Write: p = %+v\n", string(p))

	// lg.Log(Trace, "||||||| "+string(p))

	// if strings.Contains(string(b), "||||") {
	// 	fmt.Printf("|||| string(b) = %s\n", string(b))

	// 	l.l.Println(string(bytes.TrimRight(b, " \n\t")))
	// }

	return w.b.Write(p)
}

func (w *writer) WriteByte(c byte) error {
	return w.b.WriteByte(c)
}

func (w *writer) WriteString(s string) (int, error) {
	return w.b.WriteString(s)
}

// LevelWriter is the interface that wraps the LevelWrite method.
type LevelWriter interface {
	LevelWrite(level Level, p []byte) (n int, err error)
}

// LeveledWriter writes all log messages to the standard writer,
// except for log levels that are defined in the overrides map.
type LeveledWriter struct {
	standard  io.Writer
	overrides map[Level]io.Writer
}

// NewLeveledWriter returns an initialized LeveledWriter.
//
// standard will be used as the default writer for all log levels,
// except for log levels that are defined in the overrides map.
func NewLeveledWriter(standard io.Writer, overrides map[Level]io.Writer) *LeveledWriter {
	return &LeveledWriter{
		standard:  standard,
		overrides: overrides,
	}
}

// Write implements io.Writer.
func (lw *LeveledWriter) Write(p []byte) (int, error) {
	fmt.Printf("||| LeveledWriter Write: p = %+v\n", string(p))

	return lw.standard.Write(p)
}

// LevelWrite implements LevelWriter.
func (lw *LeveledWriter) LevelWrite(level Level, p []byte) (int, error) {
	fmt.Printf("||| LevelWrite Write: p = %+v\n", string(p))

	w, ok := lw.overrides[level]
	if !ok {
		w = lw.standard
	}
	return w.Write(p)
}
