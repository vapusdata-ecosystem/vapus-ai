package services

// import (
// 	"log"
// 	"net/http"
// 	"os/exec"

// 	"github.com/gorilla/websocket"
// 	"github.com/labstack/echo/v4"
// )

// // WebSocket upgrader
// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // Allow all origins (you can restrict this)
// 	},
// }

// func HandleWebSocket(c echo.Context) error {
// 	// Upgrade the connection
// 	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
// 	if err != nil {
// 		log.Println("Upgrade error:", err)
// 		return err
// 	}
// 	defer conn.Close()

// 	// Start a bash shell
// 	cmd := exec.Command("bash")
// 	cmdStdIn, _ := cmd.StdinPipe()
// 	cmdStdOut, _ := cmd.StdoutPipe()
// 	cmdStdErr, _ := cmd.StderrPipe()

// 	// Start the shell process
// 	if err := cmd.Start(); err != nil {
// 		log.Println("Command start error:", err)
// 		return err
// 	}

// 	// Stream stdout and stderr to WebSocket
// 	go func() {
// 		buf := make([]byte, 1024)
// 		for {
// 			n, err := cmdStdOut.Read(buf)
// 			if n > 0 {
// 				conn.WriteMessage(websocket.TextMessage, buf[:n])
// 			}
// 			if err != nil {
// 				break
// 			}
// 		}
// 	}()
// 	go func() {
// 		buf := make([]byte, 1024)
// 		for {
// 			n, err := cmdStdErr.Read(buf)
// 			if n > 0 {
// 				conn.WriteMessage(websocket.TextMessage, buf[:n])
// 			}
// 			if err != nil {
// 				break
// 			}
// 		}
// 	}()

// 	// Read from WebSocket and write to stdin
// 	for {
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 		cmdStdIn.Write(msg)
// 	}

// 	return nil
// }
