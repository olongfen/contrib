package contrib

import (
	"encoding/json"
	"os"
	"strings"
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
	ID         uint `gorm:"primary_key；AUTO_INCREMENT"`
	CreateTime time.Time
	UpdateTime time.Time
	DeleteTime *time.Time `sql:"index"`
}

type TimeData struct {
	CreateTime time.Time
	UpdateTime time.Time
	DeleteTime *time.Time `sql:"index"`
}

// PubGetEnvString
func PubGetEnvString(key string, defaultValue string) (ret string) {
	ret = os.Getenv(key)
	if len(ret) == 0 {
		ret = defaultValue
	}
	return
}

// PubGetEnvBool
func PubGetEnvBool(key string, defaultValue bool) (ret bool) {
	val := strings.ToLower(os.Getenv(key))
	if val == "true" {
		ret = true
	} else if val == "false" {
		ret = false
	} else {
		ret = defaultValue
	}
	return
}
