package assert

import (
	"encoding/json"
	"fmt"
	"github.com/jucardi/go-testx/testutilx"
	"github.com/jucardi/go-terminal-colors"
	"math"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	colorErrMsg     = []fmtc.Color{fmtc.White, fmtc.Italic}
	colorErrStack   = []fmtc.Color{fmtc.Cyan}
	colorFailureMsg = []fmtc.Color{fmtc.White, fmtc.Bold}
	colorErrLabels  = []fmtc.Color{fmtc.Bold, fmtc.Red}
	colorSubLabels  = []fmtc.Color{fmtc.Bold}
	colorDiff       = []fmtc.Color{fmtc.Yellow}

	onFailureCallback func()
)

// OnFailure registers a callback that will be executed if a test fails
func OnFailure(handler func()) {
	onFailureCallback = handler
}

// FailNow fails test
func FailNow(t TestingT, failureMessage string, msgAndArgs ...interface{}) bool {
	helper(t)
	Fail(t, failureMessage, msgAndArgs...)

	// We cannot extend TestingT with FailNow() and
	// maintain backwards compatibility, so we fallback
	// to panicking when FailNow is not available in
	// TestingT.
	// See issue #263

	if t, ok := t.(IFailNow); ok {
		t.FailNow()
	} else {
		panic("test failed and t is missing `FailNow()`")
	}
	return false
}

// Fail reports a failure through
func Fail(t TestingT, failureMessage string, msgAndArgs ...interface{}) bool {
	helper(t)

	var content []labeledContent
	message := messageFromMsgAndArgs(msgAndArgs...)

	if len(message) > 0 {
		content = append(content, labeledContent{
			fmtc.New().Print("Message:", colorErrLabels...).String(),
			fmtc.New().Print(message, colorErrMsg...).String()},
		)
	}

	content = append(content, labeledContent{
		fmtc.New().Print("At:", colorErrLabels...).String(),
		fmtc.New().PrintLn(strings.Join(CallerInfo(), "\n    "), colorErrStack...).PrintLn("").String()},
	)

	content = append(content, labeledContent{
		fmtc.New().Print("Error:", colorErrLabels...).String(),
		fmtc.New().Print(strings.Replace(failureMessage, "\n", "\n    ", -1), colorFailureMsg...).String()},
	)

	if n, ok := t.(IAssertLogger); ok {
		n.FailMsgf("\n\n%s", ""+labeledOutput(content...))
	} else {
		t.Errorf("\n\n%s", ""+labeledOutput(content...))
	}

	// Add test name if the Go version supports it
	if n, ok := t.(interface {
		Fail()
	}); ok {
		n.Fail()
	}

	if onFailureCallback != nil {
		onFailureCallback()
	}

	return false
}

