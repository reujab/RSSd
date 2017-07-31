package commands

type ListItem struct {
	GUID  string `json:"guid"`
	Name  string `json:"name"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

const (
	List byte = iota
	Read
)
