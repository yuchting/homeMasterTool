package bt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type BtorrentLibrary interface {
	ParseLib(infohash string) string
	GetRequestMethod() string
}

type Magnet2Bt struct {
	libs []BtorrentLibrary
}

func downloadFile(requestMethod string, url string) ([]byte, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
		Jar:     jar,
	}

	request, err := http.NewRequest(requestMethod, url, nil)

	//request.Header.Add(":authority", "itorrents.org")
	//request.Header.Add(":method", "GET")
	//request.Header.Add(":scheme", "https")
	//request.Header.Add(":path", "/torrent/B415C913643E5FF49FE37D304BBB5E6E11AD5101.torrent")
	request.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	request.Header.Add("accept-encoding", "gzip, deflate, br")
	request.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-US;q=0.7")
	request.Header.Add("cookie", "__cfduid=d858b50be09e77a0ed1f2eafd25de20261532600128; _ga=GA1.2.1107172433.1532600134; _gid=GA1.2.1435366109.1532600134; cf_clearance=8af637566d06a729bbc6e48cd26390c619ddf592-1532669773-3600")
	request.Header.Add("dnt", "1")
	request.Header.Add("referer", "https://itorrents.org/torrent/B415C913643E5FF49FE37D304BBB5E6E11AD5101.torrent")
	request.Header.Add("upgrade-insecure-requests", "1")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")

	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

const btorrentFileHeader = "d8:announce49"

func isBtorrent(data []byte) bool {

	if len(data) < len(btorrentFileHeader) {
		return false
	}

	for i, v := range []byte(btorrentFileHeader) {
		if data[i] != v {
			return false
		}
	}

	return true
}

func (bt *Magnet2Bt) DownloadTorrent(infohash string) ([]byte, error) {
	for _, v := range bt.libs {
		if data, err := downloadFile(v.GetRequestMethod(), v.ParseLib(infohash)); err == nil && isBtorrent(data) {
			return data, nil
		}
	}

	return nil, fmt.Errorf("cannot download btorrent file! %s", infohash)
}
