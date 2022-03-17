package document

import (
	"database/sql"
	"flow-documents/mysql"
)

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
