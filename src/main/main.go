package main

import(
    _"syscall"
    _"fmt"
    "os/exec"
    "log"
    "io"
    "net"
    _"os"
)

// TCP Server


func main() {
    // 1) Send the endpoint INFO to server










    // 2) Wait for shutdown command 

    l, err := net.Listen("tcp", ":8003")
    if nil != err {
        log.Fatalf("fail to bind address to 5032; err: %v", err)
    }
    defer l.Close()

    for {
        conn, err := l.Accept()
        if nil != err {
            log.Printf("fail to accept; err: %v", err)
            continue
        }
        go ConnHandler(conn)
    }
}

func ConnHandler(conn net.Conn) {
    recvBuf := make([]byte, 4096) // receive buffer: 4kB
    for {
        n, err := conn.Read(recvBuf)
        if nil != err {
            if io.EOF == err {
                log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
                return
            }
            log.Printf("fail to receive data; err: %v", err)
            return
        }

        if 0 < n {
            data := recvBuf[:n]
            log.Println(string(data))
            
            if string(data) == "shutdown"{
                Shutdown()
            }

        }
    }
}



// Shutdown the this computer
func Shutdown(){
    
    cmd := exec.Command("shutdown", "now")
    err := cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
}
