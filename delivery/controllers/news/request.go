package news

type createNewsRequest struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
	Tags  []int  `json:"tags" validate:"required"`
}

type updateNewsRequest struct {
	Title string `json:"title" validate:"omitempty"`
	Body  string `json:"body" validate:"omitempty"`
	Tags  []int  `json:"tags" validate:"required"`
}
