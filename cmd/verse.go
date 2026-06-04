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
	"fmt"

	"github.com/kmdinake/rapela/rapela"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// verseCmd represents the verse command
var verseCmd = &cobra.Command{
	Use:   "verse",
	Short: "Retrieve the bible verse of the day",
	Long: `You can retrieve the bible verse of the day by executing the following command.
	For example: 
	> rapela verse

	Which will display 
	> (en-kjv) Genesis 1:1 - In the beginning God created the heaven and the earth.`,
	Run: func(cmd *cobra.Command, args []string) {
		BibleVersionFromConfig := viper.GetString("bible-version")
		bs, err := rapela.NewBibleService(BibleVersionFromConfig)
		if err != nil {
			cobra.CheckErr(err)
		}
		verseOfTheDay := bs.GetVerseOfTheDay()
		fmt.Printf("Dumela ngwana waka! Tseya sebaka se go rapela.\nVerse Of The Day: %s\n", verseOfTheDay)
	},
}

func init() {
	rootCmd.AddCommand(verseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
