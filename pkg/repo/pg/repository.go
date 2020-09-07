package pg

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type PgRepository struct {
	db  *sql.DB
	log *zerolog.Logger
}

func NewPgRepository(db *sql.DB, log *zerolog.Logger) *PgRepository {
	return &PgRepository{
		db,
		log,
	}
}
