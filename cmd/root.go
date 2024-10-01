package cmd

import (
	"log"

	"leetstalker/internal/config"
	"leetstalker/internal/leetcode"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile config.Configuration

	rootCmd = &cobra.Command{
		Use: "leetstalker",
		Short: "Fetch profile data from LeetCode",
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				leetcode.FetchMultiple(cfgFile.Users)
			} else {
				leetcode.FetchMultiple(args)
			}
		},
	}
)

type UserData struct {
	Username	string
	Data		[]byte
	Error		error
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	cobra.OnInitialize(func () {
		if err := config.InitConfig(&cfgFile); err != nil {
			log.Fatalln(err)
		}
	})
}
