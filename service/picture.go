package service

import (
	"context"
	"database/sql"

	"github.com/sungyo4869/portfolio/model"
)

type PictureService struct {
	db *sql.DB
}

func NewPictureService(db *sql.DB) *PictureService {
	return &PictureService{
		db: db,
	}
}

func (s *PictureService) CreatePicture(ctx context.Context, card_id int64, path string) (*model.Picture, error) {
	var picture model.Picture

	const (
		insert  = `INSERT INTO pictures(card_id, path) VALUES(?, ?)`
		confirm = `SELECT id, card_id, path FROM pictures WHERE id = ?`
	)

	result, err := s.db.ExecContext(ctx, insert, card_id, path)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	row, err := s.db.QueryContext(ctx, confirm, id)
	if err != nil {
		return nil, err
	}

	err = row.Scan(
		&picture.ID,
		&picture.CardID,
		&picture.PicturePath,
	)
	if err != nil {
		return nil, err
	}

	return &picture, nil
}
func (s *PictureService) ReadPicture(ctx context.Context) ([]*model.Picture, error) {
	var pictures []*model.Picture

	const query = `SELECT id, card_id, path FROM pictures`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var picture model.Picture
		err := rows.Scan(
			&picture.ID, 
			&picture.CardID, 
			&picture.PicturePath,
		)
		if err != nil {
			return nil, err
		}

		pictures = append(pictures, &picture)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pictures, nil
}

func (s *PictureService) DeletePicture(ctx context.Context, pictureID int64) error {
	const query = `DELETE FROM pictures WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query, pictureID)
	if err != nil {
		return err
	}
	return nil
}
