package serviceregistry

// Interface to talking to external
type ServiceRegistry interface {

	// Register a service with local agent
	Register() error

	// Deregister a service with local agent
	Deregister() error
}
