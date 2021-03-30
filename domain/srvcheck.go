// Create package in v.1.0.0
// Same as entities, struct and method in domain package will used in all layer.
// srvcheck.go is file that define model as struct and abstract method of model as interface.
// Also, it declare usecase interface used as business layer.

// srvcheck domain is managing the state of the service (elasticsearch, swarm, consul, etc ...) periodically

// All model struct and interface is about service check domain
// Most importantly, it only defines and does not implement interfaces.

package domain
