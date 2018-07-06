/*
这个程序参考，基本是抄，主要是刚学，重在理解：
https://www.jianshu.com/p/eebf44c1d6db
https://www.oschina.net/translate/build-a-realtime-chat-server-with-go-and-websockets?lang=chs&page=1#
https://astaxie.gitbooks.io/build-web-application-with-golang/zh/03.2.html
https://godoc.org/github.com/gorilla/websocket#Conn.ReadMessage
*/
package main 
import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
    // "bufio"
    "fmt"
)
var clients=make(map[*websocket.Conn]bool) //连接的客户端
var broadcast=make(chan Message) //接受的消息管道
// var upgrader=websocket.Upgrader{} //获取一个普通的HTTP连接，然后将其升级为websocket
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Message struct {
	// Email string `json:"email"`
	Username string `json:"username"`
	Message string `json:"message"`
}
// type Message string 
func main() {
	//创建一个简单的文件服务器
	// fs:=http.FileServer(http.Dir("../public"))
	// http.Handle("/",fs)
	//设置访问路由，在url中/ws的就用handleConnections函数处理
	http.HandleFunc("/ws", handleConnections)
    // fmt.Println("捕获了")
	//处理用户发送来的数据，将这些数据发送到每一个客户端
	go handleMessages()
	//开启服务
	log.Println("http server started on:8080")
	err:=http.ListenAndServe("localhost:8080", nil)
	if err!=nil {
		log.Fatal("ListenAndServer:",err)
	}
}
func handleConnections(w http.ResponseWriter, r *http.Request) {
        // 得到websocket连接的客户端ws
		//调用upgrader.Upgrade(w, r, nil)获得这个连接的指针
        fmt.Println("handleConnections")
        ws, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
                log.Fatal(err)
        }
        // 确保在函数返回时将websocket关闭
        defer ws.Close()
        //将websocket加入到clients中
        clients[ws] = true
        fmt.Println("有用户连接")
         for {
                var msg Message               
                //将ws中的数据已json的类型读入到msg中
                err := ws.ReadJSON(&msg)
                // msg:=bufio.NewReader.ReadString()
                if err != nil {
                        log.Printf("error: %v", err)
                        delete(clients, ws)
                        break
                }
                //将刚接收的客户端的信息放入到消息管道中，用于向其他客户端发送
                broadcast <- msg
            }
}
func handleMessages() {
        for {
                msg := <-broadcast
                for client := range clients {
                        err := client.WriteJSON(msg)
                        // client.Write([]byte(msg + "\n"))
                        if err != nil {
                                log.Printf("error: %v", err)
                                client.Close()
                                delete(clients, client)
                        }
                }
        }
}
