package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func transmit_file(file_name string, ip_port string) {

	fmt.Println("Sending " + file_name + " to " + ip_port + "...")

	// Connect to the server
	conn, err := net.Dial("tcp", ip_port)
	if err != nil {
		fmt.Println("Error connecting: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to", conn.RemoteAddr().String())

	// Read the file into memory
	fileData, err := os.ReadFile(file_name)
	if err != nil {
		fmt.Println("Error reading file: ", err.Error())
		os.Exit(1)
	}

	// Get the file size and convert it to a string
	fileSize := strconv.Itoa(len(fileData))

	// Send the file size and name to the server
	_, err = conn.Write([]byte(fileSize + ":" + file_name + ":"))
	if err != nil {
		fmt.Println("Error sending file size and file name: ", err.Error())
		os.Exit(1)
	}

	// Send the file data to the server
	_, err = conn.Write(fileData)
	if err != nil {
		fmt.Println("Error sending file data: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("File sent")
}

func recieve_file(port_semi string) {
	// Listen for incoming connections on port 8080
	listener, err := net.Listen("tcp", port_semi)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Listening on: ", listener.Addr().String())

	for {

		// Wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Received connection from: ", conn.RemoteAddr().String())

		// Receive the file size from the client
		fileSize := make([]byte, 100)
		_, err = conn.Read(fileSize)
		if err != nil {
			fmt.Println("Error receiving file size: ", err.Error())
			os.Exit(1)
		}

		// Convert the file size to an integer
		var str_size = strings.Split(string(fileSize), ":")[0]
		var file_name = strings.Split(string(fileSize), ":")[1]

		size, err := strconv.Atoi(str_size)
		if err != nil {
			fmt.Println("Error converting file size: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Recieved: " + file_name + ", File Size: " + string(str_size) + " bytes")

		// Receive the file data from the client
		fileData := make([]byte, size)
		_, err = conn.Read(fileData)
		if err != nil {
			fmt.Println("Error receiving file data: ", err.Error())
			os.Exit(1)
		}

		// Write the file to disk
		err = os.WriteFile(file_name, fileData, 0644)
		if err != nil {
			fmt.Println("Error writing file: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("File received and saved")
	}
}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Please Input a File to be Sent and IP/Port")
		fmt.Println("e.g. <.exe name> tx 127.0.0.1 file.py")
		fmt.Println("Or Enter arg '-ip_rx' to list get IP Info of rx Machine\n")
		fmt.Println("Or Enter arg '-ip_list' to list get IP Info of rx Machine")
		os.Exit(1)
	}

	if os.Args[1] == "-ip_list" {
		cmd := exec.Command("arp", "-a")

		stdout, err := cmd.CombinedOutput()

		if err != nil {
			panic(err)
		}
		fmt.Println(string(stdout))

		os.Exit(1)
	}

	if os.Args[1] == "-ip_rx" {

		// Retrieve the network interfaces on the current machine
		ifaces, err := net.Interfaces()
		if err != nil {
			panic(err)
		}

		// Iterate through each interface's address list to extract the IP addresses
		for _, iface := range ifaces {
			fmt.Println(iface.Name + ":")
			addrs, err := iface.Addrs()
			if err != nil {
				panic(err)
			}

			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						fmt.Println(ipNet.IP.String())
					}
				}
			}
		}

		os.Exit(1)
	}

	//e.g. "tx" or "rx"
	var tx_rx = os.Args[1]

	//choosing ip port
	var port_only = ":8080"
	//e.g. "127.0.0.1"
	var ip_adr_port = os.Args[2] + port_only

	if tx_rx == "tx" {

		//e.g. file.py
		//this is only used in tx, file name is sent to rx
		if len(os.Args) < 4 {
			fmt.Println("Please Enter the File Name after the IP/Port")
			os.Exit(1)
		}
		var file_name = os.Args[3]

		transmit_file(file_name, ip_adr_port)

	} else if tx_rx == "rx" {

		recieve_file(port_only)

	}
}
