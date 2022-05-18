package db

type QueryParams struct {
	Query  []string
	Offset int
	Limit  int
	Sort   string
	Group  string
}

type UpdateParams struct {
	Update map[string]interface{}
	Query  []string
}

type DeleteParams struct {
	Query []string
	Model interface{}
}
