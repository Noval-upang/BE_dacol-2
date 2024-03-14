package validation

import (
	"encoding/json"
	"strings"

	"github.com/Hy-Iam-Noval/dacol-2/src/DB"
)

func imageValidate(f Field) bool {
	field := f.Field().String()
	if field == "" {
		return true
	}

	data := []DB.JSONImg{}
	json.Unmarshal([]byte(field), &data)

	pass := true
	for _, i := range data {
		if !pass {
			break
		}
		pass = checkExtFile(i.Name) && checkSizeFile(i.Size)
	}
	return pass
}

func checkSizeFile(field int) bool {
	// 1 mb eq 1024000
	return field > 2*1024000
}

func checkExtFile(field string) bool {
	spared := strings.Split(field, ".")
	extFile := spared[len(spared)-1]
	switch extFile {
	case "png":
	case "jpg":
	case "img":
		return true
	}
	return false
}
