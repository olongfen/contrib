package contrib

import (
	"encoding/json"
	"errors"
	"fmt"
)

// DataBody
type DataBody map[string]interface{}

// GetValueByString
func (d DataBody) GetValueByString(key string) (ret string, err error) {
	if val, ok := d[key]; !ok {
		err = fmt.Errorf(fmt.Sprintf(`"map dose't have key: %s"`, key))
		return
	} else {
		switch v := val.(type) {
		case string:
			ret = v
		default:
			err = fmt.Errorf(fmt.Sprintf(`%v type is %T`, key, val))
			return
		}
	}

	return
}

// GetValueByInt
func (d DataBody) GetValueByInt(key string) (ret int64, err error) {
	if val, ok := d[key]; !ok {
		err = fmt.Errorf(fmt.Sprintf(`"map dose't have key: %s"`, key))
		return
	} else {
		switch v := val.(type) {
		case int:
			ret = int64(v)
		case int32:
			ret = int64(v)
		case int8:
			ret = int64(v)
		case int16:
			ret = int64(v)
		case int64:
			ret = int64(v)
		default:
			err = fmt.Errorf(fmt.Sprintf(`%v type is %T`, key, val))
			return
		}
	}

	return
}

// GetValueByBool
func (d DataBody) GetValueByBool(key string) (ret bool, err error) {
	if val, ok := d[key]; !ok {
		err = errors.New(fmt.Sprintf(`"map dose't have key: %s"`, key))
		return
	} else {
		switch v := val.(type) {
		case bool:
			ret = v
		default:
			err = errors.New(fmt.Sprintf(`%v type is %T`, val, val))
			return
		}
	}

	return
}

// GetValueByFloat
func (d DataBody) GetValueByFloat(key string) (ret float64, err error) {
	if val, ok := d[key]; !ok {
		err = fmt.Errorf(fmt.Sprintf(`"map dose't have key: %s"`, key))
		return
	} else {
		switch v := val.(type) {
		case float64:
			ret = v
		case float32:
			ret = float64(v)
		default:
			err = fmt.Errorf(fmt.Sprintf(`%v type is %T`, key, val))
			return
		}
	}

	return
}

// GetValueByMAP
func (d DataBody) GetValueByMAP(key string) (ret DataBody, err error) {
	if val, ok := d[key]; !ok {
		err = fmt.Errorf(fmt.Sprintf(`"map dose't have key: %s"`, key))
		return
	} else {
		switch v := val.(type) {
		case map[string]interface{}:
			ret = v
		default:
			err = fmt.Errorf(fmt.Sprintf(`%v type is %T`, key, val))
			return
		}
	}

	return
}

func (d DataBody) MarshalToBody(val interface{}) (err error) {
	var (
		data []byte
	)
	if data, err = json.Marshal(val); err != nil {
		return
	}
	return json.Unmarshal(data, &d)
}

func (d DataBody) Unmarshal(val interface{}) (err error) {
	var (
		body []byte
	)
	if body, err = json.Marshal(d); err != nil {
		return
	}
	return json.Unmarshal(body, val)
}

func (d DataBody) JSONUnmarshalByKey(key string, val interface{}) (err error) {
	var (
		body []byte
	)
	if body, err = json.Marshal(d[key]); err != nil {
		return
	}
	return json.Unmarshal(body, val)
}
