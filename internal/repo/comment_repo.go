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
	err := r.db.QueryRowContext(ctx, `INSERT INTO comment (message, post_id, user_id, upvotes, downvotes) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at`, c.Message, c.PostID, c.UserID, c.Upvotes, c.Downvotes).Scan(&id, &createdAt, &updatedAt)
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

func (r *CommentRepo) GetByID(ctx context.Context, id int) (*models.Comment, error) {
	var c models.Comment
	var createdAt, updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `SELECT id, message, post_id, user_id, upvotes, downvotes, created_at, updated_at FROM comment WHERE id=$1`, id).Scan(&c.ID, &c.Message, &c.PostID, &c.UserID, &c.Upvotes, &c.Downvotes, &createdAt, &updatedAt)
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
	err := r.db.QueryRowContext(ctx, `UPDATE comment SET message=$1, post_id=$2, user_id=$3, upvotes=$4, downvotes=$5 WHERE id=$6 RETURNING updated_at`, c.Message, c.PostID, c.UserID, c.Upvotes, c.Downvotes, c.ID).Scan(&updatedAt)
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
	rows, err := r.db.QueryContext(ctx, `SELECT id, message, post_id, user_id, upvotes, downvotes, created_at, updated_at FROM comment ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Comment
	for rows.Next() {
		var c models.Comment
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&c.ID, &c.Message, &c.PostID, &c.UserID, &c.Upvotes, &c.Downvotes, &createdAt, &updatedAt); err != nil {
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
