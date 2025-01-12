package gosql

import (
	_sql "database/sql"
	"errors"
	"reflect"
	"strings"
)

func (c *Connection) Select(sql string, bind ...interface{}) (rs *_sql.Rows, err error) {
	if c.tx != nil {
		rs, err = c.tx.Query(sql, bind...)
	} else {
		rs, err = c.db.Query(sql, bind...)
	}
	return
}

func Select[T any](c *Connection, sql string, bind ...interface{}) ([]T, error) {
	rs, err := c.Select(sql, bind...)
	if err != nil {
		return nil, err
	}

	columns, err := rs.Columns()
	if err != nil {
		return nil, err
	}
	columnNameIndexMap := map[string]int{}
	for i, c := range columns {
		columnNameIndexMap[c] = i
	}

	result := []T{}
	for rs.Next() {
		var t T
		var row []interface{}
		if reflect.TypeOf(t).Kind() == reflect.Map {
			row = makeScanArrayForMap(&t, len(columns))
		} else if reflect.TypeOf(t).Kind() == reflect.Struct {
			row = transformStructFieldToScanArray(&t, columnNameIndexMap)
		} else {
			return nil, errors.New("return type must be array of struct or map")
		}
		err = rs.Scan(row...)
		if reflect.TypeOf(t).Kind() == reflect.Map {
			transformScanArrayToMap(&t, row, columns)
		}
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func makeScanArrayForMap[T any](t *T, length int) []interface{} {
	row := make([]interface{}, length)
	typeOfT := reflect.TypeOf(t).Elem()
	for i := 0; i < len(row); i++ {
		row[i] = reflect.New(typeOfT.Elem()).Interface()
	}
	return row
}

func transformScanArrayToMap[T any](t *T, row []interface{}, columns []string) []interface{} {
	typeOfT := reflect.TypeOf(t).Elem()
	reflect.ValueOf(t).Elem().Set(reflect.MakeMap(typeOfT))
	valueOfT := reflect.ValueOf(t).Elem()
	for i := 0; i < len(row); i++ {
		v := *row[i].(*interface{})
		switch v.(type) {
		case []byte:
			v = string(v.([]byte))
		}
		valueOfT.SetMapIndex(reflect.ValueOf(columns[i]), reflect.ValueOf(v))
	}
	return row
}

func transformStructFieldToScanArray[T any](t *T, columnNameIndexMap map[string]int) []interface{} {
	row := make([]interface{}, len(columnNameIndexMap))
	typeOfT := reflect.TypeOf(t).Elem()
	valueOfT := reflect.ValueOf(t).Elem()
	deepTransformStructFieldToScanArray(typeOfT, valueOfT, row, columnNameIndexMap)
	for i, v := range row {
		if v == nil {
			row[i] = new(interface{})
		}
	}
	return row
}

func deepTransformStructFieldToScanArray(typeOfT reflect.Type, valueOfT reflect.Value, row []interface{}, columnNameIndexMap map[string]int) {
	for i := 0; i < typeOfT.NumField(); i++ {
		field := typeOfT.Field(i)
		colName := findNameFromDbTag(field.Tag.Get("db"))
		if colName == "" {
			if field.Type.Kind() == reflect.Struct && field.Name == field.Type.Name() {
				// embedded type
				deepTransformStructFieldToScanArray(field.Type, valueOfT.Field(i), row, columnNameIndexMap)
			} else {
				continue
			}
		}
		colIndex, ok := columnNameIndexMap[colName]
		if !ok {
			continue
		}
		if field.Type.Kind() == reflect.Pointer {
			valueOfT.Field(i).Set(reflect.New(field.Type.Elem()))
		} else {
			switch valueOfT.Field(i).Kind() {
			case reflect.Map, reflect.Slice, reflect.Interface:
				if valueOfT.Field(i).IsNil() {
					valueOfT.Field(i).Set(reflect.New(field.Type).Elem())
				}
			}
		}
		row[colIndex] = valueOfT.Field(i).Addr().Interface()
	}
}

func findNameFromDbTag(tagContent string) string {
	pos := strings.Index(tagContent, ",")
	if pos >= 0 {
		return tagContent[0:pos]
	} else {
		return tagContent
	}
}
