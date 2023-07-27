package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"github.com/yscheng-gf/gioco-db-migration/helper"
	"github.com/yscheng-gf/gioco-db-migration/internal/config"
	"github.com/yscheng-gf/gioco-db-migration/models"
	"go.mongodb.org/mongo-driver/bson"
)

func NewTransLogPrecisionMigrateCmd(c *config.Conf) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "trans-log-precision",
		Short: "This is migrate postgresDB column balance digital",
		Long:  "This is migrate postgresDB column balance digital",
		Run: func(cmd *cobra.Command, args []string) {
			digital, _ := cmd.Flags().GetInt("digital")
			mongoCli := helper.InitMongo(c)
			pgDb := helper.InitPgDb(c)

			cursor, err := mongoCli.Database("GIOCO_PLUS").Collection("sys_operators").Find(context.TODO(), bson.M{"status": "1"})
			if err != nil {
				fmt.Println(err)
				return
			}
			defer cursor.Close(context.TODO())

			ops := new([]models.Operator)
			err = cursor.All(context.TODO(), ops)
			if err != nil {
				fmt.Println(err)
				return
			}

			p := mpb.New(mpb.WithWidth(64))
			bar := p.AddBar(int64(len(*ops)),
				mpb.PrependDecorators(
					decor.Name("Migrate", decor.WCSyncSpaceR),
					decor.Spinner([]string{}, decor.WCSyncSpaceR),
					decor.Percentage(decor.WCSyncSpaceR),
				),
				mpb.AppendDecorators(
					decor.OnComplete(
						decor.EwmaETA(decor.ET_STYLE_GO, 30, decor.WCSyncWidth), "done",
					),
				),
			)

			for _, op := range *ops {
				start := time.Now()
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN balance TYPE numeric(24, %d);", op.Code+"_member_wallets", digital))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN before_balance TYPE numeric(24, %d);", op.Code+"_member_transactions", digital))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN amount TYPE numeric(24, %d);", op.Code+"_member_transactions", digital))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN after_balance TYPE numeric(24, %d);", op.Code+"_member_transactions", digital))

				bar.EwmaIncrement(time.Since(start))
			}
			p.Wait()
		},
	}

	migrateCmd.Flags().IntP("digital", "d", 8, "This flag for setting Nnd decimal place.")

	return migrateCmd
}
