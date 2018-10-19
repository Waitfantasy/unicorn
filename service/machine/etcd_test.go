package machine

import (
	"go.etcd.io/etcd/clientv3"
	"log"
	"strconv"
	"testing"
)

var ip = "1.0.0.1"

func newEtcdMachine() *EtcdMachine {
	e, err := NewEtcdMachine(clientv3.Config{
		Endpoints: []string{
			"192.168.10.10:2379",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return e
}

func Test_init(t *testing.T)  {
	e := newEtcdMachine()

	defer func() {
		if err := e.Close(); err != nil {
			t.Error(err)
		}
	}()

	if err := e.putSlots(newSlots()); err != nil {
		t.Error(err)
	}
}

func TestEtcdMachine_Put(t *testing.T) {
	e := newEtcdMachine()

	defer func() {
		if err := e.Close(); err != nil {
			t.Error(err)
		}
	}()

	item, err := e.Put(ip)
	if err != nil {
		t.Error(err)
		return
	}
	if item.Ip != ip {
		t.Errorf("put ip is %s, but the item ip is %s", ip, item.Ip)
		return
	}

	if item.Id != 1 {
		t.Errorf("the item id is not 1")
		return
	}

	slots, err := e.getSlots()
	if err != nil {
		t.Error(err)
		return
	}
	if slots.Free[1] != 1 {
		t.Errorf("the Free slots calculate error")
		return
	}

	for i:= 2; i < maxFreeSlots; i++ {
		if slots.Free[i] != 0 {
			t.Errorf("the Free slots calculate error")
			return
		}
	}

	if slots.Use != 1 {
		t.Errorf("the slots Use id calculate error")
	}
}

func TestEtcdMachine_Get(t *testing.T) {
	e := newEtcdMachine()

	defer func() {
		if err := e.Close(); err != nil {
			t.Error(err)
		}
	}()

	item, err := e.Get(ip)
	if err != nil {
		t.Error(err)
		return
	}

	if item == nil {
		return
	}

	if item.Ip != ip {
		t.Errorf("put ip is %s, but the item ip is %s", ip, item.Ip)
		return
	}

	if item.Id != 1 {
		t.Errorf("the item id is not 1")
		return
	}
}

func TestEtcdMachine_Reset(t *testing.T) {
	e := newEtcdMachine()

	defer func() {
		if err := e.Close(); err != nil {
			t.Error(err)
		}
	}()

	if err := e.Reset(ip, "1.0.0.2"); err != nil {
		t.Error(err)
	}

	ip = "1.0.0.2"

	item, err := e.Get(ip)
	if err != nil {
		t.Error(err)
		return
	}

	if item.Ip != ip {
		t.Errorf("put ip is %s, but the item ip is %s", ip, item.Ip)
		return
	}

	if item.Id != 1 {
		t.Errorf("the item id is not 1")
		return
	}
}

func TestEtcdMachine_Del(t *testing.T) {
	e := newEtcdMachine()

	defer func() {
		if err := e.Close(); err != nil {
			t.Error(err)
		}
	}()

	item, err := e.Del(ip)
	if err != nil {
		t.Error(err)
		return
	}

	if item.Ip != ip {
		t.Errorf("del ip is %s, but the item ip is %s", ip, item.Ip)
		return
	}

	if item.Id != 1 {
		t.Errorf("the item id is not 1")
		return
	}

	if item, err = e.Get(ip); !(item == nil && err == nil) {
		t.Errorf("del item error")
		return
	}
}

func TestEtcdMachine_All(t *testing.T) {
	e := newEtcdMachine()

	defer func() {
		if err := e.Close(); err != nil {
			t.Error(err)
		}
	}()

	for i:= 1; i < maxFreeSlots; i++ {
		ip := "10.10.10." + strconv.Itoa(i)
		item, err := e.Put(ip)
		if err != nil {
			t.Error(err)
			return
		}
		if item.Ip != ip {
			t.Errorf("put ip is %s, get item ip is %s", ip, item.Ip)
			return
		}

		if item.Id != i {
			t.Errorf("slots free calculate error")
			return
		}
	}

	slots, err := e.getSlots()
	if err != nil {
		t.Error(err)
		return
	}

	if slots.Use != maxFreeSlots - 1 {
		t.Error("slots use calculate error")
		return
	}

	if index := slots.findFreeIndex(); index != 0 {
		t.Error("slots free calculate error")
		return
	}


	for i:= 1; i < maxFreeSlots; i++ {
		ip := "10.10.10." + strconv.Itoa(i)
		item, err := e.Del(ip)
		if err != nil {
			t.Error(err)
			return
		}
		if item.Ip != ip {
			t.Errorf("put ip is %s, get item ip is %s", ip, item.Ip)
			return
		}

		if item.Id != i {
			t.Errorf("slots free calculate error")
			return
		}
	}
}