// conrol device or miner
package ctrl

import (
	"log"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/godbus/dbus/v5"
)

type RGBColor struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

func RebootDevice() error {
	log.Println("to reboot device")
	go func() { exec.Command("reboot").Run() }()
	log.Println("sent reboot device cmd")
	return nil
}

var lightBlinking int32 = 0

func SetDeviceLightBlink(durSec uint8) error {
	colorA := RGBColor{255, 0, 0}
	colorB := RGBColor{0, 0, 255}
	interval := time.Millisecond * 200
	count := int(((time.Second * time.Duration(durSec)) / interval).Nanoseconds())

	cColor, err := DeviceLightColor()
	if err != nil {
		return err
	}
	go func() {
		canRun := atomic.CompareAndSwapInt32(&lightBlinking, 0, 1)
		if !canRun {
			log.Println("give up setting light blinking, cause it is blinking")
			return
		}
		defer func() {
			atomic.StoreInt32(&lightBlinking, 0)
			SetDeviceLightColor(*cColor)
		}()
		var err error
		for i := 1; i <= count; i++ {
			if i%2 == 0 {
				err = SetDeviceLightColor(colorA)
			} else {
				err = SetDeviceLightColor(colorB)
			}
			if err != nil {
				log.Println("set device light color error: ", err)
				break
			}
			time.Sleep(interval)
		}
	}()

	return nil
}

func SetDeviceLightColor(c RGBColor) error {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return err
	}
	defer conn.Close()
	obj := conn.Object("org.hiot.led", "/org/hiot/led")
	call := obj.Call("org.hiot.led.SetColor", 0, c.R, c.G, c.B)
	if call.Err != nil {
		return call.Err
	} else {
		return nil
	}
}

func DeviceLightColor() (*RGBColor, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	obj := conn.Object("org.hiot.led", "/org/hiot/led")
	va, err := obj.GetProperty("org.hiot.led.Color")
	v := va.Value().([]interface{})
	c := RGBColor{
		R: v[0].(uint8),
		G: v[1].(uint8),
		B: v[2].(uint8),
	}
	log.Println("get device color:", c)
	return &c, nil
}
