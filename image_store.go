package main

import "database/sql"

var globalImageStore ImageStore

const pageSize = 25

type ImageStore interface {
	Save(img *Image) error
	Find(id string) (*Image, error)
	FindAll(offset int) ([]Image, error)
	FindAllByUser(user *User, offset int) ([]Image, error)
}

type DBImageStore struct {
	db *sql.DB
}

func NewDBImageStore() ImageStore {
	return &DBImageStore{db: globalMySQLDB}
}

func (store *DBImageStore) Save(img *Image) error {
	_, err := store.db.Exec(
		`
		REPLACE INTO images
			(id, user_id, name, location, description, size, created_at)
		VALUES
			(?, ?, ?, ?, ?, ?, ?)
		`,
		img.ID,
		img.UserID,
		img.Name,
		img.Location,
		img.Description,
		img.Size,
		img.CreatedAt,
	)
	return err
}

func (store *DBImageStore) Find(id string) (*Image, error) {
	row := store.db.QueryRow(
		`
		SELECT id, user_id, name, location, description, size, created_at
		FROM images
		WHERE id = ?
		`,
		id,
	)
	img := &Image{}
	err := row.Scan(
		img.ID,
		img.UserID,
		img.Name,
		img.Description,
		img.Size,
		img.CreatedAt,
	)
	return img, err
}

func (store *DBImageStore) FindAll(offset int) ([]Image, error) {
	return store.findAllQuery(
		`
		SELECT id, user_id, name, location, description, size, created_at
		FROM images
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?
		`,
		pageSize,
		offset,
	)
}

func (store *DBImageStore) FindAllByUser(user *User, offset int) ([]Image, error) {
	return store.findAllQuery(
		`
		SELECT id, user_id, name, location, description, size, created_at
		FROM images
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?
		`,
		user.ID,
		pageSize,
		offset,
	)
}

func (store *DBImageStore) findAllQuery(query string, args ...interface{}) ([]Image, error) {
	rows, err := store.db.Query(query, args)
	if err != nil {
		return nil, err
	}

	imgs := []Image{}
	for rows.Next() {
		img := &Image{}
		err := rows.Scan(
			img.ID,
			img.UserID,
			img.Name,
			img.Description,
			img.Size,
			img.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, *img)
	}
	return imgs, nil
}
