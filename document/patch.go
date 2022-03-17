package document

import "flow-documents/mysql"

type PatchBody struct {
	Name      *string `json:"name" validate:"omitempty"`
	Url       *string `json:"url" validate:"omitempty"`
	ProjectId *uint64 `json:"project_id" validate:"omitempty,gte=1"`
}

func Patch(userId uint64, id uint64, new PatchBody) (d Document, notFound bool, err error) {
	// Get old
	old, notFound, err := Get(userId, id)
	if err != nil {
		return Document{}, false, err
	}
	if notFound {
		return Document{}, true, nil
	}

	// Set no update values
	if new.Name == nil {
		new.Name = &old.Name
	}
	if new.Url == nil {
		new.Url = &old.Url
	}
	if new.ProjectId == nil {
		new.ProjectId = &old.ProjectId
	}

	// Update row
	db, err := mysql.Open()
	if err != nil {
		return Document{}, false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("UPDATE documents SET name = ?, url = ?, project_id = ? WHERE user_id = ? AND id = ?")
	if err != nil {
		return Document{}, false, err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(new.Name, new.Url, new.ProjectId, userId, id)
	if err != nil {
		return Document{}, false, err
	}

	d.Id = id
	d.Name = *new.Name
	d.Url = *new.Url
	d.ProjectId = *new.ProjectId

	return
}
