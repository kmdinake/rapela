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

var ListVersions bool
var NewBibleVersionId string

// bibleCmd represents the bible command
var bibleCmd = &cobra.Command{
	Use:   "bible",
	Short: "To see the current bible version",
	Long: `You can retrieve the bible version currently configured by executing the following command.
	For example: 
	> rapela bible
	
	Which will display the current bible version as follows:
	> en-kjv: King James Version of the Holy Bible (KJV)
	
	You can see the list of bible versions by executing the following command.
	For example: 
	> rapela bible --list-versions or -l
	
	Which will display bible versions as follows:
	> - en-engbrent: Brenton English Septuagint (engbrent)
	  - en-oke: Targum Onkelos Etheridge (OKE)
	  - sr-Latn-srp1865: Sveta Biblija (SRP1865)
	  - en-gnv: Geneva Bible (GNV)
	  - spm-akg-mkac: Mak Osɨrisira Akaman Aghuuŋ ko Iesusɨm Mbɨsevisir Gumasi (Akg-MkAc)
	  - mgw-matumbi: Injili ya Yesu (Matumbi)
	  - ...
	`,
	Run: func(cmd *cobra.Command, args []string) {
		BibleVersionFromConfig := viper.GetString("bible-version")
		bs, err := rapela.NewBibleService(BibleVersionFromConfig)
		if err != nil {
			cobra.CheckErr(err)
		}

		if ListVersions == true {
			versions, err := bs.GetBibleVersions()
			if err != nil {
				cobra.CheckErr(err)
			}
			fmt.Println("Available Bible Versions:")
			for _, version := range versions {
				fmt.Printf("- %s\n", version)
			}
			return
		}

		if NewBibleVersionId != "" {
			fmt.Printf("Setting Bible version to: %s\n", NewBibleVersionId)
			versions, err := bs.GetBibleVersions()
			if err != nil {
				cobra.CheckErr(err)
			}

			var newVersion rapela.BibleVersion
			for _, v := range versions {
				if v.Id == NewBibleVersionId {
					newVersion = v
					break
				}
			}
			if newVersion.Id == "" {
				cobra.CheckErr(fmt.Errorf("Invalid bible version: %s", NewBibleVersionId))
				return
			}
			bs.SetBibleVersion(newVersion)

			version, err := bs.GetBibleVersion()
			if err != nil {
				panic(err)
			}
			fmt.Printf("Current Bible Version after setting: %s\n", version)
			return
		}

		version, err := bs.GetBibleVersion()
		if err != nil {
			cobra.CheckErr(err)
		}
		fmt.Printf("Current Bible Version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(bibleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bibleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bibleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	bibleCmd.Flags().BoolVarP(&ListVersions, "list-versions", "l", false, "rapela bible --list-versions or --lv")
	bibleCmd.Flags().StringVarP(&NewBibleVersionId, "set", "s", "", "rapela bible --set=<name-of-version> or -s=<name-of-version>") // todo (keo): bind this to a config value so that on startup we use this as the default value
}
