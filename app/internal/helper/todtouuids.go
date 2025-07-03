package helper

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToDTOUUIDs(ids []pgtype.UUID) []uuid.UUID {
	gUUIDs := make([]uuid.UUID, len(ids))
	for i, c := range ids {
		gUUIDs[i] = c.Bytes
	}
	return gUUIDs
}
