package ytdownload

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type writeCounter struct {
	downloaded uint64
	Total      uint64
}

//GetDownloadData gets the json resp
func GetDownloadData(queryStr string) (map[string]interface{}, error) {
	resp, err := http.Get(queryStr)
	if err != nil {
		fmt.Println("There was a problem connecting to the remote server.\nPls check your internet connection to make sure it isnt a problem at your end.")
		return nil, err
	}
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("An error occured")
		return nil, err
	}
	data := string(rb)

	params, err := url.ParseQuery(data)
	if err != nil {
		fmt.Println("Unable to parse given URL")
		return nil, err
	}

	var result map[string]interface{}

	json.Unmarshal([]byte(params["player_response"][0]), &result)

	//q := result["streamingData"].(map[string]interface{})

	return result, nil
}

//get the title of the video
func GetFilename(q map[string]interface{}) string {
	q = q["videoDetails"].(map[string]interface{})
	name := q["title"].(string)
	return name
}

//GetDownloadAudioURL get the url for downloading only audio
func GetDownloadAudioURL(q map[string]interface{}, high bool, audio bool) (string, error) {

	q = q["streamingData"].(map[string]interface{})

	h := q["adaptiveFormats"].([]interface{})

	check := h[len(h)-1].(map[string]interface{})

	if _, ok := check["url"]; !ok {
		fmt.Println("This audio cant be downloaded due to Copyright Infringment.This is generally the case with music videos")
		return "nil", errors.New("Error")
	}

	s := h[len(h)-1].(map[string]interface{})

	return s["url"].(string), nil
}

//GetDownloadURL get the url for downoading the video
func GetDownloadURL(q map[string]interface{}, high bool, audio bool) (string, error) {

	q = q["streamingData"].(map[string]interface{})

	h := q["formats"].([]interface{})
	var index int

	check := h[0].(map[string]interface{})

	if _, ok := check["url"]; !ok {
		fmt.Println("This video cant be downloaded due to Copyright Infringment.This is generally the case with music videos")
		return "nil", errors.New("Error")
	}

	if high == true && len(h) >= 2 {
		index = 1
	} else {
		if high == true {
			fmt.Println("Unfortunately 720p quality wasnt available")
		}
		index = 0
	}

	s := h[index].(map[string]interface{})

	return s["url"].(string), nil
}

//DownloadVideo actually downloads the video,calculates % downloaded and writes to file
func DownloadVideo(link string, out *os.File) error {
	counter := &writeCounter{}

	res, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	//Check server response
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	counter.Total = uint64(res.ContentLength)

	//Writer the body to file
	_, err = io.Copy(out, io.TeeReader(res.Body, counter))
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer out.Close()

	fmt.Printf("\r%s", strings.Repeat(" ", 36))
	fmt.Printf("\rDownload Completed...U can view the file at in your ~/Downloads directory\n")

	return nil
}

func CheckIsLive(q map[string]interface{}) bool {
	q = q["videoDetails"].(map[string]interface{})
	if _, ok := q["isLive"]; ok {
		if q["isLive"].(bool) {
			return true
		}
	}
	return false
}
