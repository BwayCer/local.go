package protlook

import (
	"fmt"
	"reflect"
	"strings"
)

func Print(args ...interface{}) {
	for idx, leng := 0, len(args); idx < leng; idx++ {
		printMsg("", "", reflect.ValueOf(args[idx]))
	}
}

func printMsg(preSpace string, nameArgu string, valInfo reflect.Value) {
	var name string
	if nameArgu == "" {
		name = "~~~"
	} else {
		name = nameArgu
	}

	kind := valInfo.Kind().String()

	if valInfo.CanInterface() {
		switch kind {
		case "struct":
			fmt.Printf("%s: (%s(%s))\n", preSpace+name, kind, valInfo.Type())
		case "bool":
			fmt.Printf("%s: (%s) %v\n", preSpace+name, kind, valInfo)
		case "string":
			fmt.Printf("%s: (%s) %q\n", preSpace+name, kind, valInfo)
		case "array", "slice":
			fmt.Printf("%s: (%s(%s))\n"+preSpace+"  %v\n", preSpace+name, kind, valInfo.Type(), valInfo)
		case "map":
			fmt.Printf("%s: (%s(%s))\n"+preSpace+"  %v\n", preSpace+name, kind, valInfo.Type(), valInfo)
		case "func":
			fmt.Printf("%s: (%s)\n", preSpace+name, valInfo.Type())
		case "ptr":
			fmt.Printf("%s: (%s(%s))\n", preSpace+name, kind, valInfo.Type())
		default:
			// reflect.Value.String() => {
			// 	if valInfo.CanInterface() {
			// 		string(valInfo.Interface())
			// 	} else {
			// 		string(valInfo)
			// 	}
			// }
			fmt.Printf("%s: (%s(%s)) %v\n", preSpace+name, kind, valInfo.Type(), valInfo /* valInfo.Interface() */)
		}

		handleDeepParse(preSpace, valInfo)
	} else {
		var formatValueCode string
		switch kind {
		case "string":
			formatValueCode = "%q"
			fallthrough
		default:
			if formatValueCode == "" {
				formatValueCode = "%v"
			}
			fmt.Printf("%s: (inner) (%s(%s)) "+formatValueCode+"\n", preSpace+name, kind, valInfo.Type(), valInfo)
		}
	}
}

func handleDeepParse(preSpace string, valInfo reflect.Value) {
	kind := valInfo.Kind().String()
	val := valInfo.Interface()

	ynStructType := kind == "struct"
	ynBaseType := !isBaseType(val)

	if ynStructType || ynBaseType {
		typeInfo := reflect.TypeOf(val)

		if ynStructType {
			for idx, leng := 0, typeInfo.NumField(); idx < leng; idx++ {
				subTypeInfo := typeInfo.Field(idx)
				subValInfo := valInfo.Field(idx)
				printMsg(preSpace+"  ", subTypeInfo.Name, subValInfo)
			}
		}

		if ynBaseType {
			for idx, leng := 0, typeInfo.NumMethod(); idx < leng; idx++ {
				subTypeInfo := typeInfo.Method(idx)
				subValInfo := valInfo.Method(idx)
				printMsg(preSpace+"  ", subTypeInfo.Name, subValInfo)
			}
		}
	}
}

func isBaseType(target interface{}) bool {
	baseTypeTxt := "|*bool|"
	baseTypeTxt += "*int|*uint|*int8|*uint8|*int16|*uint16|"
	baseTypeTxt += "*int32|*uint32|*int64|*uint64|*float32|*float64|"
	baseTypeTxt += "*string|"
	baseTypeTxt += "*complex64|*complex128|"
	typ := reflect.TypeOf(target).String()
	if strings.Index(typ, "*") == -1 {
		typ = "*" + typ
	}
	idx := strings.Index(baseTypeTxt, "|"+typ+"|")
	return idx != -1
}
