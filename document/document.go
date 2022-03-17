package document

type Document struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	ProjectId uint64 `json:"project_id"`
}
