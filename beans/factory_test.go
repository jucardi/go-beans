package beans_test // Using a different package name so this test represents in reality how the dependency injection would work.

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/jucardi/go-beans/beans"
	"github.com/stretchr/testify/assert"
)

// The following are implementations of the IService defined in the "example_test.go"
// in that file the best practices to use the bean factory are explained.

// TestServiceImpl is an implementation of IService that can be used as a bean.
type TestServiceImpl struct {
}

func (t *TestServiceImpl) SomeMethod() bool {
	return true
}

func (t *TestServiceImpl) GetName() string {
	return "bean1"
}

// TestServiceImpl2 is an implementation of IService that can be used as a bean.
type TestServiceImpl2 struct {
}

func (t *TestServiceImpl2) SomeMethod() bool {
	return true
}

func (t *TestServiceImpl2) GetName() string {
	return "bean2"
}

// TestServiceImpl3 is an implementation of IService that can be used as a bean.
type TestServiceImpl3 struct {
	name string
}

func (t *TestServiceImpl3) SomeMethod() bool {
	return true
}

func (t *TestServiceImpl3) GetName() string {
	return t.name
}

// Other test interface
type IOther interface {
	Name() string
}

type OtherImpl1 struct {
	name string
}

func (o *OtherImpl1) Name() string {
	return o.name
}

// Not used interface

type INotUsed interface{}

// endregion

// region Test Factory methods impl

var ComponentType *IService

func init() {
	beans.SetLogger(&beans.ConsoleLogger{})
	beans.Register(ComponentType, "default", &TestServiceImpl{})
	beans.SetPrimary(ComponentType, "default")
}

// Wraps `bean.Resolve` but automatically casts from `interface {}` to `IService`
func Resolve(name string) IService {
	return beans.Resolve(ComponentType, name).(IService)
}

// Wraps `bean.Primary` but automatically casts from `interface {}` to `IService`
func Primary() IService {
	return beans.Primary(ComponentType).(IService)
}

// endregion

func TestGet(t *testing.T) {
	assert.True(t, Resolve("default").SomeMethod())
	assert.Equal(t, "bean1", Resolve("default").GetName())
}

func TestGetPrimary(t *testing.T) {
	assert.True(t, Primary().SomeMethod())
}

func TestRegister(t *testing.T) {
	beans.Register(ComponentType, "some-name", &TestServiceImpl2{})
	assert.Equal(t, "bean1", Primary().GetName())
	assert.Equal(t, "bean2", Resolve("some-name").GetName())
}

func TestSetPrimary(t *testing.T) {
	beans.Register(ComponentType, "some-name", &TestServiceImpl2{})
	beans.SetPrimary(ComponentType, "some-name")
	assert.Equal(t, "bean2", Primary().GetName())
}

func TestSetPrimary_NotFound(t *testing.T) {
	name := "something"
	err := beans.SetPrimary(ComponentType, name)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("dependency %s not registered, unable to set as primary", name), err.Error())
}

func TestRegister_Exists(t *testing.T) {
	name := "some-name"
	beans.Register(ComponentType, name, &TestServiceImpl2{})
	err := beans.Register(ComponentType, name, &TestServiceImpl2{})
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("a dependency with name %s is already registered", name), err.Error())
}

func TestRegister_InvalidName(t *testing.T) {
	err := beans.Register(ComponentType, "", &TestServiceImpl{})
	assert.NotNil(t, err)
	assert.Equal(t, "the name cannot be empty", err.Error())
}

func TestGet_NotFound(t *testing.T) {
	s := beans.Resolve(ComponentType, "something")
	assert.Nil(t, s)
}

func TestRegister_Invalid(t *testing.T) {
	err := beans.Register(ComponentType, "some-name", "Something")
	assert.NotNil(t, err)
	assert.Equal(t, "the component type 'string' does not implement the provided type 'IService'", err.Error())
}

func TestResolve(t *testing.T) {
	something := beans.Resolve(ComponentType, "default").(IService)
	assert.True(t, something.SomeMethod())
	assert.Equal(t, "bean1", something.GetName())
}

func TestRegisterFunc(t *testing.T) {
	i := 0
	beans.RegisterFunc(ComponentType, "test-maker", func() interface{} {
		i++
		return &TestServiceImpl3{
			name: "TEST_" + strconv.Itoa(i),
		}
	})

	comp1 := beans.Resolve(ComponentType, "test-maker").(IService)
	comp2 := beans.Resolve(ComponentType, "test-maker").(IService)

	assert.Equal(t, "TEST_1", comp1.GetName())
	assert.Equal(t, "TEST_2", comp2.GetName())
}

func TestSetPrimary_TypeNotFound(t *testing.T) {
	err := beans.SetPrimary((*string)(nil), "some-random-name")
	assert.NotNil(t, err)
	assert.Equal(t, "no dependencies found for type string, unable to resolve", err.Error())
}

func TestResolve_TypeNotFound(t *testing.T) {
	val := beans.Resolve((*string)(nil), "some-random-name")
	assert.Nil(t, val)
}

func TestImplicitPrimary(t *testing.T) {
	beans.Register((*IOther)(nil), "something", &OtherImpl1{name: "primary"})
	val := beans.Primary((*IOther)(nil)).(IOther)
	assert.NotNil(t, val)
	assert.Equal(t, "primary", val.Name())
}

func TestPrimary_NotFound(t *testing.T) {
	beans.Register((*IOther)(nil), "name1", &OtherImpl1{name: "name1"})
	beans.Register((*IOther)(nil), "name2", &OtherImpl1{name: "name2"})
	val := beans.Primary((*IOther)(nil))
	assert.Nil(t, val)
}

func TestSetAllowOverrides(t *testing.T) {
	err := beans.Register((*IOther)(nil), "name", &OtherImpl1{name: "name1"})
	assert.Nil(t, err)
	err = beans.Register((*IOther)(nil), "name", &OtherImpl1{name: "name2"})
	assert.NotNil(t, err)
	comp1 := beans.Resolve((*IOther)(nil), "name").(IOther)
	assert.Equal(t, "name1", comp1.Name())
	beans.SetAllowOverrides(true)
	err = beans.Register((*IOther)(nil), "name", &OtherImpl1{name: "name2"})
	assert.Nil(t, err)
	comp2 := beans.Resolve((*IOther)(nil), "name").(IOther)
	assert.Equal(t, "name2", comp2.Name())
}

func TestExists(t *testing.T) {
	name := "something"
	beans.Register((*IOther)(nil), name, &OtherImpl1{name: "primary"})
	contains := beans.Exists((*IOther)(nil), name)
	assert.True(t, contains)
	assert.False(t, beans.Exists((*IService)(nil), name))
	assert.False(t, beans.Exists((*INotUsed)(nil), name))
}

func TestGetPrimaryName(t *testing.T) {
	name := "something"
	beans.Register((*IOther)(nil), name, &OtherImpl1{name: "primary"})
	beans.SetPrimary((*IOther)(nil), name)
	assert.Equal(t, name, beans.GetPrimaryName((*IOther)(nil)))
	assert.Equal(t, "", beans.GetPrimaryName((*INotUsed)(nil)))
}
