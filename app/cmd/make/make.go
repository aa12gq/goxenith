package make

import (
	"embed"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"goxenith/pkg/console"
	"goxenith/pkg/file"
	"goxenith/pkg/str"
	"strings"
)

type Model struct {
	TableName          string
	StructName         string
	StructNamePlural   string
	VariableName       string
	VariableNamePlural string
	PackageName        string
}

//go:embed stubs
var stubsFS embed.FS

var CmdMake = &cobra.Command{
	Use:   "make",
	Short: "Generate file and code",
}

func init() {
	// 注册 make 的子命令
	CmdMake.AddCommand(
		CmdMakeCMD,
	)
}

func makeModelFromString(name string) Model {
	model := Model{}
	model.StructName = str.Singular(strcase.ToCamel(name))
	model.StructNamePlural = str.Plural(model.StructName)
	model.TableName = str.Snake(model.StructNamePlural)
	model.VariableName = str.LowerCamel(model.StructName)
	model.PackageName = str.Snake(model.StructName)
	model.VariableNamePlural = str.LowerCamel(model.StructNamePlural)
	return model
}

func createFileFromStub(filePath string, stubName string, model Model, variables ...interface{}) {

	// 实现最后一个参数可选
	replaces := make(map[string]string)
	if len(variables) > 0 {
		replaces = variables[0].(map[string]string)
	}

	// 目标文件已存在
	if file.Exists(filePath) {
		console.Exit(filePath + " already exists!")
	}

	// 读取 stub 模板文件
	modelData, err := stubsFS.ReadFile("stubs/" + stubName + ".stub")
	if err != nil {
		console.Exit(err.Error())
	}
	modelStub := string(modelData)

	// 添加默认的替换变量
	replaces["{{VariableName}}"] = model.VariableName
	replaces["{{VariableNamePlural}}"] = model.VariableNamePlural
	replaces["{{StructName}}"] = model.StructName
	replaces["{{StructNamePlural}}"] = model.StructNamePlural
	replaces["{{PackageName}}"] = model.PackageName
	replaces["{{TableName}}"] = model.TableName

	// 对模板内容做变量替换
	for search, replace := range replaces {
		modelStub = strings.ReplaceAll(modelStub, search, replace)
	}

	// 存储到目标文件中
	err = file.Put([]byte(modelStub), filePath)
	if err != nil {
		console.Exit(err.Error())
	}

	// 提示成功
	console.Success(fmt.Sprintf("[%s] created.", filePath))
}
