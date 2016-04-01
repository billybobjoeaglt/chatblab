package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"bufio"
	"strings"
)

func main() {
	// Passphrase for private key
	/*var passprase, err = ioutil.ReadFile("./pass.key")
	if err != nil {
		panic(err)
	}*/

	err := loadConfig()
	if err != nil {
		panic(err)
	}

	go printAll(outputChannel)
	go listen()
	go func(){
		reader := bufio.NewReader(os.Stdin)
		for{
			text, _ := reader.ReadString('\n')
			text=strings.TrimSpace(text)
			if strings.Contains(text,"connect "){
				createConnection(strings.Split(text,"connect ")[1]+":8080")
			}else{
				broadcastMessage(text)
			}
		}
	}()
	// Exit capture
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-sigs
		err := saveConfig()
		if err != nil {
			panic(err)
		}
		done <- true
	}()
	<-done
	fmt.Println("Safe Exited")
}
