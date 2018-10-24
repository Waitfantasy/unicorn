package machine

import "encoding/json"

type Item struct {
	Key           string `json:"key"`
	Id            int    `json:"id"`
	Ip            string `json:"ip"`
	LastTimestamp uint64 `json:"last_timestamp"`
}

func JsonUnmarshalItem(data []byte) (*Item, error) {
	item := &Item{}
	if err := json.Unmarshal(data, item); err != nil {
		return nil, err
	}
	return item, nil
}

func JsonMarshalItem(item *Item) ([]byte, error) {
	return json.Marshal(item)
}
