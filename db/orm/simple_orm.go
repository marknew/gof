/**
 * Copyright 2015@press soft
 * name :
 * author : mark zhang
 * date : 2015-09-07 10:11
 * description :
 * history :
 */

package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jrsix/gof/log"
	"reflect"
	"strings"
)

var _ Orm = new(simpleOrm)

//it's a IOrm Implements for mysql
type simpleOrm struct {
	tableMap map[string]*TableMapMeta
	*sql.DB
	useTrace bool
}

func NewOrm(db *sql.DB) Orm {
	return &simpleOrm{
		DB:       db,
		tableMap: make(map[string]*TableMapMeta),
	}
}

func (this *simpleOrm) Version() string {
	return "1.0.2"
}

func (this *simpleOrm) err(err error) error {
	if this.useTrace && err != nil {
		log.Println("[ ORM][ !][ ERROR]:", err.Error())
	}
	return err
}

func (this *simpleOrm) debug(format string, args ...interface{}) {
	if this.useTrace {
		log.Printf(format+"\n", args...)
	}
}

func (this *simpleOrm) getTableMapMeta(t reflect.Type) *TableMapMeta {
	m, exists := this.tableMap[t.String()]
	if exists {
		return m
	}

	m = GetTableMapMeta(t)
	this.tableMap[t.String()] = m

	if this.useTrace {
		log.Println("[ ORM][ META]:", m)
	}

	return m
}

func (this *simpleOrm) getTableName(t reflect.Type) string {
	//todo: 用int做键
	v, exists := this.tableMap[t.String()]
	if exists {
		return v.TableName
	}
	return t.Name()
}

//if not defined primary key.the first key will as primary key
func (this *simpleOrm) getPKName(t reflect.Type) (pkName string, pkIsAuto bool) {
	v, exists := this.tableMap[t.String()]
	if exists {
		return v.PkFieldName, v.PkIsAuto
	}
	return GetPKName(t)
}

func (this *simpleOrm) unionField(meta *TableMapMeta, v string) string {
	if len(meta.TableName) != 0 {
		return meta.TableName + "." + v
	}
	return v
}

func (this *simpleOrm) SetTrace(b bool) {
	this.useTrace = b
}

//create a fixed table map
func (this *simpleOrm) CreateTableMap(v interface{}, tableName string) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	meta := this.getTableMapMeta(t)
	meta.TableName = tableName
	this.tableMap[t.String()] = meta
}

func (this *simpleOrm) Get(primaryVal interface{}, entity interface{}) error {
	var sql string
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return this.err(errors.New("Unaddressable of entity ,it must be a ptr"))
	}
	val = val.Elem()

	/* build sql */
	meta := this.getTableMapMeta(t)
	fieldLen = len(meta.FieldsIndex)
	fieldArr := make([]string, fieldLen)
	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = v
		scanVal[i] = &rawBytes[i]
	}

	sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s=?",
		strings.Join(fieldArr, ","),
		meta.TableName,
		meta.PkFieldName,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%+v", sql, primaryVal))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		return this.err(errors.New(err.Error() + "\n[ SQL]:" + sql))
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRow(primaryVal)
	err = row.Scan(scanVal...)
	if err != nil {
		return this.err(err)
	}
	for i := 0; i < fieldLen; i++ {
		field := val.Field(i)
		SetField(field, rawBytes[i])
	}
	return nil
}

func (this *simpleOrm) GetBy(entity interface{}, where string,
	args ...interface{}) error {

	var sql string
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return this.err(errors.New("unaddressable of entity ,it must be a ptr"))
	}

	if strings.Trim(where, "") == "" {
		return this.err(errors.New("param where can't be empty "))
	}

	val = val.Elem()

	if !val.IsValid() {
		return this.err(errors.New("not validate or not initialize."))
	}

	/* build sql */
	meta := this.getTableMapMeta(t)
	fieldLen = len(meta.FieldsIndex)
	fieldArr := make([]string, fieldLen)
	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = v
		scanVal[i] = &rawBytes[i]
	}

	sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(fieldArr, ","),
		meta.TableName,
		where,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%s - %+v", sql, where, args))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		if this.useTrace {
			log.Println("[ ORM][ ERROR]:", err.Error(), " [ SQL]:", sql)
		}
		return this.err(errors.New(err.Error() + "\n[ SQL]:" + sql))
	}
	defer stmt.Close()

	row := stmt.QueryRow(args...)
	err = row.Scan(scanVal...)

	if err != nil {
		return this.err(err)
	}

	for i := 0; i < fieldLen; i++ {
		field := val.Field(i)
		SetField(field, rawBytes[i])
	}
	return nil
}

