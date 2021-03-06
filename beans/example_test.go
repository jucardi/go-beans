package beans_test

import "github.com/jucardi/go-beans/beans"

// IService is an example interface definition of a bean.
type IService interface {
	SomeMethod() bool
	GetName() string
}

// ServiceImpl is an implementation of IService that can be used as a bean.
type ServiceImpl struct {
}

func (t *ServiceImpl) SomeMethod() bool {
	return true
}

func (t *ServiceImpl) GetName() string {
	return "bean1"
}

func (t *ServiceImpl) init() {
	// any init logic of the bean.
}

// ####################################################################
// ## Things needed to register a component into the global factory. ##
// ####################################################################

// Here we define what the bean name will be. We declare it in a constant, so if multiple implementations
// exist, it is easy to find the proper bean name so the bean can be easily retrieved. If doing multiple
// implementations, the beans name must be unique
const BeanName = "some-bean-name"

// The following like is to validate the implementation of the interface on build, so it do not fail
// in runtime if a new function was added to the interface and missed to add the implementation in
// this section.
var _ IService = (*ServiceImpl)(nil)

// Registering the bean implementation.
func init() {
	// From time to time, the implementation of a component may depend on the initialization of other
	// components or packages. Using RegisterFunc can help implement a lazy initialization of a singleton
	// bean, where the singleton will be initialize on first use rather than on load of the application.
	//
	// The (*IService)(nil) is a nil pointer to the bean interface. It is required so the factory knows
	// what is the bean type so it can properly register it. An alternative to this approach is to declare
	// pointer to the interface without assigning any value to it, so it can be passed as an argument.
	//
	// Eg.
	//        var reference *IService
	//
	//        bean.RegisterFunc(reference, BeanName, func() interface{}) {
	//
	if err := beans.RegisterFunc((*IService)(nil), BeanName, func() interface{} {
		instance := &ServiceImpl{}
		instance.init()
		return instance
	}, true); err != nil {
		// Use your logger here.
		println(err.Error())
	}

	// To the bean as primary
	if err := beans.SetPrimary((*IService)(nil), BeanName); err != nil {
		// Use your logger here.
		println(err.Error())
	}
}

// With this singleton function, the primary implementation of the bean can be accessed directly from the
// package where it was defined, so third party consumers won't have to import the bean package, making it
// transparent to utilize any implemented service.
//
// This example shows how to use a variadic arg for an optional bean name to be provided so it can retrieve
// the proper bean by the given name. If no name is provided, it will return the primary bean.
func Service(name ...string) IService {
	if len(name) == 0 {
		return beans.Primary((*IService)(nil)).(IService)
	}
	return beans.Resolve((*IService)(nil), name[0]).(IService)
}
