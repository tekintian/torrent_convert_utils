package btutils

import (
	"errors"
	"os"
	"path"
	"regexp"
	"strconv"
)

// convert magnet link to torrent file
func MagnetToTorrent(magnetLink string, savePath string) (fileName string, err error) {
	var torrentFilename string
	// Compile a couple of regexp we need
	var validMagnet = regexp.MustCompile(`xt=urn(\.[0-9]*)?:btih:([^&/]+)`)
	var displayName = regexp.MustCompile(`dn=([^&/]+)`)
	var infoHash = regexp.MustCompile(`btih:([^&/]+)`)

	if validMagnet.MatchString(magnetLink) {
		if displayName.MatchString(magnetLink) {
			torrentFilename = displayName.FindString(magnetLink)
		} else if infoHash.MatchString(magnetLink) {
			//torrentFilename = infoHash.FindString(magnetLink)
			torrentFilename = validMagnet.FindString(magnetLink)
		} else {
			return "", errors.New("Format of magnet URI not supported!")
		}

		// split at '='
		fileName = regexp.MustCompile(`=`).Split(torrentFilename, -1)[1]
		if len(fileName) == 0 {
			xt := validMagnet.FindString(magnetLink)
			fileName = regexp.MustCompile(`:`).Split(xt, -1)[2]
		}
		// Add torrent extension
		fileName = fileName + ".torrent"
	} else {
		// not a valid magnet URI given
		return "", errors.New("The magnet URI is not correct or unparseable!")
	}

	// Create torrent file and output magnet link string to it
	f, err := os.Create(path.Join(savePath, fileName))
	if err != nil {
		return "", err
	}
	_, wErr := f.WriteString("d10:magnet-uri" + strconv.Itoa(len(magnetLink)) + ":" + magnetLink + "e")
	if wErr != nil {
		return "", err
	}
	defer f.Close()

	return fileName, err
}
