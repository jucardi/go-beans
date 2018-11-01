package beans

import (
	"fmt"
	"os"
)

// IBeanLogger is a simple interface that represents a logger so any logger that matches this interface
// can be hooked into the bean factory by doing:  bean.SetLogger(log)
//
// For example, to attach github.com/sirupsen/logrus:
//
//        bean.SetLogger(logrus.StandardLogger())
//
type IBeanLogger interface {
	// Error logs a message at level Error.
	Error(args ...interface{})

	// Errorf logs a message at level Error.
	Errorf(format string, args ...interface{})
}

type emptyLogger struct {
}

func (e *emptyLogger) Error(args ...interface{}) {
}

func (e *emptyLogger) Errorf(format string, args ...interface{}) {
}

// ConsoleLogger is a simple implementation of a logger that will print messages to the console
type ConsoleLogger struct {
}

// Error logs a message at level Error.
func (c *ConsoleLogger) Error(args ...interface{}) {
	argx := append([]interface{}{"[ ERROR ] "}, args)
	fmt.Fprintln(os.Stderr, argx...)
}

// Errorf logs a message at level Error.
func (c *ConsoleLogger) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[ ERROR ] "+format, args...)
	fmt.Fprintln(os.Stderr)
}
