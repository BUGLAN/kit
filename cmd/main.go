package main

import (
	"fmt"
	"github.com/BUGLAN/kit/cmd/b"
	"github.com/spf13/cobra"
	"log"
)

var cmd = &cobra.Command{
	Use:   "kit",
	Short: "Kit is a command-line application for coding",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the kit version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kit v0.1.4")
	},
}

func init() {
	cmd.AddCommand(versionCmd)
	cmd.AddCommand(b.BsonParseCmd)
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Panic(err)
	}
}
