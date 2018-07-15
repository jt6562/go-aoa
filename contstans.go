package aoa

import (
	"errors"
)

var (
	// Errors
	ErrNoDevice              error = errors.New("No device")
	ErrorFailedToGetProtocol       = errors.New("Failed to get protocol")
	ErrorNotSupport                = errors.New("Not support this function")
)

const (
	MODE_ACCESSORY uint8 = 1
	MODE_AUDIO           = 2
)

const (
	ACCESSORY_VENDOR_ID = 0x18D1

	PRODUCT_ID_ACCESSORY           = 0x2D00
	PRODUCT_ID_ACCESSORY_ADB       = 0x2D01
	PRODUCT_ID_AUDIO               = 0x2D02
	PRODUCT_ID_AUDIO_ADB           = 0x2D03
	PRODUCT_ID_ACCESSORY_AUDIO     = 0x2D04
	PRODUCT_ID_ACCESSORY_AUDIO_ADB = 0x2D05

	ACCESSORY_STRING_MANUFACTURER = 0
	ACCESSORY_STRING_MODEL        = 1
	ACCESSORY_STRING_DESCRIPTION  = 2
	ACCESSORY_STRING_VERSION      = 3
	ACCESSORY_STRING_URI          = 4
	ACCESSORY_STRING_SERIAL       = 5
)

// requestType
const (
	USB_DIR_OUT     uint8 = 0
	USB_DIR_IN            = 0x80
	USB_TYPE_VENDOR       = (0x02 << 5)
)

// requests
const (
	ACCESSORY_GET_PROTOCOL        uint8 = 51
	ACCESSORY_SEND_STRING               = 52
	ACCESSORY_START                     = 53
	ACCESSORY_REGISTER_HID              = 54
	ACCESSORY_UNREGISTER_HID            = 55
	ACCESSORY_SET_HID_REPORT_DESC       = 56
	ACCESSORY_SEND_HID_EVENT            = 57
	ACCESSORY_SET_AUDIO_MODE            = 58
)
