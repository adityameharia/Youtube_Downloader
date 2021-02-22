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
	"errors"
	"fmt"
	"os"
	"regexp"

	ytdownload "yt-downloader/cmd/ytdownload"

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
		err := download(args, high, audio)
		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolP("hd", "d", false, "Downloads video in 720p if available")
	downloadCmd.Flags().BoolP("audio", "a", false, "Downloads only audio")
}

// func (wc writeCounter) printProgress() {
// 	// Clear the line by using a character return to go back to the start and remove
// 	// the remaining characters by filling it with spaces
// 	fmt.Printf("\r%s", strings.Repeat(" ", 36))

// 	progress := (wc.downloaded * 100) / wc.Total

// 	fmt.Print("\rDownloading... " + fmt.Sprint(progress) + "% complete")
// }

// func onSigInt(path string) {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt)
// 	go func() {
// 		<-c
// 		err := os.Remove(path)
// 		if err != nil {
// 			fmt.Println("Unable to delete the file created")
// 		}
// 		os.Exit(1)
// 	}()
// }

func download(args []string, high bool, audio bool) error {

	match, _ := regexp.Match(`^(https?\:\/\/)?(www\.youtube\.com|youtu\.?be)\/.+$`, []byte(args[0]))

	if !match {
		fmt.Println("Pls enter a valid url")
		return errors.New("bad url")
	}

	// u, err := url.Parse(args[0])
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil
	// }

	// var id string
	// if u.Host == "youtu.be" {
	// 	num := strings.LastIndex(args[0], "/")
	// 	id = args[0][num+1:]
	// } else {

	// 	par, _ := url.ParseQuery(u.RawQuery)

	// 	id = par["v"][0]
	// }

	id, err := ytdownload.GetID(args[0])
	if err != nil {
		fmt.Println("Unable to extract ID from the link")
		return err
	}

	//url at which we send our request
	queryStr := "https://www.youtube.com/get_video_info?video_id=" + id + "&el=embedded&eurl=https://youtube.googleapis.com/v/" + id + "&sts=18333"

	// home, err := os.UserHomeDir()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// home += "/Downloads/"
	// onSigInt(home + args[1])
	// if _, err := os.Stat(home + args[1]); err == nil {

	// 	fmt.Println("The given filename already exists")
	// 	return nil
	// }

	out, err := ytdownload.CreateFile(args[1])
	if err != nil {
		return err
	}

	q, err := ytdownload.GetDownloadData(queryStr)
	if err != nil {
		return err
	}
	//var s map[string]interface{}

	var link string

	if audio == false {
		// h := q["formats"].([]interface{})
		// var index int

		// check := h[0].(map[string]interface{})

		// if _, ok := check["url"]; !ok {
		// 	fmt.Println("This video cant be downloaded as it requires YouTube Premium")
		// 	return nil
		// }

		// if high == true && len(h) >= 2 {
		// 	index = 1
		// } else {
		// 	if high == true {
		// 		fmt.Println("Unfortunately 720p quality wasnt available")
		// 	}
		// 	index = 0
		// }

		// s = h[index].(map[string]interface{})

		link, err = ytdownload.GetDownloadURL(q, high, audio)
		if err != nil {
			return err
		}

	} else {

		link, err = ytdownload.GetDownloadAudioURL(q, high, audio)
		if err != nil {
			return err
		}

		// h := q["adaptiveFormats"].([]interface{})

		// check := h[len(h)-1].(map[string]interface{})

		// if _, ok := check["url"]; !ok {
		// 	fmt.Println("This audio cant be downloaded as it requires YouTube Premium")
		// 	return nil
		// }

		// s = h[len(h)-1].(map[string]interface{})

	}
	// y := spinner.New(spinner.CharSets[0], 100*time.Millisecond)
	// y.Prefix = "Downloading the video in ~/Downloads: "
	// y.Start()
	// defer y.Stop()
	// Get the data

	err = ytdownload.DownloadVideo(link, out)
	if err != nil {
		return err
	}

	return nil

}
