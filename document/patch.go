package document

import (
	"flow-documents/mysql"
	"strings"
)

type PatchBody struct {
	Name      *string `json:"name" validate:"omitempty"`
	Url       *string `json:"url" validate:"omitempty"`
	ProjectId *uint64 `json:"project_id" validate:"omitempty,gte=1"`
}

func Patch(userId uint64, id uint64, new PatchBody) (d Document, notFound bool, err error) {
	// Get old
	d, notFound, err = Get(userId, id)
	if err != nil {
		return Document{}, false, err
	}
	if notFound {
		return Document{}, true, nil
	}

	// Generate query
	queryStr := "UPDATE schemes SET"
	var queryParams []interface{}
	if new.Name != nil {
		queryStr += " name = ?,"
		queryParams = append(queryParams, new.Name)
		d.Name = *new.Name
	}
	if new.Url != nil {
		queryStr += " url = ?,"
		queryParams = append(queryParams, new.Url)
		d.Url = *new.Url
	}
	if new.ProjectId != nil {
		queryStr += " project_id = ?"
		queryParams = append(queryParams, new.ProjectId)
		d.ProjectId = *new.ProjectId
	}
	queryStr = strings.TrimRight(queryStr, ",")
	queryStr += " WHERE user_id = ? AND id = ?"
	queryParams = append(queryParams, userId, id)

	// Update row
	db, err := mysql.Open()
	if err != nil {
		return Document{}, false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare(queryStr)
	if err != nil {
		return Document{}, false, err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(queryParams...)
	if err != nil {
		return Document{}, false, err
	}

	return
}
