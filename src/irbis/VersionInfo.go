package irbis

import "strconv"

type VersionInfo struct {
	Organization     string
	Version          string
	MaxClients       int
	ConnectedClients int
}

func (version *VersionInfo) Parse(lines []string) {
	if len(lines) == 3 {
		version.Version = lines[0]
		version.ConnectedClients, _ = strconv.Atoi(lines[1])
		version.MaxClients, _ = strconv.Atoi(lines[2])
	} else {
		version.Organization = lines[0]
		version.Version = lines[1]
		version.ConnectedClients, _ = strconv.Atoi(lines[2])
		version.MaxClients, _ = strconv.Atoi(lines[3])
	}
}

func (version *VersionInfo) String() string {
	return version.Version
}
