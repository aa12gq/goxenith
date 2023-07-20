package make

import (
	"github.com/spf13/cobra"
	"os"
)

var CmdMakeSchema = &cobra.Command{
	Use:   "schema",
	Short: "Crate scheme file, example: make schema user",
	Run:   runMakeSchema,
	Args:  cobra.ExactArgs(1),
}

func runMakeSchema(cmd *cobra.Command, args []string) {
	model := makeModelFromString(args[0])

	// 保证 "app/models/schema/" 目录存在
	dir := "app/models/schema/"
	os.MkdirAll(dir, os.ModePerm)
	createFileFromStub(dir+model.PackageName+".go", "schema/schema", model)
}
