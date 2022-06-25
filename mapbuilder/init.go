package mapbuilder

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/crackeer/gopkg/util"
	"github.com/tidwall/gjson"
)

func init() {

	gjson.AddModifier("remove", func(jsonStr, arg string) string {
		if len(arg) < 1 {
			return jsonStr
		}
		keys := strings.Split(arg, ",")
		mapData := map[string]interface{}{}
		if err := json.Unmarshal([]byte(jsonStr), &mapData); err != nil {
			return jsonStr
		}

		for _, k := range keys {
			delete(mapData, k)
		}

		bs, _ := json.Marshal(mapData)
		return string(bs)
	})
	gjson.AddModifier("listremove", func(jsonStr, arg string) string {
		if len(arg) < 1 {
			return jsonStr
		}
		keys := strings.Split(arg, ",")
		listData := []map[string]interface{}{}
		decoder := json.NewDecoder(bytes.NewReader([]byte(jsonStr)))
		decoder.UseNumber()
		if err := decoder.Decode(&listData); err != nil {
			return jsonStr
		}

		list2 := []map[string]interface{}{}
		for _, item := range listData {
			for _, k := range keys {
				delete(item, k)
			}
			list2 = append(list2, item)
		}

		bs, _ := json.Marshal(list2)
		return string(bs)
	})
	gjson.AddModifier("default", func(jsonStr, arg string) string {
		if len(jsonStr) < 1 {
			return arg
		}
		return ""
	})
	gjson.AddModifier("omitnil", func(jsonStr, arg string) string {
		if len(jsonStr) < 1 || jsonStr == "\"\"" || jsonStr == "{}" || jsonStr == "[]" {
			return nilString
		}
		return jsonStr
	})

	gjson.AddModifier("convert", func(jsonStr, arg string) string {

		container := map[string]interface{}{
			"key": jsonStr,
		}
		switch arg {
		case "bool":
			val := util.GetBooleanValFromMap(container, "key")
			byteData, _ := util.Marshal(val)
			return string(byteData)
		case "int":
			val := util.GetInt64ValFromMap(container, "key")
			byteData, _ := util.Marshal(val)
			return string(byteData)
		}
		return jsonStr
	})

	gjson.AddModifier("validnil", func(jsonStr, arg string) string {
		tmp := map[string]interface{}{}
		if err := json.Unmarshal([]byte(jsonStr), &tmp); err != nil {
			return nilString
		}

		if _, ok := tmp[arg]; !ok {
			return nilString
		}
		return jsonStr
	})
}
