package mysql

import (
	"fmt"
	"gitlab.tessan.com/data-center/tessan-erp-common/convert"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"time"
)

// Condition MySQL 搜索条件
type Condition struct {
	Symbol      OperatorSymbol `json:"symbol"`       // 操作符
	ColumnName  string         `json:"column_name"`  // 相关字段
	ColumnValue string         `json:"column_value"` // 字段值
}

// OperatorSymbol 。。。
type OperatorSymbol string

//goland:noinspection ALL
const (
	EQUAL              OperatorSymbol = "="
	LIKE               OperatorSymbol = "like"
	IN                 OperatorSymbol = "in"
	NULL               OperatorSymbol = "is null"
	NOT_NULL           OperatorSymbol = "is not null"
	GREATER_THEN       OperatorSymbol = ">"
	GREATER_THEN_EQUAL OperatorSymbol = ">="
	LESS_THAN          OperatorSymbol = "<"
	LESS_THAN_EQUAL    OperatorSymbol = "<="
	NOT_EQUAL          OperatorSymbol = "<>"
)

type Selected string

const (
	COUNT        Selected = "count(*)"
	ALL          Selected = "*"
	DISTINCT_ALL Selected = " DISTINCT * "
)

//goland:noinspection GoSnakeCaseUsage
const (
	BASIC_MODEL_CREATED_AT  = "created_at"
	BASIC_MODEL_UPDATED_AT  = "updated_at"
	BASIC_MODEL_DELETED_AT  = "deleted_at"
	BASIC_MODEL_PRIMARY_KEY = "id"
)

// CalcMysqlBatchSize 批量插入时计算长度的
func CalcMysqlBatchSize(data interface{}) int {
	count := countStructFields(reflect.ValueOf(data))
	return 60000 / count
}

func countStructFields(v reflect.Value) int {
	var count int
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.Type.Kind() == reflect.Struct && isInnerStruct(field) {
			num := countStructFields(v.Field(i))
			count += num
		} else if !strings.Contains(field.Tag.Get("gorm"), "foreignKey") {
			count += 1
		}
	}
	return count
}

func isInnerStruct(structField reflect.StructField) bool {
	if structField.Type == reflect.TypeOf(time.Time{}) || structField.Type == reflect.TypeOf(&time.Time{}) {
		return false
	}
	if strings.Contains(structField.Tag.Get("gorm"), "foreignKey") {
		return false
	}

	return true
}

// sqlGeneration sql生成器
type sqlGeneration struct {
	sql string
}

// InitSqlGeneration 初始化
func InitSqlGeneration(model schema.Tabler, selected Selected) *sqlGeneration {
	var s sqlGeneration
	s.sql = fmt.Sprintf("select %s from %s where 1=1 ", selected, model.TableName())
	return &s
}

// AddConditions 添加 where 条件，都是 and
func (s *sqlGeneration) AddConditions(symbol OperatorSymbol, conditions ...string) *sqlGeneration {
	length := len(conditions)
	for i := range conditions {
		if i%2 != 0 {
			continue
		}
		if i+1 >= length {
			break
		}
		if conditions[i+1] != "" {
			tempSql := ""
			switch symbol {
			case LIKE:
				tempSql = fmt.Sprintf(" and %s %s '%%%s%%'", conditions[i], symbol, conditions[i+1])
			case IN:
				tempSql = fmt.Sprintf(" and %s %s (%s)", conditions[i], symbol, conditions[i+1])
			case EQUAL:
				tempSql = fmt.Sprintf(" and %s %s '%s'", conditions[i], symbol, conditions[i+1])
			case NULL:
				tempSql = fmt.Sprintf(" and %s %s ", conditions[i], symbol)
			case NOT_NULL:
				tempSql = fmt.Sprintf(" and %s %s ", conditions[i], symbol)
			case NOT_EQUAL:
				tempSql = fmt.Sprintf(" and %s %s '%s'", conditions[i], symbol, conditions[i+1])

			}
			s.sql += tempSql
		}
	}
	return s
}

