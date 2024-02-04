package component

// Disposable component means that
// we don't want to store state or use dynamic actions.
// On action call you'll get explicit error about wrong usage.
type Disposable struct {
	Name
}

// Marshal for disposable returns "disposable" string.
// This flag will help to detect wrong action usage.
func (*Disposable) Marshal(src any) string {
	return "disposable"
}

// Unmarshal for disposable returns nothing.
func (*Disposable) Unmarshal(dst any, str string) {
	return
}
