package document

import (
	"database/sql"
	"flow-documents/mysql"
)

type Document struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	ProjectId uint64 `json:"project_id"`
}

type Post struct {
	Name      string `json:"name" validate:"required"`
	Url       string `json:"url" validate:"required"`
	ProjectId uint64 `json:"project_id" validate:"required,gte=1"`
}

type Patch struct {
	Name      *string `json:"name" validate:"omitempty"`
	Url       *string `json:"url" validate:"omitempty"`
	ProjectId *uint64 `json:"project_id" validate:"omitempty,gte=1"`
}

func Get(userId uint64, id uint64) (d Document, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT name, url, project_id FROM documents WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, id)
	if err != nil {
		return
	}

	// TODO: uint64に対応
	var (
		name      string
		url       string
		projectId uint64
	)
	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	err = rows.Scan(&name, &url, &projectId)
	if err != nil {
		return Document{}, false, err
	}

	d.Id = id
	d.Name = name
	d.Url = url
	d.ProjectId = projectId

	return
}

func Insert(userId uint64, post Post) (d Document, err error) {
	// Insert DB
	db, err := mysql.Open()
	if err != nil {
		return Document{}, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO documents (user_id, name, url, project_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return Document{}, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(userId, post.Name, post.Url, post.ProjectId)
	if err != nil {
		return Document{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Document{}, err
	}

	d.Id = uint64(id)
	d.Name = post.Name
	d.Url = post.Url

	return
}

func Update(userId uint64, id uint64, new Patch) (d Document, notFound bool, err error) {
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

func Delete(userId uint64, id uint64) (notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("DELETE FROM documents WHERE user_id = ? AND id = ?")
	if err != nil {
		return false, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(userId, id)
	if err != nil {
		return false, err
	}
	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affectedRowCount == 0 {
		// Not found
		return true, nil
	}

	return false, nil
}

func GetList(userId uint64, projectId *uint64) (documents []Document, err error) {
	// Generate query
	queryStr := "SELECT id, name, url, project_id FROM documents WHERE user_id = ?"
	if projectId != nil {
		queryStr += " AND project_id = ?"
	}

	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	// TODO: sort
	stmtOut, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	var rows *sql.Rows
	if projectId == nil {
		rows, err = stmtOut.Query(userId)
	} else {
		rows, err = stmtOut.Query(userId, *projectId)
	}
	if err != nil {
		return
	}

	for rows.Next() {
		// TODO: uint64に対応
		var (
			id        uint64
			name      string
			url       string
			projectId uint64
		)
		err = rows.Scan(&id, &name, &url, &projectId)
		if err != nil {
			return
		}

		documents = append(documents, Document{id, name, url, projectId})
	}

	return
}
