package music

import "strings"

type (
	Playlist struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Songs Songs  `json:"songs"`
	}
	Artist struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Songs Songs  `json:"songs"`
	}
	Album struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Songs Songs  `json:"songs"`
	}
	Lyric struct {
		Lyric string `json:"lyric"`
		Trans string `json:"trans"`
	}
	Song struct {
		Id      string    `json:"id"`
		Name    string    `json:"name"`
		Picture string    `json:"picture"`
		Artist  []*Artist `json:"artist"`
		Album   *Album    `json:"album"`
	}
	Songs    []*Song
	SongLink struct {
		Id   string `json:"id"`
		Br   int64  `json:"br"`
		Size int64  `json:"size"`
		URL  string `json:"url"`
	}
)

func (s *Song) GetArtist() string {
	authors := make([]string, len(s.Artist))
	for i, art := range s.Artist {
		authors[i] = art.Name
	}
	return strings.Join(authors, "/")
}
