package contrib

import (
	"encoding/json"
	"time"
)

// JSONMarshalMust must to  marshal json
func JSONMarshalMust(v interface{}) (ret []byte) {
	var (
		err error
	)
	if ret, err = json.Marshal(v); err == nil {
		return
	} else {
		ret = []byte("{}")
	}
	return
}

// ModelBase
type ModelBase struct {
	ID          uint `gorm:"unique_index；AUTO_INCREMENT"`
	CreatedTime time.Time
	UpdatedTime time.Time
	DeletedTime *time.Time `sql:"index"`
}
