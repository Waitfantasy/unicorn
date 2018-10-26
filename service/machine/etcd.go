package machine

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"hash/crc32"
	"strconv"
	"time"
)

type EtcdMachine struct {
	prefixKey string
	slotsKey  string
	timeout   int
	cli       *clientv3.Client
}

func NewEtcdMachine(cfg clientv3.Config, timeout int) (*EtcdMachine, error) {
	clientv3.New(cfg)
	if cli, err := clientv3.New(cfg); err != nil {
		return nil, err
	} else {
		return &EtcdMachine{
			prefixKey: "/unicorn_machine_items/",
			slotsKey:  "unicorn_machine_slots",
			timeout:   timeout,
			cli:       cli,
		}, nil
	}
}

func (e *EtcdMachine) Close() error {
	return e.cli.Close()
}

func (e *EtcdMachine) getSlots() (*slots, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(e.timeout)*time.Second)
	res, err := e.cli.Get(ctx, e.slotsKey)
	if err != nil {
		return nil, err
	}

	for _, kv := range res.Kvs {
		if string(kv.Key) == e.slotsKey {
			return jsonUnmarshalSlots(kv.Value)
		}
	}

	slots := newSlots()
	return slots, nil
}

func (e *EtcdMachine) putSlots(s *slots) error {
	b, err := jsonMarshalSlots(s)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(e.timeout)*time.Second)
	_, err = e.cli.Put(ctx, e.slotsKey, string(b))
	return err
}

func (e *EtcdMachine) All() ([]*Item, error) {
	items := make([]*Item, 0)

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(e.timeout)*time.Second)
	res, err := e.cli.Get(ctx, e.prefixKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range res.Kvs {
		if item, err := JsonUnmarshalItem(kv.Value); err == nil {
			item.Key = string(kv.Key)
			items = append(items, item)
		}
	}

	return items, nil
}

func (e *EtcdMachine) Get(ip string) (*Item, error) {
	key := e.key(ip)
	return e.get(key)
}

func (e *EtcdMachine) get(key string) (*Item, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(e.timeout)*time.Second)
	res, err := e.cli.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	for _, kv := range res.Kvs {
		if string(kv.Key) == key {
			return JsonUnmarshalItem(kv.Value)
		}
	}

	return nil, nil
}

func (e *EtcdMachine) Put(ip string) (*Item, error) {
	key := e.key(ip)
	item, err := e.get(key)
	if err != nil {
		return nil, err
	}

	if item != nil {
		return nil, fmt.Errorf("the machine ip %s already exists in verify", ip)
	}

	slots, err := e.getSlots()
	if err != nil {
		return nil, err
	}

	if slots.Use > MaxMachine {
		return nil, errors.New("no machine id available in slots")
	}

	index := slots.findFreeIndex()
	item = &Item{
		Key: key,
		Id:  index,
		Ip:  ip,
	}

	if err = e.PutItem(item); err != nil {
		return nil, err
	}

	slots.Use++
	slots.Free[index] = 1
	for ; ; {
		if err = e.putSlots(slots); err == nil {
			break
		}
	}
	return item, nil
}

func (e *EtcdMachine) PutItem(item *Item) error {
	b, err := JsonMarshalItem(item)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(e.timeout)*time.Second)
	_, err = e.cli.Put(ctx, item.Key, string(b))
	return err
}

func (e *EtcdMachine) Del(ip string) (*Item, error) {
	key := e.key(ip)
	item, err := e.get(key)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, fmt.Errorf("no machine item can be delete by ip: %s", ip)
	}

	slots, err := e.getSlots()
	if err != nil {
		return nil, err
	}

	slots.Free[item.Id] = 0
	slots.Use--
	if err = e.putSlots(slots); err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(e.timeout)*time.Second)
	if _, err = e.cli.Delete(ctx, key); err != nil {
		return nil, err
	}

	return item, nil
}

func (e *EtcdMachine) Reset(oldIp, newIp string) error {
	if oldIp == newIp {
		return nil
	}

	newKey := e.key(newIp)
	item, err := e.get(newKey)
	if err != nil {
		return err
	}

	if item != nil {
		return fmt.Errorf("the machine ip %s already exists in verify", newIp)
	}

	oldKey := e.key(oldIp)
	if item, err = e.get(oldKey); err != nil {
		return err
	}

	if item == nil {
		return fmt.Errorf("the machine ip %s not exists in etcd", oldIp)
	}

	// delete old machine by old key
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(e.timeout)*time.Second)
	if _, err = e.cli.Delete(ctx, oldKey); err != nil {
		return err
	}

	item.Key = newKey
	item.Ip = newIp
	if err = e.PutItem(item); err != nil {
		return err
	}

	return nil
}

func (e *EtcdMachine) key(ip string) string {
	uint32Key := crc32.ChecksumIEEE([]byte(ip))
	return e.prefixKey + strconv.Itoa(int(uint32Key))
}
