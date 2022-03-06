package tags

type TagRequest struct {
	Name string `json:"name" validate:"required"`
}
