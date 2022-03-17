package document

import "flow-documents/mysql"

type PostBody struct {
	Name      string `json:"name" validate:"required"`
	Url       string `json:"url" validate:"required"`
	ProjectId uint64 `json:"project_id" validate:"required,gte=1"`
}

func Post(userId uint64, post PostBody) (d Document, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO documents (user_id, name, url, project_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(userId, post.Name, post.Url, post.ProjectId)
	if err != nil {
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	d.Id = uint64(id)
	d.Name = post.Name
	d.Url = post.Url
	return
}
