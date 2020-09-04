package tzdb

import (
	"github.com/moisespsena-go/time-helpers"
	"strings"
	"time"
)

type LocationCity string

func (this LocationCity) Region() string {
	pos := strings.IndexRune(string(this), '/')
	return string(this)[pos+1:]
}

func (this LocationCity) City() string {
	pos := strings.IndexRune(string(this), '/')
	return string(this)[0:pos]
}

func (this LocationCity) Location() *time.Location {
	if this == "" {
		return time.Local
	}
	return Db.ByName[string(this)]
}

func (this LocationCity) String() string {
	return string(this)
}

func (this LocationCity) Label() string {
	if this == "" {
		return ""
	}
	return string(this) + " (" + time_helpers.LocToGmtS(this.Location()) + ")"
}
