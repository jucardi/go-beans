package testx

import (
	"github.com/jucardi/go-testx/assert"
	"time"
)

// OnFailure registers a callback that will be executed if a test fails
func OnFailure(handler func()) {
	assert.OnFailure(handler)
}

// FailNow fails test
func FailNow(failureMessage string, msgAndArgs ...interface{}) bool {
	return assert.FailNow(currentCtx(), failureMessage, msgAndArgs...)
}

// Fail reports a failure through
func Fail(failureMessage string, msgAndArgs ...interface{}) bool {
	return assert.Fail(currentCtx(), failureMessage, msgAndArgs...)
}

// ShouldImplement asserts that an object is implemented by the specified interface.
//
//    ShouldImplement((*MyInterface)(nil), new(MyObject))
//
func ShouldImplement(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	return assert.Implements(currentCtx(), interfaceObject, object, msgAndArgs...)
}

// IsType asserts that the specified objects are of the same type.
func IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	return assert.IsType(currentCtx(), expectedType, object, msgAndArgs...)
}

// ShouldEqual asserts that two objects are equal.
//
//    ShouldEqual(123, 123)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func ShouldEqual(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.Equal(currentCtx(), expected, actual, msgAndArgs ...)
}

// EqualValues asserts that two objects are equal or convertable to the same types
// and equal.
//
//    EqualValues(uint32(123), int32(123))
//
func ShouldEqualValues(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.EqualValues(currentCtx(), expected, actual, msgAndArgs...)
}

// ShouldMatchExactly asserts that two objects are equal in value and type.
//
//    ShouldMatchExactly(int32(123), int64(123))
//
func ShouldMatchExactly(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.Exactly(currentCtx(), expected, actual, msgAndArgs...)
}

// ShouldNotBeNil asserts that the specified object is not nil.
//
//    ShouldNotBeNil(err)
//
func ShouldNotBeNil(object interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotNil(currentCtx(), object, msgAndArgs...)
}

// ShouldBeNil asserts that the specified object is nil.
//
//    ShouldBeNil(err)
//
func ShouldBeNil(object interface{}, msgAndArgs ...interface{}) bool {
	return assert.Nil(currentCtx(), object, msgAndArgs...)
}

// ShouldBeEmpty asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//    ShouldBeEmpty(obj)
//
func ShouldBeEmpty(object interface{}, msgAndArgs ...interface{}) bool {
	return assert.Empty(currentCtx(), object, msgAndArgs...)
}

// ShouldNotBeEmpty asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if ShouldNotBeEmpty(obj) {
//     ShouldBeEqual("two", obj[1])
//  }
//
func ShouldNotBeEmpty(object interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotEmpty(currentCtx(), object, msgAndArgs...)
}

// ShouldLen asserts that the specified object has specific length.
// ShouldLen also fails if the object has a type that len() not accept.
//
//    ShouldLen(mySlice, 3)
//
func ShouldLen(object interface{}, length int, msgAndArgs ...interface{}) bool {
	return assert.Len(currentCtx(), object, length, msgAndArgs...)
}

// ShouldBeTrue asserts that the specified value is true.
//
//    ShouldBeTrue(myBool)
//
func ShouldBeTrue(value bool, msgAndArgs ...interface{}) bool {
	return assert.True(currentCtx(), value, msgAndArgs...)
}

// ShouldBeFalse asserts that the specified value is false.
//
//    ShouldBeFalse(myBool)
//
func ShouldBeFalse(value bool, msgAndArgs ...interface{}) bool {
	return assert.False(currentCtx(), value, msgAndArgs...)
}

// ShouldNotBeEqual asserts that the specified values are NOT equal.
//
//    ShouldNotBeEqual(obj1, obj2)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func ShouldNotBeEqual(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotEqual(currentCtx(), expected, actual, msgAndArgs...)
}

// ShouldContain asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    ShouldContain("Hello World", "World")
//    ShouldContain(["Hello", "World"], "World")
//    ShouldContain({"Hello": "World"}, "Hello")
//
func ShouldContain(s, contains interface{}, msgAndArgs ...interface{}) bool {
	return assert.Contains(currentCtx(), s, contains, msgAndArgs...)
}

// ShouldNotContain asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    ShouldNotContain("Hello World", "Earth")
//    ShouldNotContain(["Hello", "World"], "Earth")
//    ShouldNotContain({"Hello": "World"}, "Earth")
//
func ShouldNotContain(s, contains interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotContains(currentCtx(), s, contains, msgAndArgs...)
}

// ShouldHaveSubset asserts that the specified list(array, slice...) contains all
// elements given in the specified subset(array, slice...).
//
//    ShouldHaveSubset([1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]")
//
func ShouldHaveSubset(list, subset interface{}, msgAndArgs ...interface{}) (ok bool) {
	return assert.Subset(currentCtx(), list, subset, msgAndArgs...)
}

