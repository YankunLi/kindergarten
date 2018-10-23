package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	. "fly"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"utils"
)

type Command struct {
	CMD string
}

type CommandInterface interface {
	do() error
}

func (c *Command) do() error {
	//do
}

func Usage() {
	fmt.Println("fly {command} [option]")
	fmt.Println("Options:")
	fmt.Println("    -s_file         The path of file to send")
	fmt.Println("    -d_file_name    The name of destination file")
	fmt.Println("    -config         config file")
	fmt.Println("    -log_file       log file")
	fmt.Println("    -remote_server  remote server ip and port")
	return
}

func InitConfig(flyConf *ConfType) {
	flyConf.SrcFilePath = flag.String("s_file", "./test.file", "file to send")
	var f = flag.String("file", "./test.file", "file to send")
	flyConf.DesFileName = flag.String("d_file_name", "", "the name of destination file")
	flyConf.ConfigFile = flag.String("conf", "/etc/fly/fly.conf", "config file")
	flyConf.LogFile = flag.String("log_file", "/var/log/fly/fly.log", "fly log file")
	flyConf.RemoteSer = flag.String("remote_server", "127.0.0.1:6789", "remote server")
	flag.Parse()
	flyConf.CMD = flag.Arg(0)
	//flag.PrintDefaults()
	fmt.Println(*f)
}

func InitLog(logger **log.Logger, logFile *string) {
	fd, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("open log file fail")
		panic("exit")
	}
	//init log

	*logger = log.New(fd, "[INFO]", log.Llongfile|log.LstdFlags)
}

func doCommand(flyCtx *FlyContext, cmd *Command) {
	logger := flyCtx.Config.Logger

	switch cmd.CMD {
	case "sendfile":
		SendFile(flyCtx)
	default:
		logger.Println("conn't found this command : %s", cmd.CMD)
	}
}

func SendFile(flyCtx *FlyContext) {
	client := flyCtx.Client
	config := flyCtx.Config
	logger := config.Logger
	srcFilePath := config.SrcFilePath
	destFileName := config.DesFileName

	//build msg head
	fileInfo, err := os.Stat(*srcFilePath)
	if err != nil {
		logger.Println("get file infor fail os.Stat")
		panic("get file infor fail os.Stat")
	}

	var msg_head MsgHeadT
	msg_head.Tag = 1
	msg_head.MType = 2
	msg_head.Length = uint32(fileInfo.Size())
	StringToArr(*srcFilePath, msg_head.SrcName[:])
	if *destFileName == "" {
		*destFileName = *srcFilePath + ".recive"
	}
	StringToArr(*destFileName, msg_head.DesName[:])

	//send request head
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, msg_head)
	client.Conn.Write(buf.Bytes())

	//send request body
	fmt.Println("src file ", *srcFilePath)
	fd, err := os.Open(*srcFilePath)
	if err != nil {
		logger.Printf("Open file %s fail", *srcFilePath)
		panic("Open file fail")
	}
	defer fd.Close()

	rbuf := make([]byte, 4096)
	for {
		n, err := fd.Read(rbuf)
		if err != nil && err != io.EOF {
			logger.Println("Read file fail")
			panic("C: Read file fail")
		}

		if n == 0 {
			break
		}

		fmt.Println("C: send data : ", string(rbuf[:n]))
		client.Conn.Write(rbuf[:n])
	}
	fmt.Println("send file finished")
	logger.Println("send file finished")

	var toRead, ackHeadSize = 7, 7
	var readIndex = 0
	for {
		n, err := client.Conn.Read(rbuf[readIndex:toRead])
		if err != nil {
			logger.Println("receive request ack fail")
		}
		fmt.Println("C: receive data length", n)
		readIndex += n
		if readIndex == 7 {
			break
		}
	}

	fmt.Println("C: Decode ack binary")

	var ack MsgACKT

	bBuf := bytes.NewBuffer(rbuf[:ackHeadSize])
	binary.Read(bBuf, binary.LittleEndian, &ack)

	n, err := client.Conn.Read(rbuf)
	if err != nil {
		logger.Println("receive ack fail")
	}
	fmt.Printf("ack code: %d, message: %s\n", ack.Code, rbuf[:n])

	logger.Printf("receive ack msg, code : %d, message: %s", ack.Code, rbuf)

}

func main() {
	//parse command parameter
	var flyCtx FlyContext
	var config ConfType

	flyCtx.Config = &config

	InitConfig(&config)

	utils.Show()

	//init log
	InitLog(&config.Logger, config.LogFile)

	var logger *log.Logger = config.Logger

	logger.Printf("Start connect remote server : %s", *config.RemoteSer)

	//connect server
	var client NetClient
	flyCtx.Client = &client
	var err error
	client.Conn, err = net.Dial("tcp", *config.RemoteSer)
	if err != nil {
		logger.Printf("Connect remote server %s fail", *config.RemoteSer)
		logger.SetPrefix("[debug] ")
		logger.Printf("[debug] Connect server %s fail", *config.RemoteSer)
		return
	}

	logger.Println("Connecting server successfully")

	defer client.Conn.Close()

	//build command and do it

	// get file infor
	if config.CMD == "" {
		logger.Println("cann't found command")
		Usage()
	}
	var cmd Command
	cmd.CMD = config.CMD

	doCommand(&flyCtx, &cmd)

	return
}

func StringToArr(str string, s []byte) (err error) {
	if len([]rune(str)) > len(s) {
		return errors.New("string too long")
	}

	temp := []byte(str)
	for i := 0; i < len(temp); i++ {
		s[i] = temp[i]
	}

	return
}

//func ReadAll(fd *File, buf []byte, toRead uint32) {
//
//}
