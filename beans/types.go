package beans

// IResolveHandler defines an optional contract for dependencies to trigger a function on a successful get/resolve
//
// The implementation of this interface is optional, the beans manager will verify if this contract is implemented
// in a bean and trigger the handler(s) accordingly
//
type IResolveHandler interface {
	// OnResolve is a handler that gets triggered when a bean is resolved/fetched from the beans manager.
	OnResolve()
}

// IFirstTimeResolveHandler defines an optional contract for dependencies to trigger a function on the first time of
// a successful get/resolve. This function will not get triggered on sub-sequent resolved. It is meant for lazy
// initialization of components where said initialization is only desired when the component is used.
//
// The implementation of this interface is optional, the beans manager will verify if this contract is implemented
// in a bean and trigger the handler(s) accordingly
//
type IFirstTimeResolveHandler interface {
	OnFirstTimeResolve()
}

// ErrorCallback defines a function callback that can be registered when errors occur
type ErrorCallback func(error)

// MessageCallback defines a function callback that can be registered to broadcast messages
type MessageCallback func(string)

// ILogCallback defines the contract for logging callback assignments for this package
type ILogCallback interface {
	// SetErrorCallback sets a callback that will be triggered when an error occurs
	SetErrorCallback(callback ErrorCallback)
	// SetInfoCallback sets a callback that will be triggered when an info message is generated
	SetInfoCallback(callback MessageCallback)
	// SetDebugCallback sets a callback that will be triggered when a debug message is generated
	SetDebugCallback(callback MessageCallback)
}

// ILogger defines the contract for a full logger that can be used by this package
type ILogger interface {
	// Error logs an error to the logger
	Error(err error)
	// Info logs an info message to the logger
	Info(msg string)
	// Debug logs a debug message to the logger
	Debug(msg string)
}