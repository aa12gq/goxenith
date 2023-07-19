package cmd

import (
	"github.com/spf13/cobra"
	"goxenith/pkg/helpers"
	"os"
)

// Env 存储全局选项 --env 的值
var Env string

func RegisterGlobalFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVarP(&Env, "env", "e", "", "load .env file, example: --env=testing will use .env.testing file")
}

func RegisterDefaultCmd(rootCmd *cobra.Command, subCmd *cobra.Command) {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	firstArg := helpers.FirstElement(os.Args[1:])
	if err == nil && cmd.Use == rootCmd.Use && firstArg != "-h" && firstArg != "--help" {
		args := append([]string{subCmd.Use}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
}

// CobraCommandUsageTemplate Cobra 命令使用帮助模板
func CobraCommandUsageTemplate() string {
	return `{{if .Version}}版本:
{{.CommandPath}} {{.Version}} {{end}}{{if .Annotations.buildDateTime}}构建时间 {{.Annotations.buildDateTime}}
{{end}}{{if .Annotations.copyright}}
版权:
{{.Annotations.copyright}}
{{end}}{{if or .Runnable .HasAvailableSubCommands}}
运行方式:{{end}}{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
{{.CommandPath}} [子命令 [子命令参数]]{{end}}{{if gt (len .Aliases) 0}}
别名:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

示例:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

可用子命令:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

其他命令:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

参数:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

全局参数:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

其他帮助主题:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

使用 "{{.CommandPath}} [子命令] --help" 查看命令的详细帮助说明.{{end}}
`
}

// CobraHelpCommand 自定义的cobra帮助命令
func CobraHelpCommand(rootCmd *cobra.Command) *cobra.Command {
	helpCmd := &cobra.Command{
		Use:   "help",
		Short: "查看一个命令的帮助信息",
		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				c.Printf("未知命令 %#q\n", args)
				cobra.CheckErr(c.Root().Usage())
			} else {
				cmd.InitDefaultHelpFlag()
				cmd.InitDefaultVersionFlag()
				cobra.CheckErr(cmd.Help())
			}
		},
	}
	helpCmd.Long = "查看一个命令的帮助信息. 格式: " + rootCmd.Name() + " [子命令 [子命令]...], 如: " + rootCmd.Name() + " " + helpCmd.Name() + " start serve management_server"
	rootCmd.PersistentFlags().BoolP("help", "h", false, "查看命令的帮助信息")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "查看版本信息")
	return helpCmd
}
