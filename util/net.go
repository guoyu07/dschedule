/* *
 * Author: liubin8@staff.weibo.com
 * Date: 2015.12.28
 * Refer: hashicorp/consul/command/agent/
 * */

package util

import (
	"fmt"
	"net"
	"strings"
)

// Refer: hashicorp/consul/command/agent/config.go
// unixSocketAddr tests if a given address describes a domain socket,
// and returns the relevant path part of the string if it is.
func UnixSocketAddr(addr string) (string, bool) {
	if !strings.HasPrefix(addr, "unix://") {
		return "", false
	}
	return strings.TrimPrefix(addr, "unix://"), true
}

// Refer: hashicorp/consul/command/agent/config.go
// ClientListener is used to format a listener for a
// port on a ClientAddr
func ClientListener(addr string, port int) (net.Addr, error) {
	if path, ok := UnixSocketAddr(addr); ok {
		return &net.UnixAddr{Name: path, Net: "unix"}, nil
	}
	ip := net.ParseIP(addr)
	if ip == nil {
		return nil, fmt.Errorf("Failed to parse IP: %v", addr)
	}
	return &net.TCPAddr{IP: ip, Port: port}, nil
}

// Refer: hashicorp/consul/command/agent/util.go

// FilePermissions is an interface which allows a struct to set
// ownership and permissions easily on a file it describes.
type FilePermissions interface {
	// User returns a user ID or user name
	User() string

	// Group returns a group ID. Group names are not supported.
	Group() string

	// Mode returns a string of file mode bits e.g. "0644"
	Mode() string
}

// setFilePermissions handles configuring ownership and permissions settings
// on a given file. It takes a path and any struct implementing the
// FilePermissions interface. All permission/ownership settings are optional.
// If no user or group is specified, the current user/group will be used. Mode
// is optional, and has no default (the operation is not performed if absent).
// User may be specified by name or ID, but group may only be specified by ID.
func SetFilePermissions(path string, p FilePermissions) error {
	/* TODO open
	   	var err error
	   	uid, gid := os.Getuid(), os.Getgid()

	   	if p.User() != "" {
	   		if uid, err = strconv.Atoi(p.User()); err == nil {
	   			goto GROUP
	   		}

	   		// Try looking up the user by name
	   		if u, err := user.Lookup(p.User()); err == nil {
	   			uid, _ = strconv.Atoi(u.Uid)
	   			goto GROUP
	   		}

	   		return fmt.Errorf("invalid user specified: %v", p.User())
	   	}

	   GROUP:
	   	if p.Group() != "" {
	   		if gid, err = strconv.Atoi(p.Group()); err != nil {
	   			return fmt.Errorf("invalid group specified: %v", p.Group())
	   		}
	   	}
	   	if err := os.Chown(path, uid, gid); err != nil {
	   		return fmt.Errorf("failed setting ownership to %d:%d on %q: %s",
	   			uid, gid, path, err)
	   	}

	   	if p.Mode() != "" {
	   		mode, err := strconv.ParseUint(p.Mode(), 8, 32)
	   		if err != nil {
	   			return fmt.Errorf("invalid mode specified: %v", p.Mode())
	   		}
	   		if err := os.Chmod(path, os.FileMode(mode)); err != nil {
	   			return fmt.Errorf("failed setting permissions to %d on %q: %s",
	   				mode, path, err)
	   		}
	   	}
	*/

	return nil
}

//http://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
