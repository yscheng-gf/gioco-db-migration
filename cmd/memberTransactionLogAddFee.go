package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"github.com/yscheng-gf/gioco-db-migration/helper"
	"github.com/yscheng-gf/gioco-db-migration/internal/config"
)

func NewTransLogFeeMigrateCmd(c *config.Conf) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trans-log-add-fee",
		Short: "This is migrate member transaction log add fee columns",
		Long:  "This is migrate member transaction log add fee columns",
		Run: func(cmd *cobra.Command, args []string) {
			mongo := helper.InitMongo(c)
			pgDb := helper.InitPgDb(c)

			// 取出所有 operator
			ops, err := helper.FetchAllOperator(mongo)
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
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN fee_open character varying default 'N'", op.Code+"_member_transactions"))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN fee_rate numeric(24, 2) ", op.Code+"_member_transactions"))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN fixed_fee numeric(24, 2) ", op.Code+"_member_transactions"))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN fee numeric(24, 2) ", op.Code+"_member_transactions"))
				bar.EwmaIncrement(time.Since(start))
			}
			p.Wait()
		},
	}

	return cmd
}
