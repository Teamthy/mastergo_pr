package repository

import (
	"backend/internal/models"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, onboarding_status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query, user.Email, user.PasswordHash, user.OnboardingStatus).
		Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, onboarding_status, is_verified FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.OnboardingStatus, &user.IsVerified)

	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *UserRepository) UpdateOnboardingStatus(ctx context.Context, email string, status models.OnboardingStep) error {

	query := `UPDATE users SET onboarding_status = $1, is_verified = TRUE WHERE email = $2`
	_, err := r.db.Exec(ctx, query, status, email)
	return err
}
func (r *UserRepository) MarkUserVerified(ctx context.Context, email string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET is_verified = true WHERE email = $1`,
		email,
	)
	return err
}
func (r *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user := &models.User{}

	query := `
		SELECT id, email, password_hash, full_name, onboarding_status, is_verified, created_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, userID).
		Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.OnboardingStatus,
			&user.IsVerified,
			&user.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateProfile(
	ctx context.Context,
	userID uuid.UUID,
	fullName string,
	status models.OnboardingStep,
) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET full_name = $1,
		    onboarding_status = $2
		WHERE id = $3
	`, fullName, status, userID)

	return err
}
