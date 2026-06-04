/*
Copyright © 2026 Keoagile Dinake kmdinake@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kmdinake/rapela/rapela"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rapela",
	Short: "Daily prompt to take a moment to pray 🙇🏾",
	Long: `Rapela - Daily prompt to take a moment to pray 🙇🏾

	How It Works
	By default when you open your terminal, it will display a bible verse for the day. For example: (en-kjv) Genesis 1:1 - In the beginning God created the heaven and the earth.

	Should you clear your terminal then you can retrieve the bible verse of the day by executing the following command:
	rapela or rapela verse

	To see the list of bible versions:
	rapela bible --list-versions or --lv

	To set a bible version:
	rapela bible --set=<name-of-version>
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { // PersistentPreRunE is called after flags are parsed but before the command's RunE function is called. 
		return initializeConfig(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		bibleVersionFromConfig := viper.GetString("bible-version")
		bs, err := rapela.NewBibleService(bibleVersionFromConfig)
		if err != nil {
			cobra.CheckErr(err)
		}
		verseOfTheDay := bs.GetVerseOfTheDay()
		fmt.Printf("Dumela ngwana waka! Tseya sebaka se go rapela.\nVerse Of The Day: %s\n", verseOfTheDay)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rapela.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initializeConfig(cmd *cobra.Command) error {
	// setup Viper to use environment variables
	viper.SetEnvPrefix("RAPELA")

	// allow fo nested keys in environment variables (e.g. RAPELA_BIBLE_VERSION)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	viper.AutomaticEnv()

	// handle the configuration file
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		
		// only panic if we can't get the home directory
		cobra.CheckErr(err)

		// search for a config file with the name "config" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName(".rapela")
		viper.SetConfigType("yaml")
	}

	// read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		// it's okay if the config file doesn't exist
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	// bind Cobra flags to Viper
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	// fmt.Println("Configuration initialized. Using config file: ", viper.ConfigFileUsed())  // debug-only
	return nil
}