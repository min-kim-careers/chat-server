package helper

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToDBUUIDs(ids []uuid.UUID) []pgtype.UUID {
	pgUUIDs := make([]pgtype.UUID, len(ids))
	for i, u := range ids {
		pgUUIDs[i] = pgtype.UUID{Bytes: u, Valid: true}
	}
	return pgUUIDs
}
