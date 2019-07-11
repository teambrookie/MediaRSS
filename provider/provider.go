package provider

import "github.com/teambrookie/mediarss/db"

type MediaProvider interface {
	UnseenMedias(string) ([]db.Media, error)
}
