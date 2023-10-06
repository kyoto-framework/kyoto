package component

type Disposable struct {
	Name
}

func (*Disposable) Marshal() string {
	return "disposable"
}

func (*Disposable) Unmarshal(str string) {
	return
}
