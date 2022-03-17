package document

import (
	"flow-documents/mysql"
)

func GetList(userId uint64, projectId *uint64) (documents []Document, err error) {
	// Generate query
	queryStr := "SELECT id, name, url, project_id FROM documents WHERE user_id = ?"
	queryParams := []interface{}{userId}
	if projectId != nil {
		queryStr += " AND project_id = ?"
		queryParams = append(queryParams, *projectId)
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

	rows, err := stmtOut.Query(queryParams...)
	if err != nil {
		return
	}

	for rows.Next() {
		d := Document{}
		err = rows.Scan(&d.Id, &d.Name, &d.Url, &d.ProjectId)
		if err != nil {
			return
		}
		documents = append(documents, d)
	}

	return
}
