package beans

import (
	"errors"
	"fmt"
	"reflect"
)

type dependencyCollection struct {
	primary   string
	instances map[string]func() interface{}
}

var (
	allowOverrides = false
	dependencies   map[reflect.Type]*dependencyCollection
	log            IBeanLogger
)

func init() {
	dependencies = make(map[reflect.Type]*dependencyCollection)
	log = &emptyLogger{}
}

// SetLogger sets an IBeanLogger implementation to use as a logger for the bean factory.
func SetLogger(logger IBeanLogger) {
	log = logger
}

// SetAllowOverrides is normally used when testing. It allows a registered bean to be overwritten by another implementation, like a Mock
func SetAllowOverrides(allow bool) {
	allowOverrides = allow
}

// Resolve resolves the bean type by the given name.
//
// The 'interfaceRef' is a reference pointer to the bean interface so the bean factory knows what is the bean type, and
// it can properly register it. It can be a nil pointer to the interface. There are 2 ways of passing this a nil pointer.
// Let us assume we have a bean interface called IService. To obtain a nil pointer to the interface:
//
// Option 1:   (*IService)(nil)
//
//   Eg.   bean.Resolve((*IService)(nil), beanName)
//
// Option 2:   var reference *IService
//
//   Eg.   bean.Resolve(reference, beanName)
//
func Resolve(ref interface{}, name string) interface{} {
	return Get(getType(ref), name)
}

// Primary resolves the bean type using the registered bean set as primary
//
// The 'interfaceRef' is a reference pointer to the bean interface so the bean factory knows what is the bean type, and
// it can properly register it. It can be a nil pointer to the interface. There are 2 ways of passing this a nil pointer.
// Let us assume we have a bean interface called IService. To obtain a nil pointer to the interface:
//
// Option 1:   (*IService)(nil)
//
//   Eg.   bean.Primary((*IService)(nil))
//
// Option 2:   var reference *IService
//
//   Eg.   bean.Primary(reference)
//
func Primary(ref interface{}) interface{} {
	return GetPrimary(getType(ref))
}

// Get gets the the instance by the specified name.
func Get(t reflect.Type, name string) interface{} {
	if !containsType(dependencies, t) {
		log.Errorf("no dependencies found for type %s, unable to resolve.", t.Name())
		return nil
	}

	if name == "" {
		if dependencies[t].primary != "" {
			return Get(t, dependencies[t].primary)
		}

		if len(dependencies[t].instances) == 1 {
			for _, fn := range dependencies[t].instances {
				return fn()
			}
		}

		log.Error(fmt.Sprintf("No primary dependency found for type '%s'", t.Name()))
		return nil
	}

	if fn, ok := dependencies[t].instances[name]; ok {
		return fn()
	}

	log.Errorf("dependency %s not registered, unable to resolve.", name)
	return nil
}

// GetPrimary gets the primary dependency registered in this factory instance, same as primary but with a reflect.Type
func GetPrimary(t reflect.Type) interface{} {
	return Get(t, "")
}

// RegisterFuncByType registers a bean function retriever into the factory,
// the function could return a singleton instance or could also be used for a constructor.
func RegisterFuncByType(t reflect.Type, name string, fn func() interface{}) error {
	if name == "" {
		return errors.New("the name cannot be empty")
	}

	if !containsType(dependencies, t) {
		dependencies[t] = &dependencyCollection{
			instances: make(map[string]func() interface{}),
		}
	}

	if _, ok := dependencies[t].instances[name]; ok && !allowOverrides {
		return fmt.Errorf("a dependency with name %s is already registered", name)
	}

	dependencies[t].instances[name] = fn
	return nil
}

// RegisterFunc registers a bean function retriever into the factory.
// the function could return a singleton instance or could also be used for a constructor.
//
// The 'interfaceRef' is a reference pointer to the bean interface so the bean factory knows what is the bean type, and
// it can properly register it. It can be a nil pointer to the interface. There are 2 ways of passing this a nil pointer.
// Let us assume we have a bean interface called IService. To obtain a nil pointer to the interface:
//
// Option 1:   (*IService)(nil)
//
//   Eg.   bean.RegisterFunc((*IService)(nil), beanName, func() interface{}) { ...
//
// Option 2:   var reference *IService
//
//   Eg.   bean.RegisterFunc(reference, beanName, func() interface{}) { ...
//
func RegisterFunc(interfaceRef interface{}, name string, fn func() interface{}) error {
	return RegisterFuncByType(getType(interfaceRef), name, fn)
}

// RegisterByType registers a bean singleton instance into the factory.
func RegisterByType(t reflect.Type, name string, component interface{}) error {
	if ct := reflect.TypeOf(component); !ct.Implements(t) {
		return fmt.Errorf("the component type '%s' does not implement the provided type '%s'", ct.Name(), t.Name())
	}

	return RegisterFuncByType(t, name, func() interface{} { return component })
}

// Register registers a bean singleton instance into the factory.
//
// The 'interfaceRef' is a reference pointer to the bean interface so the bean factory knows what is the bean type, and
// it can properly register it. It can be a nil pointer to the interface. There are 2 ways of passing this a nil pointer.
// Let us assume we have a bean interface called IService. To obtain a nil pointer to the interface:
//
// Option 1:   (*IService)(nil)
//
//   Eg.   bean.Register((*IService)(nil), beanName, instance)
//
// Option 2:   var reference *IService
//
//   Eg.   bean.Register(reference, beanName, instance)
//
func Register(interfaceRef interface{}, name string, component interface{}) error {
	return RegisterByType(getType(interfaceRef), name, component)
}

// SetPrimaryByType sets the primary bean name to be used.
func SetPrimaryByType(t reflect.Type, name string) error {
	if !containsType(dependencies, t) {
		return fmt.Errorf("no dependencies found for type %s, unable to resolve", t.Name())
	}

	if _, ok := dependencies[t].instances[name]; ok {
		dependencies[t].primary = name
		return nil
	}

	return fmt.Errorf("dependency %s not registered, unable to set as primary", name)
}

// SetPrimary sets the primary bean name to be used.
//
// The 'interfaceRef' is a reference pointer to the bean interface so the bean factory knows what is the bean type, and
// it can properly register it. It can be a nil pointer to the interface. There are 2 ways of passing this a nil pointer.
// Let us assume we have a bean interface called IService. To obtain a nil pointer to the interface:
//
// Option 1:   (*IService)(nil)
//
//   Eg.   bean.SetPrimary((*IService)(nil), beanName)
//
// Option 2:   var reference *IService
//
//   Eg.   bean.SetPrimary(reference, beanName)
//
func SetPrimary(interfaceRef interface{}, name string) error {
	return SetPrimaryByType(getType(interfaceRef), name)
}

func containsType(c map[reflect.Type]*dependencyCollection, key reflect.Type) bool {
	if _, ok := c[key]; ok {
		return ok
	}

	return false
}

func getType(obj interface{}) reflect.Type {
	return reflect.TypeOf(obj).Elem()
}
