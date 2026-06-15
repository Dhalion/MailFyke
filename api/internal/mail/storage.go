package mail

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

func Store(pool *pgxpool.Pool, email *ParsedEmail, orgID string) error {
	// TODO: implement DB insert into emails table
	return nil
}
