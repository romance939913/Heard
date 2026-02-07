package repo

import (
	"context"
	"database/sql"

	"github.com/brennanromance/heard/internal/models"
)

type PostRepo struct{ db *sql.DB }

func NewPostRepo(db *sql.DB) *PostRepo { return &PostRepo{db: db} }

func (r *PostRepo) Create(ctx context.Context, pModel *models.Post) error {
	var id int
	var createdAt, updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `INSERT INTO post (title, description, company_id, user_id, likes) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at`, pModel.Title, pModel.Description, pModel.CompanyID, pModel.UserID, pModel.Likes).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return err
	}
	pModel.ID = id
	if createdAt.Valid {
		pModel.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		pModel.UpdatedAt = updatedAt.Time
	}
	return nil
}

func (r *PostRepo) ToggleLike(ctx context.Context, userID, postID int) (bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var id int
	err = tx.QueryRowContext(ctx, `SELECT id FROM post_likes WHERE user_id=$1 AND post_id=$2`, userID, postID).Scan(&id)
	if err == sql.ErrNoRows {
		if err = tx.QueryRowContext(ctx, `INSERT INTO post_likes (user_id, post_id) VALUES ($1,$2) RETURNING id`, userID, postID).Scan(&id); err != nil {
			tx.Rollback()
			return false, err
		}
		if err = tx.Commit(); err != nil {
			return false, err
		}
		return true, nil
	}
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM post_likes WHERE id=$1`, id); err != nil {
		tx.Rollback()
		return false, err
	}
	if err = tx.Commit(); err != nil {
		return false, err
	}
	return false, nil
}

func (r *PostRepo) GetByID(ctx context.Context, id int) (*models.Post, error) {
	var p models.Post
	var createdAt, updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `SELECT id, title, description, company_id, user_id, likes, created_at, updated_at FROM post WHERE id=$1`, id).Scan(&p.ID, &p.Title, &p.Description, &p.CompanyID, &p.UserID, &p.Likes, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		p.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		p.UpdatedAt = updatedAt.Time
	}
	return &p, nil
}

func (r *PostRepo) Update(ctx context.Context, pModel *models.Post) error {
	var updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `UPDATE post SET title=$1, description=$2, company_id=$3, user_id=$4, likes=$5 WHERE id=$6 RETURNING updated_at`, pModel.Title, pModel.Description, pModel.CompanyID, pModel.UserID, pModel.Likes, pModel.ID).Scan(&updatedAt)
	if err != nil {
		return err
	}
	if updatedAt.Valid {
		pModel.UpdatedAt = updatedAt.Time
	}
	return nil
}

func (r *PostRepo) Delete(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM comment WHERE post_id=$1`, id); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM post WHERE id=$1`, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *PostRepo) List(ctx context.Context) ([]*models.Post, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, description, company_id, user_id, likes, created_at, updated_at FROM post ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Post
	for rows.Next() {
		var p models.Post
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.CompanyID, &p.UserID, &p.Likes, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		if createdAt.Valid {
			p.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			p.UpdatedAt = updatedAt.Time
		}
		out = append(out, &p)
	}
	return out, nil
}
