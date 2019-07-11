package handlers

import (
	"net/http"
	"time"

	"fmt"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/teambrookie/mediarss/db"
	"github.com/teambrookie/mediarss/provider"
)

type rssHandler struct {
	store         db.MediaStore
	mediaProvider provider.MediaProvider
}

func (h *rssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["user"]
	if token == "" {
		http.Error(w, "token must be set in query params", http.StatusNotAcceptable)
		return
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "ShowRSS by binou",
		Link:        &feeds.Link{Href: "https://github.com/TeamBrookie/showrss"},
		Description: "A list of torrent for your unseen Betaseries episodes",
		Author:      &feeds.Author{Name: "Fabien Foerster", Email: "fabienfoerster@gmail.com"},
		Created:     now,
	}
	medias, err := h.mediaProvider.UnseenMedias(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, m := range medias {
		media, err := h.store.GetMedia(m.ID, db.FOUND)
		if media.Magnet == "" || err != nil {
			continue
		}
		description := fmt.Sprintf("<a href='%s'>Link</a>", media.TorrentURL)
		item := &feeds.Item{
			Title:       media.Name,
			Link:        &feeds.Link{Href: media.TorrentURL},
			Description: description,
			Created:     media.LastUpdate,
		}
		feed.Add(item)
	}

	w.Header().Set("Content-Type", "text/xml")
	feed.WriteRss(w)
	return

}

func RSSHandler(store db.MediaStore, mediaProvider provider.MediaProvider) http.Handler {
	return &rssHandler{
		store:         store,
		mediaProvider: mediaProvider,
	}
}
