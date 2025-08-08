package main

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
)

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// CircularIDAllocator ID分配器
type CircularIDAllocator struct {
	MinID     int64
	MaxID     int64
	CurrentID int64
	IP        string
	Port      int
	IPHash    int64
	Offset    int64
	Lock      sync.Mutex
}

// NewCircularIDAllocator 创建ID分配器
func NewCircularIDAllocator(mockIP string, port int) *CircularIDAllocator {
	ip := mockIP
	if ip == "" {
		ip = getLocalIP()
	}

	ipHash := hashIP(ip)
	offset := (ipHash * (1 << 32)) + (int64(port) * (1 << 16))

	return &CircularIDAllocator{
		MinID:     0,
		MaxID:     math.MaxInt64,
		CurrentID: 0,
		IP:        ip,
		Port:      port,
		IPHash:    ipHash,
		Offset:    offset,
	}
}

// Allocate 分配新ID
func (a *CircularIDAllocator) Allocate() int64 {
	a.Lock.Lock()
	defer a.Lock.Unlock()

	allocatedID := a.CurrentID + a.Offset
	a.CurrentID++

	if a.CurrentID >= (1 << 32) {
		a.CurrentID = 0
		fmt.Println("ID counter reset", "currentID", a.CurrentID)
	}

	fmt.Println("Allocated new room ID",
		"allocatedID", allocatedID,
		"currentID", a.CurrentID,
		"offset", a.Offset,
		"ip", a.IP,
		"port", a.Port)

	return allocatedID
}

// hashIP IP地址哈希
func hashIP(ip string) int64 {
	parts := strings.Split(ip, ".")
	hash := int64(0)

	for _, part := range parts {
		num, _ := strconv.ParseInt(part, 10, 64)
		hash = (hash*256 + num) % (1 << 32)
	}

	return hash
}

func main() {
	a := NewCircularIDAllocator("", 0)
	fmt.Println(a.Allocate())
	fmt.Println(a.Allocate())
	fmt.Println(a.Allocate())

	a.Allocate()
	a.Allocate()

}
