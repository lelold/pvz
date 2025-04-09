package model

import (
	"time"
)

type PVZ struct {
	ID               string    `json:"id" format:"uuid"`
	RegistrationDate time.Time `json:"registrationDate" format:"date-time"`
	City             string    `json:"city" enum:"Москва,Санкт-Петербург,Казань"`
}

type Reception struct {
	ID       string    `json:"id" format:"uuid"`
	DateTime time.Time `json:"dateTime" format:"date-time"`
	PVZID    string    `json:"pvzId" format:"uuid"`
	Status   string    `json:"status" enum:"in_progress,close"`
}

type Product struct {
	ID          string    `json:"id" format:"uuid"`
	DateTime    time.Time `json:"dateTime" format:"date-time"`
	Type        string    `json:"type" enum:"электроника,одежда,обувь"`
	ReceptionID string    `json:"receptionId" format:"uuid"`
}
