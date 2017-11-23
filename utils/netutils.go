package utils

import "net"

func GetIP() string {
    inters, err := net.Interfaces()
    if err != nil {
        return "error"
    }
    for _, inter := range inters {
        if inter.Name == "lo" {
            continue
        }
        addrs, err := inter.Addrs()
        if err != nil {
            continue
        }
        for _, addr := range addrs {
            ipnet, ok := addr.(*net.IPNet)
            if ok && !ipnet.IP.IsLoopback() {
                if ipnet.IP.To4() != nil {
                    return ipnet.IP.String()
                }
            }
        }
    }
    panic("Unable to determine local IP address (non loopback). Exiting.")
}
