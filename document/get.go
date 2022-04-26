package document

import "flow-documents/mysql"

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
	defer rows.Close()

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	err = rows.Scan(&d.Name, &d.Url, &d.ProjectId)
	if err != nil {
		return
	}

	d.Id = id
	return
}