// Implements asserts that an object is implemented by the specified interface.
//
//    assert.Implements(t, (*MyInterface)(nil), new(MyObject))
func Implements(t TestingT, interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	helper(t)
	interfaceType := reflect.TypeOf(interfaceObject).Elem()

	if object == nil {
		return Fail(t, fmt.Sprintf("Cannot check if nil implements %v", interfaceType), msgAndArgs...)
	}
	if !reflect.TypeOf(object).Implements(interfaceType) {
		return Fail(t, fmt.Sprintf("%T must implement %v", object, interfaceType), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// IsType asserts that the specified objects are of the same type.
func IsType(t TestingT, expectedType interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	helper(t)

	if !testutilx.ObjectsAreEqual(reflect.TypeOf(object), reflect.TypeOf(expectedType)) {
		return Fail(t, fmt.Sprintf("Object expected to be of type %v, but was %v", reflect.TypeOf(expectedType), reflect.TypeOf(object)), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// Equal asserts that two objects are equal.
//
//    assert.Equal(t, 123, 123)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func Equal(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	helper(t)

	if err := validateEqualArgs(expected, actual); err != nil {
		return Fail(t, fmtc.New().Print(
			fmt.Sprintf("Invalid operation: %#v == %#v (%s)", expected, actual, err),
			colorFailureMsg...).String(), msgAndArgs...)
	}

	if !testutilx.ObjectsAreEqual(expected, actual) {
		diff := diff(expected, actual)
		if diff != "" {
			diff = diff + "\n"
		}
		expected, actual = formatUnequalValues(expected, actual)
		return Fail(t,
			fmtc.New().
				PrintLn("Not equal", colorFailureMsg...).
				PrintLn("").
				Print("expected: ", colorSubLabels...).PrintLn(expected).
				Print("actual  : ", colorSubLabels...).PrintLn(actual).
				PrintLn("").
				Print(diff).
				String(),
			msgAndArgs...)
	}

	addAssert(t)
	return true
}

// EqualValues asserts that two objects are equal or convertable to the same types
// and equal.
//
//    assert.EqualValues(t, uint32(123), int32(123))
func EqualValues(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	helper(t)

	if !testutilx.ObjectsAreEqualValues(expected, actual) {
		diff := diff(expected, actual)
		expected, actual = formatUnequalValues(expected, actual)
		return Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s%s", expected, actual, diff), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// Exactly asserts that two objects are equal in value and type.
//
//    assert.Exactly(t, int32(123), int64(123))
func Exactly(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	helper(t)

	aType := reflect.TypeOf(expected)
	bType := reflect.TypeOf(actual)

	if aType != bType {
		return Fail(t, fmt.Sprintf("Types expected to match exactly\n\t%v != %v", aType, bType), msgAndArgs...)
	}

	addAssert(t)
	return Equal(t, expected, actual, msgAndArgs...)
}

// NotNil asserts that the specified object is not nil.
//
//    assert.NotNil(t, err)
func NotNil(t TestingT, object interface{}, msgAndArgs ...interface{}) bool {
	helper(t)
	if !isNil(object) {
		addAssert(t)
		return true
	}
	return Fail(t, "Expected value not to be nil.", msgAndArgs...)
}

// Nil asserts that the specified object is nil.
//
//    assert.Nil(t, err)
func Nil(t TestingT, object interface{}, msgAndArgs ...interface{}) bool {
	helper(t)
	if isNil(object) {
		addAssert(t)
		return true
	}
	return Fail(t, fmt.Sprintf("Expected nil, but got: %#v", object), msgAndArgs...)
}

// Empty asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  assert.Empty(t, obj)
func Empty(t TestingT, object interface{}, msgAndArgs ...interface{}) bool {
	helper(t)

	pass := isEmpty(object)
	if !pass {
		return Fail(t, fmt.Sprintf("Should be empty, but was %v", object), msgAndArgs...)
	}

	addAssert(t)
	return pass
}

// NotEmpty asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if assert.NotEmpty(t, obj) {
//    assert.Equal(t, "two", obj[1])
//  }
func NotEmpty(t TestingT, object interface{}, msgAndArgs ...interface{}) bool {
	helper(t)

	pass := !isEmpty(object)
	if !pass {
		return Fail(t, fmt.Sprintf("Should NOT be empty, but was %v", object), msgAndArgs...)
	}

	addAssert(t)
	return pass
}

// Len asserts that the specified object has specific length.
// Len also fails if the object has a type that len() not accept.
//
//    assert.Len(t, mySlice, 3)
func Len(t TestingT, object interface{}, length int, msgAndArgs ...interface{}) bool {
	helper(t)
	ok, l := getLen(object)
	if !ok {
		return Fail(t,
			fmtc.New().
				PrintLn("'Len()' is not available for object", colorFailureMsg...).
				PrintLn(object, colorDiff...).
				String(),
			msgAndArgs...)
	}

	if l != length {
		return Fail(t,
			fmtc.New().
				PrintLn("Length mismatch for", colorFailureMsg...).
				PrintLn("").
				Print("    > ").PrintLn(object, colorDiff...).
				PrintLn("").
				Print("expected: ", colorSubLabels...).PrintLn(length).
				Print("actual  : ", colorSubLabels...).PrintLn(l).
				PrintLn("").
				String(),
			msgAndArgs...)
	}
	addAssert(t)
	return true
}

// True asserts that the specified value is true.
//
//    assert.True(t, myBool)
func True(t TestingT, value bool, msgAndArgs ...interface{}) bool {
	helper(t)
	if h, ok := t.(interface {
		Helper()
	}); ok {
		h.Helper()
	}

	if !value {
		return Fail(t, fmtc.New().PrintLn("Should be true", colorFailureMsg...).String(), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// False asserts that the specified value is false.
//
//    assert.False(t, myBool)
func False(t TestingT, value bool, msgAndArgs ...interface{}) bool {
	helper(t)

	if value {
		return Fail(t, fmtc.New().PrintLn("Should be false", colorFailureMsg...).String(), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// NotEqual asserts that the specified values are NOT equal.
//
//    assert.NotEqual(t, obj1, obj2)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func NotEqual(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	helper(t)
	if err := validateEqualArgs(expected, actual); err != nil {
		return Fail(t, fmt.Sprintf("Invalid operation: %#v != %#v (%s)",
			expected, actual, err), msgAndArgs...)
	}

	if testutilx.ObjectsAreEqual(expected, actual) {
		return Fail(t,
			fmtc.New().
				PrintLn("Should not be equal", colorFailureMsg...).
				PrintLn("").
				Print("actual  : ", colorSubLabels...).PrintLn(actual).
				PrintLn("").
				String(),
			msgAndArgs...)
	}

	addAssert(t)
	return true
}

// Contains asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    assert.Contains(t, "Hello World", "World")
//    assert.Contains(t, ["Hello", "World"], "World")
//    assert.Contains(t, {"Hello": "World"}, "Hello")
func Contains(t TestingT, s, contains interface{}, msgAndArgs ...interface{}) bool {
	helper(t)

	ok, found := includeElement(s, contains)
	if !ok {
		return Fail(t, fmt.Sprintf("\"%s\" could not be applied builtin len()", s), msgAndArgs...)
	}
	if !found {
		return Fail(t, fmt.Sprintf("\"%s\" does not contain \"%s\"", s, contains), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    assert.NotContains(t, "Hello World", "Earth")
//    assert.NotContains(t, ["Hello", "World"], "Earth")
//    assert.NotContains(t, {"Hello": "World"}, "Earth")
func NotContains(t TestingT, s, contains interface{}, msgAndArgs ...interface{}) bool {
	helper(t)

	ok, found := includeElement(s, contains)
	if !ok {
		return Fail(t, fmt.Sprintf("\"%s\" could not be applied builtin len()", s), msgAndArgs...)
	}
	if found {
		return Fail(t, fmt.Sprintf("\"%s\" should not contain \"%s\"", s, contains), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// Subset asserts that the specified list(array, slice...) contains all
// elements given in the specified subset(array, slice...).
//
//    assert.Subset(t, [1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]")
func Subset(t TestingT, list, subset interface{}, msgAndArgs ...interface{}) (ok bool) {
	helper(t)
	if subset == nil {
		return true // we consider nil to be equal to the nil set
	}

	subsetValue := reflect.ValueOf(subset)
	defer func() {
		if e := recover(); e != nil {
			ok = false
		}
	}()

	listKind := reflect.TypeOf(list).Kind()
	subsetKind := reflect.TypeOf(subset).Kind()

	if listKind != reflect.Array && listKind != reflect.Slice {
		return Fail(t, fmt.Sprintf("%q has an unsupported type %s", list, listKind), msgAndArgs...)
	}

	if subsetKind != reflect.Array && subsetKind != reflect.Slice {
		return Fail(t, fmt.Sprintf("%q has an unsupported type %s", subset, subsetKind), msgAndArgs...)
	}

	for i := 0; i < subsetValue.Len(); i++ {
		element := subsetValue.Index(i).Interface()
		ok, found := includeElement(list, element)
		if !ok {
			return Fail(t, fmt.Sprintf("\"%s\" could not be applied builtin len()", list), msgAndArgs...)
		}
		if !found {
			return Fail(t, fmt.Sprintf("\"%s\" does not contain \"%s\"", list, element), msgAndArgs...)
		}
	}

	addAssert(t)
	return true
}

// NotSubset asserts that the specified list(array, slice...) contains not all
// elements given in the specified subset(array, slice...).
//
//    assert.NotSubset(t, [1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]")
func NotSubset(t TestingT, list, subset interface{}, msgAndArgs ...interface{}) (ok bool) {
	helper(t)
	if subset == nil {
		return Fail(t, fmt.Sprintf("nil is the empty set which is a subset of every set"), msgAndArgs...)
	}

	subsetValue := reflect.ValueOf(subset)
	defer func() {
		if e := recover(); e != nil {
			ok = false
		}
	}()

	listKind := reflect.TypeOf(list).Kind()
	subsetKind := reflect.TypeOf(subset).Kind()

	if listKind != reflect.Array && listKind != reflect.Slice {
		return Fail(t, fmt.Sprintf("%q has an unsupported type %s", list, listKind), msgAndArgs...)
	}

	if subsetKind != reflect.Array && subsetKind != reflect.Slice {
		return Fail(t, fmt.Sprintf("%q has an unsupported type %s", subset, subsetKind), msgAndArgs...)
	}

	for i := 0; i < subsetValue.Len(); i++ {
		element := subsetValue.Index(i).Interface()
		ok, found := includeElement(list, element)
		if !ok {
			return Fail(t, fmt.Sprintf("\"%s\" could not be applied builtin len()", list), msgAndArgs...)
		}
		if !found {
			return true
		}
	}

	addAssert(t)
	return Fail(t, fmt.Sprintf("%q is a subset of %q", subset, list), msgAndArgs...)
}

// ElementsMatch asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// assert.ElementsMatch(t, [1, 3, 2, 3], [1, 3, 3, 2])
func ElementsMatch(t TestingT, listA, listB interface{}, msgAndArgs ...interface{}) bool {
	helper(t)
	if isEmpty(listA) && isEmpty(listB) {
		addAssert(t)
		return true
	}

	aKind := reflect.TypeOf(listA).Kind()
	bKind := reflect.TypeOf(listB).Kind()

	if aKind != reflect.Array && aKind != reflect.Slice {
		return Fail(t,
			fmtc.New().
				PrintLn(fmt.Sprintf("%q has an unsupported type %s", listA, aKind)).
				String(),
			msgAndArgs...)
	}

	if bKind != reflect.Array && bKind != reflect.Slice {
		return Fail(t,
			fmtc.New().
				PrintLn(fmt.Sprintf("%q has an unsupported type %s", listB, aKind)).
				String(),
			msgAndArgs...)
	}

	aValue := reflect.ValueOf(listA)
	bValue := reflect.ValueOf(listB)

	aLen := aValue.Len()
	bLen := bValue.Len()

	if aLen != bLen {
		return Fail(t,
			fmtc.New().
				PrintLn(fmt.Sprintf("lengths mismatch: %d != %d", aLen, bLen)).
				String(),
			msgAndArgs...)
	}

	// Mark indexes in bValue that we already used
	visited := make([]bool, bLen)
	for i := 0; i < aLen; i++ {
		element := aValue.Index(i).Interface()
		found := false
		for j := 0; j < bLen; j++ {
			if visited[j] {
				continue
			}
			if testutilx.ObjectsAreEqual(bValue.Index(j).Interface(), element) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			return Fail(t,
				fmtc.New().
					Print("Element                 ").PrintLn(element, colorDiff...).
					Print("appears more times in   ").PrintLn(aValue, colorDiff...).
					Print("than                    ").PrintLn(bValue, colorDiff...).
					String(),
				msgAndArgs...)
		}
	}

	addAssert(t)
	return true
}

// Condition uses a EvalFunc to assert a complex condition.
func Condition(t TestingT, comp EvalFunc, msgAndArgs ...interface{}) bool {
	helper(t)
	result := comp()
	if !result {
		Fail(t, "Condition failed!", msgAndArgs...)
	}
	addAssert(t)
	return result
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//   assert.Panics(t, func(){ GoCrazy() })
func Panics(t TestingT, f PanicTestFunc, msgAndArgs ...interface{}) bool {
	helper(t)

	if funcDidPanic, panicValue := didPanic(f); !funcDidPanic {
		return Fail(t, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", f, panicValue), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// PanicsWithValue asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//   assert.PanicsWithValue(t, "crazy error", func(){ GoCrazy() })
func PanicsWithValue(t TestingT, expected interface{}, f PanicTestFunc, msgAndArgs ...interface{}) bool {
	helper(t)

	funcDidPanic, panicValue := didPanic(f)
	if !funcDidPanic {
		return Fail(t, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", f, panicValue), msgAndArgs...)
	}
	if panicValue != expected {
		return Fail(t, fmt.Sprintf("func %#v should panic with value:\t%#v\n\tPanic value:\t%#v", f, expected, panicValue), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   assert.NotPanics(t, func(){ RemainCalm() })
func NotPanics(t TestingT, f PanicTestFunc, msgAndArgs ...interface{}) bool {
	helper(t)

	if funcDidPanic, panicValue := didPanic(f); funcDidPanic {
		return Fail(t, fmt.Sprintf("func %#v should not panic\n\tPanic value:\t%v", f, panicValue), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// WithinDuration asserts that the two times are within duration delta of each other.
//
//   assert.WithinDuration(t, time.Now(), time.Now(), 10*time.Second)
func WithinDuration(t TestingT, expected, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) bool {
	helper(t)

	dt := expected.Sub(actual)
	if dt < -delta || dt > delta {
		return Fail(t, fmt.Sprintf("Max difference between %v and %v allowed is %v, but difference was %v", expected, actual, delta, dt), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// InDelta asserts that the two numerals are within delta of each other.
//
// 	 assert.InDelta(t, math.Pi, (22 / 7.0), 0.01)
func InDelta(t TestingT, expected, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	helper(t)

	af, aok := toFloat(expected)
	bf, bok := toFloat(actual)

	if !aok || !bok {
		return Fail(t, fmt.Sprintf("Parameters must be numerical"), msgAndArgs...)
	}

	if math.IsNaN(af) {
		return Fail(t, fmt.Sprintf("Expected must not be NaN"), msgAndArgs...)
	}

	if math.IsNaN(bf) {
		return Fail(t, fmt.Sprintf("Expected %v with delta %v, but was NaN", expected, delta), msgAndArgs...)
	}

	dt := af - bf
	if dt < -delta || dt > delta {
		return Fail(t, fmt.Sprintf("Max difference between %v and %v allowed is %v, but difference was %v", expected, actual, delta, dt), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// InDeltaSlice is the same as InDelta, except it compares two slices.
func InDeltaSlice(t TestingT, expected, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	helper(t)
	if expected == nil || actual == nil ||
		reflect.TypeOf(actual).Kind() != reflect.Slice ||
		reflect.TypeOf(expected).Kind() != reflect.Slice {
		return Fail(t, fmt.Sprintf("Parameters must be slice"), msgAndArgs...)
	}

	actualSlice := reflect.ValueOf(actual)
	expectedSlice := reflect.ValueOf(expected)

	for i := 0; i < actualSlice.Len(); i++ {
		result := InDelta(t, actualSlice.Index(i).Interface(), expectedSlice.Index(i).Interface(), delta, msgAndArgs...)
		if !result {
			return result
		}
	}

	addAssert(t)
	return true
}

// InDeltaMapValues is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func InDeltaMapValues(t TestingT, expected, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	helper(t)
	if expected == nil || actual == nil ||
		reflect.TypeOf(actual).Kind() != reflect.Map ||
		reflect.TypeOf(expected).Kind() != reflect.Map {
		return Fail(t, "Arguments must be maps", msgAndArgs...)
	}

	expectedMap := reflect.ValueOf(expected)
	actualMap := reflect.ValueOf(actual)

	if expectedMap.Len() != actualMap.Len() {
		return Fail(t, "Arguments must have the same number of keys", msgAndArgs...)
	}

	for _, k := range expectedMap.MapKeys() {
		ev := expectedMap.MapIndex(k)
		av := actualMap.MapIndex(k)

		if !ev.IsValid() {
			return Fail(t, fmt.Sprintf("missing key %q in expected map", k), msgAndArgs...)
		}

		if !av.IsValid() {
			return Fail(t, fmt.Sprintf("missing key %q in actual map", k), msgAndArgs...)
		}

		if !InDelta(
			t,
			ev.Interface(),
			av.Interface(),
			delta,
			msgAndArgs...,
		) {
			return false
		}
	}

	addAssert(t)
	return true
}

// InEpsilon asserts that expected and actual have a relative error less than epsilon
func InEpsilon(t TestingT, expected, actual interface{}, epsilon float64, msgAndArgs ...interface{}) bool {
	helper(t)
	actualEpsilon, err := calcRelativeError(expected, actual)
	if err != nil {
		return Fail(t, err.Error(), msgAndArgs...)
	}
	if actualEpsilon > epsilon {
		return Fail(t, fmt.Sprintf("Relative error is too high: %#v (expected)\n"+
			"        < %#v (actual)", epsilon, actualEpsilon), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// InEpsilonSlice is the same as InEpsilon, except it compares each value from two slices.
func InEpsilonSlice(t TestingT, expected, actual interface{}, epsilon float64, msgAndArgs ...interface{}) bool {
	helper(t)
	if expected == nil || actual == nil ||
		reflect.TypeOf(actual).Kind() != reflect.Slice ||
		reflect.TypeOf(expected).Kind() != reflect.Slice {
		return Fail(t, fmt.Sprintf("Parameters must be slice"), msgAndArgs...)
	}

	actualSlice := reflect.ValueOf(actual)
	expectedSlice := reflect.ValueOf(expected)

	for i := 0; i < actualSlice.Len(); i++ {
		result := InEpsilon(t, actualSlice.Index(i).Interface(), expectedSlice.Index(i).Interface(), epsilon)
		if !result {
			return result
		}
	}

	addAssert(t)
	return true
}

/*
	Errors
*/

// NoError asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.NoError(t, err) {
//	   assert.Equal(t, expectedObj, actualObj)
//   }
func NoError(t TestingT, err error, msgAndArgs ...interface{}) bool {
	helper(t)

	if err == nil || !reflect.ValueOf(err).IsValid() || (reflect.ValueOf(err).Kind() == reflect.Ptr && reflect.ValueOf(err).IsNil()) {
		addAssert(t)
		return true
	}

	return Fail(t, fmt.Sprintf("Received unexpected error:\n%+v", err), msgAndArgs...)
}

// Error asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.Error(t, err) {
//	   assert.Equal(t, expectedError, err)
//   }
func Error(t TestingT, err error, msgAndArgs ...interface{}) bool {
	helper(t)

	if err == nil || !reflect.ValueOf(err).IsValid() || (reflect.ValueOf(err).Kind() == reflect.Ptr && reflect.ValueOf(err).IsNil()) {
		return Fail(t, "An error is expected but got nil or invalid error.", msgAndArgs...)
	}

	addAssert(t)
	return true
}

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   assert.EqualError(t, err,  expectedErrorString)
func EqualError(t TestingT, theError error, errString string, msgAndArgs ...interface{}) bool {
	helper(t)
	if !Error(t, theError, msgAndArgs...) {
		return false
	}
	expected := errString
	actual := theError.Error()
	// don't need to use deep equals here, we know they are both strings
	if expected != actual {
		return Fail(t, fmt.Sprintf("Error message not equal:\n"+
			"expected: %q\n"+
			"actual  : %q", expected, actual), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// Regexp asserts that a specified regexp matches a string.
//
//  assert.Regexp(t, regexp.MustCompile("start"), "it's starting")
//  assert.Regexp(t, "start...$", "it's not starting")
func Regexp(t TestingT, rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	helper(t)

	match := matchRegexp(rx, str)

	if !match {
		return Fail(t, fmt.Sprintf("Expect \"%v\" to match \"%v\"", str, rx), msgAndArgs...)
	}

	addAssert(t)
	return true
}

// NotRegexp asserts that a specified regexp does not match a string.
//
//  assert.NotRegexp(t, regexp.MustCompile("starts"), "it's starting")
//  assert.NotRegexp(t, "^start", "it's not starting")
func NotRegexp(t TestingT, rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	helper(t)
	match := matchRegexp(rx, str)

	if match {
		return Fail(t, fmt.Sprintf("Expect \"%v\" to NOT match \"%v\"", str, rx), msgAndArgs...)
	}

	addAssert(t)
	return true

}

// Zero asserts that i is the zero value for its type.
func Zero(t TestingT, i interface{}, msgAndArgs ...interface{}) bool {
	helper(t)
	if i != nil && !reflect.DeepEqual(i, reflect.Zero(reflect.TypeOf(i)).Interface()) {
		return Fail(t, fmt.Sprintf("Should be zero, but was %v", i), msgAndArgs...)
	}
	addAssert(t)
	return true
}

// NotZero asserts that i is not the zero value for its type.
func NotZero(t TestingT, i interface{}, msgAndArgs ...interface{}) bool {
	helper(t)
	if i == nil || reflect.DeepEqual(i, reflect.Zero(reflect.TypeOf(i)).Interface()) {
		return Fail(t, fmt.Sprintf("Should not be zero, but was %v", i), msgAndArgs...)
	}
	addAssert(t)
	return true
}

// FileExists checks whether a file exists in the given path. It also fails if the path points to a directory or there is an error when trying to check the file.
func FileExists(t TestingT, path string, msgAndArgs ...interface{}) bool {
	helper(t)
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Fail(t, fmt.Sprintf("unable to find file %q", path), msgAndArgs...)
		}
		return Fail(t, fmt.Sprintf("error when running os.Lstat(%q): %s", path, err), msgAndArgs...)
	}
	if info.IsDir() {
		return Fail(t, fmt.Sprintf("%q is a directory", path), msgAndArgs...)
	}
	addAssert(t)
	return true
}

// DirExists checks whether a directory exists in the given path. It also fails if the path is a file rather a directory or there is an error checking whether it exists.
func DirExists(t TestingT, path string, msgAndArgs ...interface{}) bool {
	helper(t)
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Fail(t, fmt.Sprintf("unable to find file %q", path), msgAndArgs...)
		}
		return Fail(t, fmt.Sprintf("error when running os.Lstat(%q): %s", path, err), msgAndArgs...)
	}
	if !info.IsDir() {
		return Fail(t, fmt.Sprintf("%q is a file", path), msgAndArgs...)
	}
	addAssert(t)
	return true
}

// JSONEq asserts that two JSON strings are equivalent.
//
//  assert.JSONEq(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
func JSONEq(t TestingT, expected string, actual string, msgAndArgs ...interface{}) bool {
	helper(t)
	var expectedJSONAsInterface, actualJSONAsInterface interface{}

	if err := json.Unmarshal([]byte(expected), &expectedJSONAsInterface); err != nil {
		return Fail(t, fmt.Sprintf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error()), msgAndArgs...)
	}

	if err := json.Unmarshal([]byte(actual), &actualJSONAsInterface); err != nil {
		return Fail(t, fmt.Sprintf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", actual, err.Error()), msgAndArgs...)
	}

	result := Equal(t, expectedJSONAsInterface, actualJSONAsInterface, msgAndArgs...)
	if result {
		addAssert(t)
	}
	return result
}
