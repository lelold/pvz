package model

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PVZID    uuid.UUID `json:"pvzId"`
	Status   string    `json:"status"`
}
