package assert

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
}

type IAssertLogger interface {
	FailMsgf(format string, args ...interface{})
}

type IAssertsCounter interface {
	Increment()
}

type IFailNow interface {
	FailNow()
}

type IHelper interface {
	Helper()
}

// ComparisonAssertionFunc is a common function prototype when comparing two values.  Can be useful
// for table driven tests.
type ComparisonAssertionFunc func(TestingT, interface{}, interface{}, ...interface{}) bool

// ValueAssertionFunc is a common function prototype when validating a single value.  Can be useful
// for table driven tests.
type ValueAssertionFunc func(TestingT, interface{}, ...interface{}) bool

// BoolAssertionFunc is a common function prototype when validating a bool value.  Can be useful
// for table driven tests.
type BoolAssertionFunc func(TestingT, bool, ...interface{}) bool

// ErrorAssertionFunc is a common function prototype when validating an error value.  Can be useful
// for table driven tests.
type ErrorAssertionFunc func(TestingT, error, ...interface{}) bool

// EvalFunc a custom function that returns true on success and false on failure
type EvalFunc func() (success bool)

// PanicTestFunc defines a func that should be passed to the assert.Panics and assert.NotPanics
// methods, and represents a simple func that takes no arguments, and returns nothing.
type PanicTestFunc func()
