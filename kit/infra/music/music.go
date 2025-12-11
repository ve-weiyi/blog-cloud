package music

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"
)

type (
	Server interface {
		Search(string) ([]*Song, error)
		Song(string) (*Song, error)
		SongLink(string) (*SongLink, error)
		Lyric(string) (*Lyric, error)
		Album(string) (*Album, error)
		Artist(string) (*Artist, error)
		Playlist(string) (*Playlist, error)
	}

	Music struct {
		servers map[string]Server
	}
)

var (
	ErrIsExists   = errors.New("server is exists")
	ErrNotFound   = errors.New("server is not found")
	ErrNotSupport = errors.New("server is not support")
)

func (m *Music) Search(server string, id string) ([]*Song, error) {
	srv, ok := m.servers[server]
	if !ok {
		return nil, ErrNotFound
	}
	return srv.Search(id)
}

func (m *Music) Song(server string, id string) (*Song, error) {
	srv, ok := m.servers[server]
	if !ok {
		return nil, ErrNotFound
	}
	return srv.Song(id)
}

func (m *Music) SongLink(server string, id string) (*SongLink, error) {
	srv, ok := m.servers[server]
	if !ok {
		return nil, ErrNotFound
	}
	return srv.SongLink(id)
}

func (m *Music) Album(server string, id string) (*Album, error) {
	srv, ok := m.servers[server]
	if !ok {
		return nil, ErrNotFound
	}
	return srv.Album(id)
}

func (m *Music) Lyric(server string, id string) (*Lyric, error) {
	srv, ok := m.servers[server]
	if !ok {
		return nil, ErrNotFound
	}
	return srv.Lyric(id)
}

func (m *Music) Artist(server string, id string) (*Artist, error) {
	srv, ok := m.servers[server]
	if !ok {
		return nil, ErrNotFound
	}
	return srv.Artist(id)
}

func (m *Music) Playlist(server string, id string) (*Playlist, error) {
	srv, ok := m.servers[server]
	if !ok {
		return nil, ErrNotFound
	}
	return srv.Playlist(id)
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
