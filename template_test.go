package table_gen

import (
	"testing"
)

func TestGenerator_Gen(t *testing.T) {
	config := &Config{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "",
	}
	g, err := NewGenerator(config)
	if err != nil {
		t.Errorf("have err during NewGenerator, err is %s", err)
		return
	}
	genInfo := &GenInfo{
		Schema:       "elliot_test",
		Table:        "test_table",
		ExportFolder: "",
		TemplatePath: "struct_gen_test_template",
		FileName:     "test_gen.go",
		PackageName:  "table_gen",
		StructName:   "TestGenStruct",
	}
	isSuccess, err := g.Gen(genInfo)
	if err != nil {
		t.Errorf("have err during Gen file, err is %s", err)
		return
	}
	t.Logf("does it gen file successully?  %v", isSuccess)
}
