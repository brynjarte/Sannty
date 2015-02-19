package UDP 

import (
	"fmt"
	"strconv"
	"time"
	
)

type Process struct {
	Master bool
	Backup bool
	SequenceNumber int
}

func ProcessPair(p Process, rec_channel chan UDPMessage){
	
	for{
		timeChannel := make(chan bool,1)
		if(p.Master){
			for i:=p.SequenceNumber; i<p.SequenceNumber+6; i++ {
				msg := UDPMessage{strconv.Itoa(i),i}
				sendUdpMessage(msg)
				fmt.Println(i)
			}
		return
		}

		if(p.Backup){
			if(p.SequenceNumber == 0){
				go recieveUdpMessage(rec_channel)//DERSOM KANALEN ER UNBUFFERED SÅ TRENGE ME IKKJE SJEKKD DA.
			}
			for{
				go timeOut(timeChannel)
						
				select{
					case recievedMsg := <- rec_channel: 
						p.SequenceNumber = recievedMsg.MessageNumber	
				
					case <-timeChannel:
						fmt.Println("timeout")
						p.Backup = false
						p.Master = true
						p1 := Process{false,true,p.SequenceNumber}
						go ProcessPair(p1, rec_channel) // spawn new backup
						fmt.Println("NEW MASTAH")
						break
					}
				if p.Master{
					break
				}
	
			}
		}

	}
}	


func timeOut(timeChannel chan bool){
	time.Sleep(3*time.Second)
	timeChannel <- true
}

func primary(){
	
	msg := UDPMessage{"I'm alive",1}
	for{
		sendUdpMessage(msg)
		time.Sleep(1*time.Second)
	}
}

func Backup(rec_channel chan UDPMessage){
	timeChannel := make(chan bool,1)
	go recieveUdpMessage(rec_channel)
	for{
		go timeOut(timeChannel)
		select{
			case <-rec_channel: 
				fmt.Println("KONTAKT MED MASTER")	
			case <-timeChannel:
				go primary()
				fmt.Println("Startar ny master")
			}
	}
}









