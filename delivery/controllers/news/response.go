package news

type newsResponse struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Status string   `json:"status"`
	Tags   []string `json:"tags"`
}
