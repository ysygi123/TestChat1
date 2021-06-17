package mysql

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

func init()  {
	var err error
	DB, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/mywork1")
	if err != nil {
		fmt.Println(err)
		panic("safd")
	}
	DB.SetMaxOpenConns(500)
	DB.SetMaxIdleConns(500)
}

func GetOneRow(rows *sql.Rows) (map[string]string, error) {
	col, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	vals := make([][]byte, len(col))
	scans := make([]interface{}, len(col))
	for k := range col{
		scans[k] = &vals[k]
	}
	result := make(map[string]string)

	for rows.Next(){
		err := rows.Scan(scans...)
		if err != nil {
			fmt.Println(err)
		}
		for k, v := range vals {
			key := col[k]
			result[key] = string(v)
		}
	}
	return result, nil
}

