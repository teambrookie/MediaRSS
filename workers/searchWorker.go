package workers

import (
	"log"
	"time"

	"github.com/qopher/go-torrentapi"
	"github.com/teambrookie/mediarss/db"
	"github.com/teambrookie/mediarss/filter"
)

type SearchFunc func(db.Media) (torrentapi.TorrentResults, error)

const apiRateLimit = 2 * time.Second

//SearchWorker is a function who search
func SearchWorker(mediaToSearch <-chan db.Media, store db.MediaStore, searchFunc SearchFunc, config filter.Config) {
	for media := range mediaToSearch {
		time.Sleep(apiRateLimit)
		torrents, err := searchFunc(media)
		if err != nil {
			log.Printf("Error processing : %s", media.Name)
		}
		torrents = filter.FilterOutDeadTorrents(torrents)
		torrents = filter.Filter(config.Categories, torrents)
		torrent := filter.BestTorrent(torrents)

		//if no torrent is found continue
		if torrent == (torrentapi.TorrentResult{}) {
			log.Printf("%s NOT FOUND", media.Name)
			continue
		}

		media.AddTorrentInfo(torrent)
		log.Println(media.Filename)
		err = store.AddMedia(media, db.FOUND)
		if err != nil {
			log.Printf("Error adding media : %s to FOUND collection", media.Name)
		}
		err = store.DeleteMedia(media.ID, db.NOTFOUND)
		if err != nil {
			log.Printf("Error removing media : %s from NOTFOUND collection", media.Name)
		}
	}
}
