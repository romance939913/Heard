package repo

import (
	"context"
	"database/sql"

	"github.com/brennanromance/heard/internal/models"
)

type CompanyRepo struct {
	db *sql.DB
}

func NewCompanyRepo(db *sql.DB) *CompanyRepo { return &CompanyRepo{db: db} }

func (r *CompanyRepo) Create(ctx context.Context, c *models.Company) error {
	var id int
	err := r.db.QueryRowContext(ctx, `INSERT INTO company (name, description, parent_company_id, industry, sub_industry, headquarters, date_incorporated, user_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`, c.Name, c.Description, c.ParentCompanyID, c.Industry, c.SubIndustry, c.Headquarters, c.DateIncorporated, c.UserID).Scan(&id)
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *CompanyRepo) GetByID(ctx context.Context, id int) (*models.Company, error) {
	var c models.Company
	var sub sql.NullString
	var hq sql.NullString
	var dt sql.NullTime
	var uid sql.NullInt32
	err := r.db.QueryRowContext(ctx, `SELECT id, name, description, parent_company_id, industry, sub_industry, headquarters, date_incorporated, user_id FROM company WHERE id=$1`, id).Scan(&c.ID, &c.Name, &c.Description, &c.ParentCompanyID, &c.Industry, &sub, &hq, &dt, &uid)
	if err != nil {
		return nil, err
	}
	if sub.Valid {
		s := sub.String
		c.SubIndustry = &s
	}
	if hq.Valid {
		h := hq.String
		c.Headquarters = &h
	}
	if dt.Valid {
		d := dt.Time.Format("2006-01-02")
		c.DateIncorporated = &d
	}
	if uid.Valid {
		u := int(uid.Int32)
		c.UserID = &u
	}
	return &c, nil
}

func (r *CompanyRepo) Update(ctx context.Context, c *models.Company) error {
	_, err := r.db.ExecContext(ctx, `UPDATE company SET name=$1, description=$2, parent_company_id=$3, industry=$4, sub_industry=$5, headquarters=$6, date_incorporated=$7, user_id=$8 WHERE id=$9`, c.Name, c.Description, c.ParentCompanyID, c.Industry, c.SubIndustry, c.Headquarters, c.DateIncorporated, c.UserID, c.ID)
	return err
}

func (r *CompanyRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM company WHERE id=$1`, id)
	return err
}

func (r *CompanyRepo) List(ctx context.Context) ([]*models.Company, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, description, parent_company_id, industry, sub_industry, headquarters, date_incorporated, user_id FROM company ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Company
	for rows.Next() {
		var c models.Company
		var sub sql.NullString
		var hq sql.NullString
		var dt sql.NullTime
		var uid sql.NullInt32
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.ParentCompanyID, &c.Industry, &sub, &hq, &dt, &uid); err != nil {
			return nil, err
		}
		if sub.Valid {
			s := sub.String
			c.SubIndustry = &s
		}
		if hq.Valid {
			h := hq.String
			c.Headquarters = &h
		}
		if dt.Valid {
			d := dt.Time.Format("2006-01-02")
			c.DateIncorporated = &d
		}
		if uid.Valid {
			u := int(uid.Int32)
			c.UserID = &u
		}
		out = append(out, &c)
	}
	return out, nil
}
