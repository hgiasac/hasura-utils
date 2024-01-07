package graphql

// UpdateManyInput represents the common structure of update_many mutations
type UpdateManyInput struct {
	Where map[string]any `json:"where"`
	Set   map[string]any `json:"_set,omitempty"`
	Inc   map[string]any `json:"_inc,omitempty"`
}

// AffectedRowsOutput represent affected rows response
type AffectedRowsOutput struct {
	AffectedRows int `graphql:"affected_rows" json:"affected_rows"`
}
