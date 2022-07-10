package dao

import "database/sql"

var DB *sql.DB

func InitMysql(dns string) (err error) {
	database, err := sql.Open("mysql", dns)
	if err != nil {
		return err
	}
	DB = database
	err = DB.Ping()
	if err != nil {
		return err
	}

	return nil
}
