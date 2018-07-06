/*
这个程序参考：
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
var clients=make(map[*websocket.Conn]bool) //Connected client
var broadcast=make(chan Message) //Accepted message pipeline
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
	
	//Set access route
	http.HandleFunc("/ws", handleConnections)
	//Process the data sent by the user and send the data to each client
	go handleMessages()
	//Open service
	log.Println("http server started on:8080")
	err:=http.ListenAndServe("localhost:8080", nil)
	if err!=nil {
		log.Fatal("ListenAndServer:",err)
	}
}
func handleConnections(w http.ResponseWriter, r *http.Request) {
		//Call upgrader.Upgrade(w, r, nil) to get the pointer to this connection
        fmt.Println("handleConnections")
        ws, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
                log.Fatal(err)
        }
        // Make sure the websocket is closed when the function returns
        defer ws.Close()
        //Add websocket to clients
        clients[ws] = true
        fmt.Println("有用户连接")
         for {
                var msg Message               
                //Read the type of json data in ws into msg
                err := ws.ReadJSON(&msg)
                // msg:=bufio.NewReader.ReadString()
                if err != nil {
                        log.Printf("error: %v", err)
                        delete(clients, ws)
                        break
                }
                //Put the information of the client just received into the message pipeline for sending to other clients
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
