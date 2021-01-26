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
	"strings"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <youtube video url> <filename>",
	Short: "Download Youtube videos to your ~/Downloads directory with the given file",
	Args:  cobra.MinimumNArgs(2),
	Long:  `Download Youtube videos to your ~/Downloads directory with the given file`,
	Run: func(cmd *cobra.Command, args []string) {
		high, _ := cmd.Flags().GetBool("high")
		audio, _ := cmd.Flags().GetBool("audio")
		//fmt.Println("download called")
		downloadVideo(args, high, audio)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().Bool("high", false, "Downloads video in 720p if available")
	downloadCmd.Flags().BoolP("audio", "a", false, "Downloads only audio")
}

type writeCounter struct {
	downloaded uint64
	Total      uint64
}

func downloadVideo(args []string, high bool, audio bool) error {

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

	var id string
	if u.Host == "youtu.be" {
		num := strings.LastIndex(args[0], "/")
		id = args[0][num+1:]
	} else {

		par, _ := url.ParseQuery(u.RawQuery)

		id = par["v"][0]
	}
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

	var s map[string]interface{}

	if audio == false {
		h := q["formats"].([]interface{})
		var index int

		check := h[0].(map[string]interface{})

		if _, ok := check["url"]; !ok {
			fmt.Println("This video cant be downloaded as it requires YouTube Premium")
			return nil
		}

		if high == true && len(h) >= 2 {
			index = 1
		} else {
			if high == true {
				fmt.Println("Unfortunately 720p quality wasnt available")
			}
			index = 0
		}

		s = h[index].(map[string]interface{})

	} else {
		h := q["adaptiveFormats"].([]interface{})

		check := h[len(h)-1].(map[string]interface{})

		if _, ok := check["url"]; !ok {
			fmt.Println("This audio cant be downloaded as it requires YouTube Premium")
			return nil
		}

		s = h[len(h)-1].(map[string]interface{})

	}
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

	// y := spinner.New(spinner.CharSets[0], 100*time.Millisecond)
	// y.Prefix = "Downloading the video in ~/Downloads: "
	// y.Start()
	// defer y.Stop()
	// Get the data

	counter := &writeCounter{}
	res, err := http.Get(s["url"].(string))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	counter.Total = uint64(res.ContentLength)

	//Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	//Writer the body to file
	io.Copy(out, io.TeeReader(res.Body, counter))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("\r%s", strings.Repeat(" ", 36))
	fmt.Printf("\rDownload Completed...U can view the file at %v%v\n", home, args[1])
	return nil

}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.downloaded += uint64(n)
	wc.printProgress()
	return n, nil
}

func (wc writeCounter) printProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 36))

	progress := (wc.downloaded * 100) / wc.Total
	//str := fmt.Sprintf("%f", progress)
	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Print("\rDownloading... " + fmt.Sprint(progress) + "% complete")
}
