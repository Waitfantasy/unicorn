package machine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"hash/crc32"
	"strconv"
	"time"
)

const (
	MinId         = 1
	MaxId         = 1 << 10
	FreeSlotCount = MaxId + 1
)

type machineExtraData struct {
	FreeSlots [FreeSlotCount]int `json:"free_slots"`
	LastId    int                `json:"last_id"`
}

func NewMachineExtraData() *machineExtraData {
	return &machineExtraData{
		LastId:    0,
		FreeSlots: [FreeSlotCount]int{0: 1024,},
	}
}

func (extra *machineExtraData) findSlotFreeIndex() int {
	for i := 1; i < MaxId; i++ {
		if extra.FreeSlots[i] == 0 {
			return i
		}
	}

	return 0
}

func (extra *machineExtraData) encode() ([]byte, error) {
	return json.Marshal(extra)
}

type machineItem struct {
	Id               int    `json:"id"`
	Ip               string `json:"ip"`
	Key              string `json:"-"`
	CreatedTimestamp int64  `json:"created_timestamp"`
	UpdatedTimestamp int64  `json:"updated_timestamp"`
}

func newMachineItem() *machineItem {
	ts := time.Now().Unix()
	return &machineItem{
		CreatedTimestamp: ts,
		UpdatedTimestamp: ts,
	}
}

func (item *machineItem) withIpId(ip string, id int) *machineItem {
	item.Ip = ip
	item.Id = id
	return item
}

func (item *machineItem) encode() ([]byte, error) {
	return json.Marshal(item)
}

func (item *machineItem) FormatCreatedTime() string {
	return time.Unix(item.CreatedTimestamp, 0).Format("2006-01-02 15:04:05")
}

func (item *machineItem) FormatUpdatedTime() string {
	return time.Unix(item.UpdatedTimestamp, 0).Format("2006-01-02 15:04:05")
}



type EtcdConfig struct {
	Endpoints []string
}

type Service struct {
	cfg       clientv3.Config
	cli       *clientv3.Client
	prefixKey string
	extraKey  string
}

func NewService(cfg clientv3.Config) *Service {
	return &Service{
		cfg:       cfg,
		prefixKey: "/unicorn_machine_items/",
		extraKey:  "unicorn_machine_extra",
	}
}

func (s *Service) EtcdConnection() error {
	if cli, err := clientv3.New(s.cfg); err != nil {
		return err
	} else {
		s.cli = cli
		return nil
	}
}

func (s *Service) GetMachineItemList() ([]*machineItem, error) {
	items := make([]*machineItem, 0 )
	res, err := s.cli.Get(context.Background(), s.prefixKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range res.Kvs {
		if item, err := s.extractMachineItem(kv.Value); err == nil {
			item.Key = string(kv.Key)
			items = append(items, item)
		}
	}

	return items, nil
}

func (s *Service) GetMachineItem(key string) (*machineItem, error) {
	res, err := s.cli.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	for _, kv := range res.Kvs {
		if string(kv.Key) == key {
			return s.extractMachineItem(kv.Value)
		}
	}

	return nil, nil
}

func (s *Service) PutMachineItem(ip string) (*machineItem, error) {
	key := s.MachineKey(ip)
	item, err := s.GetMachineItem(key)
	if err != nil {
		return nil, err
	}

	if item != nil {
		return nil, fmt.Errorf("the machine ip %s already exists in etcd", ip)
	}

	extraData, err := s.GetExtraData()
	if err != nil {
		return nil, err
	}

	if extraData.LastId >= MaxId {
		if index := extraData.findSlotFreeIndex(); index == 0 {
			return nil, errors.New("no machine id available in free slot")
		} else {
			item := newMachineItem().withIpId(ip, index)
			if err = s.putMachineItem(key, item); err != nil {
				return nil, err
			}

			extraData.FreeSlots[index] = 1
			for ; ; {
				if err = s.PutExtraData(extraData); err == nil {
					break
				}
			}
			return item, nil
		}
	} else {
		extraData.LastId++
		item := newMachineItem().withIpId(ip, extraData.LastId)
		if err = s.putMachineItem(key, item); err != nil {
			return nil, err
		}
		extraData.FreeSlots[extraData.LastId] = 1
		for ; ; {
			if err = s.PutExtraData(extraData); err == nil {
				break
			}
		}
		return item, nil
	}
}

func (s *Service) putMachineItem(key string, machineItem *machineItem) (error) {
	b, err := machineItem.encode()
	if err != nil {
		return err
	}

	_, err = s.cli.Put(context.Background(), key, string(b))
	return err
}

func (s *Service) DelMachineItem(ip string) (*machineItem, error) {
	key := s.MachineKey(ip)
	item, err := s.GetMachineItem(key)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, fmt.Errorf("no machine item can be delete by ip: %s", ip)
	}

	extraData, err := s.GetExtraData()
	if err != nil {
		return nil, err
	}

	extraData.FreeSlots[item.Id] = 0

	if err = s.PutExtraData(extraData); err != nil {
		return nil, err
	}

	if err = s.delMachineItem(key); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) delMachineItem(key string) error {
	_, err := s.cli.Delete(context.Background(), key)
	return err
}

func (s *Service) GetExtraData() (*machineExtraData, error) {
	res, err := s.cli.Get(context.Background(), s.extraKey)
	if err != nil {
		return nil, err
	}

	for _, kv := range res.Kvs {
		if string(kv.Key) == s.extraKey {
			return s.extractMachineExtraData(kv.Value)
		}
	}

	extraData := NewMachineExtraData()
	b, err := extraData.encode()
	if err != nil {
		return nil, err
	}

	if _, err = s.cli.Put(context.Background(), s.extraKey, string(b)); err != nil {
		return nil, err
	}

	return extraData, nil
}

func (s *Service) PutExtraData(extraData *machineExtraData) error {
	b, err := extraData.encode()
	if err != nil {
		return err
	}

	_, err = s.cli.Put(context.Background(), s.extraKey, string(b))
	return err
}

func (s *Service) MachineKey(ip string) string {
	uint32Key := crc32.ChecksumIEEE([]byte(ip))
	return s.prefixKey + strconv.Itoa(int(uint32Key))
}

func (s *Service) extractMachineItem(data []byte) (*machineItem, error) {
	machineItem := &machineItem{}
	if err := json.Unmarshal(data, machineItem); err != nil {
		return nil, err
	}
	return machineItem, nil
}

func (s *Service) extractMachineExtraData(data []byte) (*machineExtraData, error) {
	extraData := &machineExtraData{}
	if err := json.Unmarshal(data, extraData); err != nil {
		return nil, err
	}
	return extraData, nil
}
