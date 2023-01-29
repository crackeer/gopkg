package database

import (
	"gorm.io/gorm"
)

// GetSQLiteTables
//
//	@param db
//	@return []map
func GetSQLiteTables(db *gorm.DB) []map[string]interface{} {
	list := []map[string]interface{}{}
	db.Table("sqlite_master").Where(map[string]interface{}{
		"type": "table",
	}).Find(&list)
	return list
}

// ExecSQL
//
//	@param db
//	@return []map
func ExecSQL(db *gorm.DB, sql string) error {
	return db.Exec(sql).Error
}
