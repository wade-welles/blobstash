/*

Package luautil implements utility for gopher-lua.

*/
package luautil

import (
	"encoding/json"

	"github.com/yuin/gopher-lua"
)

// TableToMap convert a `*lua.LTable` to a `map[string]interface{}`
func TableToMap(table *lua.LTable) map[string]interface{} {
	return tomap(table, map[*lua.LTable]bool{})
}

func tomap(table *lua.LTable, visited map[*lua.LTable]bool) map[string]interface{} {
	res := map[string]interface{}{}
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		switch converted := value.(type) {
		case lua.LBool:
			res[key.String()] = converted
		case lua.LChannel:
			panic("no channel")
		case lua.LNumber:
			res[key.String()] = converted
		case *lua.LFunction:
			panic("no function")
		case *lua.LNilType:
			res[key.String()] = converted
		case *lua.LState:
			panic("no LState")
		case lua.LString:
			res[key.String()] = converted
		case *lua.LTable:
			var arr []interface{}
			obj := map[string]interface{}{}

			if visited[converted] {
				panic("nested table")
			}
			visited[converted] = true

			converted.ForEach(func(k lua.LValue, v lua.LValue) {
				_, numberKey := k.(lua.LNumber)
				// if numberKey, then convert to a slice of interface
				subtable, istable := v.(*lua.LTable)
				if numberKey {
					if istable {
						arr = append(arr, tomap(subtable, visited))
					} else {
						arr = append(arr, v)
					}
				} else {
					if istable {
						obj[k.(lua.LString).String()] = tomap(subtable, visited)
					} else {
						obj[k.(lua.LString).String()] = v
					}
				}
			})
			if len(arr) > 0 {
				res[key.String()] = arr
			} else {
				res[key.String()] = obj
			}
		}
	})
	return res
}

// Convert a Lua table to JSON
func ToJSON(table *lua.LTable) []byte {
	js, err := json.Marshal(TableToMap(table))
	if err != nil {
		panic(err)
	}
	return js
}

// Convert the JSON to a Lua object ready to be pushed
// Adapted from https://github.com/layeh/gopher-json/blob/master/util.go (Public domain)
func FromJSON(L *lua.LState, js []byte) lua.LValue {
	var res interface{}
	if err := json.Unmarshal(js, &res); err != nil {
		panic(err)
	}
	return fromJSON(L, res)
}

func fromJSON(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(fromJSON(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), fromJSON(L, item))
		}
		return tbl
	}
	return lua.LNil
}