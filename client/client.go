// app2.go
// client.go
package main

import (
	"context"
	"fmt"
	"log"

	netLib "github.com/shirou/gopsutil/v3/net"
	"go.uber.org/zap"
)

func GetPIDFromPort(_ context.Context, logger *zap.Logger, port int) (uint32, error) {
	logger.Debug("Getting pid using port", zap.Int("port", port))

	connections, err := netLib.Connections("inet")
	if err != nil {
		return 0, err
	}

	for _, conn := range connections {
		if conn.Status == "LISTEN" && conn.Laddr.Port == uint32(port) {
			if conn.Pid > 0 {
				return uint32(conn.Pid), nil
			}
			return 0, fmt.Errorf("pid %d is out of bounds", conn.Pid)
		}
	}

	// If we get here, no process was found using the given port
	return 0, fmt.Errorf("no process found using port %d", port)
}

func main() {
	// Wait for a brief moment to ensure the server is up
	// fmt.Println("Client process started, making request to the server...")

	// // Make an HTTP request to the server started in the main process
	// resp, err := http.Get("http://localhost:8080/ping")
	// if err != nil {
	// 	log.Fatalf("Failed to reach server: %v", err)
	// }
	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalf("Failed to read response body: %v", err)
	// }

	// fmt.Println("Response from server:", string(body))

	pid, err := GetPIDFromPort(context.Background(), zap.NewExample(), 8090)
	if err != nil {
		log.Fatalf("Failed to get pid: %v", err)
	}
	fmt.Println("PID of the process using port 8090:", pid)

}
