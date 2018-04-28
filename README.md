# go-beans
## Components manager for easy dependency injection. Inspired in the Java Spring Framework beans.

`go-beans` is a components manager / components factory that easily allows dependency injection and easily switch these dependencies with other implementations of the components (beans)

For example, let's assume we wanted to create an Alerts handler to send alerts on events in a system. There are multiple mechanisms to send out an alert (Email, SMS, chat message, push notification, RSS, etc.).

We could define a base interface

```Go
type IAlertHandler interface {
    Send(alert *Alert) error
}
```

And then have multiple implementations of this interface to send an alert via a specific delivery mechanism.

```Go
type EmailAlertHandler struct {
    // some fields here if needed
}

func (e *EmailAlertHandler) Send(alert *Alert) error {
    // handle some logic to send the SMTP message
    return nil
}
```


```Go
type SMSAlertHandler intestructrface {
    // some fields here if needed
}

func (s *SMSAlertHandler) Send(alert *Alert) error {
    // handle some logic to send the SMS message
    return nil
}
```

Using the `go-beans` factory, we could register these beans with a name (Eg. "email" and "sms"), and control which one of them will be used with a simple configuration.

```Go
// Alert trigger handler

func GetHandler() IAlertHandler {
    // name of the handler to be used. Eg "email"
    alertType := config.AlertType
    return beans.Resolve((*IAlertHandler)(nil), alertType).(IAlertHandler)
}


func sendAlert(alert *Alert) {
    err := GetHandler().Send(alert)
    if err != nil {
        log.Error(err)
    }
}
```

Additionally, when doing integration testing, the beans factory allows to overwrite a bean that goes by a given name, or a bean that is set as primary, so even using the real bean name, a mock implementation can be hooked into the application, so it can still be tested as a whole.

## How to create a bean?

For the previous example, this would be an implementation of the `email` alert handler

**`email.go`**
```Go
package email

import "some-path-to-the-alerts-package/alerts"

// The struct that represents the email handler implementation.
type EmailAlertHandler struct {
}

// The method to be implemented
func (e *EmailAlertHandler) Send(alert *Alert) error {
    return nil
}

func (t *ServiceImpl) init() {
    // any init logic of the bean.
}

// Here we define what the bean name will be. We declare it in a constant, so if multiple implementations
// exist, it is easy to find the proper bean name so the bean can be easily retrieved. If doing multiple
// implementations, the beans name must be unique
const BeanName = "email"

// The following like is to validate the implementation of the interface on build, so it do not fail in
// runtime if a new function was added to the interface and missed to add the implementation in
// this section.
var _ IAlertHandler = (*EmailAlertHandler)(nil)

// In this example we'll use a singleton instance to be registered as the bean.
var instance *EmailAlertHandler

// Registering the bean implementation.
func init() {

    // In this example we do a lazy construct of a singleton bean, where the singleton will be
    // initialized on first use rather than on load of the application.
    //
    // The (*IAlertHandler)(nil) is a nil pointer to the bean interface. It is required so the
    // factory knows what is the bean type so it can properly register it. An alternative to this
    // approach is to declare pointer to the interface without assigning any value to it, so it
    // can be passed as an argument.
    //
    // Eg.
    //        var reference *IAlertHandler
    //
    //        bean.RegisterFunc(reference, BeanName, func() interface{}) {
    //
    beans.RegisterFunc((*IAlertHandler)(nil), BeanName, func() interface{} {
        if instance != nil {
            return instance
    }

        instance = &EmailAlertHandler{}
        instance.init()
        return instance
    })
}
```

In the `alerts` package we could have the following function that will simply retrieve the instance to be used
so other services will not have to talk to the beans factory directly
```Go
package alerts

// With this singleton function, the primary implementation of the bean can be accessed directly from
// the package where it was defined, so third party consumers won't have to import the bean package,
// making it transparent to utilize any implemented service.
func Get() IAlertHandler {
    return beans.Resolve((*IAlertHandler)(nil), config.AlertType).(IAlertHandler)
}
```

## Replacing an existing bean with a mock for testing

In a test file (or a test utils file that could be used across other consuming services), implement a struct
of the interface to be used as a mock. In the following example, the mock implementation also has a logic to mock responses and validate calls. This can be extended to be more

**`alerts_test.go`**
```Go
package alerts

/* =========== Mock implementation ============ */

type MockAlertHandler struct {
    sendResponse error
}

func (m *MockAlertHandler) Send(alert *Alert) error {
    return m.sendResponse
}

// Additional method useful to add a expected response with one of the interface methods are called
func (m *MockAlertHandler) WhenSend(expectedResponse error) {
    m.sendResponse = expectedResponse
}

var (
    // Ensures the struct implements the interface on compile time, to prevent failures in runtime
    _ IAlertHandler = (*MockAlertHandler)(nil)

    // The mock instance
    mock *MockAlertHandler
)

// To obtain and register a mock implementation.
func Mock() *MockAlert {
    if mock != nil {
        return mock
    }

    mock = &MockAlertHandler{}
    beans.RegisterFunc((*IAlertHandler)(nil), "email", func() interface{} {
        return mock
    })

	return mock
}

/* =========== End of Mock implementation ============ */

/* =========== Tests begin here ============ */

// Initialize tests.
func init() {
    // Setting AllowOverrides to 'true' will allow the beans factory to replace existing beans
    // with mock implementations.
    beans.SetAllowOverrides(true)
}

func TestSomething(t *testing.T) {
    mock := Mock()
    mock.WhenAlert(errors.New("some error"))

    assert.Equal(t, errors.New("some error"), mock.Send(nil))
}
```

## Getting started.

To start using this package

```
go get github.com/jucardi/go-beans
```