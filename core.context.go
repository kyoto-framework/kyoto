package kyoto

import "sync"

// Context used as scoped temporary store for data
var (
	context   = map[Page]map[string]interface{}{}
	contextrw = &sync.RWMutex{}
)

// SetContext of a page, page based key value store using a string as a key and an interface to store the data
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

// GetContext of the specified page, return data stored by key
func GetContext(p Page, key string) interface{} {
	contextrw.RLock()
	defer contextrw.RUnlock()
	return context[p][key]
}

// DelContext delete the page's context data by key
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
