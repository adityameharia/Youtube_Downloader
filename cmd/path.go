/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// pathCmd represents the path command
var pathCmd = &cobra.Command{
	Use:   "path <starting wiki page> <ending wiki page>",
	Short: "Find a random path between two wikipedia pages",
	Long: `Finds a random path between two wikipedia pages using recursion.

Since wikipedia blocks multiple request from one ip addess.
So we had to limit the number of request we are sending per sec and hence
having a path greater than 3-4 jumps might take more than 10mins`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		GetPath(args)
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pathCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pathCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GetPath(args []string) error {
	start := args[0]
	end := args[1]
	limit := make(chan int, 250)

	pageFound := make(chan []string)

	var path []string

	set := make(map[string]bool)

	findPage(start, end, path, pageFound, limit, set)
	found := <-pageFound
	found = append([]string{start}, found...)
	for _, urlAndNextStep := range found {
		fmt.Println(urlAndNextStep)
	}
	return nil
}

func getURL(url string) (pageData string) {
	res, err := http.Get(url)

	if err != nil {
		//fmt.Println(err)
		return
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	return (string(bytes))
}

func urls(html string) (wikiURLs []string) {
	wikiURLs = make([]string, 0)
	wikipediaURLPrefix := "https://en.wikipedia.org"

	untrimmedAnchors := strings.Split(html, "<a href=\"")[1:]
	for _, untrimmedAnchor := range untrimmedAnchors {
		url := untrimmedAnchor[:strings.Index(untrimmedAnchor, "\"")]

		if strings.HasPrefix(url, "/wiki") && !strings.Contains(url, ":") {
			fullURL := wikipediaURLPrefix + url
			name := untrimmedAnchor[strings.Index(untrimmedAnchor, ">")+1 : strings.Index(untrimmedAnchor, "<")]
			if !strings.Contains(name, "<") || name == "Read" {
				wikiURLs = append(wikiURLs, fullURL)
			}
		}
	}

	return wikiURLs
}

func findPage(start string, end string, path []string, pageFound chan []string, limit chan int, set map[string]bool) {
	limit <- 1
	pageHTML := getURL(start)
	urls := urls(pageHTML)

	temp := make([]string, len(path))
	copy(temp, path)
	time.Sleep(2 * time.Second)
	for _, url := range urls {
		if set[url] {
			return
		}
		//fmt.Println(i)
		select {
		case <-pageFound:
			return
		default:
			temp = append(path, url)
			if url == end {
				pageFound <- temp
				close(pageFound)
				return
			}
			if len(temp) < 6 {
				go findPage(url, end, temp, pageFound, limit, set)
			}

		}
	}
	<-limit
	//fmt.Println(ch)
}
