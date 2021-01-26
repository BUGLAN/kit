package t

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

var BsonTimeStampCmd = &cobra.Command{
	Use:   "ts",
	Short: "timestamp get the ISO time",
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			parseTimeStamp(arg)
		}
	},
}

func parseTimeStamp(raw string) {
	ts, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		fmt.Printf("[ERROR] %s is not right format, eg: 1609434000", raw)
	}

	t := time.Unix(ts, 0)
	fmt.Printf("ts: %s, %s\n", raw, t.String())
}
