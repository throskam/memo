package lib

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func fromPGUUID(ID pgtype.UUID) uuid.NullUUID {
	return uuid.NullUUID{
		UUID:  ID.Bytes,
		Valid: ID.Valid,
	}
}

func toPGUUID(ID uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: ID,
		Valid: true,
	}
}

func toNullablePGUUID(ID uuid.NullUUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: ID.UUID,
		Valid: ID.Valid,
	}
}
