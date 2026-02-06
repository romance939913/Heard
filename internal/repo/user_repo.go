package repo

import (
    "context"
    "database/sql"

    "github.com/brennanromance/heard/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct{ db *sql.DB }

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	// Hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var id int
	err = r.db.QueryRowContext(ctx, `INSERT INTO users (username, email, password) VALUES ($1,$2,$3) RETURNING id`, u.Username, u.Email, string(hashedPassword)).Scan(&id)
	if err != nil {
		return err
	}
	u.ID = id
	u.Password = "" // Don't return the password
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx, `SELECT id, username, email, password FROM users WHERE id=$1`, id).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx, `SELECT id, username, email, password FROM users WHERE username=$1`, username).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx, `SELECT id, username, email, password FROM users WHERE email=$1`, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// VerifyPassword checks if the provided password matches the hashed password
func (r *UserRepo) VerifyPassword(ctx context.Context, userID int, password string) error {
	var hashedPassword string
	err := r.db.QueryRowContext(ctx, `SELECT password FROM users WHERE id=$1`, userID).Scan(&hashedPassword)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (r *UserRepo) Update(ctx context.Context, u *models.User) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET username=$1, email=$2, password=$3 WHERE id=$4`, u.Username, u.Email, u.Password, u.ID)
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}

func (r *UserRepo) List(ctx context.Context) ([]*models.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, username, email, password FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password); err != nil {
			return nil, err
		}
		out = append(out, &u)
	}
	return out, nil
}