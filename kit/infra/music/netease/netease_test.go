package netease

import (
	"crypto/tls"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

// https://www.jianshu.com/p/97e35fa456ce
var ns = Netease{
	client: newClient(),
}

func newClient() *http.Client {
	transport := &http.Transport{
		IdleConnTimeout:       120 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 60 * time.Second,
		}).Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}

func Test_Search(t *testing.T) {
	resp, err := ns.Search("许嵩")
	if err != nil {
		t.Error(err)
	}

	t.Log(jsonconv.AnyToJsonIndent(resp))
}

func Test_Song(t *testing.T) {
	resp, err := ns.Song("167873")
	if err != nil {
		t.Error(err)
	}

	t.Log(jsonconv.AnyToJsonIndent(resp))
}

func Test_SongLink(t *testing.T) {
	resp, err := ns.SongLink("167873")
	if err != nil {
		t.Error(err)
	}

	t.Log(jsonconv.AnyToJsonIndent(resp))
}

func Test_Lyric(t *testing.T) {
	resp, err := ns.Lyric("167873")
	if err != nil {
		t.Error(err)
	}

	t.Log(jsonconv.AnyToJsonIndent(resp))
}

func Test_Album(t *testing.T) {
	resp, err := ns.Album("16953")
	if err != nil {
		t.Error(err)
	}

	t.Log(jsonconv.AnyToJsonIndent(resp))
}

func Test_Artist(t *testing.T) {
	resp, err := ns.Artist("5771")
	if err != nil {
		t.Error(err)
	}

	t.Log(jsonconv.AnyToJsonIndent(resp))
}

func Test_Playlist(t *testing.T) {
	resp, err := ns.Playlist("5771")
	if err != nil {
		t.Error(err)
	}

	t.Log(jsonconv.AnyToJsonIndent(resp))
}

func Test_ALL(t *testing.T) {
	var resp any
	var err error

	resp, err = ns.Search("许嵩")
	if err != nil {
		t.Error(err)
	}
	t.Log(jsonconv.AnyToJsonIndent(resp))

	resp, err = ns.Song("167873")
	if err != nil {
		t.Error(err)
	}
	t.Log(jsonconv.AnyToJsonIndent(resp))

	resp, err = ns.SongLink("167873")
	if err != nil {
		t.Error(err)
	}
	t.Log(jsonconv.AnyToJsonIndent(resp))

	resp, err = ns.Lyric("167873")
	if err != nil {
		t.Error(err)
	}
	t.Log(jsonconv.AnyToJsonIndent(resp))

	resp, err = ns.Album("167873")
	if err != nil {
		t.Error(err)
	}
	t.Log(jsonconv.AnyToJsonIndent(resp))

	resp, err = ns.Artist("5771")
	if err != nil {
		t.Error(err)
	}
	t.Log(jsonconv.AnyToJsonIndent(resp))

	resp, err = ns.Playlist("167873")
	if err != nil {
		t.Error(err)
	}
	t.Log(jsonconv.AnyToJsonIndent(resp))
}
