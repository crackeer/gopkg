package table

import (
	"errors"
	"fmt"

	"github.com/crackeer/gopkg/util"
	"gorm.io/gorm"
)

// Table
type Table struct {
	Name    string
	DB      *gorm.DB
	OrderBy string
	Driver  string

	primaryKey string
}

// NewTable
//  @param name
//  @param db
//  @return *Table
func NewTable(db *gorm.DB, name string) (*Table, error) {
	if db == nil {
		return nil, errors.New("db nil")
	}
	return &Table{
		Name: name,
		DB:   db,
	}, nil
}

// SetPrimaryKey ...
//  @receiver table
//  @param key
func (table *Table) SetPrimaryKey(key string) {
	table.primaryKey = key
}

// SetDriver
//  @receiver table
//  @param key
func (table *Table) SetDriver(key string) {
	table.Driver = key
}

// GetPrimaryKey
//  @receiver table
//  @return string
func (table *Table) GetPrimaryKey() string {
	if len(table.primaryKey) > 0 {
		return table.primaryKey
	}
	list := []map[string]interface{}{}
	sql := fmt.Sprintf("pragma table_info('%s')", table.Name)
	table.DB.Raw(sql).Scan(&list)
	for _, item := range list {
		mapLoader := util.LoadMap(item)
		if mapLoader.GetInt64("pk", 0) > 0 {
			return mapLoader.GetString("name", "")
		}
	}

	return "id"
}

// Create
//  @receiver table
//  @param data
//  @return map
func (table *Table) Create(data map[string]interface{}) (map[string]interface{}, error) {
	err := table.DB.Table(table.Name).Create(&data).Error
	return data, err
}

// Update
//  @receiver table
//  @param id
//  @param data
//  @return int64
func (table *Table) Update(where map[string]interface{}, data map[string]interface{}) int64 {
	return table.DB.Table(table.Name).Where(where).Updates(data).RowsAffected
}

// Query
//  @receiver table
//  @param query
//  @return []map
func (table *Table) Query(query map[string]interface{}, limit int) []map[string]interface{} {
	list := []map[string]interface{}{}
	table.DB.Table(table.Name).Where(query).Limit(limit).Find(&list)
	return list
}

// GetPageList
//  @receiver table
//  @param query
//  @param page
//  @param pageSize
//  @return []map
func (table *Table) GetPageList(query map[string]interface{}, page, pageSize int64) []map[string]interface{} {
	list := []map[string]interface{}{}
	offset := (page - 1) * pageSize
	table.DB.Table(table.Name).Where(query).Offset(int(offset)).Order("id desc").Limit(int(pageSize)).Find(&list)
	return list
}

// Count
//  @receiver table
//  @param query
//  @return int64
func (table *Table) Count(query map[string]interface{}) int64 {
	var count int64
	table.DB.Table(table.Name).Where(query).Count(&count)
	return count
}

// Delete
//  @receiver table
//  @param primaryKey
//  @param value
//  @return int64
func (table *Table) Delete(primaryKey string, value interface{}) int64 {
	sql := fmt.Sprintf("delete from %s where %s=?", table.Name, primaryKey)
	return table.DB.Exec(sql, value).RowsAffected
}

// Distinct
//  @receiver table
//  @param field
//  @param where
//  @return []interface{}
func (table *Table) Distinct(field string, where map[string]interface{}) []interface{} {
	list := []interface{}{}
	table.DB.Table(table.Name).Select(fmt.Sprintf("distinct(%s) as %s", field, field)).Where(where).Pluck(field, &list)
	return list
}
