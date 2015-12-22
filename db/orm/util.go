/**
 * Copyright 2013 @press soft.
 * name :
 * author : mark zhang
 * date : 2013-10-22 21:43
 * description :
 * history :
 */

package orm

import (
	"database/sql"
	//"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 获取表元数据
func GetTableMapMeta(t reflect.Type) *TableMapMeta {
	ixs, maps := GetFields(t)
	pkName, pkIsAuto := GetPKName(t)
	m := &TableMapMeta{
		TableName:     t.Name(),
		PkFieldName:   pkName,
		PkIsAuto:      pkIsAuto,
		FieldsIndex:   ixs,
		FieldMapNames: maps,
	}
	return m
}

//if not defined primary key.the first key will as primary key
func GetPKName(t reflect.Type) (pkName string, pkIsAuto bool) {
	var ti int = t.NumField()

	ffc := func(f reflect.StructField) (string, bool) {
		if f.Tag != "" {
			var isAuto bool
			var fieldName string

			if ia := f.Tag.Get("auto"); ia == "yes" || ia == "1" {
				isAuto = true
			}

			if fieldName = f.Tag.Get("db"); fieldName != "" {
				return fieldName, isAuto
			}
			return f.Name, isAuto
		}
		return f.Name, false
	}

	for i := 0; i < ti; i++ {
		f := t.Field(i)
		if f.Tag != "" {
			pk := f.Tag.Get("pk")
			if pk == "1" || pk == "yes" {
				return ffc(f)
			}
		}
	}

	return ffc(t.Field(0))
}

// 获取实体的字段
func GetFields(t reflect.Type) (ixs []int, mapNames []string) {
	ixs = []int{}
	mapNames = []string{}

	fNum := t.NumField()
	var fmn string

	for i := 0; i < fNum; i++ {
		f := t.Field(i)
		if f.Tag != "" {
			fmn = f.Tag.Get("db")
			if fmn == "-" || fmn == "_" || len(fmn) == 0 {
				continue
			}
		}
		if fmn == "" {
			fmn = f.Name
		}
		mapNames = append(mapNames, fmn)
		ixs = append(ixs, i)
		fmn = ""
	}

	return ixs, mapNames
}

func SetField(field reflect.Value, d []byte) {
	//fmt.Println(field.Kind().String())
	//fmt.Println(string(d))
	if field.IsValid() {

		//fmt.Println(reflect.ValueOf(string(d)))
		switch field.Type().Kind() {
		// case template.HTML:
		// 	field.Set(reflect.ValueOf(template.HTML(d)))
		case reflect.String:
			field.Set(reflect.ValueOf(string(d)))

		case reflect.Int:
			val, err := strconv.ParseInt(string(d), 10, 0)
			if err == nil {
				field.Set(reflect.ValueOf(int(val)))
			}
		case reflect.Int32:
			val, err := strconv.ParseInt(string(d), 10, 32)
			if err == nil {
				field.Set(reflect.ValueOf(val))
			}
		case reflect.Int64:
			val, err := strconv.ParseInt(string(d), 10, 64)
			if err == nil {
				field.Set(reflect.ValueOf(val))
			}

		case reflect.Float32:
			val, err := strconv.ParseFloat(string(d), 32)
			if err == nil {
				field.Set(reflect.ValueOf(float32(val)))
			}

		case reflect.Float64:
			val, err := strconv.ParseFloat(string(d), 64)
			if err == nil {
				field.Set(reflect.ValueOf(val))
			}

		case reflect.Bool:
			strVal := string(d)
			val := strings.ToLower(strVal) == "true" || strVal == "1"
			field.Set(reflect.ValueOf(val))
			return

		case reflect.Slice:
			field.Set(reflect.ValueOf([]byte(d)))
			return

			//接口类型
		case reflect.Struct:
			//fmt.Println(reflect.TypeOf(time.Now()), field.Type())
			if reflect.TypeOf(time.Now()) == field.Type() {
				t, err := time.Parse("2006-01-02 15:04:05", string(d))
				if err == nil {
					field.Set(reflect.ValueOf(t.Local()))
				}
			}
			return
		}

	}
}

//遍历所有列，并得到参数及列名
func ItrFieldForSave(meta *TableMapMeta, val *reflect.Value, includePk bool) (params []interface{}, fieldArr []string) {
	var isSet bool
	for i, k := range meta.FieldMapNames {

		if !includePk && meta.PkIsAuto &&
			meta.FieldMapNames[i] == meta.PkFieldName {
			continue
		}

		field := val.Field(i)
		isSet = false

		switch field.Type().Kind() {
		case reflect.String:
			if field.String() != "" {
				isSet = true
				if val.Kind() == reflect.Ptr {
					params = append(params, field.String())
				} else {
					params = append(params, field.String())
				}
			}
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			//if field.Int() != 0 {
			isSet = true
			params = append(params, field.Int())
			//}
		case reflect.Float32, reflect.Float64:
			//if v := field.Float(); v != 0 {
			isSet = true
			params = append(params, field.Float())
			//}

		case reflect.Bool:
			strVal := field.String()
			val := strings.ToLower(strVal) == "true" || strVal == "1"
			field.Set(reflect.ValueOf(val))
			break

		case reflect.Struct:
			v := field.Interface()
			switch v.(type) {
			case time.Time:
				if v.(time.Time).Year() > 1 {
					isSet = true
					params = append(params, v.(time.Time))
				}
			}
		}

		if isSet {
			fieldArr = append(fieldArr, k)
		}
	}
	return params, fieldArr
}

func ItrField(meta *TableMapMeta, val *reflect.Value, includePk bool) (params []interface{}, fieldArr []string) {
	var isSet bool
	for i, k := range meta.FieldMapNames {

		if !includePk && meta.PkIsAuto &&
			meta.FieldMapNames[i] == meta.PkFieldName {
			continue
		}

		field := val.Field(i)
		isSet = false

		switch field.Type().Kind() {
		case reflect.String:
			if field.String() != "" {
				isSet = true
				if val.Kind() == reflect.Ptr {
					params = append(params, field.String())
				} else {
					params = append(params, field.String())
				}
			}
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				isSet = true
				params = append(params, field.Int())
			}
		case reflect.Float32, reflect.Float64:
			if v := field.Float(); v != 0 {
				isSet = true
				params = append(params, field.Float())
			}

			//		case reflect.Bool:
			//			val := strings.ToLower(strVal) == "true" || strVal == "1"
			//			field.Set(reflect.ValueOf(val))
			//			break

		case reflect.Struct:
			v := field.Interface()
			switch v.(type) {
			case time.Time:
				if v.(time.Time).Year() > 1 {
					isSet = true
					params = append(params, v.(time.Time))
				}
			}
		}

		if isSet {
			fieldArr = append(fieldArr, k)
		}
	}
	return params, fieldArr
}

