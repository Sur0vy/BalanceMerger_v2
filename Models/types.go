package Models

import (
	"path/filepath"
	"time"
)

type Sources struct {
	Balance string
	Journal string
	Card    string
}

func (s *Sources) GetOutFileName(onlyFileName bool) string {
	var ret string
	if !onlyFileName {
		ret = filepath.Dir(s.Journal) + "/"
	}
	ret = ret + "out_" + time.Now().Format("02-01-2006_15-04") + ".xlsx"
	return ret
}
