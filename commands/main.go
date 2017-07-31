package commands

type ListItem struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

const (
	List byte = iota
	Read
)
