package workers

import (
	"log"
	"time"

	"github.com/teambrookie/mediarss/db"
	"github.com/teambrookie/mediarss/provider"
)

func RefreshWorker(limiter <-chan time.Time, users []string, store db.MediaStore, mediaProvider provider.MediaProvider, mediaToSearch chan<- db.Media) {
	for {
		<-limiter
		log.Println("Refresh started")
		for _, user := range users {
			log.Printf("Refreshing for user %s\n", user)
			medias, err := mediaProvider.UnseenMedias(user)
			if err != nil {
				log.Printf("Error retriving medias for user %s : %s\n", user, err)
			}
			for _, media := range medias {
				if media, _ := store.GetMedia(media.ID, db.FOUND); (media == db.Media{}) {
					err := store.AddMedia(media, db.NOTFOUND)
					if err != nil {
						log.Printf("Error adding medias to database: %s", err)
					}
				}
			}
		}

		log.Println("Passing not found episodes to the search worker")
		notFounds, err := store.GetCollection(db.NOTFOUND)
		if err != nil {
			log.Printf("Error retriving unfound episodes from db : %s\n", err)
		}
		for _, media := range notFounds {
			mediaToSearch <- media
		}

	}
}
