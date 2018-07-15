package aoa

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/google/gousb"
)

type (
	Accessory struct {
		usbctx      *gousb.Context
		device      *gousb.Device
		AoAProtocol uint16
		mode        uint8

		accessory *AccessoryMode
		audio     *AudioMode
	}

	AccessoryMode struct {
		epIn  *gousb.InEndpoint
		epOut *gousb.OutEndpoint
		done  func()
	}

	AudioMode struct {
		epIn *gousb.InEndpoint
		done func()
	}
)

func NewAccessory() *Accessory {
	return &Accessory{
		usbctx: gousb.NewContext(),
	}
}

func (acc *Accessory) Close() {
	if acc.accessory != nil {
		acc.accessory.done()
	}

	if acc.audio != nil {
		acc.audio.done()
	}

	if acc.device != nil {
		acc.device.Close()
	}

	acc.usbctx.Close()
}

func (a *AccessoryMode) Read(b []byte) (int, error) {
	return a.epIn.Read(b)
}

func (a *AccessoryMode) Write(b []byte) (int, error) {
	return a.epOut.Write(b)
}

func (acc *Accessory) findAndroidDevice(onlyAccessory bool) (dev *gousb.Device, err error) {
	devs, err := acc.usbctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		if desc.Class == gousb.ClassHub {
			return false
		}

		// fmt.Printf("DescClass: %v: %+v, %+v, %v\n", desc, desc.Speed, desc.Class, desc.SubClass)
		if !onlyAccessory {
			return true
		}

		if desc.Vendor == ACCESSORY_VENDOR_ID {
			switch desc.Product {
			case PRODUCT_ID_ACCESSORY, PRODUCT_ID_ACCESSORY_ADB, PRODUCT_ID_AUDIO, PRODUCT_ID_AUDIO_ADB, PRODUCT_ID_ACCESSORY_AUDIO, PRODUCT_ID_ACCESSORY_AUDIO_ADB:
				return true
			}
		}

		return false
	})

	if err != nil {
		return nil, err
	}

	if len(devs) == 0 {
		return nil, ErrNoDevice
	}

	for _, d := range devs {
		d.ControlTimeout = time.Duration(time.Second * 2)
		protoVer, err := getProtocol(d)
		if err != nil || protoVer < 1 {
			d.Close()
			continue
		}
		fmt.Println("Found a accessory", d, protoVer, err)

		// Only use the first
		dev = d
		acc.AoAProtocol = protoVer
		break
	}

	return
}

func (acc *Accessory) SwitchToAccessoryMode(mode uint8) (err error) {
	dev, err := acc.findAndroidDevice(false)
	if err != nil {
		return err
	}
	acc.mode = mode

	// Switch to accessory mode
	fmt.Println("Switching a android to accessory mode", dev)
	if mode&MODE_ACCESSORY > 0 {
		sendString(dev, ACCESSORY_SEND_STRING, 0, ACCESSORY_STRING_MANUFACTURER, []byte("megvii\x00"))
		sendString(dev, ACCESSORY_SEND_STRING, 0, ACCESSORY_STRING_MODEL, []byte("WorkerStorage\x00"))
		sendString(dev, ACCESSORY_SEND_STRING, 0, ACCESSORY_STRING_DESCRIPTION, []byte("Beehive Worker Storage Service\x00"))
		sendString(dev, ACCESSORY_SEND_STRING, 0, ACCESSORY_STRING_VERSION, []byte("1.0.0\x00"))
		sendString(dev, ACCESSORY_SEND_STRING, 0, ACCESSORY_STRING_URI, []byte("https://collect.zzcrowd.com\x00"))
		sendString(dev, ACCESSORY_SEND_STRING, 0, ACCESSORY_STRING_SERIAL, []byte("1234567890\x00"))
	}

	if mode&MODE_AUDIO > 0 && acc.AoAProtocol > 1 {
		sendString(dev, ACCESSORY_SET_AUDIO_MODE, 1, 0, nil)
	}

	err = startAccessoryMode(dev)
	if err != nil {
		return err
	}

	// Switch accessory mode, android usb vendor_id and product_id will change
	dev.Close()
	time.Sleep(time.Second)

	dev, err = acc.findAndroidDevice(true)
	if err != nil {
		return err
	}
	dev.SetAutoDetach(true)

	acc.device = dev

	return nil
}

func (acc *Accessory) OpenAcessoryInterface() (*AccessoryMode, error) {
	if acc.device == nil {
		return nil, ErrNoDevice
	}

	if acc.mode&MODE_ACCESSORY == 0 {
		return nil, ErrorNotSupport
	}
	intf, done, err := acc.device.DefaultInterface()
	if err != nil {
		return nil, err
	}

	a := &AccessoryMode{done: done}

	for _, ei := range intf.Setting.Endpoints {
		if ei.Direction == gousb.EndpointDirectionIn {
			a.epIn, err = intf.InEndpoint(ei.Number)
			if err != nil {
				return nil, err
			}
		}
		if ei.Direction == gousb.EndpointDirectionOut {
			a.epOut, err = intf.OutEndpoint(ei.Number)
			if err != nil {
				return nil, err
			}
		}
	}

	acc.accessory = a

	return acc.accessory, nil
}

// libusb operation
func getProtocol(dev *gousb.Device) (uint16, error) {
	if dev == nil {
		return 0, ErrNoDevice
	}

	var data = make([]byte, 2)
	n, err := dev.Control(USB_DIR_IN|USB_TYPE_VENDOR, ACCESSORY_GET_PROTOCOL, 0, 0, data)
	if err != nil {
		return 0, err
	}
	if n != 2 {
		return 0, ErrorFailedToGetProtocol
	}

	return binary.LittleEndian.Uint16(data), nil
}

func sendString(dev *gousb.Device, request uint8, val, idx uint16, data []byte) error {
	// TODO, use interface{} instead of []byte, if string, convert to []byte
	if dev == nil {
		return ErrNoDevice
	}

	_, err := dev.Control(USB_DIR_OUT|USB_TYPE_VENDOR, request, val, idx, data)
	return err
}

func startAccessoryMode(dev *gousb.Device) error {
	if dev == nil {
		return ErrNoDevice
	}

	_, err := dev.Control(USB_DIR_OUT|USB_TYPE_VENDOR, ACCESSORY_START, 0, 0, nil)
	return err
}
