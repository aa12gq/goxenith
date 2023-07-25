package migrate

import (
	"context"
	"github.com/spf13/cobra"
	cmd2 "goxenith/app/cmd"
	"goxenith/app/models"
	"goxenith/bootstrap"
	"goxenith/dao"
	"goxenith/pkg/config"
	"time"
)

func RunUp(cmd *cobra.Command, args []string) {
	var err error
	var mCtx context.Context

	config.InitConfig(cmd2.Env)
	bootstrap.SetupLogger()

	timeout, _ := cmd.Flags().GetUint("timeout")
	var cancel context.CancelFunc
	if timeout > 0 {
		mCtx, cancel = context.WithTimeout(context.TODO(), time.Second*time.Duration(timeout))
		defer cancel()
	} else {
		mCtx = context.Background()
	}
	debug, _ := cmd.Flags().GetBool("verbose")
	dropColumn, _ := cmd.Flags().GetBool("drop-column")
	dropIndex, _ := cmd.Flags().GetBool("drop-index")
	foreignKey, _ := cmd.Flags().GetBool("create-foreign-key")

	daoIn, cleanFun, err := dao.NewDAO()
	if err != nil {
		panic(err)
	} else {
		defer cleanFun()
	}

	if err := models.Migrate(mCtx, daoIn, &models.MigrateOptions{
		Debug:            debug,
		DropColumn:       dropColumn,
		DropIndex:        dropIndex,
		CreateForeignKey: foreignKey,
	}); err != nil {
		panic(err)
	}
}
