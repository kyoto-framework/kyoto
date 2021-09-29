package ssc

import "sync"

// Context used as scoped temporary store for data
var context = map[Page]map[string]interface{}{}
var contextLock = &sync.RWMutex{}

func SetContext(p Page, key string, value interface{}) {
	contextLock.Lock()
	defer contextLock.Unlock()
	space, ok := context[p]
	if !ok {
		space = map[string]interface{}{}
	}
	space[key] = value
	context[p] = space
}

func GetContext(p Page, key string) interface{} {
	contextLock.RLock()
	defer contextLock.RUnlock()
	return context[p][key]
}

func DelContext(p Page, key string) {
	contextLock.Lock()
	defer contextLock.Unlock()
	space, ok := context[p]
	if !ok {
		space = map[string]interface{}{}
	}
	if key != "" {
		delete(space, key)
	} else {
		delete(context, p)
		return
	}
	context[p] = space
}
