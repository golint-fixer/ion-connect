package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv()
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	fi, err := os.Open(usr.HomeDir + "/.ionchannel/credentials")
	if err != nil {
		panic(err)
	}
	viper.SetDefault("endpoint_url", "https://api.ionchannel.io/")
	viper.SetConfigType("props")
	viper.SetEnvPrefix("ionchannel")
	err = viper.ReadConfig(fi) // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(err)
	}

}

//RootCmd - Root command container for ion-connect
var RootCmd = &cobra.Command{
	Use:   "ion-connect",
	Short: "Ion Connect is awesome!",
	Long:  `Ion connect is awesome with more words`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside rootCmd Run with args: %v\n", args)
	},
}

// Execute runs the command defined for the root
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
