package main

import (
	"strconv"

	pexu_dom "github.com/AnimusPEXUS/wasmtools/dom"
	"github.com/AnimusPEXUS/wasmtools/widgetcollection/select00"
)

func Select00FromMapUIntString(
	document *pexu_dom.Document,
	data map[uint]string,
	preselected uint,
	cb func(),
) *select00.Select00 {

	correct_map := make(map[string]string)

	for k, v := range data {
		correct_map[strconv.Itoa(int(k))] = v
	}

	ret := select00.NewSelect00(
		document,
		correct_map,
		strconv.Itoa(int(preselected)),
		cb,
	)

	return ret
}
