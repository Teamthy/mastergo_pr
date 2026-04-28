package repository

import (
	"context"
	"time"

	"backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	firstName,
	lastName,
	email,
	passwordHash string,
) (*models.User, error) {

	user := models.NewUser(firstName, lastName, email, passwordHash)

	query := `
	INSERT INTO users (
		id,
		first_name,
		last_name,
		email,
		password_hash,
		onboarding_status,
		email_verified,
		created_at,
		updated_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.OnboardingStatus,
		user.EmailVerified,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `
	SELECT 
		id,
		first_name,
		last_name,
		email,
		password_hash,
		phone,
		address,
		onboarding_status,
		email_verified,
		last_login_at,
		created_at,
		updated_at
	FROM users
	WHERE email=$1
	`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.Address,
		&user.OnboardingStatus,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User

	query := `
	SELECT 
		id,
		first_name,
		last_name,
		email,
		password_hash,
		phone,
		address,
		onboarding_status,
		email_verified,
		last_login_at,
		created_at,
		updated_at
	FROM users
	WHERE id=$1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.Address,
		&user.OnboardingStatus,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) MarkUserVerified(ctx context.Context, email string) error {
	query := `
	UPDATE users
	SET email_verified = true,
	    onboarding_status = $2,
	    updated_at = $3
	WHERE email = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		email,
		models.StepEmailVerified,
		time.Now(),
	)

	return err
}

// Update updates user's last_login_at timestamp
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
	UPDATE users
	SET
		last_login_at = $2,
		updated_at = $3
	WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.LastLoginAt,
		time.Now(),
	)

	return err
}

func (r *UserRepository) UpdateProfile(
	ctx context.Context,
	userID uuid.UUID,
	firstName,
	lastName,
	phone,
	address string,
	status models.OnboardingStep,
) error {

	query := `
	UPDATE users
	SET
		first_name = $2,
		last_name = $3,
		phone = $4,
		address = $5,
		onboarding_status = $6,
		updated_at = $7
	WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		firstName,
		lastName,
		phone,
		address,
		status,
		time.Now(),
	)

	return err
}
