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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

// downloadAudioCmd represents the downloadAudio command
var downloadAudioCmd = &cobra.Command{
	Use:   "downloadAudio <youtube url> <fileName>",
	Short: "Download Youtube video as a podcast",
	Long:  `It downloads only the audio of a youtube video`,
	Run: func(cmd *cobra.Command, args []string) {
		downloadAudio(args)
	},
}

func init() {
	rootCmd.AddCommand(downloadAudioCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadAudioCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadAudioCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func downloadAudio(args []string) error {
	match, _ := regexp.Match(`^(https?\:\/\/)?(www\.youtube\.com|youtu\.?be)\/.+$`, []byte(args[0]))

	if !match {
		fmt.Println("Pls enter a valid url")
		return nil
	}
	u, err := url.Parse(args[0])
	if err != nil {
		fmt.Println(err)
		return nil
	}

	par, _ := url.ParseQuery(u.RawQuery)

	id := par["v"][0]
	queryStr := "https://www.youtube.com/get_video_info?video_id=" + id + "&el=embedded&eurl=https://youtube.googleapis.com/v/" + id + "&sts=18333"

	resp, _ := http.Get(queryStr)
	rb, _ := ioutil.ReadAll(resp.Body)
	data := string(rb)

	params, err := url.ParseQuery(data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(params["player_response"][0]), &result)
	q := result["streamingData"].(map[string]interface{})
	h := q["adaptiveFormats"].([]interface{})

	check := h[len(h)-1].(map[string]interface{})

	if _, ok := check["url"]; !ok {
		fmt.Println("This audio cant be downloaded as it requires YouTube Premium")
		return nil
	}

	s := h[len(h)-1].(map[string]interface{})

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return err
	}

	home += "/Downloads/"

	out, err := os.Create(filepath.Join(home, filepath.Base(args[1])))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer out.Close()

	// Get the data
	res, _ := http.Get(s["url"].(string))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	fmt.Println("Downloading the audio in ~/Downloads")

	//Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	//Writer the body to file
	io.Copy(out, res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
