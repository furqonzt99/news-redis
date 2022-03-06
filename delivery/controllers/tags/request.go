package tags

type tagRequest struct {
	Name string `json:"name" validate:"required"`
}
