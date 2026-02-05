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
    err := r.db.QueryRowContext(ctx, `INSERT INTO comment (message, post_id, user_id, upvotes, downvotes) VALUES ($1,$2,$3,$4,$5) RETURNING id`, c.Message, c.PostID, c.UserID, c.Upvotes, c.Downvotes).Scan(&id)
    if err != nil { return err }
    c.ID = id
    return nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id int) (*models.Comment, error) {
    var c models.Comment
    err := r.db.QueryRowContext(ctx, `SELECT id, message, post_id, user_id, upvotes, downvotes FROM comment WHERE id=$1`, id).Scan(&c.ID, &c.Message, &c.PostID, &c.UserID, &c.Upvotes, &c.Downvotes)
    if err != nil { return nil, err }
    return &c, nil
}

func (r *CommentRepo) Update(ctx context.Context, c *models.Comment) error {
    _, err := r.db.ExecContext(ctx, `UPDATE comment SET message=$1, post_id=$2, user_id=$3, upvotes=$4, downvotes=$5 WHERE id=$6`, c.Message, c.PostID, c.UserID, c.Upvotes, c.Downvotes, c.ID)
    return err
}

func (r *CommentRepo) Delete(ctx context.Context, id int) error {
    _, err := r.db.ExecContext(ctx, `DELETE FROM comment WHERE id=$1`, id)
    return err
}

func (r *CommentRepo) List(ctx context.Context) ([]*models.Comment, error) {
    rows, err := r.db.QueryContext(ctx, `SELECT id, message, post_id, user_id, upvotes, downvotes FROM comment ORDER BY id`)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []*models.Comment
    for rows.Next() {
        var c models.Comment
        if err := rows.Scan(&c.ID, &c.Message, &c.PostID, &c.UserID, &c.Upvotes, &c.Downvotes); err != nil { return nil, err }
        out = append(out, &c)
    }
    return out, nil
}
