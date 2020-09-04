package tzdb

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/moisespsena-go/assetfs/assetfsapi"
	iocommon "github.com/moisespsena-go/io-common"
)

const ZonesZipFile = "golib/time/zoneinfo.zip"

var (
	Db DB
)

type DB struct {
	mu     sync.Mutex
	ByName map[string]*time.Location
	All    []*time.Location
}

func (d *DB) Load(fs assetfsapi.Interface) (err error) {
	if d.ByName != nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.ByName != nil {
		return
	}

	d.ByName = map[string]*time.Location{}

	var (
		asset assetfsapi.FileInfo
		r     io.ReadCloser
		zr    *zip.Reader
	)
	if asset, err = fs.AssetInfo(ZonesZipFile); err != nil {
		return err
	}

	if r, err = asset.Reader(); err != nil {
		return err
	}
	defer r.Close()

	if zr, err = zip.NewReader(&zipreader{r.(iocommon.ReadSeekCloser)}, asset.Size()); err != nil {
		return
	}

	var loc *time.Location

	for _, file := range zr.File {
		if !file.FileInfo().IsDir() {
			if loc, err = loadFile(file); err != nil {
				return
			}
			d.All = append(d.All, loc)
			d.ByName[file.Name] = loc
		}
	}
	return
}

func (d *DB) Names() (names []string) {
	for _, loc := range d.All {
		names = append(names, loc.String())
	}
	return
}

func (d *DB) Pair() (pairs [][]string) {
	for _, loc := range d.All {
		pairs = append(pairs, []string{loc.String(), LocationCity(loc.String()).Label()})
	}
	return
}

func (d *DB) Structs() (pairs []struct {
	Value, Label string
	Loc          *time.Location
}) {
	for _, loc := range d.All {
		pairs = append(pairs, struct {
			Value, Label string
			Loc          *time.Location
		}{loc.String(), LocationCity(loc.String()).Label(), loc})
	}
	return
}

func (d *DB) Labels() (s []string) {
	for _, loc := range d.All {
		s = append(s, LocationCity(loc.String()).Label())
	}
	return
}

func (d *DB) Strings() (s []string) {
	for _, loc := range d.All {
		s = append(s, loc.String())
	}
	return
}

type zipreader struct {
	iocommon.ReadSeekCloser
}

func (r *zipreader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("negative offset")
	}

	if _, err = r.Seek(off, io.SeekStart); err != nil {
		return
	}

	return r.Read(p)
}

func loadFile(f *zip.File) (loc *time.Location, err error) {
	var (
		r    io.ReadCloser
		data []byte
	)

	if r, err = f.Open(); err != nil {
		return
	}

	defer r.Close()

	if data, err = ioutil.ReadAll(r); err != nil {
		return
	}

	return time.LoadLocationFromTZData(f.Name, data)
}
