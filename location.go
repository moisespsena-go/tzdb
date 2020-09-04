package tzdb

import (
	"time"
)

var Sys = NewLocationGetterFunc(func() *time.Location {
	return time.Local
})

type Location interface {
	Location() *time.Location
}

type locationGetter struct {
	loc *time.Location
}

func (this locationGetter) Location() *time.Location {
	return this.loc
}

func NewLocationGetter(loc *time.Location) Location {
	return locationGetter{loc}
}

type LocationGetterFunc func() *time.Location

func (this LocationGetterFunc) Location() *time.Location {
	return this()
}

func NewLocationGetterFunc(f func() *time.Location) Location {
	return LocationGetterFunc(f)
}
