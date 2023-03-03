package res

type LogIndexShowFoldsRes struct {
	Title    string                  `json:"title"`
	Id       string                  `json:"id"`
	PathId   string                  `json:"path"`
	Children []*LogIndexShowFoldsRes `json:"children"`
}
