package models

import "time"

type Vaccination struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Drug      string    `json:"drug"`
	DrugID    int32     `json:"drug_id"`
	Dose      int32     `json:"dose"`
	AppliedAt time.Time `json:"date"`
}
