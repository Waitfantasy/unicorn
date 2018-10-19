package machine

import "encoding/json"

type Item struct {
	Key           string
	Id            int
	Ip            string
	LastTimestamp uint64
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
