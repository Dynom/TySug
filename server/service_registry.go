package server

import (
	"sync"
)

// ServiceRegistry holds (opinionated) service objects for handling specific requests
type ServiceRegistry struct {
	services map[string]Service
	lock     *sync.RWMutex
}

// NewServiceRegistry creates a new registry
func NewServiceRegistry() ServiceRegistry {
	return ServiceRegistry{
		services: make(map[string]Service),
		lock:     &sync.RWMutex{},
	}
}

// Register registers a new service with a list of references
func (rlh *ServiceRegistry) Register(listName string, svc Service) *ServiceRegistry {
	rlh.lock.Lock()
	rlh.services[listName] = svc
	rlh.lock.Unlock()

	return rlh
}

// GetServiceForList returns a service able to handle a specific list
func (rlh ServiceRegistry) GetServiceForList(name string) Service {
	rlh.lock.RLock()
	svc := rlh.services[name]
	rlh.lock.RUnlock()

	return svc
}

// HasServiceForList returns true only if a service has been registered for a certain list
func (rlh ServiceRegistry) HasServiceForList(name string) bool {
	rlh.lock.RLock()
	_, ok := rlh.services[name]
	rlh.lock.RUnlock()

	return ok
}
