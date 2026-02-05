package repo

import (
    "context"
    "database/sql"

    "github.com/brennanromance/heard/internal/models"
)

type CompanyRepo struct{
    db *sql.DB
}

func NewCompanyRepo(db *sql.DB) *CompanyRepo { return &CompanyRepo{db: db} }

func (r *CompanyRepo) Create(ctx context.Context, c *models.Company) error {
    var id int
    err := r.db.QueryRowContext(ctx, `INSERT INTO company (name, description, parent_company_id, sector) VALUES ($1,$2,$3,$4) RETURNING id`, c.Name, c.Description, c.ParentCompanyID, c.Sector).Scan(&id)
    if err != nil { return err }
    c.ID = id
    return nil
}

func (r *CompanyRepo) GetByID(ctx context.Context, id int) (*models.Company, error) {
    var c models.Company
    err := r.db.QueryRowContext(ctx, `SELECT id, name, description, parent_company_id, sector FROM company WHERE id=$1`, id).Scan(&c.ID, &c.Name, &c.Description, &c.ParentCompanyID, &c.Sector)
    if err != nil { return nil, err }
    return &c, nil
}

func (r *CompanyRepo) Update(ctx context.Context, c *models.Company) error {
    _, err := r.db.ExecContext(ctx, `UPDATE company SET name=$1, description=$2, parent_company_id=$3, sector=$4 WHERE id=$5`, c.Name, c.Description, c.ParentCompanyID, c.Sector, c.ID)
    return err
}

func (r *CompanyRepo) Delete(ctx context.Context, id int) error {
    _, err := r.db.ExecContext(ctx, `DELETE FROM company WHERE id=$1`, id)
    return err
}

func (r *CompanyRepo) List(ctx context.Context) ([]*models.Company, error) {
    rows, err := r.db.QueryContext(ctx, `SELECT id, name, description, parent_company_id, sector FROM company ORDER BY id`)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []*models.Company
    for rows.Next() {
        var c models.Company
        if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.ParentCompanyID, &c.Sector); err != nil { return nil, err }
        out = append(out, &c)
    }
    return out, nil
}