func (this *simpleOrm) GetByQuery(entity interface{}, sql string,
	args ...interface{}) error {
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return this.err(errors.New("Unaddressable of entity ,it must be a ptr"))
	}

	val = val.Elem()

	/* build sql */
	meta := this.getTableMapMeta(t)
	fieldLen = len(meta.FieldsIndex)
	fieldArr := make([]string, fieldLen)
	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = this.unionField(meta, v)
		scanVal[i] = &rawBytes[i]
	}

	if strings.Index(sql, "*") != -1 {
		sql = strings.Replace(sql, "*", strings.Join(fieldArr, ","), 1)
	}

	if this.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s", sql))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)

	if err != nil {
		return this.err(errors.New(err.Error() + "\n[ SQL]:" + sql))
	}
	defer stmt.Close()

	row := stmt.QueryRow(args...)
	err = row.Scan(scanVal...)

	if err != nil {
		return this.err(err)
	}

	for i := 0; i < fieldLen; i++ {
		field := val.Field(i)
		SetField(field, rawBytes[i])
	}
	return nil
}

//Select more than 1 entity list
//@to : referenced queried entity list
//@entity : query condition
//@where : other condition
func (this *simpleOrm) Select(to interface{}, where string, args ...interface{}) error {
	return this.selectBy(to, where, false, args...)
}

func (this *simpleOrm) SelectByQuery(to interface{}, sql string, args ...interface{}) error {
	return this.selectBy(to, sql, true, args...)
}

//解决结构体的字段顺序、个数、大小写必须与查询语句相同
// query rows
func (this *simpleOrm) selectBy(to interface{}, sql string, fullSql bool, args ...interface{}) error {
	var fieldLen int
	var eleIsPtr bool // 元素是否为指针

	toVal := reflect.Indirect(reflect.ValueOf(to))
	toTyp := reflect.TypeOf(to).Elem()

	if toTyp.Kind() == reflect.Ptr {
		toTyp = toTyp.Elem()
	}

	if toTyp.Kind() != reflect.Slice {
		return this.err(errors.New("目标必须是切片类型"))
	}

	baseType := toTyp.Elem()
	//fmt.Println("basetype", baseType)
	if baseType.Kind() == reflect.Ptr {
		eleIsPtr = true
		baseType = baseType.Elem()
	}

	/* build sql */
	meta := this.getTableMapMeta(baseType)

	//fmt.Println("meta", meta)

	fieldLen = len(meta.FieldMapNames)

	fieldArr := make([]string, fieldLen)
	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = this.unionField(meta, v)
		scanVal[i] = &rawBytes[i]
	}
	/*
		if fullSql {
			if strings.Index(sql, "*") != -1 {
				sql = strings.Replace(sql, "*", strings.Join(fieldArr, ","), 1)
			}
		} else {
			where := sql
			if len(where) == 0 {
				sql = fmt.Sprintf("SELECT %s FROM %s",
					strings.Join(fieldArr, ","),
					meta.TableName)
			} else {
				// 此时,sql为查询条件
				sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s",
					strings.Join(fieldArr, ","),
					meta.TableName,
					where)
			}
		}
	*/
	if this.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s [ Params] - %+v", sql, args))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		return this.err(errors.New(fmt.Sprintf("%s - [ SQL]: %s- [Args]:%+v", err.Error(), sql, args)))
	}

	defer stmt.Close()
	rows, err := stmt.Query(args...)

	if err != nil {
		return this.err(errors.New(err.Error() + "\n[ SQL]:" + sql))
	}

	defer rows.Close()

	/* 用反射来对输出结果复制 */
	toArr := toVal
	results := [](map[string][]byte){} //数据切片
	results = RowsToMap(rows)
	var fieldname string
	//results = RowsToMap(rows) //读入SQL Result into map
	//http://127.0.0.1:1003/udo/mobile/api/login.html?Usr_Login=admin&usr_Password=63a9f0ea7bb98050796b649e85481845
	//http://127.0.0.1:1003/udo/mobile/api/getGrnMainList.html?dayfrom=2014-01-01&dayto=2015-09-09
	//fmt.Println(results)

	//fmt.Println("v", v)
	//fieldpos表示映射字段在结构体中的顺序
	for _, result := range results {
		e := reflect.New(baseType)
		v := e.Elem()
		for i, fieldpos := range meta.FieldsIndex {
			//fmt.Println(meta.FieldMapNames[i], result[meta.FieldMapNames[i]])
			//fmt.Println(result)
			fieldname = strings.ToLower(meta.FieldMapNames[i])
			SetField(v.Field(fieldpos), result[fieldname])
			//SetField(v.FieldByName(field), result[field])
		}

		if eleIsPtr {
			toArr = reflect.Append(toArr, e)
		} else {
			toArr = reflect.Append(toArr, v)
		}
	}

	/*for i, fi := range meta.FieldsIndex {
		SetField(v.Field(fi), rawBytes[i])
	}
	*/
	//fmt.Println("columns", results)

	/*
		for rows.Next() {
			e := reflect.New(baseType)
			v := e.Elem()

			if err = rows.Scan(scanVal...); err != nil {
				break
			}
			//fmt.Println("scanVal", scanVal)

			for i, fi := range meta.FieldsIndex {
				SetField(v.Field(fi), rawBytes[i])
			}

			fmt.Println(v)
			if eleIsPtr {
				toArr = reflect.Append(toArr, e)
			} else {
				toArr = reflect.Append(toArr, v)
			}
			fmt.Println(toArr)
		}
	*/
	toVal.Set(toArr)

	return this.err(err)
}

