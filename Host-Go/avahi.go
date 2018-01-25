package main

import (
	"os/exec"
	"strconv"
	"strings"
)

func collect_with_avahi(path string) (string, error) {
	out, err := exec.Command(path, "-t", "-c", "-r", "-p", "_http._tcp").Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		fields := strings.Split(line, ";")
		if len(fields) < 9 || fields[0] != "=" {
			// consider only complete entries
			continue
		}
		if strings.ToLower(fields[2]) != "ipv4" {
			// ignore IPv6 addresses
			continue
		}
		if strings.Contains(fields[7], ":") {
			// seems to be an IPv6 address
			continue
		}
		port, err := strconv.Atoi(fields[8])
		if err != nil || port < 1 || port > 65535 {
			// missing or invalid port
			continue
		}

		txt := []string{}
		if len(fields) >= 9 {
			list := strings.Split(fields[9], `" "`)
			for _, elem := range list {
				elem = strings.Trim(elem, `"`)
				elem = strings.Replace(elem, "\\032", " ", -1)
				txt = append(txt, elem)
			}
		}
		add_server(fields[3], fields[6], fields[7], port, txt)
	}

	return "Go (Avahi)", nil
}