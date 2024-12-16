package models

type Item struct {
	ID   string
	Name string
	Desc string
	Code string
}

func (i Item) Title() string       { return i.Name }
func (i Item) Description() string { return i.Desc }
func (i Item) FilterValue() string { return i.Name }
