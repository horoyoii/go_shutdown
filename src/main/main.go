package main

import(
    _"syscall"
    "fmt"
    "os/exec"
    "log"
    "io"
    "net"
    _"os"
    "bytes"
    "net/http"
)

// TCP Server


func main() {
    // 1) Send the endpoint INFO to server

    reqBody := bytes.NewBufferString(GetLocalIP())
    fmt.Println(GetLocalIP())

    recv, err := http.Post("http://34.225.204.24:8004/api/v1/notebook/turnon", "text/plain", reqBody) 
    fmt.Println("Send turnon signal to server")
    if err != nil{
        fmt.Println(err.Error())
    }

    fmt.Println(recv)
    fmt.Println("recieve Response")


    // 2) Wait for shutdown command 
    fmt.Println("local program server is running")

    http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request){
        fmt.Fprintf(w, "pong %s", "pong")
    })

    http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request){
        fmt.Fprintf(w, "ok %s","shutdown")
        Shutdown()
    })
   

    http.ListenAndServe(":8003", nil)


//    fmt.Println("wait for signal from server")
//    l, err := net.Listen("tcp", ":8003")
//    if nil != err {
//        log.Fatalf("fail to bind address to 5032; err: %v", err)
//    }
//    defer l.Close()
//
//    for {
//        conn, err := l.Accept()
//        if nil != err {
//            log.Printf("fail to accept; err: %v", err)
//            continue
//        }
//        go ConnHandler(conn)
//    }
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

func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

// Shutdown the this computer
func Shutdown(){
    
    cmd := exec.Command("shutdown", "now")
    err := cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
}
