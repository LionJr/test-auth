package auth

import "context"

type PgRepo interface {
	GetUserIdByEmail(ctx context.Context, email string) (string, error)
	UpdateTokenHash(ctx context.Context, userId, tokenHash string) error
	GetEmailByUserId(ctx context.Context, userId string) (string, error)
	CheckRefreshTokenHash(ctx context.Context, userId, tokenHash string) (bool, error)
}