func (this *simpleOrm) Delete(entity interface{}, where string,
	args ...interface{}) (effect int64, err error) {
	var sql string

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	/* build sql */
	meta := this.getTableMapMeta(t)

	if where == "" {
		return 0, errors.New("Unknown condition")
	}

	sql = fmt.Sprintf("DELETE FROM %s WHERE %s",
		meta.TableName,
		where,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%s", sql, args))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		return 0, this.err(errors.New(fmt.Sprintf("[ ORM][ ERROR]:%s [ SQL]:%s", err.Error(), sql)))

	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	var rowNum int64 = 0
	if err == nil {
		rowNum, err = result.RowsAffected()
	}
	if err != nil {
		return rowNum, this.err(errors.New(err.Error() + "\n[ SQL]" + sql))
	}
	return rowNum, nil
}

func (this *simpleOrm) DeleteByPk(entity interface{}, primary interface{}) (err error) {
	var sql string
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	/* build sql */
	meta := this.getTableMapMeta(t)

	sql = fmt.Sprintf("DELETE FROM %s WHERE %s=?",
		meta.TableName,
		meta.PkFieldName,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%s", sql, primary))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		return this.err(errors.New(fmt.Sprintf("[ ORM][ ERROR]:%s \n [ SQL]:%s", err.Error(), sql)))

	}
	defer stmt.Close()

	_, err = stmt.Exec(primary)
	if err != nil {
		return this.err(errors.New(err.Error() + "\n[ SQL]" + sql))
	}
	return nil
}

func (this *simpleOrm) Save(primaryKey interface{}, entity interface{}) (rows int64, lastInsertId int64, err error) {
	var sql string
	//var condition string
	//var fieldLen int

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	val := reflect.Indirect(reflect.ValueOf(entity))

	/* build sql */
	meta := this.getTableMapMeta(t)
	//fieldLen = len(meta.FieldNames)
	params, fieldArr := ItrFieldForSave(meta, &val, false)

	//insert
	if primaryKey == nil {
		var pArr = make([]string, len(fieldArr))
		for i := range pArr {
			pArr[i] = "?"
		}

		sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", meta.TableName,
			strings.Join(fieldArr, ","),
			strings.Join(pArr, ","),
		)

		if this.useTrace {
			log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%+v", sql, params))
		}

		/* query */
		stmt, err := this.DB.Prepare(sql)
		if err != nil {
			return 0, 0, this.err(errors.New("[ ORM][ ERROR]:" + err.Error() + "\n[ SQL]" + sql))
		}
		defer stmt.Close()

		result, err := stmt.Exec(params...)
		var rowNum int64 = 0
		var lastInsertId int64 = 0
		if err == nil {
			rowNum, err = result.RowsAffected()
			lastInsertId, _ = result.LastInsertId()
			return rowNum, lastInsertId, err
		}
		return rowNum, lastInsertId, this.err(errors.New(err.Error() + "\n[ SQL]" + sql))
	} else {
		//update model

		var setCond string

		for i, k := range fieldArr {
			if i == 0 {
				setCond = fmt.Sprintf("%s = ?", k)
			} else {
				setCond = fmt.Sprintf("%s,%s = ?", setCond, k)
			}
		}

		sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", meta.TableName,
			setCond,
			meta.PkFieldName,
		)

		/* query */
		stmt, err := this.DB.Prepare(sql)
		if err != nil {
			return 0, 0, this.err(errors.New("[ ORM][ ERROR]:" + err.Error() + " [ SQL]" + sql))
		}
		defer stmt.Close()

		params = append(params, primaryKey)

		if this.useTrace {
			log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%+v", sql, params))
		}

		result, err := stmt.Exec(params...)
		var rowNum int64 = 0
		if err == nil {
			rowNum, err = result.RowsAffected()
			return rowNum, 0, err
		}
		return rowNum, 0, this.err(errors.New(err.Error() + "\n[ SQL]" + sql))
	}
}
