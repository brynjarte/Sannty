
package main


import(
	"fmt"
	"net"
	"time"
	"encoding/json"
)

var boolvar bool

type UDPMessage struct{
	Message string
	MessageNumber int 
}


func Recieve(){
	boolvar := true
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", ":25555")
	recievesock,_ := net.ListenUDP("udp", raddr)
	for(boolvar)  {
		mlen , _,_ := recievesock.ReadFromUDP(buffer)

		
	
		var rec_msg UDPMessage
		json.Unmarshal(buffer[:mlen], &rec_msg)
	
		fmt.Println(rec_msg.Message, rec_msg.MessageNumber)	
	}
}

func Send(){
	
	baddr,err := net.ResolveUDPAddr("udp", "129.241.187.255:25555")
	sendSock, err := net.DialUDP("udp", nil ,baddr) // connection
	send_msg := UDPMessage{"Det er mottat!",1}
	time.Sleep(1*time.Second)
	buf,_ := json.Marshal(send_msg)
	_,err = sendSock.Write(buf)
	if err != nil{
		panic(err)
	}
	
}

func main(){
	
	go Recieve()
	time.Sleep(1*time.Second)
	Send()
	time.Sleep(100*time.Second)
	boolvar = false
		
}