// ShouldNotHaveSubset asserts that the specified list(array, slice...) contains not all
// elements given in the specified subset(array, slice...).
//
//    ShouldNotHaveSubset([1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]")
//
func ShouldNotHaveSubset(list, subset interface{}, msgAndArgs ...interface{}) (ok bool) {
	return assert.NotSubset(currentCtx(), list, subset, msgAndArgs...)
}

// ShouldMatchElements asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
//    ShouldMatchElements([1, 3, 2, 3], [1, 3, 3, 2])
//
func ShouldMatchElements(listA, listB interface{}, msgAndArgs ...interface{}) bool {
	return assert.ElementsMatch(currentCtx(), listA, listB, msgAndArgs...)
}

// Eval uses a EvalFunc to assert a complex condition.
func Eval(f assert.EvalFunc, msgAndArgs ...interface{}) bool {
	return assert.Condition(currentCtx(), f, msgAndArgs...)
}

// ShouldPanic asserts that the code inside the specified PanicTestFunc panics.
//
//   ShouldPanic(func(){ GoCrazy() })
//
func ShouldPanic(f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	return assert.Panics(currentCtx(), f, msgAndArgs...)
}

// ShouldPanicWithValue asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//   ShouldPanicWithValue("crazy error", func(){ GoCrazy() })
//
func ShouldPanicWithValue(expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	return assert.PanicsWithValue(currentCtx(), expected, f, msgAndArgs...)
}

// ShouldNotPanic asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   ShouldNotPanic(func(){ RemainCalm() })
//
func ShouldNotPanic(f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	return assert.NotPanics(currentCtx(), f, msgAndArgs...)
}

// WithinDuration asserts that the two times are within duration delta of each other.
//
//     WithinDuration(time.Now(), time.Now(), 10*time.Second)
//
func WithinDuration(expected, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) bool {
	return assert.WithinDuration(currentCtx(), expected, actual, delta, msgAndArgs...)
}

// InDelta asserts that the two numerals are within delta of each other.
//
//    InDelta(math.Pi, (22 / 7.0), 0.01)
//
func InDelta(expected, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	return assert.InDelta(currentCtx(), expected, actual, delta, msgAndArgs...)
}

// ShouldNotError asserts that a function returned no error (i.e. `nil`).
//
//    actualObj, err := SomeFunction()
//    if ShouldNotError(err) {
//        ShouldEqual(expectedObj, actualObj)
//    }
//
func ShouldNotError(err error, msgAndArgs ...interface{}) bool {
	return assert.NoError(currentCtx(), err, msgAndArgs...)
}

// ShouldError asserts that a function returned an error (i.e. not `nil`).
//
//    actualObj, err := SomeFunction()
//    if ShouldError(err) {
//        ShouldEqual(expectedError, err)
//    }
//
func ShouldError(err error, msgAndArgs ...interface{}) bool {
	return assert.Error(currentCtx(), err, msgAndArgs...)
}

// ShouldEqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//    actualObj, err := SomeFunction()
//    ShouldEqualError(err,  expectedErrorString)
//
func ShouldEqualError(err error, errString string, msgAndArgs ...interface{}) bool {
	return assert.EqualError(currentCtx(), err, errString, msgAndArgs...)
}

// ShouldMatchRegexp asserts that a specified regexp matches a string.
//
//    ShouldMatchRegexp(regexp.MustCompile("start"), "it's starting")
//    ShouldMatchRegexp("start...$", "it's not starting")
//
func ShouldMatchRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	return assert.Regexp(currentCtx(), rx, str, msgAndArgs...)
}

// ShouldNotMatchRegexp asserts that a specified regexp does not match a string.
//
//    ShouldNotMatchRegexp(regexp.MustCompile("starts"), "it's starting")
//    ShouldNotMatchRegexp("^start", "it's not starting")
//
func ShouldNotMatchRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotRegexp(currentCtx(), rx, str, msgAndArgs...)
}

// ShouldZero asserts that i is the zero value for its type.
func ShouldZero(i interface{}, msgAndArgs ...interface{}) bool {
	return assert.Zero(currentCtx(), i, msgAndArgs...)
}

// NotZero asserts that i is not the zero value for its type.
func ShouldNotZero(i interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotZero(currentCtx(), i, msgAndArgs...)
}

// JSONEq asserts that two JSON strings are equivalent.
//
//  JsonShouldEq(`{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
//
func JsonShouldEq(expected string, actual string, msgAndArgs ...interface{}) bool {
	return assert.JSONEq(currentCtx(), expected, actual, msgAndArgs...)
}
