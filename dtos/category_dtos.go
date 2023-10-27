package dtos

type CategoryProducts struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
