package main

import (
	"fmt"
	"github.com/BUGLAN/kit/cmd/b"
	"github.com/BUGLAN/kit/cmd/t"
	"github.com/spf13/cobra"
	"log"
)

var cmd = &cobra.Command{
	Use:   "cmd",
	Short: "Cmd is a command-line application for coding",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the cmd version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cmd v0.1.6")
	},
}

func init() {
	cmd.AddCommand(
		versionCmd,
		b.BsonParseCmd,
		t.BsonTimeStampCmd,
	)
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Panic(err)
	}
}
