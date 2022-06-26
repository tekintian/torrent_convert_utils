# A tiny .Torrent file Magnet convert utils lib

convert torrent file to magnet;
convert magnet to torrent file

support convert the remote .torrent file to magnet link

## Install

```sh
go get github.com/tekintian/torrent_convert_utils
```

## Usage

```go
package demo

import (
	btutils "github.com/tekintian/torrent_convert_utils"
)

func main(){
	//convert local .torrent file to magnet
	magnetLink,err:=btutils.TorrentToMagnet("yourpath/youfile.torrent")
	fmt.Println(magnetLink)
	fmt.Println(err)

	// convert remote .torrent file to magnet
	magnetLink2, err:=btutils.RemoteTorrentToMagnet("http://file.com/youfile.torrent")
	fmt.Println(magnetLink2)
	fmt.Println(err)

  // convert magnet link string to torrent file
	fileName, err:=btutils.MagnetToTorrent("magnet:?xt=urn:btih:3A5F88EB1F2ECCAEC4424416AC5DCA3B615CE515","public/upload")
	fmt.Println(fileName) // your entry file paht is public/upload/{fileName}
	fmt.Println(err)

}
```



