package driver

import (
	"sync"

	"github.com/kubeedge/mapper-generator/pkg/common"
)

// CustomizedDev is the customized device configuration and client information.
type CustomizedDev struct {
	Instance         common.DeviceInstance
	CustomizedClient *CustomizedClient
}

type CustomizedClient struct {
	// TODO add some variables to help you better implement device drivers
	deviceMutex sync.Mutex
	ProtocolCommonConfig
	ProtocolConfig
}

type ProtocolConfig struct {
	SlaveID float64 `json:"slaveID"`
}

type ProtocolCommonConfig struct {
	Com
	// TODO: add your Common data according to configmap
	CommonCustomizedValues `json:"customizedValues"`
}

type Com struct {
	SerialPort string `json:"serialPort"`
	DataBits   int    `json:"dataBits"`
	BaudRate   int    `json:"baudRate"`
	Parity     string `json:"parity"`
	StopBits   int    `json:"stopBits"`
}

type CommonCustomizedValues struct {
	SerialType string `json:"serialType"`
	// TODO: add your CommonCustomizedValues according to configmap
}
type VisitorConfig struct {
	Name           string
	Register       string  `json:"register"`
	Offset         uint16  `json:"offset"`
	Limit          int     `json:"limit"`
	Scale          float64 `json:"scale,omitempty"`
	IsSwap         bool    `json:"isSwap,omitempty"`
	IsRegisterSwap bool    `json:"isRegisterSwap,omitempty"`
}
