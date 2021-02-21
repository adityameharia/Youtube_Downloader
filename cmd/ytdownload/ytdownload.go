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

func GetDownloadData(queryStr string) (map[string]interface{}, error) {
	resp, _ := http.Get(queryStr)
	rb, _ := ioutil.ReadAll(resp.Body)
	data := string(rb)

	params, err := url.ParseQuery(data)
	if err != nil {
		fmt.Println("Unable to parse given URL")
		return nil, err
	}

	var result map[string]interface{}

	json.Unmarshal([]byte(params["player_response"][0]), &result)

	q := result["streamingData"].(map[string]interface{})

	return q, nil
}

func GetDownloadAudioUrl(q map[string]interface{}, high bool, audio bool) (string, error) {
	h := q["adaptiveFormats"].([]interface{})

	check := h[len(h)-1].(map[string]interface{})

	if _, ok := check["url"]; !ok {
		fmt.Println("This audio cant be downloaded due to Copyright Infringment")
		return "nil", errors.New("Error")
	}

	s := h[len(h)-1].(map[string]interface{})

	return s["url"].(string), nil
}

func GetDownloadUrl(q map[string]interface{}, high bool, audio bool) (string, error) {
	h := q["formats"].([]interface{})
	var index int

	check := h[0].(map[string]interface{})

	if _, ok := check["url"]; !ok {
		fmt.Println("This audio cant be downloaded due to Copyright Infringment")
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
	fmt.Printf("\rDownload Completed...U can view the file at in your /Downloads directory\n")

	return nil
}
