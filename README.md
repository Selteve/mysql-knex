# Go-MySQL-Knex

Go-MySQL-Knex 是一个 Go 语言的查询构建器,灵感来自于 Node.js 的 Knex.js 库。它提供了一种简单、直观的方式来构建和执行 SQL 查询,而无需直接编写 SQL 语句。

## 安装

使用 go get 安装包:

```bash
go get -u github.com/Selteve/mysql-knex
```

## 初始化

在使用查询构建器之前,您需要初始化数据库连接:

```go
import (
"github.com/your_username/go-mysql-knex/db"
"github.com/your_username/go-mysql-knex/types"
)
func main() {
dbOptions := &types.DB{
Host: "localhost:3306",
User: "your_username",
Password: "your_password",
DBName: "your_database",
}
db.Init(dbOptions)
defer db.Close()
// 您的代码...
}
```

## 使用示例

### 查询单条记录

```go
qb := db.NewQueryBuilder("users")
qb.Select("id", "name").Where(map[string]interface{}{"username": "17710620417"}).First()
```

### 使用select

```go
qb := db.NewQueryBuilder("users")
qb.Select("username", "password").Where(map[string]interface{}{"username": "17710620417"}).First()
``` 

### 查询多条记录

```go
qb := db.NewQueryBuilder("users")
qb.Select("id", "name").Where("id", ">", 1).Get()
```

### 插入记录

```go
qb := db.NewQueryBuilder("users")
qb.Insert(map[string]interface{}{
"name": "John Doe",
"email": "john.doe@example.com",
})
```

### 更新记录

```go
qb := db.NewQueryBuilder("users")
qb.Update(map[string]interface{}{
"name": "John Smith",
}).Where("id", "=", 1)
```

### 删除记录

```go
qb := db.NewQueryBuilder("users")
qb.Delete().Where("id", "=", 1)
```

## 可用方法

- `DB(tableName string)`: 指定要操作的表名
- `Select(columns ...string)`: 选择特定的列
- `Where(conditions map[string]interface{})`: 添加 WHERE 条件
- `OrderBy(column string, direction string)`: 添加排序
- `Limit(limit int)`: 限制结果数量
- `Offset(offset int)`: 设置结果偏移量
- `First()`: 获取第一条匹配的记录
- `Get()`: 获取所有匹配的记录
- `Insert(data map[string]interface{})`: 插入新记录
- `Update(data map[string]interface{})`: 更新记录
- `Delete()`: 删除记录

## 注意事项

- 此包使用参数化查询来防止 SQL 注入,但在使用时仍需注意安全性。
- 确保在应用程序退出前调用 `db.Close()` 以正确关闭数据库连接。
- 此包目前主要支持 MySQL 数据库,如需支持其他数据库可能需要进行修改。

## 贡献

欢迎提交 issues 和 pull requests 来帮助改进这个包。

## 许可证

此项目采用 MIT 许可证。详情请见 [LICENSE](LICENSE) 文件.


