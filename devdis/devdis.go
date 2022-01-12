package devdis

// device discovery

import (
	"fmt"
	"log"
	"net"

	"github.com/godbus/dbus/v5"
	"github.com/holoplot/go-avahi"
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
			log.Println("error getting interfaces", err)
			return err
		}
		for i, itf := range l {
			log.Println("net interface", i+1, itf.Name)
			if itf.Name == config.Config().LanDevIntface {
				// TODO: check if avahi interface index is correct
				dis.netInterfaceIndex = int32(i) + 1
				log.Println("interface index:", i+1)
			}
		}

		// starts
		go func() {
			err := dis.Register()
			if err != nil {
				log.Println("Discovery register error", err)
			}
		}()
		go func() {
			err := dis.Browse()
			if err != nil {
				log.Println("Discovery browse error", err)
			}
		}()
		log.Println("mdns inited discovery")
	} else {
		log.Println("do not init device discovery repeadly")
	}
	return nil
}

func (d *DevDiscovery) Register() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Printf("Cannot get system bus: %v\n", err)
		return err
	}

	a, err := avahi.ServerNew(conn)
	if err != nil {
		log.Printf("Avahi new failed: %v\n", err)
		return err
	}

	eg, err := a.EntryGroupNew()
	if err != nil {
		log.Printf("EntryGroupNew() failed: %v\n", err)
		return err
	}

	hostname, err := a.GetHostName()
	if err != nil {
		log.Printf("GetHostName() failed: %v\n", err)
		return err
	}

	fqdn, err := a.GetHostNameFqdn()
	if err != nil {
		log.Printf("GetHostNameFqdn() failed: %v\n", err)
		return err
	}
	txt := [][]byte{[]byte("cap=hm-diag"), []byte("other=xxx")}
	err = eg.AddService(avahi.InterfaceUnspec, avahi.ProtoInet, 0, hostname, service, "local", fqdn, 80, txt)
	if err != nil {
		log.Printf("AddService() failed: %v\n", err)
		return err
	}

	err = eg.Commit()
	if err != nil {
		log.Printf("Commit() failed: %v\n", err)
		return err
	}

	log.Println("Entry published. Hit ^C to exit.")

	for {
		select {}
	}
	// TODO: retry
}

func (d *DevDiscovery) Browse() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Printf("Cannot get system bus: %v\n", err)
		return err
	}

	server, err := avahi.ServerNew(conn)
	if err != nil {
		log.Printf("Avahi new failed: %v\n", err)
		return err
	}

	host, err := server.GetHostName()
	if err != nil {
		log.Printf("GetHostName() failed: %v\n", err)
		return err
	}
	log.Println("GetHostName()", host)

	fqdn, err := server.GetHostNameFqdn()
	if err != nil {
		log.Printf("GetHostNameFqdn() failed: %v\n", err)
		return err
	}
	log.Println("GetHostNameFqdn()", fqdn)

	s, err := server.GetAlternativeHostName(host)
	if err != nil {
		log.Printf("GetAlternativeHostName() failed: %v\n", err)
		return err
	}
	log.Println("GetAlternativeHostName()", s)

	i, err := server.GetAPIVersion()
	if err != nil {
		log.Printf("GetAPIVersion() failed: %v\n", err)
		return err
	}
	log.Println("GetAPIVersion()", i)

	sb, err := server.ServiceBrowserNew(avahi.InterfaceUnspec, avahi.ProtoUnspec, service, "local", 0)
	if err != nil {
		log.Printf("ServiceBrowserNew() failed: %v\n", err)
		return err
	}

	var service avahi.Service

	for {
		select {
		case service = <-sb.AddChannel:
			log.Printf("Discovery service NEW: %#v\n", service)
			service, err := server.ResolveService(service.Interface, service.Protocol, service.Name,
				service.Type, service.Domain, avahi.ProtoUnspec, 0)
			if err == nil {
				log.Println("Discovery service RESOLVED >>", service.Address)
				if service.Interface == d.netInterfaceIndex {
					key := fmt.Sprintf("%s.%s.%s", service.Name, service.Type, service.Domain)
					d.Services[key] = service
					log.Printf("Discovered service ADDED: %#v\n", service)
				}
			}
		case service = <-sb.RemoveChannel:
			log.Println("Discovery sevice REMOVE: ", service)
			if service.Interface == d.netInterfaceIndex {
				key := fmt.Sprintf("%s.%s.%s", service.Name, service.Type, service.Domain)
				delete(d.Services, key)
				log.Println("Discovery service REMOVED: ", service)
			}
		}
	}

	// TODO: retry
}
