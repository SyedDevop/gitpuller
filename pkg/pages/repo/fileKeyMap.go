package repo

import "github.com/charmbracelet/bubbles/key"

type FileKeyMap struct {
	Open       key.Binding
	SelectItem key.Binding
	Conform    key.Binding
	SelectAll  key.Binding
	DSelectAll key.Binding
	UpDir      key.Binding
}

func NewFileKeyMap() *FileKeyMap {
	km := new(FileKeyMap)

	km.Open = key.NewBinding(
		key.WithKeys(
			"enter",
		),
		key.WithHelp(
			"enter",
			"Open",
		),
	)

	km.SelectItem = key.NewBinding(
		key.WithKeys(
			" ",
		),
		key.WithHelp(
			"space",
			"Select",
		),
	)

	km.UpDir = key.NewBinding(
		key.WithKeys(
			"backspace",
			"b",
		),
		key.WithHelp(
			"⌫/b",
			"Go up a directory",
		),
	)
	km.Conform = key.NewBinding(
		key.WithKeys(
			"y",
		),
		key.WithHelp(
			"y",
			"Conform the Selected Items",
		),
	)

	km.SelectAll = key.NewBinding(
		key.WithKeys(
			"a",
			"A",
		),
		key.WithHelp(
			"a/A",
			"Select All the File/Dir",
		),
	)
	km.DSelectAll = key.NewBinding(
		key.WithKeys(
			"d",
			"D",
		),
		key.WithHelp(
			"d/D",
			"D Select All the File/Dir",
		),
	)
	return km
}
