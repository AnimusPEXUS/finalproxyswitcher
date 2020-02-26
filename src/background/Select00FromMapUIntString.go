package main

import (
	"sort"
	"strconv"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
	"github.com/AnimusPEXUS/wasmtools/widgetcollection"
)

func Select00FromMapUIntString(
	etc *elementtreeconstructor.ElementTreeConstructor,
	data map[uint]string,
	preselected uint,
	cb func(),
) *widgetcollection.Select00 {

	keys := []int{}

	for k, _ := range data {
		keys = append(keys, int(k))
	}

	sort.Ints(keys)

	correct_input := make([][2]string, 0)

	for _, k := range keys {
		correct_input = append(correct_input, [2]string{strconv.Itoa(int(k)), data[uint(k)]})
	}

	ret := widgetcollection.NewSelect00(
		etc,
		correct_input,
		strconv.Itoa(int(preselected)),
		cb,
	)

	return ret
}
