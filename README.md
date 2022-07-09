## GenStructByTable

​	GenStructByTable是一个简单的代码生成工具，具体工作是讲mysql中的表生成成go中的结构体。以下是简单的demo：

在生成代码之前你可能需要配置一个简单的template文件,以此来代表你需要的生成的代码的模版，详情可以参考golang text/template的语法。

````
package {{.PackageName}}


type {{.StructName}} struct {
    {{- range $i, $v := .Meta }}
        {{$v.CamelName}} {{$v.DataTypeInGo}}
    {{- end }}
}
````

写好template之后可以参考以下代码生成代码：

````go
func TestGenerator_Gen(t *testing.T) {
	config := &Config{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "", // 我才不告诉你我的密码呢。
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
````

