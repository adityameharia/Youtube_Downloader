/*
Package cmd Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

// linksCmd represents the links command
var linksCmd = &cobra.Command{
	Use:   "links",
	Short: "Prints all the links on a HTML page",
	Long:  `Parses the HTML response to find all the a tags on a particular page`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		links(args[0])
	},
}

func init() {
	rootCmd.AddCommand(linksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// linksCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// linksCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func links(u string) error {
	doc, err := goquery.NewDocument(u)

	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		fmt.Println(href)

	})
	return nil
}
