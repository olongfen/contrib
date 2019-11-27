package middle

import "encoding/json"

// JSONMarshalMust must to  marshal json
func JSONMarshalMust(v interface{})(ret []byte)  {
	var(
		err error
	)
	if ret ,err = json.Marshal(v);err==nil{
		return
	}else {
		ret =[]byte("{}")
	}
	return
}
