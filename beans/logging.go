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
	_, _ = fmt.Fprintln(os.Stderr, args...)
}

// Errorf logs a message at level Error.
func (c *ConsoleLogger) Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
}
