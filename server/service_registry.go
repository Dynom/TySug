package server

import (
	"sync"
)

type ServiceRegistry struct {
	services map[string]Service
	lock     *sync.RWMutex
}

func NewServiceRegistry() ServiceRegistry {
	return ServiceRegistry{
		services: make(map[string]Service),
		lock:     &sync.RWMutex{},
	}
}

func (rlh *ServiceRegistry) Register(listName string, svc Service) *ServiceRegistry {
	rlh.lock.Lock()
	defer rlh.lock.Unlock()

	rlh.services[listName] = svc
	return rlh
}

func (rlh ServiceRegistry) GetServiceForList(name string) Service {
	rlh.lock.RLock()
	defer rlh.lock.RUnlock()

	return rlh.services[name]
}

func (rlh ServiceRegistry) HasServiceForList(name string) bool {
	rlh.lock.RLock()
	defer rlh.lock.RUnlock()

	_, ok := rlh.services[name]
	return ok
}
