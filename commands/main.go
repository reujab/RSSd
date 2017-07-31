package commands

type ListItem struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

const (
	List byte = iota
)
