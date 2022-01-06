package smode

func SetContext(p Page, key string, value interface{}) {
	cmap[p].Context.Set(key, value)
}

func GetContext(p Page, key string) interface{} {
	return cmap[p].Context.Get(key)
}

func DelContext(p Page, key string) {
	cmap[p].Context.Del(key)
}
