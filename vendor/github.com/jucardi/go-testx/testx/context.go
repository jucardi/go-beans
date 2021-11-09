package testx

import (
	"fmt"
	"github.com/jucardi/go-testx/assert"
	"github.com/jucardi/go-logger-lib/log"
	"github.com/jucardi/go-terminal-colors"
	"io"
	"strings"
	"testing"
)

const singleIndent = "    "

var (
	previous    = 0
	contextPile = []*context{newCtx(nil)}

	// Ensure implementation of these interfaces on build.
	_ assert.TestingT        = (*context)(nil)
	_ assert.IAssertsCounter = (*context)(nil)
	_ assert.IHelper         = (*context)(nil)
	_ assert.IFailNow        = (*context)(nil)
)

type context struct {
	assertions int
	failures   int
	level      log.Level
	offsetLn   int
	rows       int
	currentCol int
	t          *testing.T
}

func newCtx(t *testing.T) *context {
	row, _ := getSize()
	return &context{assertions: 0, t: t, rows: row}
}

func currentCtx() *context {
	return contextPile[len(contextPile)-1]
}

func (c *context) SprintIndent(str string, colors ...fmtc.Color) string {
	return c.Sprint(c.doIndent(str), colors...)
}

func (c *context) PrintIndent(str string, colors ...fmtc.Color) *context {
	return c.Print(c.doIndent(str), colors...)
}

func (c *context) Fprint(w io.Writer, str string, colors ...fmtc.Color) *context {
	_, _ = fmtc.WithColors(colors...).Fprint(w, str)
	return c
}

func (c *context) Println(str string, colors ...fmtc.Color) *context {
	_, _ = fmtc.WithColors(colors...).Println(c.doIndent(str))
	return c
}

func (c *context) Print(str string, colors ...fmtc.Color) *context {
	_, _ = fmtc.WithColors(colors...).Print(str)
	return c
}

func (c *context) Sprint(str string, colors ...fmtc.Color) string {
	return fmtc.New().Print(str, colors...).String()
}

func (c *context) FailMsgf(format string, args ...interface{}) {
	str := "✘"
	c.Print(str, fmtc.Red, fmtc.Bold)
	c.failures++
	msg := fmt.Sprintf(format, args...)
	c.PrintIndent(msg)
}

func (c *context) Increment() {
	str := "✔"
	c.Print(str, fmtc.Green, fmtc.Bold)
	c.assertions++
}

func (c *context) FailNow() {
	c.t.FailNow()
}

func (c *context) Helper() {
	c.t.Helper()
}

func (c *context) Fail() {
	c.t.Fail()
}

func (c *context) indent() string {
	return strings.Repeat(singleIndent, len(contextPile)-1)
}

func (c *context) doIndent(str string) string {
	return c.indent() + strings.Replace(str, "\n", "\n"+c.indent(), -1)
}
