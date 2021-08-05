package ssc

import "sync"

// Context used as scoped temporary store for data
var context = map[Page]map[string]interface{}{}
var contextlock = &sync.Mutex{}

func SetContext(p Page, key string, value interface{}) {
	contextlock.Lock()
	space, ok := context[p]
	if !ok {
		space = map[string]interface{}{}
	}
	space[key] = value
	context[p] = space
	contextlock.Unlock()
}

func GetContext(p Page, key string) interface{} {
	return context[p][key]
}

func DelContext(p Page, key string) {
	contextlock.Lock()
	space, ok := context[p]
	if !ok {
		space = map[string]interface{}{}
	}
	if key != "" {
		delete(space, key)
	} else {
		space = map[string]interface{}{}
	}
	context[p] = space
	contextlock.Unlock()
}
