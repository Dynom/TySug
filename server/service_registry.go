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
	defer rlh.lock.Unlock()

	rlh.services[listName] = svc
	return rlh
}

// GetServiceForList returns a service able to handle a specific list
func (rlh ServiceRegistry) GetServiceForList(name string) Service {
	rlh.lock.RLock()
	defer rlh.lock.RUnlock()

	return rlh.services[name]
}

// HasServiceForList returns true only if a service has been registered for a certain list
func (rlh ServiceRegistry) HasServiceForList(name string) bool {
	rlh.lock.RLock()
	defer rlh.lock.RUnlock()

	_, ok := rlh.services[name]
	return ok
}
