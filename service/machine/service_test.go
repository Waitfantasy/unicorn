package machine

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"strconv"
	"testing"
)

func TestConnection(t *testing.T) {
	service := newService()
	err := service.EtcdConnection()
	if err != nil {
		t.Error(err)
	}
}

func TestService_GetExtraData(t *testing.T) {
	service := newService()
	service.EtcdConnection()
	extraData, err := service.GetExtraData()
	if err != nil {
		t.Error(err)
	} else {
		if extraData.FreeSlots[0] != MaxId {
			t.Error("get machine data free slots first slot error")
		}

		if len(extraData.FreeSlots) - 1 != MaxId {
			t.Error("get machine data free slots length error")
		}
	}
}

func TestService_PutMachineItem(t *testing.T) {
	service := newService()
	err := service.EtcdConnection()
	if err != nil {
		t.Error(err)
	}


	for i := 1; i <= 1024; i++ {
		_, err := service.PutMachineItem("1.0.0." + strconv.Itoa(i))
		if err != nil {
			t.Error(err)
			return
		}
	}
	extraData, _ := service.GetExtraData()
	for i := 1; i <= 1024; i++ {
		if extraData.FreeSlots[i] != 1 {
			t.Error("TestService_PutMachineNode free slots error")
			return
		}
	}

	if extraData.LastId != 1024 {
		t.Error("TestService_PutMachineNode last id error")
		return
	}

	_, err = service.PutMachineItem("1.0.0.1025")
	if err == nil {
		t.Error("TestService_PutMachineNode free slot full error")
		return
	}

	if item, err := service.DelMachineItem("1.0.0.10"); err != nil {
		t.Error(err)
		return
	} else {
		extraData, _ := service.GetExtraData()
		if extraData.FreeSlots[item.Id] != 0 {
			t.Error("TestService_PutMachineNode del not after free slot update error")
		}
	}

	_, err = service.PutMachineItem("192.168.10.10")

	if err != nil {
		t.Error("TestService_PutMachineNode slot free error")
		return
	}

	_, err = service.PutMachineItem("1.0.0.1025")
	if err == nil {
		t.Error("TestService_PutMachineNode free slot full error")
		return
	}

}

func TestService_GetMachineItemList(t *testing.T) {
	service := newService()
	service.EtcdConnection()
	items, err := service.GetMachineItemList()
	if err != nil {
		t.Error(err)
	} else  {
		for _, item := range items {
			fmt.Println(item)
		}
	}
}


func newService() *Service {
	return NewService(clientv3.Config{
		Endpoints: []string{
			"192.168.10.10:2379",
		},
	})
}
