package main

import (
	"log"
	"net/http"

	"github.com/ve-weiyi/ve-blog-golang/pkg/plugins/music"
)

func main() {
	handler := music.MusicHandler("/music/")

	http.HandleFunc("/music/", handler)

	log.Println("Music server starting on :8080")
	log.Println("Endpoints:")
	log.Println("  GET /music/search?keyword=xxx")
	log.Println("  GET /music/song?id=xxx")
	log.Println("  GET /music/song/link?id=xxx")
	log.Println("  GET /music/lyric?id=xxx")
	log.Println("  GET /music/album?id=xxx")
	log.Println("  GET /music/artist?id=xxx")
	log.Println("  GET /music/playlist?id=xxx")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
