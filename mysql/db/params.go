package db

type QueryParams struct {
	SelectFiled []string
	Query       []string
	Not         []string
	Or          []string
	Offset      int
	Limit       int
	Sort        string
	Group       string
}

type UpdateParams struct {
	Update map[string]interface{}
	Query  []string
}

type DeleteParams struct {
	Query []string
	Model interface{}
}

type UpsertParams struct {
	ConflictFiled []string    // 冲突字段
	UpdateFiled   []string    // 更新字段
	Model         interface{} // model
}

type UpsertFiledDefaultParams struct {
	ConflictFiled []string               // 冲突字段
	Update        map[string]interface{} // 更新字段
	Model         interface{}            // model
}

type DistinctParams struct {
	DistinctFiled []string
	QueryParams
}

type JoinsParams struct {
	Joins []string
	QueryParams
}
