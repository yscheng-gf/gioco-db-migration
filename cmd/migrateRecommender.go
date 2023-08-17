package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yscheng-gf/gioco-db-migration/helper"
	"github.com/yscheng-gf/gioco-db-migration/internal/config"
	"github.com/yscheng-gf/gioco-db-migration/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMigrateRecommenderCmd(c *config.Conf) *cobra.Command {
	var operators *[]string
	cmd := &cobra.Command{
		Use:   "recommender",
		Short: "This is migrate member recommender code to recommender account",
		Long:  "This is migrate member recommender code to recommender account",
		Run: func(cmd *cobra.Command, args []string) {
			mongo := helper.InitMongo(c)

			opts := options.Find().
				SetProjection(bson.M{
					"account":        1,
					"recommend_code": 1,
					"recommender":    1,
				})
			for _, opCode := range *operators {
				members := new([]models.OpMembers)
				cur, err := mongo.Database("GIOCO_PLUS").
					Collection(opCode+"_members").
					Find(cmd.Context(), bson.M{}, opts)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer cur.Close(cmd.Context())
				cur.All(cmd.Context(), members)

				// mapping recommendMap
				recommendMap := make(map[string]string)
				for _, member := range *members {
					recommendMap[member.RecommendCode] = member.Account
				}

				for _, member := range *members {
					if member.Recommender == "" {
						continue
					}
					recommender, ok := recommendMap[member.Recommender]
					if !ok {
						continue
					}
					_, err := mongo.Database("GIOCO_PLUS").
						Collection(opCode+"_members").
						UpdateOne(
							cmd.Context(),
							bson.M{"account": member.Account},
							bson.M{
								"$set": bson.M{
									"recommender": recommender,
								},
							})
					if err != nil {
						fmt.Println(err)
						return
					}
				}
				fmt.Printf("Operator: %s done.", opCode)
			}

		},
	}

	operators = cmd.Flags().StringSliceP("operators", "p", []string{}, "operators e.g. -p iphd,xbwl")
	return cmd
}
