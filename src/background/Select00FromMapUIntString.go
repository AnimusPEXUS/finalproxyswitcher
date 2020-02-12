package main

import (
	"sort"
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

	keys := []int{}

	for k, _ := range data {
		keys = append(keys, int(k))
	}

	sort.Ints(keys)

	correct_input := make([][2]string, 0)

	for _, k := range keys {
		correct_input = append(correct_input, [2]string{strconv.Itoa(int(k)), data[uint(k)]})
	}

	ret := select00.NewSelect00(
		document,
		correct_input,
		strconv.Itoa(int(preselected)),
		cb,
	)

	return ret
}