func (s *sqlGeneration) AddGroupBy(columnName string) *sqlGeneration {
	s.sql += fmt.Sprintf(" group by %s ", columnName)
	return s
}
func (s *sqlGeneration) AddOrderBy(columnName string) *sqlGeneration {
	s.sql += fmt.Sprintf(" order by %s ", columnName)
	return s
}

func (s *sqlGeneration) AddSuffixOther(text string) *sqlGeneration {
	s.sql += text
	return s
}

func (s *sqlGeneration) HarvestSql() string {
	return s.sql
}

// BatchSqlGeneration 批量sql
type BatchSqlGeneration struct {
	sqlMap map[string]*sqlGeneration
}

func InitBatchSqlGeneration() *BatchSqlGeneration {
	return &BatchSqlGeneration{
		make(map[string]*sqlGeneration),
	}
}

func (b *BatchSqlGeneration) AddSqlGeneration(name string, s *sqlGeneration) *BatchSqlGeneration {
	b.sqlMap[name] = s
	return b
}

// AddConditions 添加条件
func (b *BatchSqlGeneration) AddConditions(symbol OperatorSymbol, conditions ...string) *BatchSqlGeneration {
	for i := range b.sqlMap {
		b.sqlMap[i].AddConditions(symbol, conditions...)
	}
	return b
}

// AddSuffixOther 在尾部追加sql
func (b *BatchSqlGeneration) AddSuffixOther(text string) *BatchSqlGeneration {
	for i := range b.sqlMap {
		b.sqlMap[i].AddSuffixOther(text)
	}
	return b
}

// HarvestSql 根据名字获取对应sql
func (b *BatchSqlGeneration) HarvestSql(name string) string {
	if v, ok := b.sqlMap[name]; ok {
		return v.HarvestSql()
	}
	return ""
}

// Harvest 获取sqlGeneration对象
func (b *BatchSqlGeneration) Harvest(name string) *sqlGeneration {
	if v, ok := b.sqlMap[name]; ok {
		return v
	}
	return nil
}

// HarvestAllSql 获取所有Sql
func (b *BatchSqlGeneration) HarvestAllSql() []string {
	var sqlList []string
	for key := range b.sqlMap {
		sqlList = append(sqlList, b.sqlMap[key].HarvestSql())
	}
	return sqlList
}

func AddConditions(symbol OperatorSymbol, conditions ...string) string {
	var sql = ""
	length := len(conditions)
	for i := range conditions {
		if i%2 != 0 {
			continue
		}
		if i+1 >= length {
			break
		}
		if conditions[i+1] != "" {
			col := convert.TransSql(conditions[i])
			data := convert.TransSql(conditions[i+1])
			tempSql := ""
			switch symbol {
			case LIKE:
				tempSql = fmt.Sprintf(" and %s %s '%%%s%%'", col, symbol, data)
			case IN:
				tempSql = fmt.Sprintf(" and %s %s (%s)", col, symbol, data)
			case EQUAL, NOT_EQUAL, GREATER_THEN, GREATER_THEN_EQUAL, LESS_THAN, LESS_THAN_EQUAL:
				tempSql = fmt.Sprintf(" and %s %s '%s'", col, symbol, data)
			case NULL, NOT_NULL:
				tempSql = fmt.Sprintf(" and %s %s ", col, symbol)

			}
			sql += tempSql
		}
	}
	return sql
}

// BrotherSQL 同级SQl
type BrotherSQL struct {
	ConditionText    string
	Condition        bool
	ConditionContent any
}

const outset = " and ( "

func Brother(symbol string, b ...BrotherSQL) ([]any, string) {
	var res = make([]any, 0, len(b))
	text := outset
	for i, sql := range b {

		if !sql.Condition {
			continue
		}
		if i != 0 && text != outset {
			text += " " + symbol + " "
		}
		text += sql.ConditionText
		res = append(res, sql.ConditionContent)
	}
	if text == outset {
		return nil, ""
	}
	text = text + ")"

	return res, text
}
