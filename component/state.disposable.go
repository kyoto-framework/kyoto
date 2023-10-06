package component

//
type Disposable struct {
	Name
}

func (*Disposable) Marshal(state any) string {
	return "disposable"
}

func (*Disposable) Unmarshal(str string, state any) {
	return
}
