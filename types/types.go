package types

import "time"

//easyjson:json
type DosingState struct {
	LastDose time.Time `json:"lastDose"`
}
