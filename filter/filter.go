package filter

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/qopher/go-torrentapi"
)

type Category struct {
	Type      string   `json:"type"`
	Optionnal bool     `json:"optionnal"`
	Keywords  []string `json:"keywords"`
}

type Config struct {
	Categories []Category
}

func LoadConfig(configURL, configPath string) (Config, error) {
	config := Config{}
	if configURL != "" {
		resp, err := http.Get(configURL)
		if err != nil {
			return Config{}, err
		}
		err = json.NewDecoder(resp.Body).Decode(&config)
		if err != nil {
			return Config{}, err
		}
		return config, nil
	}
	if configPath == "" {
		configPath = "./config.json"
	}
	data, _ := os.Open(configPath)
	err := json.NewDecoder(bufio.NewReader(data)).Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil

}

func chooseFilterFunc(catType string) func(string, torrentapi.TorrentResults) torrentapi.TorrentResults {
	if catType == "include" {
		return include
	}
	return exclude

}

func filterCat(category Category, torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	results := torrents
	filter := chooseFilterFunc(category.Type)

	for _, keyword := range category.Keywords {
		results = filter(keyword, torrents)
		if results != nil {
			break
		}
	}
	if category.Optionnal && (results == nil) {
		return torrents
	}
	return results

}

func Filter(categories []Category, torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	results := torrents
	for _, cat := range categories {
		tmp := filterCat(cat, results)
		results = tmp
	}
	return results
}

func include(keyword string, torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	var results torrentapi.TorrentResults
	keyword = strings.ToLower(keyword)
	for _, t := range torrents {
		var filename = strings.ToLower(t.Download)
		if strings.Contains(filename, keyword) {
			results = append(results, t)
		}
	}
	return results
}

func exclude(keyword string, torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	var results torrentapi.TorrentResults
	keyword = strings.ToLower(keyword)
	for _, t := range torrents {
		var filename = strings.ToLower(t.Download)
		if !strings.Contains(filename, keyword) {
			results = append(results, t)
		}
	}
	return results
}

func BestTorrent(torrents torrentapi.TorrentResults) torrentapi.TorrentResult {
	bt := torrentapi.TorrentResult{}
	for _, t := range torrents {
		if (bt == torrentapi.TorrentResult{}) {
			bt = t
			continue
		}
		if (t.Seeders / (1 + t.Leechers)) > (bt.Seeders / (1 + bt.Leechers)) {
			bt = t
		}
	}
	return bt
}

func FilterOutDeadTorrents(torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	var res torrentapi.TorrentResults
	for _, t := range torrents {
		if t.Seeders > 0 || t.Leechers > 0 {
			res = append(res, t)
		}
	}
	return res
}
