package dto

type Pagination struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func Empty() map[string]interface{} {
	return map[string]interface{}{}
}
