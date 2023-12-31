package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"goxenith/app/cmd"
	"goxenith/app/cmd/make"
	"goxenith/app/cmd/migrate"
	"goxenith/app/cmd/serve"
	_ "goxenith/app/models/ent/runtime"
	"goxenith/bootstrap"
	btsConig "goxenith/config"
	"goxenith/pkg/config"
	"goxenith/pkg/console"
	"os"
)

func init() {
	btsConig.Initialize()
}

func main() {

	var rootCmd = &cobra.Command{
		Use:  "goxenith",
		Long: `默认运行 "serve". 使用 "-h" 参数查看所有可用命令.`,
		PersistentPreRun: func(command *cobra.Command, args []string) {
			config.InitConfig(cmd.Env)
			bootstrap.SetupLogger()
			bootstrap.SetupRedis()
			bootstrap.SetupDB()
			bootstrap.SetupCache()
		},
	}
	var CmdServe = &cobra.Command{
		Use:   "serve",
		Short: "服务启动",
		Long:  `启动服务。提供服务及相关接口`,
		RunE:  serve.RunWeb,
		Args:  cobra.NoArgs,
	}
	var CmdMigrate = &cobra.Command{
		Use:              "migrate",
		Short:            "数据迁移",
		Long:             `数据迁移。将models下的schema同步到数据库，并进行相关数据初始化`,
		Run:              migrate.Migrate,
		Args:             cobra.MaximumNArgs(1),
		PersistentPreRun: func(*cobra.Command, []string) {},
	}
	rootCmd.AddCommand(
		CmdServe,
		CmdMigrate,
		make.CmdMake,
		cmd.CmdCache,
	)
	{
		CmdMigrate.PersistentFlags().UintP("timeout", "t", 0, "迁移执行超时时间，单位：秒。大于等于0的整数，等于0时，永不超时。")
		CmdMigrate.PersistentFlags().BoolP("verbose", "V", false, "显示详细日志，如：打印sql日志等。")
		CmdMigrate.PersistentFlags().Bool("drop-column", false, "是否删除原有字段")
		CmdMigrate.PersistentFlags().Bool("drop-index", false, "是否删除原有索引")
		CmdMigrate.PersistentFlags().Bool("foreign-key", true, "是否创建外键")
	}
	cmd.RegisterDefaultCmd(rootCmd, CmdServe)
	cmd.RegisterGlobalFlags(rootCmd)
	rootCmd.SetHelpCommand(cmd.CobraHelpCommand(rootCmd))
	rootCmd.SetUsageTemplate(cmd.CobraCommandUsageTemplate())
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