//转换为字典数组
//参考：http://my.oschina.net/nowayout/blog/143278
func RowsToMarshalMap(rows *sql.Rows) (rowsMap []map[string]interface{}) {
	rowsMap = [](map[string]interface{}){} //数据切片
	var tmpInt int = 0                     //序列
	columns, _ := rows.Columns()           //列名

	//定义数组，数组的类型为[]byte
	var values []interface{} = make([]interface{}, len(columns))
	var rawBytes [][]byte = make([][]byte, len(values))

	for v := range values {
		values[v] = &rawBytes[v]
	}

	for rows.Next() {
		rows.Scan(values...)

		if len(rowsMap) == tmpInt {
			rowsMap = append(rowsMap, make(map[string]interface{}))
		}

		for i, v := range columns {
			rowsMap[tmpInt][v] = string(rawBytes[i])
		}

		tmpInt++
	}
	return rowsMap
}

//http://127.0.0.1:1003/udo/mobile/api/getAccountsReceiveDetail.html?custNo=A001&dayFrom=2014-09-08&dayTo=2015-09-08
func RowsToMap(rows *sql.Rows) (rowsMap []map[string][]byte) {
	rowsMap = [](map[string][]byte){} //数据切片
	var colIndex int = 0              //序列
	columns, _ := rows.Columns()      //列名
	//fmt.Println(columns)
	//定义数组，数组的类型为[]byte
	var values []interface{} = make([]interface{}, len(columns))
	var rawBytes [][]byte = make([][]byte, len(values))

	for v := range values {
		values[v] = &rawBytes[v]
	}

	for rows.Next() {
		rows.Scan(values...)
		//for index, val := range values {
		//	fmt.Println(index, val)
		//}
		if len(rowsMap) == colIndex {
			rowsMap = append(rowsMap, make(map[string][]byte))
		}

		for i, v := range columns {
			rowsMap[colIndex][strings.ToLower(v)] = rawBytes[i]
			//fmt.Println(v + "===>" + string(rawBytes[i]))
		}
		colIndex++
	}

	//fmt.Println(rowsMap)
	return rowsMap
}

func RowToMap(rows *sql.Rows) map[string][]byte {
	rowMap := make(map[string][]byte)
	columns, _ := rows.Columns() //列名
	if rows.Next() {
		row := rows
		//数据
		//定义数组，数组的类型为[]byte
		var values []interface{} = make([]interface{}, len(columns))
		var rawBytes [][]byte = make([][]byte, len(values))
		for v := range values {
			values[v] = &rawBytes[v]
		}
		row.Scan(values...)
		for i, v := range columns {
			rowMap[v] = rawBytes[i]
			//fmt.Println(v + "===>" + string(rawBytes[i]))
		}
	}
	return rowMap
}
