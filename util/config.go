// 使用LUA脚本返回配置参数
package util

import (
	"path"
	"fmt"
	"github.com/yuin/gopher-lua"
)

//const configPath = `./config`
const configPath = `../config`
var configTbl *lua.LTable

func init() {
	l := lua.NewState()
	err := l.DoFile(path.Join(configPath, "config.lua"))
	if err != nil {
		fmt.Println("Config Load Error")
		panic(err)
	}
	config := l.Get(-1)

	if tbl, ok := config.(*lua.LTable);ok {
		configTbl = tbl
	}
}

func GetConfig(section, key string) string {
	val := configTbl.RawGetString(section)
	if !lua.LVIsFalse(val) {
		t := val.(*lua.LTable)
		if key == "" {
			return t.String()
		}
		if lua.LVIsFalse(t.RawGetString(key)) {
			panic("Config " + section + "|" + key + " Missing")
		}
		return t.RawGetString(key).String()
	}
	panic("Config " + section + " Missing")
}