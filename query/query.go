package query

const DEFAULT_LIMIT = 50

// Query
type Query struct {
	Limit           int                    `form:"limit" json:"limit"`
	Page            int                    `form:"page" json:"page"`
	SearchCondition map[string]interface{} `json:"searchCondition"`
	Range           *OrmSearchRange
	Sort            interface{} `json:"sort,omitempty"` // 可选信息.刚才用户提交的排序信息
}

//
type OrmSearchRange struct {
	Key   string
	Value interface{}
}

// Default
func (q *Query) Default() {
	if q == nil {
		q = new(Query)
	}
	if q.Limit == 0 {
		q.Limit = DEFAULT_LIMIT
	}
	if q.Page == 0 {
		q.Page = 1
	}

	q.Page = (q.Page - 1) * q.Limit
	return
}
