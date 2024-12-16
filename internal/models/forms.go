package models

type FormsResponse struct {
	Data []Form `json:"data"`
}

func (f *Form) ToItem() Item {
	return Item{
		ID:   f.ID,
		Name: f.Name,
		Desc: f.Code,
	}
}

func FormsToItems(forms []Form) []Item {
	items := []Item{}
	for _, form := range forms {
		items = append(items, form.ToItem())
	}
	return items
}
