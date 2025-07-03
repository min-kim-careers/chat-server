package helper

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToDTOUUID(id pgtype.UUID) uuid.UUID {
	return id.Bytes
}
