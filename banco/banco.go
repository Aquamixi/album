package banco

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Conectar() (*sql.DB, error) {
	urlConn := "erick:@/recordings?charset=utf8&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", urlConn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
