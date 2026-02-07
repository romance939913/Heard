package repo

import (
	"context"
	"database/sql"

	"github.com/brennanromance/heard/internal/models"
)

type CommentRepo struct{ db *sql.DB }

func NewCommentRepo(db *sql.DB) *CommentRepo { return &CommentRepo{db: db} }

func (r *CommentRepo) Create(ctx context.Context, c *models.Comment) error {
	var id int
	var createdAt, updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `INSERT INTO comment (message, post_id, user_id, likes) VALUES ($1,$2,$3,$4) RETURNING id, created_at, updated_at`, c.Message, c.PostID, c.UserID, c.Likes).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return err
	}
	c.ID = id
	if createdAt.Valid {
		c.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		c.UpdatedAt = updatedAt.Time
	}
	return nil
}

func (r *CommentRepo) ToggleLike(ctx context.Context, userID, commentID int) (bool, error) {
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
	err = tx.QueryRowContext(ctx, `SELECT id FROM comment_likes WHERE user_id=$1 AND comment_id=$2`, userID, commentID).Scan(&id)
	if err == sql.ErrNoRows {
		// insert
		if err = tx.QueryRowContext(ctx, `INSERT INTO comment_likes (user_id, comment_id) VALUES ($1,$2) RETURNING id`, userID, commentID).Scan(&id); err != nil {
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
	// exists -> delete
	if _, err = tx.ExecContext(ctx, `DELETE FROM comment_likes WHERE id=$1`, id); err != nil {
		tx.Rollback()
		return false, err
	}
	if err = tx.Commit(); err != nil {
		return false, err
	}
	return false, nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id int) (*models.Comment, error) {
	var c models.Comment
	var createdAt, updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `SELECT id, message, post_id, user_id, likes, created_at, updated_at FROM comment WHERE id=$1`, id).Scan(&c.ID, &c.Message, &c.PostID, &c.UserID, &c.Likes, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		c.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		c.UpdatedAt = updatedAt.Time
	}
	return &c, nil
}

func (r *CommentRepo) Update(ctx context.Context, c *models.Comment) error {
	var updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `UPDATE comment SET message=$1, post_id=$2, user_id=$3, likes=$4 WHERE id=$5 RETURNING updated_at`, c.Message, c.PostID, c.UserID, c.Likes, c.ID).Scan(&updatedAt)
	if err != nil {
		return err
	}
	if updatedAt.Valid {
		c.UpdatedAt = updatedAt.Time
	}
	return nil
}

func (r *CommentRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM comment WHERE id=$1`, id)
	return err
}

func (r *CommentRepo) List(ctx context.Context) ([]*models.Comment, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, message, post_id, user_id, likes, created_at, updated_at FROM comment ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Comment
	for rows.Next() {
		var c models.Comment
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&c.ID, &c.Message, &c.PostID, &c.UserID, &c.Likes, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		if createdAt.Valid {
			c.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			c.UpdatedAt = updatedAt.Time
		}
		out = append(out, &c)
	}
	return out, nil
}
