package bt

import "testing"

type TestLib struct {
}

func (lib *TestLib) ParseLib(infohash string) string {
	return "https://itorrents.org/torrent/" + infohash + ".torrent"
}
func (lib *TestLib) GetRequestMethod() string {
	return "GET"
}

func TestDownloadTorrent(T *testing.T) {
	m2bt := Magnet2Bt{
		libs: []BtorrentLibrary{
			&TestLib{},
		},
	}

	data, err := m2bt.DownloadTorrent("B415C913643E5FF49FE37D304BBB5E6E11AD5101")
	if err != nil {
		T.Error(err)
	}

	T.Logf("download data length:%d", len(data))
}
