package devdis

// device discovery

import (
	"fmt"
	"net"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/holoplot/go-avahi"
	"github.com/kpango/glg"
	"xdt.com/hm-diag/config"
)

type Dev struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

type DevDiscovery struct {
	Services          map[string]avahi.Service
	netInterfaceIndex int32
}

const (
	service = "_hummingbird._tcp"
)

var dis *DevDiscovery = nil

func Discovery() *DevDiscovery {
	return dis
}

func Services() []Dev {
	res := make([]Dev, 0, len(dis.Services))
	for _, v := range dis.Services {
		res = append(res,
			Dev{
				Name:    v.Name,
				Host:    v.Host,
				Address: v.Address,
				Port:    v.Port,
			})
	}
	return res
}

func Init() error {
	if dis == nil {
		// init
		dis = &DevDiscovery{Services: make(map[string]avahi.Service)}
		l, err := net.Interfaces()
		if err != nil {
			glg.Error("error getting interfaces", err)
			return err
		}
		for i, itf := range l {
			glg.Info("net interface", i+1, itf.Name)
			if itf.Name == config.Config().LanDevIntface {
				dis.netInterfaceIndex = int32(i) + 1
				glg.Debug("interface index:", i+1)
			}
		}

		// starts
		go func() {
			err := dis.Register()
			if err != nil {
				glg.Error("Discovery register error >>>>>>", err)
			} else {
				glg.Info("Discovery register finished")
			}
		}()
		go func() {
			err := dis.Browse()
			if err != nil {
				glg.Error("Discovery browse error >>>>>>", err)
			} else {
				glg.Info("Discovery browse end")
			}
		}()
		glg.Info("mdns inited discovery")

		// retry register after some time
		// when hostname is duplicated, avahi will retry different hostname
		// assume avahi has got a stable hostname after some time
		time.Sleep(time.Minute * 5)
		glg.Info("Discovery register retry")
		err = dis.Register()
		if err != nil {
			glg.Error("Discovery register retry error >>>>>>", err)
		} else {
			glg.Info("Discovery register retry finished")
		}

	} else {
		glg.Info("do not init device discovery repeadly")
	}
	return nil
}

func (d *DevDiscovery) Register() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		glg.Errorf("Cannot get system bus: %v\n", err)
		return err
	}

	a, err := avahi.ServerNew(conn)
	if err != nil {
		glg.Errorf("Avahi new failed: %v\n", err)
		return err
	}

	eg, err := a.EntryGroupNew()
	if err != nil {
		glg.Errorf("EntryGroupNew() failed: %v\n", err)
		return err
	}

	hostname, err := a.GetHostName()
	if err != nil {
		glg.Errorf("GetHostName() failed: %v\n", err)
		return err
	}

	fqdn, err := a.GetHostNameFqdn()
	if err != nil {
		glg.Errorf("GetHostNameFqdn() failed: %v\n", err)
		return err
	}
	txt := [][]byte{[]byte("cap=hm-diag"), []byte("other=xxx")}
	glg.Infof("Discovery registering hostname:%s service:%s fqdn:%s txt:%s\n",
		hostname, service, fqdn, txt)
	err = eg.AddService(avahi.InterfaceUnspec, avahi.ProtoInet, 0, hostname, service, "local", fqdn, 80, txt)
	if err != nil {
		glg.Errorf("AddService() failed: %v\n", err)
		return err
	}

	err = eg.Commit()
	if err != nil {
		glg.Errorf("Commit() failed: %v\n", err)
		return err
	}

	glg.Info("Discovery Entry published.")

	return nil
}

func (d *DevDiscovery) Browse() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		glg.Errorf("Cannot get system bus: %v\n", err)
		return err
	}

	server, err := avahi.ServerNew(conn)
	if err != nil {
		glg.Errorf("Avahi new failed: %v\n", err)
		return err
	}

	host, err := server.GetHostName()
	if err != nil {
		glg.Errorf("GetHostName() failed: %v\n", err)
		return err
	}
	glg.Info("GetHostName()", host)

	fqdn, err := server.GetHostNameFqdn()
	if err != nil {
		glg.Errorf("GetHostNameFqdn() failed: %v\n", err)
		return err
	}
	glg.Info("GetHostNameFqdn()", fqdn)

	s, err := server.GetAlternativeHostName(host)
	if err != nil {
		glg.Errorf("GetAlternativeHostName() failed: %v\n", err)
		return err
	}
	glg.Info("GetAlternativeHostName()", s)

	i, err := server.GetAPIVersion()
	if err != nil {
		glg.Errorf("GetAPIVersion() failed: %v\n", err)
		return err
	}
	glg.Info("GetAPIVersion()", i)

	sb, err := server.ServiceBrowserNew(avahi.InterfaceUnspec, avahi.ProtoUnspec, service, "local", 0)
	if err != nil {
		glg.Errorf("ServiceBrowserNew() failed: %v\n", err)
		return err
	}

	var service avahi.Service

	for {
		select {
		case service = <-sb.AddChannel:
			glg.Infof("Discovery service NEW: %#v\n", service)
			glg.Infof("Discovery net interface index compare  service=expect %d=%d\n",
				service.Interface, d.netInterfaceIndex)
			if service.Interface == d.netInterfaceIndex {
				service, err := server.ResolveService(service.Interface, service.Protocol, service.Name,
					service.Type, service.Domain, avahi.ProtoUnspec, 0)
				key := fmt.Sprintf("%s.%s.%s", service.Name, service.Type, service.Domain)
				if err == nil {
					glg.Info("Discovery service RESOLVED >>", service.Address)
					// remove the service which has the same address
					for k, v := range d.Services {
						if v.Address == service.Address {
							delete(d.Services, k)
						}
					}
					d.Services[key] = service
					glg.Infof("Discovered service ADDED: %#v\n", service)
				} else {
					glg.Error("Discovered service RESOLVE ERROR:", err)
				}
			}
		case service = <-sb.RemoveChannel:
			glg.Infof("Discovery sevice REMOVE: ", service)
			glg.Infof("Discovery net interface index compare  service=expect %d=%d\n",
				service.Interface, d.netInterfaceIndex)
			if service.Interface == d.netInterfaceIndex {
				key := fmt.Sprintf("%s.%s.%s", service.Name, service.Type, service.Domain)
				delete(d.Services, key)
				glg.Infof("Discovery service REMOVED: ", service)
			}
		}
	}

	// TODO: retry
}
