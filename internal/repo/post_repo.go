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
    err := r.db.QueryRowContext(ctx, `INSERT INTO post (title, description, company_id, user_id, upvotes, downvotes) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`, pModel.Title, pModel.Description, pModel.CompanyID, pModel.UserID, pModel.Upvotes, pModel.Downvotes).Scan(&id)
    if err != nil { return err }
    pModel.ID = id
    return nil
}

func (r *PostRepo) GetByID(ctx context.Context, id int) (*models.Post, error) {
    var p models.Post
    err := r.db.QueryRowContext(ctx, `SELECT id, title, description, company_id, user_id, upvotes, downvotes FROM post WHERE id=$1`, id).Scan(&p.ID, &p.Title, &p.Description, &p.CompanyID, &p.UserID, &p.Upvotes, &p.Downvotes)
    if err != nil { return nil, err }
    return &p, nil
}

func (r *PostRepo) Update(ctx context.Context, pModel *models.Post) error {
    _, err := r.db.ExecContext(ctx, `UPDATE post SET title=$1, description=$2, company_id=$3, user_id=$4, upvotes=$5, downvotes=$6 WHERE id=$7`, pModel.Title, pModel.Description, pModel.CompanyID, pModel.UserID, pModel.Upvotes, pModel.Downvotes, pModel.ID)
    return err
}

func (r *PostRepo) Delete(ctx context.Context, id int) error {
    _, err := r.db.ExecContext(ctx, `DELETE FROM post WHERE id=$1`, id)
    return err
}

func (r *PostRepo) List(ctx context.Context) ([]*models.Post, error) {
    rows, err := r.db.QueryContext(ctx, `SELECT id, title, description, company_id, user_id, upvotes, downvotes FROM post ORDER BY id`)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []*models.Post
    for rows.Next() {
        var p models.Post
        if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.CompanyID, &p.UserID, &p.Upvotes, &p.Downvotes); err != nil { return nil, err }
        out = append(out, &p)
    }
    return out, nil
}
