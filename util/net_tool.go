package util

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type ReceiveBytes uint64
type TransmitBytes uint64

func UploadDownloadFlow(dev string) (ReceiveBytes, TransmitBytes, error) {
	down, up, err := TotalFlowByDevice(dev)
	if err != nil {
		panic(err)
		return 0, 0, err
	}
	//	fmt.Println(down, up)
	time.Sleep(time.Second * 1)

	down2, up2, err := TotalFlowByDevice(dev)
	if err != nil {
		panic(err)
		return 0, 0, err
	}
	//fmt.Println(down2, up2)

	return down2 - down, up2 - up, nil
}

func TotalFlowByDevice(dev string) (ReceiveBytes, TransmitBytes, error) {
	devInfo, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		panic(err)
		return 0, 0, err
	}

	var receive int = -1
	var transmit int = -1

	var receiveBytes uint64
	var transmitBytes uint64

	lines := strings.Split(string(devInfo), "\n")
	for _, line := range lines {
		//fmt.Println(line)
		if strings.Contains(line, dev) {
			i := 0
			fields := strings.Split(line, ":")
			for _, field := range fields {
				if strings.Contains(field, dev) {
					i = 1
				} else {
					values := strings.Fields(field)
					for _, value := range values {
						//logger.Debug(value)
						if receive == i {
							bytes, _ := strconv.ParseInt(value, 10, 64)
							receiveBytes = uint64(bytes)
						} else if transmit == i {
							bytes, _ := strconv.ParseInt(value, 10, 64)
							transmitBytes = uint64(bytes)
						}
						i++
					}
				}
			}
		} else if strings.Contains(line, "face") {
			index := 0
			tag := false
			fields := strings.Split(line, "|")
			for _, field := range fields {
				if strings.Contains(field, "face") {
					index = 1
				} else if strings.Contains(field, "bytes") {
					values := strings.Fields(field)
					for _, value := range values {
						//logger.Debug(value)
						if strings.Contains(value, "bytes") {
							if !tag {
								tag = true
								receive = index
							} else {
								transmit = index
							}
						}
						index++
					}
				}
			}
		}
	}
	//	fmt.Println("receive_bytes :", receiveBytes)
	//	fmt.Println("transmit_bytes :", transmitBytes)

	return ReceiveBytes(receiveBytes), TransmitBytes(transmitBytes), nil
}

// func main() {
// 	dev := flag.String("dev", "eth0", "net device name")
// 	flag.Parse()
// 	for {
// 		d, u, _ := UploadDownloadFlow(*dev)
// 		fmt.Println("download:", humanize.Bytes(uint64(d))+"/s", "upload:", humanize.Bytes(uint64(u))+"/s")
// 		time.Sleep(time.Millisecond * 500)
// 	}

// }
