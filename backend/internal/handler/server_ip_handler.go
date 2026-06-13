package handler

import (
	"io"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"regonline-backend/internal/response"

	"github.com/gin-gonic/gin"
)

type ServerIPHandler struct {
}

func NewServerIPHandler() *ServerIPHandler {
	return &ServerIPHandler{}
}

type ServerIPResponse struct {
	PublicIP  string `json:"public_ip"`
	LocalIP   string `json:"local_ip"`
	BestIP    string `json:"best_ip"`
	Port      int    `json:"port"`
	URL       string `json:"url"`
	IPType    string `json:"ip_type"`
}

var publicIPServices = []string{
	"https://api.ipify.org",
	"https://ifconfig.me/ip",
	"https://icanhazip.com",
	"https://checkip.amazonaws.com",
}

func (h *ServerIPHandler) GetServerIP(c *gin.Context) {
	result := ServerIPResponse{
		Port: 5000,
	}

	publicIP := detectPublicIP()
	if publicIP != "" {
		result.PublicIP = publicIP
		result.BestIP = publicIP
		result.IPType = "public"
		if isIPv4(publicIP) {
			result.IPType = "public_ipv4"
		} else {
			result.IPType = "public_ipv6"
		}
	}

	localIPs := detectLocalIPs()

	if result.BestIP == "" {
		bestLocal := findBestLocalIP(localIPs)
		if bestLocal != "" {
			result.BestIP = bestLocal
			result.IPType = "private_ipv4"
		}
	}

	if result.BestIP != "" {
		if isIPv6(result.BestIP) {
			result.URL = "http://[" + result.BestIP + "]:" + itoa(result.Port)
		} else {
			result.URL = "http://" + result.BestIP + ":" + itoa(result.Port)
		}
	}

	if len(localIPs) > 0 {
		result.LocalIP = localIPs[0]
	}

	response.Success(c, gin.H{
		"public_ip": result.PublicIP,
		"local_ip":  result.LocalIP,
		"best_ip":   result.BestIP,
		"port":      result.Port,
		"url":       result.URL,
		"ip_type":   result.IPType,
	})
}

func detectPublicIP() string {
	client := &http.Client{Timeout: 3 * time.Second}

	type result struct {
		ip  string
		err error
	}

	ch := make(chan result, len(publicIPServices))
	for _, svc := range publicIPServices {
		go func(url string) {
			resp, err := client.Get(url)
			if err != nil {
				ch <- result{err: err}
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(io.LimitReader(resp.Body, 64))
			if err != nil {
				ch <- result{err: err}
				return
			}
			ip := strings.TrimSpace(string(body))
			ip = regexp.MustCompile(`[^\d\.:a-fA-F]`).ReplaceAllString(ip, "")
			if net.ParseIP(ip) != nil {
				ch <- result{ip: ip}
			} else {
				ch <- result{err: io.EOF}
			}
		}(svc)
	}

	var ipv4, ipv6 string
	for i := 0; i < len(publicIPServices); i++ {
		r := <-ch
		if r.err != nil {
			continue
		}
		if ipv4 == "" && isIPv4(r.ip) {
			ipv4 = r.ip
		}
		if ipv6 == "" && isIPv6(r.ip) {
			ipv6 = r.ip
		}
	}

	if ipv4 != "" {
		return ipv4
	}
	if ipv6 != "" {
		return ipv6
	}
	return ""
}

func detectLocalIPs() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var ips []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP
			if ip == nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
				continue
			}
			if ip4 := ip.To4(); ip4 != nil {
				ips = append(ips, ip4.String())
			}
		}
	}

	sort.Slice(ips, func(i, j int) bool {
		pi := isPrivateIP(ips[i])
		pj := isPrivateIP(ips[j])
		if pi != pj {
			return !pi
		}
		return ips[i] < ips[j]
	})

	return ips
}

func findBestLocalIP(ips []string) string {
	for _, ip := range ips {
		if !isPrivateIP(ip) {
			return ip
		}
	}
	for _, ip := range ips {
		return ip
	}
	return ""
}

var privateRanges = []*net.IPNet{
	{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
	{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
	{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
	{IP: net.IPv4(169, 254, 0, 0), Mask: net.CIDRMask(16, 32)},
}

func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return true
	}
	for _, r := range privateRanges {
		if r.Contains(ip) {
			return true
		}
	}
	return false
}

func isIPv4(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil && ip.To4() != nil
}

func isIPv6(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil && ip.To4() == nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}