package menu

type Menu struct {
	Items []*Item
}

type Item struct {
	Name     string
	Icon     string
	Link     string
	Children []*Item
}
