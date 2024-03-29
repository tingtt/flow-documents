package document

import "flow-documents/mysql"

func DeleteAll(userId uint64) (err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("DELETE FROM documents WHERE user_id = ?")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(userId)
	if err != nil {
		return
	}

	return
}
