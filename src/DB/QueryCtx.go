package DB

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

type SkipColumnT = map[int]bool

func Select(table string, selected string, condition ...string) DBConn {
	var conditionCtx string
	ctx := fmt.Sprintf("SELECT %s FROM %s", selected, table)
	if len(condition) > 1 {
		conditionCtx = " WHERE " + strings.Join(condition, " AND ")
	} else {
		conditionCtx = " WHERE " + condition[0]
	}

	return initDb(ctx + conditionCtx)
}

func Delete(table string, condition ...string) DBConn {
	ctx := fmt.Sprintf("DELETE FROM %s WHERE ", table)
	if len(condition) == 0 {
		ctx = ctx + condition[0]
	} else {
		ctx = ctx + strings.Join(condition, " AND ")
	}
	return initDb(ctx)
}

func getFieldName(structNames []string) string {
	lowerList := []string{}
	for _, i := range structNames {
		low := strings.ToLower(i[:1]) + i[1:]
		if low != "password2" {
			lowerList = append(lowerList, low)
		}
	}

	return strings.Join(lowerList, ",")
}

func getValue(structName []string, structField []interface{}) string {
	var (
		fieldNames string
		append     = func(format string, value interface{}) {
			if fieldNames == "" {
				fieldNames = fmt.Sprintf(format, value)
			} else {
				fieldNames =
					fieldNames + "," +
						fmt.Sprintf(format, value)
			}
		}
	)
	for k, i := range structField {
		if structName[k] == "Password2" {
			continue
		}

		if i == nil || i == "" {
			append("%s", "null")
		} else if reflect.TypeOf(i).Name() == "string" {
			append(`'%s'`, i)
		} else {
			append("%v", i)
		}

	}
	return fieldNames
}

func Create[T User | Product | Selling](table string, data T) DBConn {
	structNames := structs.Names(data)[1:]
	structValue := structs.Values(data)[1:]
	ctx := fmt.Sprintf(`
		INSERT INTO 
			%s(%s) 
			VALUES (%s)`,
		table, getFieldName(structNames), getValue(structNames, structValue))

	return initDb(ctx)
}

func Query(query string) DBConn {
	return initDb(query)
}
