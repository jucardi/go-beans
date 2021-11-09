package testx

import (
	"fmt"
	"github.com/jucardi/go-iso8601"
	"github.com/jucardi/go-logger-lib/log"
	"github.com/jucardi/go-terminal-colors"
	"io"
	"os"
	"time"
)

type levelScheme struct {
	writer    io.Writer
	level     []fmtc.Color
	timestamp []fmtc.Color
	message   []fmtc.Color
}

var (
	_ log.ILogger = (*context)(nil)

	cfgDebug = levelScheme{
		writer:    os.Stdout,
		level:     []fmtc.Color{fmtc.BgDarkGray, fmtc.White, fmtc.Bold},
		timestamp: []fmtc.Color{fmtc.Gray, fmtc.Bold},
		message:   []fmtc.Color{fmtc.Gray},
	}
	cfgInfo = levelScheme{
		writer:    os.Stdout,
		level:     []fmtc.Color{fmtc.BgBlue, fmtc.White, fmtc.Bold},
		timestamp: []fmtc.Color{fmtc.Cyan, fmtc.Bold},
		message:   []fmtc.Color{fmtc.Cyan},
	}
	cfgWarn = levelScheme{
		writer:    os.Stderr,
		level:     []fmtc.Color{fmtc.BgYellow, fmtc.Black, fmtc.Bold},
		timestamp: []fmtc.Color{fmtc.Yellow, fmtc.Bold},
		message:   []fmtc.Color{fmtc.Yellow},
	}
	cfgError = levelScheme{
		writer:    os.Stderr,
		level:     []fmtc.Color{fmtc.BgRed, fmtc.White, fmtc.Bold},
		timestamp: []fmtc.Color{fmtc.Red, fmtc.Bold},
		message:   []fmtc.Color{fmtc.Red},
	}
	config = map[log.Level]levelScheme{
		log.DebugLevel: cfgDebug,
		log.InfoLevel:  cfgInfo,
		log.WarnLevel:  cfgWarn,
		log.ErrorLevel: cfgError,
		log.FatalLevel: cfgError,
		log.PanicLevel: cfgError,
	}
)

func Log() log.ILogger {
	return currentCtx()
}

func (c *context) Name() string {
	return "Convey-Context-Logger"
}

func (c *context) SetLevel(level log.Level) {
	c.level = level
}

func (c *context) GetLevel() log.Level {
	return c.level
}

func (c *context) Debug(args ...interface{}) {
	c.logEntry(log.DebugLevel, fmt.Sprint(args...))
}

func (c *context) Debugf(format string, args ...interface{}) {
	c.logEntry(log.DebugLevel, fmt.Sprintf(format, args...))
}

func (c *context) Info(args ...interface{}) {
	c.logEntry(log.InfoLevel, fmt.Sprint(args...))
}

func (c *context) Infof(format string, args ...interface{}) {
	c.logEntry(log.InfoLevel, fmt.Sprintf(format, args...))
}

func (c *context) Warn(args ...interface{}) {
	c.logEntry(log.WarnLevel, fmt.Sprint(args...))
}

func (c *context) Warnf(format string, args ...interface{}) {
	c.logEntry(log.WarnLevel, fmt.Sprintf(format, args...))
}

func (c *context) Error(args ...interface{}) {
	c.logEntry(log.ErrorLevel, fmt.Sprint(args...))
}

func (c *context) Errorf(format string, args ...interface{}) {
	c.logEntry(log.ErrorLevel, fmt.Sprintf(format, args...))
}

func (c *context) Fatal(args ...interface{}) {
	c.logEntry(log.FatalLevel, fmt.Sprint(args...))
}

func (c *context) Fatalf(format string, args ...interface{}) {
	c.logEntry(log.FatalLevel, fmt.Sprintf(format, args...))
}

func (c *context) Panic(args ...interface{}) {
	c.logEntry(log.PanicLevel, fmt.Sprint(args...))
}

func (c *context) Panicf(format string, args ...interface{}) {
	c.logEntry(log.PanicLevel, fmt.Sprintf(format, args...))
}

func (c *context) SetFormatter(formatter log.IFormatter) {
	panic("not supported yet")
}

func (c *context) logEntry(level log.Level, str string) {
	if c.level > level {
		return
	}

	cfg := config[level]

	c.Fprint(cfg.writer, "\n").
		Fprint(cfg.writer, c.indent()+"  ").
		Fprint(cfg.writer, fmt.Sprintf(" %s ", level), cfg.level...).
		Fprint(cfg.writer, fmt.Sprintf(" %s ", iso8601.TimeToString(time.Now(), "HH:mm:ss")), cfg.timestamp...).
		Fprint(cfg.writer, "- ").
		Fprint(cfg.writer, str, cfg.message...)

	if level == log.FatalLevel {
		c.FailNow()
	}
	if level == log.PanicLevel {
		panic("failure")
	}
}
