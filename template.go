package table_gen

import (
	"os"
	"text/template"
)

var TypeConvertConfig = map[string]string{
	"int":       "int",
	"tinyint":   "int",
	"smallint":  "int",
	"bigint":    "int64",
	"char":      "string",
	"varchar":   "string",
	"longtext":  "string",
	"text":      "string",
	"tinytext":  "string",
	"date":      "time.Time",
	"datetime":  "time.Time",
	"timestamp": "time.Time",
	"double":    "float64",
	"float":     "float64",
}

type TemplateData struct {
	GenPath     string
	PackageName string
	StructName  string
	Meta        []*TemplateMetaData
}

type TemplateMetaData struct {
	*tableInfo
	CamelName    string
	DataTypeInGo string
}

func convertTableInfoToMeta(infos []*tableInfo) (metas []*TemplateMetaData, err error) {
	if len(infos) == 0 {
		return nil, err
	}
	for _, info := range infos {
		meta := &TemplateMetaData{
			tableInfo:    info,
			DataTypeInGo: TypeConvertConfig[info.DataType],
			CamelName:    camelString(info.ColumnName),
		}
		metas = append(metas, meta)
	}
	return metas, nil
}

func genCodeByTemplate(genPath string, templatePath string, data *TemplateData) (isSuccess bool, err error) {
	file, err := os.OpenFile(genPath, os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	if err != nil {
		return false, err
	}
	tem, _ := template.ParseFiles(templatePath)
	err = tem.Execute(file, data)
	if err != nil {
		return false, err
	}
	return true, nil
}

func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
