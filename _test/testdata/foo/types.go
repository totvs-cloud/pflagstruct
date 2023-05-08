package foo

import (
	"time"

	"github.com/angelokurtis/flagstruct/test/testdata/bar"
)

type Baz struct {
	Id      string    `json:"Id"`
	Name    string    `json:"Name"`
	Quux    *bar.Quux `json:"Quux"`
	Status  string    `json:"Status"`
	Created time.Time `json:"Created"`
	Error   string    `json:"Error,omitempty"`
	Tags    []*Tag    `json:"Tags,omitempty"`
}

type Tag struct {
	Key   string `json:"Id"`
	Value string `json:"Name"`
}

type Corge bar.Quux

type Grault Baz
