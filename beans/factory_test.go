package beans_test // Using a different package name so this test represents in reality how the dependency injection would work.

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/jucardi/go-beans/beans"
	. "github.com/jucardi/go-testx/testx"
)

// The following are implementations of the IService defined in the "example_test.go"
// in that file the best practices to use the bean factory are explained.

// TestServiceImpl is an implementation of IService that can be used as a bean.
type TestServiceImpl struct {
	resolveCount          int
	firstTimeResolveCount int
}

func (t *TestServiceImpl) SomeMethod() bool {
	return true
}

func (t *TestServiceImpl) GetName() string {
	return "bean1"
}

func (t *TestServiceImpl) OnResolve() {
	t.resolveCount++
}

func (t *TestServiceImpl) OnFirstTimeResolve() {
	t.firstTimeResolveCount++
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

// region Test Factory methods impl

var ComponentType *IService


func before() *TestServiceImpl {
	ret := &TestServiceImpl{}
	beans.SetAllowOverrides(true)
	ShouldNotError(beans.Clear())
	beans.SetAllowOverrides(false)
	ShouldNotError(beans.Register(ComponentType, "default", ret))
	ShouldNotError(beans.SetPrimary(ComponentType, "default"))
	return ret
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
	Convey("Testing Get", t, func() {
		Convey("Sucessful Get", t, func() {
			before()
			ShouldBeTrue(Resolve("default").SomeMethod())
			ShouldEqual("bean1", Resolve("default").GetName())
		})
		Convey("Bean not found", t, func() {
			before()
			ShouldBeNil(beans.Resolve(ComponentType, "something"))
		})
	})
}

func TestGetPrimary(t *testing.T) {
	Convey("Testing GetPrimary", t, func() {
		Convey("Success explicit", t, func() {
			before()
			ShouldBeTrue(Primary().SomeMethod())
		})
		Convey("Success implicit", t, func() {
			before()
			ShouldNotError(beans.Register((*IOther)(nil), "something", &OtherImpl1{name: "primary"}))
			val := beans.Primary((*IOther)(nil)).(IOther)
			ShouldNotBeNil(val)
			ShouldEqual("primary", val.Name())
		})
		Convey("Failure, not found", t, func() {
			before()
			ShouldNotError(beans.Register((*IOther)(nil), "name1", &OtherImpl1{name: "name1"}))
			ShouldNotError(beans.Register((*IOther)(nil), "name2", &OtherImpl1{name: "name2"}))
			val := beans.Primary((*IOther)(nil))
			ShouldBeNil(val)
		})
	})
}

func TestOnResolveHandlers(t *testing.T) {
	Convey("Testing OnResolve handlers", t, func() {
		instance := before()
		Convey("Validating resolve handlers are invoked successfully the expected amount of times", t, func() {
			ShouldEqual(0, instance.resolveCount)
			ShouldEqual(0, instance.firstTimeResolveCount)

			Resolve("default")
			ShouldEqual(1, instance.resolveCount)
			ShouldEqual(1, instance.firstTimeResolveCount)

			Resolve("default")
			ShouldEqual(2, instance.resolveCount)
			ShouldEqual(1, instance.firstTimeResolveCount)

			Resolve("default")
			ShouldEqual(3, instance.resolveCount)
			ShouldEqual(1, instance.firstTimeResolveCount)
		})
	})
}

func TestRegister(t *testing.T) {
	Convey("Testing Register", t, func() {
		Convey("Success", t, func() {
			before()
			ShouldNotError(beans.Register(ComponentType, "some-name", &TestServiceImpl2{}))
			ShouldEqual("bean1", Primary().GetName())
			ShouldEqual("bean2", Resolve("some-name").GetName())
		})
		Convey("Failure, already exists", t, func() {
			before()
			name := "some-name"
			ShouldNotError(beans.Register(ComponentType, name, &TestServiceImpl2{}))
			err := beans.Register(ComponentType, name, &TestServiceImpl2{})
			ShouldError(err)
			ShouldEqual(fmt.Sprintf("a dependency with name %s is already registered", name), err.Error())
		})
		Convey("Failure, invalid name", t, func() {
			before()
			err := beans.Register(ComponentType, "", &TestServiceImpl{})
			ShouldError(err)
			ShouldEqual("the name cannot be empty", err.Error())
		})
		Convey("Failure, type mismatch", t, func() {
			before()
			err := beans.Register(ComponentType, "some-name", "Something")
			ShouldError(err)
			ShouldEqual("the component type 'string' does not implement the provided type 'IService'", err.Error())
		})
	})
}

func TestSetPrimary(t *testing.T) {
	Convey("Testing SetPrimary", t, func() {
		Convey("Success", t, func() {
			before()
			ShouldNotError(beans.Register(ComponentType, "some-name", &TestServiceImpl2{}))
			ShouldNotError(beans.SetPrimary(ComponentType, "some-name"))
			ShouldNotBeEqual("bean2", Primary().GetName())
			ShouldNotError(beans.SetPrimary(ComponentType, "some-name", true))
			ShouldEqual("bean2", Primary().GetName())
		})
		Convey("Failure, bean not found", t, func() {
			before()
			name := "something"
			err := beans.SetPrimary(ComponentType, name)
			ShouldError(err)
			ShouldEqual(fmt.Sprintf("dependency %s not registered, unable to set as primary", name), err.Error())
		})
		Convey("Failure, type not found", t, func() {
			before()
			err := beans.SetPrimary((*string)(nil), "some-random-name")
			ShouldError(err)
			ShouldEqual("no dependencies found for type string, unable to resolve", err.Error())
		})
	})
}

func TestResolve(t *testing.T) {
	Convey("Testing Resolve", t, func() {
		Convey("Successful", t, func() {
			before()
			something := beans.Resolve(ComponentType, "default").(IService)
			ShouldBeTrue(something.SomeMethod())
			ShouldEqual("bean1", something.GetName())
		})
		Convey("Failure, type not found", t, func() {
			before()
			val := beans.Resolve((*string)(nil), "some-random-name")
			ShouldBeNil(val)
		})
	})
}

func TestRegisterFunc(t *testing.T) {
	Convey("Testing RegisterFunc", t, func() {
		Convey("Success", t, func() {
			before()
			i := 0
			ShouldNotError(
				beans.RegisterFunc(ComponentType, "test-maker", func() interface{} {
					i++
					return &TestServiceImpl3{
						name: "TEST_" + strconv.Itoa(i),
					}
				}),
			)

			comp1 := beans.Resolve(ComponentType, "test-maker").(IService)
			comp2 := beans.Resolve(ComponentType, "test-maker").(IService)

			ShouldEqual("TEST_1", comp1.GetName())
			ShouldEqual("TEST_2", comp2.GetName())
		})
	})
}

func TestSetAllowOverrides(t *testing.T) {
	Convey("Testing SetAllowOverrides", t, func() {
		before()
		Convey("While allow override is false", t, func() {
			ShouldNotError(beans.Register((*IOther)(nil), "name", &OtherImpl1{name: "name1"}))
			ShouldError(beans.Register((*IOther)(nil), "name", &OtherImpl1{name: "name2"}))
			comp1 := beans.Resolve((*IOther)(nil), "name").(IOther)
			ShouldEqual("name1", comp1.Name())
		})
		Convey("While allow override is true", t, func() {
			beans.SetAllowOverrides(true)
			ShouldNotError(beans.Register((*IOther)(nil), "name", &OtherImpl1{name: "name2"}))
			comp2 := beans.Resolve((*IOther)(nil), "name").(IOther)
			ShouldEqual("name2", comp2.Name())
		})
	})
}

func TestExists(t *testing.T) {
	Convey("Testing Exists", t, func() {
		before()
		name := "something"
		ShouldNotError(beans.Register((*IOther)(nil), name, &OtherImpl1{name: "primary"}))
		ShouldBeTrue(beans.Exists((*IOther)(nil), name))
		ShouldBeFalse(beans.Exists((*IService)(nil), name))
		ShouldBeFalse(beans.Exists((*INotUsed)(nil), name))
	})
}

func TestGetPrimaryName(t *testing.T) {
	Convey("Testing GetPrimaryName", t, func() {
		before()
		name := "something"
		ShouldNotError(beans.Register((*IOther)(nil), name, &OtherImpl1{name: "primary"}))
		ShouldNotError(beans.SetPrimary((*IOther)(nil), name))
		ShouldEqual(name, beans.GetPrimaryName((*IOther)(nil)))
		ShouldEqual("", beans.GetPrimaryName((*INotUsed)(nil)))
	})
}
