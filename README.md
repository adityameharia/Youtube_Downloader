# Youtube Downloader

A snap package written using golang and cobra, capable of downloading audio and video data from youtube links.

## Installation

The package is available on [Snap](https://snapcraft.io/yt-downloader). If you are on any linux distribution, just run the following command to install

```
$ sudo snap install yt-downloader
```

## Usage

You can easily download non-copyrighted videos 

```
USAGE
    $ yt-downloader download LINK

ARGUMENTS
  LINK          Link of the youtube video to be downloaded

Options
    -d, --hd to download videos in 720p
    -a, --audio to download only audio 
    -n, --name to specify the name of the file
```

## License

Licensed under the [Apache 2.0](LICENSE) License. All Rights Reserved.

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)

