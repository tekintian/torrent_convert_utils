package btutils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/tekintian/metainfo"
)

// get magnet link string from torrent io reader
func GetMagnetLinkFromTReader(r io.Reader) (string, error) {
	mi, err := metainfo.Load(r)
	if err != nil {
		return "", err
	}
	info, err := mi.UnmarshalInfo()
	if err != nil {
		return "", err
	}
	infoHash := mi.HashInfoBytes()
	magnet := mi.Magnet(&infoHash, &info)
	hash := magnet.InfoHash.AsString()
	fmt.Println(hash)
	//这里仅返回 magnetLink 其他信息不返回
	magnetLink := fmt.Sprintf("magnet:?xt=urn:btih:%s", magnet.InfoHash.HexString())
	return magnetLink, err
}

// convert torrent file to magent string
func TorrentToMagnet(torrentFile string) (string, error) {
	file, err := os.Open(torrentFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return GetMagnetLinkFromTReader(file)
}

//remote torrent url to magnet
func RemoteTorrentToMagnet(torrentUrl string) (string, error) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	req, _ := http.NewRequest("GET", torrentUrl, nil)
	// req.Header.Set("User-Agent", fake.GetRandUa())
	// req.Header.Set("X-FORWARDED-FOR", randIp)
	// req.Header.Set("CLIENT-IP", randIp)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	magnetLink := ""
	if resp.StatusCode == 200 {
		// ioutil.TempFile 返回文件和错误
		// 在临时文件目录创建临时文件 文件名称格式 torrent-*.torrent
		file, err := ioutil.TempFile(os.TempDir(), "torrent-*.torrent")
		if err != nil {
			return "", err
		}
		// 确保程序结束时删除临时文件
		defer os.Remove(file.Name())
		//将文件保存到临时文件
		wlen, err := io.Copy(file, resp.Body)
		if err != nil {
			return "", err
		}
		//等待文件下载完毕
		for {
			if resp.ContentLength == wlen {
				//从临时文件读取
				tfile, err := os.Open(file.Name())
				if err != nil {
					return "", err
				}
				defer tfile.Close()
				// body, err := ioutil.ReadAll(resp.Body)
				magnetLink, err = GetMagnetLinkFromTReader(tfile)
				break
			}
		}

	} else {
		err=errors.New(fmt.Sprintf("http request fail, Code:%d, Msg: %s", resp.StatusCode, resp.Status))
	}
	return magnetLink, err
}
