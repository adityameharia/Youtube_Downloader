# Youtube_Downloader

A snap package written using golang and cobra which download youtbe videos,download videos as podcast.

Uses goRoutines,channels,streams etc.

* [Usage](#usage)
* [Commands](#commands)

# Usage
<!-- usage -->
```sh-session
$ sudo snap install yt-downloader
$ huffman-compressor --help [COMMAND]
USAGE
  $ huffman-compressor COMMAND
...
```
#Commands

## `yt-downloader download LINK FILENAME`

Download non copyrighted youtube videos
```
USAGE
    $yt-downloader download LINK FILENAME

ARGUMENTS
  LINK          Link of the youtube video to be downloaded
  FILENAME      Output file name

Options
    -d, --hd=hd to download videos in 720p
    -a, --audio=audio to download only audio 
```

_See code: [cmd/download.go](https://github.com/adityameharia/Youtube_Downloader/blob/main/cmd/download.go)_