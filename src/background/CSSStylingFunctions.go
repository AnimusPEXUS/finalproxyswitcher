package main

import (
	"github.com/AnimusPEXUS/wasmtools/dom/elementtreeconstructor"
)

func applyEditorStyle(ed *elementtreeconstructor.ElementMutator) {
	ed.
		SetStyle("border", "none").
		SetStyle("border-left", "1px black solid").
		// SetStyle("margin-top", "5px").
		// SetStyle("margin-botom", "5px").
		SetStyle("width", "50px")
}

func applyBorder(ed *elementtreeconstructor.ElementMutator) {
	ed.
		SetStyle("border", "1px black solid").
		SetStyle("border-radius", "3px")
}

func applyMarginRight(ed *elementtreeconstructor.ElementMutator) {
	ed.
		SetStyle("margin-right", "1px")
}

func applyProxyEditorSpanStyle(ed *elementtreeconstructor.ElementMutator) {
	ed.
		SetStyle("padding-top", "1px").
		SetStyle("padding-bottom", "1px").
		SetStyle("white-space", "nowrap")
	applyBorder(ed)
	applyMarginRight(ed)
}

func applyButtonStyle(ed *elementtreeconstructor.ElementMutator) {
	applyBorder(ed)
	applyMarginRight(ed)
}

func applyAStyle(ed *elementtreeconstructor.ElementMutator) {
	ed.
		SetStyle("color", "blue").
		SetStyle("cursor", "pointer").
		SetStyle("text-decoration", "underline")
}

func applyBlackRoundedBoxInRuleEditor(ed *elementtreeconstructor.ElementMutator) {
	ed.
		SetStyle("border", "1px dotted black").
		SetStyle("border-left", "3px solid orange").
		SetStyle("border-radius", "3px").
		SetStyle("padding", "3px")
}

func applySpanChangedAsterisk(ed *elementtreeconstructor.ElementMutator) {
	ed.
		SetStyle("display", "none").
		SetStyle("color", "red")
}
