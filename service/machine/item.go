package machine

import "encoding/json"

type Item struct {
	Id            int
	Ip            string
	LastTimestamp uint64
}

func jsonUnmarshalItem(data []byte) (*Item, error) {
	item := &Item{}
	if err := json.Unmarshal(data, item); err != nil {
		return nil, err
	}
	return item, nil
}

func jsonMarshalItem(item *Item) ([]byte, error) {
	return json.Marshal(item)
}