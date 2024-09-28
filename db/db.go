package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	Types "github.com/Selteve/mysql-knex/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// 声明全局数据库变量
var database *sqlx.DB

// 初始化
func Init(Options *Types.DB) {
	Connect(Options)
}

// 连接数据库
func Connect(Options *Types.DB) {
	var err error
	database, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", Options.User, Options.Password, Options.Host, Options.DBName))
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	// 检查数据库连接
	if err = database.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	fmt.Println("Connected to DB successfully")
}

// 关闭数据库连接
func Close() {
	if database != nil {
		err := database.Close()
		if err != nil {
			log.Fatalf("Error closing database: %v", err)
		}
		fmt.Println("Database connection closed")
	}
}

// QueryBuilder 结构体
type QueryBuilder struct {
	table       string
	selectCols  []string
	whereClause []string
	whereArgs   []interface{}
	orderByCol  string
	orderByDir  string
	limitValue  int
	offsetValue int
}

// DB 函数,类似于 Knex 中的用法
func DB(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

// NewQueryBuilder 函数
func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

// Where 方法 (支持两种写法)
func (qb *QueryBuilder) Where(args ...interface{}) *QueryBuilder {
	switch len(args) {
	case 1:
		// 处理 map[string]interface{} 参数
		if conditions, ok := args[0].(map[string]interface{}); ok {
			for column, value := range conditions {
				qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s = ?", column))
				qb.whereArgs = append(qb.whereArgs, value)
			}
		}
	case 3:
		// 处理 (column string, operator string, value interface{}) 参数
		if column, ok := args[0].(string); ok {
			if operator, ok := args[1].(string); ok {
				qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s %s ?", column, operator))
				qb.whereArgs = append(qb.whereArgs, args[2])
			}
		}
	}
	return qb
}

// OrderBy 方法
func (qb *QueryBuilder) OrderBy(column string, direction string) *QueryBuilder {
	qb.orderByCol = column
	qb.orderByDir = direction
	return qb
}

// Limit 方法
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limitValue = limit
	return qb
}

// Offset 方法
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offsetValue = offset
	return qb
}

// Select 方法
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.selectCols = columns
	return qb
}

// buildSelectQuery 方法
func (qb *QueryBuilder) buildSelectQuery() (string, []interface{}) {
	query := "SELECT "
	if len(qb.selectCols) > 0 {
		query += strings.Join(qb.selectCols, ", ")
	} else {
		query += "*"
	}
	query += " FROM " + qb.table

	if len(qb.whereClause) > 0 {
		query += " WHERE " + strings.Join(qb.whereClause, " AND ")
	}

	if qb.orderByCol != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", qb.orderByCol, qb.orderByDir)
	}

	if qb.limitValue > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limitValue)
	}

	if qb.offsetValue > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offsetValue)
	}

	return query, qb.whereArgs
}

// First 方法
func (qb *QueryBuilder) First() (map[string]interface{}, error) {
	qb.Limit(1)
	query, args := qb.buildSelectQuery()
	row := database.QueryRowx(query, args...)
	result := make(map[string]interface{})
	err := row.MapScan(result)
	if err != nil {
		return nil, err
	}
	return qb.convertResult(result), nil
}

// Get 方法
func (qb *QueryBuilder) Get() ([]map[string]interface{}, error) {
	query, args := qb.buildSelectQuery()
	rows, err := database.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		result := make(map[string]interface{})
		err := rows.MapScan(result)
		if err != nil {
			return nil, err
		}
		results = append(results, qb.convertResult(result))
	}
	return results, nil
}

// Insert 方法
func (qb *QueryBuilder) Insert(data map[string]interface{}) (sql.Result, error) {
	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	placeholders := make([]string, 0, len(data))

	for column, value := range data {
		columns = append(columns, column)
		values = append(values, value)
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		qb.table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	return database.Exec(query, values...)
}

// Update 方法
func (qb *QueryBuilder) Update(data map[string]interface{}) (sql.Result, error) {
	setClause := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))

	for column, value := range data {
		setClause = append(setClause, fmt.Sprintf("%s = ?", column))
		values = append(values, value)
	}

	query := fmt.Sprintf("UPDATE %s SET %s", qb.table, strings.Join(setClause, ", "))

	if len(qb.whereClause) > 0 {
		query += " WHERE " + strings.Join(qb.whereClause, " AND ")
		values = append(values, qb.whereArgs...)
	}

	return database.Exec(query, values...)
}

// Delete 方法
func (qb *QueryBuilder) Delete() (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s", qb.table)

	if len(qb.whereClause) > 0 {
		query += " WHERE " + strings.Join(qb.whereClause, " AND ")
	}

	return database.Exec(query, qb.whereArgs...)
}

func (qb *QueryBuilder) convertResult(result map[string]interface{}) map[string]interface{} {
	for key, value := range result {
		switch v := value.(type) {
		case []byte:
			result[key] = string(v)
		}
	}
	return result
}

func NewDB(table string) *Types.DBName {
	return &Types.DBName{Table: table}
}