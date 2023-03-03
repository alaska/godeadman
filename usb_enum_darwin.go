package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
	"strings"
)

var deviceRE = regexp.MustCompile(`^(.*):$`)
var locationRE = regexp.MustCompile(`^Location ID: (.*)$`)

// TODO: Call library directly
func enumerateDevices() ([]device, error) {
	if err := checkExe("system_profiler"); err != nil {
		return nil, err
	}

	out, err := exec.Command("system_profiler", "SPUSBDataType").Output()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))

	var (
		deviceList []device
		deviceName string
	)

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		switch {
		case locationRE.MatchString(line):
			m := locationRE.FindStringSubmatch(line)[1]
			if deviceName != "" {
				deviceList = append(deviceList, device{deviceName, m})
			}
			deviceName = ""
		case deviceRE.MatchString(line):
			deviceName = deviceRE.FindStringSubmatch(line)[1]
		}
	}
	return deviceList, nil
}
