package bar

import (
	"time"

	"github.com/apirator/apirator/api/v1alpha1"
)

type Quux struct {
	Id     string `json:"Id"`
	Name   string `json:"Name"`
	Quuz   Quuz   `json:"Quuz"`
	Status string `json:"Status"`
}

type Quuz struct {
	Id          string               `json:"Id"`
	Name        string               `json:"Name"`
	APIMockSpec v1alpha1.APIMockSpec `json:"APIMockSpec"`
	Created     time.Time            `json:"Created"`
}
