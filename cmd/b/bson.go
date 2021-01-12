package b

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/cobra"
)

var BsonParseCmd = &cobra.Command{
	Use:   "bson",
	Short: "bson get the time",
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			parseObjectIdTime(arg)
		}
	},
}

func parseObjectIdTime(hex string) {
	if !bson.IsObjectIdHex(hex) {
		fmt.Printf("hex: %s is not right ObjectId\n", hex)
		return
	}

	fmt.Printf("hex: %s, %v\n", hex, bson.ObjectIdHex(hex).Time())
}
