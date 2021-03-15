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
	Use:   "download <youtube video url>",
	Short: "Download Youtube videos to your ~/Downloads directory",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Download Youtube videos to your ~/Downloads directory`,
	Run: func(cmd *cobra.Command, args []string) {
		high, _ := cmd.Flags().GetBool("high")
		audio, _ := cmd.Flags().GetBool("audio")
		filename, _ := cmd.Flags().GetString("name")
		//fmt.Println("download called")
		err := download(args, high, audio, filename)
		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolP("hd", "d", false, "Downloads video in 720p if available")
	downloadCmd.Flags().BoolP("audio", "a", false, "Downloads only audio")
	downloadCmd.Flags().StringP("name", "n", "", "enter filename")
}

func download(args []string, high bool, audio bool, filename string) error {

	match, _ := regexp.Match(`^(https?\:\/\/)?(www\.youtube\.com|youtu\.?be)\/.+$`, []byte(args[0]))

	if !match {
		fmt.Println("Pls enter a valid url")
		return errors.New("bad url")
	}

	id, err := ytdownload.GetID(args[0])
	if err != nil {
		fmt.Println("Unable to extract ID from the link")
		return err
	}

	//url at which we send our request
	queryStr := "https://www.youtube.com/get_video_info?video_id=" + id + "&el=embedded&eurl=https://youtube.googleapis.com/v/" + id + "&sts=18333"

	q, err := ytdownload.GetDownloadData(queryStr)
	if err != nil {
		return err
	}

	if filename == "" {
		filename = ytdownload.GetFilename(q)
	}

	out, path, err := ytdownload.CreateFile(filename)
	if err != nil {
		return err
	}

	// q, err := ytdownload.GetDownloadData(queryStr)
	// if err != nil {
	// 	err := os.Remove(path)
	// 	if err != nil {
	// 		fmt.Println("Unable to delete the file created")
	// 	}
	// 	return err
	// }
	//var s map[string]interface{}

	var link string

	if audio == false {

		link, err = ytdownload.GetDownloadURL(q, high, audio)
		if err != nil {
			err := os.Remove(path)
			if err != nil {
				fmt.Println("Unable to delete the file created")
			}
			return err
		}

	} else {

		link, err = ytdownload.GetDownloadAudioURL(q, high, audio)
		if err != nil {
			err := os.Remove(path)
			if err != nil {
				fmt.Println("Unable to delete the file created")
			}
			return err
		}

	}

	err = ytdownload.DownloadVideo(link, out)
	if err != nil {
		err := os.Remove(path)
		if err != nil {
			fmt.Println("Unable to delete the file created")
		}
		return err
	}

	return nil

}
