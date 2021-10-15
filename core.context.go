package kyoto

import "sync"

// Context used as scoped temporary store for data
var context = map[Page]map[string]interface{}{}
var contextrw = &sync.RWMutex{}

func SetContext(p Page, key string, value interface{}) {
	contextrw.Lock()
	defer contextrw.Unlock()
	space, ok := context[p]
	if !ok {
		space = map[string]interface{}{}
	}
	space[key] = value
	context[p] = space
}

func GetContext(p Page, key string) interface{} {
	contextrw.RLock()
	defer contextrw.RUnlock()
	return context[p][key]
}

func DelContext(p Page, key string) {
	contextrw.Lock()
	defer contextrw.Unlock()
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
