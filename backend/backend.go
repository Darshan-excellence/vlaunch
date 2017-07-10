package backend

import (
	"errors"
	"os"

	"github.com/lebauce/vlaunch/config"
)

var DeviceNotFound = errors.New("Could not find device")

func FindDevice() (string, error) {
	if device := config.GetConfig().GetString("device"); device != "" {
		return device, nil
	}

	if uuid := config.GetConfig().GetString("device_uuid"); uuid != "" {
		if device, err := FindDeviceByUUID(uuid); err == nil {
			return device, nil
		}
	}

	if executable, err := os.Executable(); err == nil {
		if device, err := FindDeviceByPath(executable); err == nil {
			return device, nil
		}
	}

	return "", DeviceNotFound
}