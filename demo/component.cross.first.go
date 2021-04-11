package main

import "github.com/yuriizinets/go-ssc"

type ComponentCrossFirst struct{}

func (*ComponentCrossFirst) Init() {

}

func (*ComponentCrossFirst) Async() error {
	return nil
}

func (*ComponentCrossFirst) AfterAsync() {

}

func (c *ComponentCrossFirst) Actions() ssc.ActionsMap {
	return ssc.ActionsMap{}
}
