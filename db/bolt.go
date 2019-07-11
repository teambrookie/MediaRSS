package db

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

//BoltMediaStore implement the MediaStore interface using Bolt as the database
type BoltMediaStore struct {
	db *bolt.DB
}

//Possible collection for the DB
const (
	FOUND    = "FOUND"
	NOTFOUND = "NOTFOUND"
	SECKEY   = "SECKEY"
)

func createBucket(db *bolt.DB, name string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		return err
	})
	return err
}

//Init open a new database connection
// and create bucket if neeeded
func Init(dbName string) (*BoltMediaStore, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	if err = createBucket(db, NOTFOUND); err != nil {
		return nil, err
	}
	if err = createBucket(db, FOUND); err != nil {
		return nil, err
	}
	if err = createBucket(db, SECKEY); err != nil {
		return nil, err
	}

	return &BoltMediaStore{db}, nil
}

//Close close the database connection
func (store *BoltMediaStore) Close() error {
	return store.db.Close()
}

//GetCollection retrieve a collection of media
func (store *BoltMediaStore) GetCollection(collection string) ([]Media, error) {
	var medias []Media
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var m Media
			json.Unmarshal(v, &m)
			medias = append(medias, m)
		}
		return nil
	})
	return medias, err
}

//GetMedia retrieve a specified media
func (store *BoltMediaStore) GetMedia(mediaID string, collection string) (Media, error) {
	var media Media
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		v := b.Get([]byte(mediaID))
		json.Unmarshal(v, &media)
		return nil
	})
	return media, err
}

//GetMediaInfo return the Media object given a filename
// usefull for sorting purposes
func (store *BoltMediaStore) GetMediaInfo(filename string) (Media, error) {
	var media Media
	err := store.db.View(func(tx *bolt.Tx) error {
		secKeyBucket := tx.Bucket([]byte(SECKEY))
		primaryKey := secKeyBucket.Get([]byte(filename))
		b := tx.Bucket([]byte(FOUND))
		v := b.Get([]byte(primaryKey))
		json.Unmarshal(v, &media)
		return nil

	})
	return media, err
}

//AddMedia add a media in a specific collection
func (store *BoltMediaStore) AddMedia(media Media, collection string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		if v := b.Get([]byte(media.ID)); v != nil {
			return errors.New("Media already exists")
		}
		encoded, err := json.Marshal(media)
		if err != nil {
			return err
		}
		err = b.Put([]byte(media.ID), encoded)
		if err != nil {
			return err
		}
		if media.Filename != "" {
			secKeyBucket := tx.Bucket([]byte(SECKEY))
			return secKeyBucket.Put([]byte(media.Filename), []byte(media.ID))
		}
		return nil
	})
	return err
}

//UpdateMedia update a spacific media in a specific collection
func (store *BoltMediaStore) UpdateMedia(media Media, collection string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		encoded, err := json.Marshal(media)
		if err != nil {
			return err
		}
		if err = b.Put([]byte(media.ID), encoded); err != nil {
			return err
		}
		if media.Filename != "" {
			secKeyBucket := tx.Bucket([]byte(SECKEY))
			return secKeyBucket.Put([]byte(media.Filename), []byte(media.ID))
		}
		return nil

	})
	return err
}

//DeleteMedia delete a media from a specific collection
func (store *BoltMediaStore) DeleteMedia(mediaID string, collection string) error {
	m, err := store.GetMedia(mediaID, collection)
	if err != nil {
		return err
	}

	err = store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		if err := b.Delete([]byte(mediaID)); err != nil {
			return err
		}
		secKeyBucket := tx.Bucket([]byte(SECKEY))
		return secKeyBucket.Delete([]byte(m.Filename))
	})
	return err
}
