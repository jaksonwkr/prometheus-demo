package helloworld

// HelloWorld interface that specify service port to application adapter
type HelloWorld interface {
	SayHelloWorld() (string, error)
}

// Service is the service for helloworld
type Service struct{}

// SayHelloWorld implements HelloWorld interface
func (svc Service) SayHelloWorld() (string, error) {
	return "Hello World!", nil
}

// New creates a new service
func New() (*Service, error) {
	svc := &Service{}
	return svc, nil
}
