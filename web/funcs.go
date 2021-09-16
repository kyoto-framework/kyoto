package main

import (
	"fmt"
	"html/template"

	"github.com/yuriizinets/go-common"
	"github.com/yuriizinets/go-ssc"
)

func tfuncs() template.FuncMap {
	f := ssc.Funcs()
	f["fprice"] = func(price int) string {
		return fmt.Sprintf("%v", price/100)
	}
	common.TFMAttach(&f)
	return f
}
