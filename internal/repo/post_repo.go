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
	err := r.db.QueryRowContext(ctx, `INSERT INTO post (title, description, company_id, user_id, upvotes, downvotes) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`, pModel.Title, pModel.Description, pModel.CompanyID, pModel.UserID, pModel.Upvotes, pModel.Downvotes).Scan(&id, &createdAt, &updatedAt)
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

func (r *PostRepo) GetByID(ctx context.Context, id int) (*models.Post, error) {
	var p models.Post
	var createdAt, updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `SELECT id, title, description, company_id, user_id, upvotes, downvotes, created_at, updated_at FROM post WHERE id=$1`, id).Scan(&p.ID, &p.Title, &p.Description, &p.CompanyID, &p.UserID, &p.Upvotes, &p.Downvotes, &createdAt, &updatedAt)
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
	err := r.db.QueryRowContext(ctx, `UPDATE post SET title=$1, description=$2, company_id=$3, user_id=$4, upvotes=$5, downvotes=$6 WHERE id=$7 RETURNING updated_at`, pModel.Title, pModel.Description, pModel.CompanyID, pModel.UserID, pModel.Upvotes, pModel.Downvotes, pModel.ID).Scan(&updatedAt)
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
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, description, company_id, user_id, upvotes, downvotes, created_at, updated_at FROM post ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Post
	for rows.Next() {
		var p models.Post
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.CompanyID, &p.UserID, &p.Upvotes, &p.Downvotes, &createdAt, &updatedAt); err != nil {
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
