package netease

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/music"
)

type Netease struct {
	client *http.Client
}

func (self *Netease) request(method, url string, data map[string]string) ([]byte, error) {
	var body io.Reader

	if method != "GET" && data != nil {
		body = fromData(data)
	}

	fmt.Println("url:", url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if method == "GET" && data != nil {
		q := req.URL.Query()
		for k, v := range data {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	headers := map[string]string{
		"Referer":         "https://music.163.com",
		"Cookie":          "appver=8.2.30; os=iPhone OS; osver=15.0; EVNSM=1.0.0; buildver=2206; channel=distribution; machineid=iPhone13.3",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_5, AppleWebKit/605.1.15 (KHTML, like Gecko,",
		"X-Real-IP":       randomIP(),
		"Accept":          "*/*",
		"Accept-Encoding": "gzip,deflate,sdch",
		"Accept-Language": "zh-CN,zh;q=0.8,gl;q=0.6,zh-TW;q=0.4",
		"Connection":      "keep-alive",
		"Content-Type":    "application/x-www-form-urlencoded",
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("result is null")
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	return ioutil.ReadAll(gz)
}

func (self *Netease) weapi(data map[string]string) map[string]string {
	text, err := json.Marshal(data)
	if err != nil {
		return data
	}

	const (
		iv        = "0102030405060708"
		presetKey = "0CoJUm6Qyw8W8jud"
		publicKey = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDgtQn2JZ34ZC28NWYpAUd98iZ37BUrX/aKzmFbt7clFSs6sXqHauqKWqdtLkF2KexO40H1YTX8z2lSgBBOAxLsvaklV8k4cBFK9snQXE9/DDaFt6Rr7iVZMldczhC0JNgTz+SHXT6CBHuX3e9SdB1Ua44oncaTWz7OBGLbCiK45wIDAQAB\n-----END PUBLIC KEY-----"
		charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	)

	secret := randomBytes(16, charset)

	text = base64Bytes(aesEncrypt(text, []byte(presetKey), []byte(iv)))
	text = base64Bytes(aesEncrypt(text, secret, []byte(iv)))

	enck := hexBytes(rsaEncrypt(reverseBytes(secret), []byte(publicKey)))
	return map[string]string{
		"params":    string(text),
		"encSecKey": string(enck),
	}
}

func (self *Netease) toSong(result gjson.Result) *music.Song {
	song := &music.Song{
		Id:      result.Get("id").String(),
		Name:    result.Get("name").String(),
		Picture: result.Get("al.picUrl").String(),
	}
	song.Picture = strings.ReplaceAll(song.Picture, "http://", "https://")
	if arts := result.Get("ar").Array(); len(arts) > 0 {
		artist := make([]*music.Artist, len(arts))
		for i, art := range arts {
			artist[i] = &music.Artist{
				Id:   art.Get("id").String(),
				Name: art.Get("name").String(),
			}
		}
		song.Artist = artist
	}
	return song
}

func (self *Netease) Song(id string) (*music.Song, error) {
	b, _ := json.Marshal([]any{map[string]string{
		"id": id,
	}})
	data := map[string]string{
		"c": string(b),
	}

	res, err := self.request("POST", "https://music.163.com/weapi/v3/song/detail/", self.weapi(data))
	if err != nil {
		return nil, err
	}
	result := gjson.ParseBytes(res).Get("songs.0")
	return self.toSong(result), nil
}

func (self *Netease) SongLink(id string) (*music.SongLink, error) {
	data := map[string]string{
		"br":  fmt.Sprintf("%d", 320*1000),
		"ids": fmt.Sprintf(`["%s"]`, id),
	}

	res, err := self.request("POST", "http://music.163.com/weapi/song/enhance/player/url", self.weapi(data))
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(res)

	link := &music.SongLink{
		Id:   result.Get("data.0.id").String(),
		URL:  result.Get("data.0.url").String(),
		Br:   result.Get("data.0.br").Int(),
		Size: result.Get("data.0.size").Int(),
	}
	link.URL = strings.ReplaceAll(link.URL, "http://", "https://")
	return link, nil
}

func (self *Netease) Album(id string) (*music.Album, error) {
	data := map[string]string{
		"id":            id,
		"total":         "true",
		"offset":        "0",
		"limit":         "1000",
		"ext":           "true",
		"private_cloud": "true",
	}

	res, err := self.request("POST", fmt.Sprintf(`https://music.163.com/weapi/v1/album/%s`, id), self.weapi(data))
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(res)

	album := &music.Album{
		Id:   result.Get("album.id").String(),
		Name: result.Get("album.name").String(),
	}
	if sgs := result.Get("songs").Array(); len(sgs) > 0 {
		songs := make([]*music.Song, len(sgs))
		for i, sg := range sgs {
			songs[i] = self.toSong(sg)
		}
		album.Songs = songs
	}
	return album, nil
}

func (self *Netease) Lyric(id string) (*music.Lyric, error) {
	data := map[string]string{
		"id": id,
		"os": "linux",
		"lv": "-1",
		"kv": "-1",
		"tv": "-1",
	}

	res, err := self.request("POST", "https://music.163.com/weapi/song/lyric", self.weapi(data))
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(res)

	lrc := &music.Lyric{
		Lyric: result.Get("lrc.lyric").String(),
		Trans: result.Get("tlyric.lyric").String(),
	}
	return lrc, nil
}

func (self *Netease) Artist(id string) (*music.Artist, error) {
	data := map[string]string{
		"id":            id,
		"top":           "50",
		"ext":           "true",
		"prevate_cloud": "true",
	}

	res, err := self.request("POST", fmt.Sprintf("https://music.163.com/weapi/v1/artist/%s", id), self.weapi(data))
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(res)

	ins := &music.Artist{
		Id:   result.Get("artist.id").String(),
		Name: result.Get("artist.name").String(),
	}
	if hots := result.Get("hotSongs").Array(); len(hots) > 0 {
		songs := make([]*music.Song, len(hots))
		for i, hot := range hots {
			songs[i] = self.toSong(hot)
		}
		ins.Songs = songs
	}
	return ins, nil
}

func (self *Netease) Playlist(id string) (*music.Playlist, error) {
	data := map[string]string{
		"id": id,
		"n":  "1000",
		"s":  "0",
		"t":  "0",
	}

	res, err := self.request("POST", "https://music.163.com/weapi/v6/playlist/detail", self.weapi(data))
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(res)

	ins := &music.Playlist{
		Id:   result.Get("playlist.id").String(),
		Name: result.Get("playlist.name").String(),
	}
	if tracks := result.Get("playlist.tracks").Array(); len(tracks) > 0 {
		songs := make([]*music.Song, len(tracks))
		for i, track := range tracks {
			songs[i] = self.toSong(track)
		}
		ins.Songs = songs
	}
	return ins, nil
}

func (self *Netease) Search(keyword string) ([]*music.Song, error) {
	data := map[string]string{
		"s":      keyword,
		"type":   "1",
		"limit":  "10",
		"total":  "true",
		"offset": "0",
	}

	res, err := self.request("POST", "https://music.163.com/weapi/cloudsearch/pc", self.weapi(data))
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(res)

	songs := make([]*music.Song, 0)
	for _, song := range result.Get("result.songs").Array() {
		songs = append(songs, self.toSong(song))
	}
	return songs, nil
}
