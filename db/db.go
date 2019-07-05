package db

import (
	"time"
)

//Media is a generic type
type Media struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Filename     string       `json:"filename"`
	Magnet       string       `json:"magnet"`
	TorrentURL   string       `json:"torrent_url"`
	Seeds        int          `json:"seeds"`
	Leechs       int          `json:"leechs"`
	ReleaseDate  time.Time    `json:"release_date"`
	LastUpdate   time.Time    `json:"last_update"`
	LastAccess   time.Time    `json:"last_access"`
	SearchParams SearchParams `json:"search_params"`
}

type SearchParams struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

//MediaStore define the interface for retriving media
type MediaStore interface {
	GetCollection(collection string) ([]Media, error)
	GetMedia(mediaID string, collection string) (Media, error)
	GetMediaInfo(filename string) (Media, error)
	AddMedia(media Media, collection string) error
	UpdateMedia(media Media, collection string) error
	DeleteMedia(mediaID string, collection string) error
	Close() error
}
