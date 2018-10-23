package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	//	"time"
	. "fly"
)

func main() {
	fmt.Println(os.Getppid())

	if os.Getppid() == 1 {
		execFilePath, _ := filepath.Abs(os.Args[0])
		cmd := exec.Command(execFilePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()

		return
	} else {
		ls, err := net.Listen("tcp", "0.0.0.0:6789")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ls.Close()
		fmt.Println("S: tcp server listen on 6789")

		for {
			fmt.Println("S: waiting for new connection")
			con, err := ls.Accept()
			if err != nil {
				return
			}
			go handleConnection(con)
		}
		// 	for {
		// 		fmt.Println("log loop")
		// 		time.Sleep(2 * time.Second)
		// 	}
	}

}

func handleConnection(c net.Conn) {
	fmt.Println("S: start receive client data")

	buffer := make([]byte, 518)
	var msg_head MsgHeadT
	c.Read(buffer)
	buf := bytes.NewBuffer(buffer)
	binary.Read(buf, binary.LittleEndian, &msg_head)

	fd, err := os.Create("./test.file.recv")
	if err != nil {
		fmt.Println("failllllll")
		fmt.Println(string(msg_head.DesName[:]))
		return
		panic("S: open file fail")
	}
	defer fd.Close()

	var data_size = msg_head.Length
	var rbuf = make([]byte, 1024)
	if data_size > 0 {
		for {
			n, err := c.Read(rbuf)
			if err != nil {
				panic("S: Read data from network fail")
			}
			fmt.Println("S: receive : ", string(rbuf[:n]))
			r, err := fd.Write(rbuf[:n])
			r = 0
			if err != nil {
				fmt.Println("S: write data to file fail ", r)
			}
			data_size -= uint32(n)
			if data_size < 0 {
				panic("S: receive data fail")
			}
			if data_size == 0 {
				break
			}

		}
		fd.Sync()
		//		time.Sleep(5 * time.Second)
		var ack MsgACKT
		ack.Tag = 1
		ack.MType = 1
		ack.Code = 12
		var ackMessageLength uint32 = uint32(7 + len([]byte("receive finish")))
		ack.Length = ackMessageLength

		var retMessage = "Receive finish!"

		aBuf := new(bytes.Buffer)
		binary.Write(aBuf, binary.LittleEndian, ack)
		binary.Write(aBuf, binary.LittleEndian, retMessage)
		c.Write(append(aBuf.Bytes(), []byte(retMessage)...))
		fmt.Println("S: receive data finished and send ack to client")

	}

	return

	//	for {
	c.Read(buffer)
	c.Write(buffer)

	fmt.Println("S: finish receive client data")
	//	}
}
