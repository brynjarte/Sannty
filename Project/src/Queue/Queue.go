package Queue

import (
	"Source"
	//"driver"
	//"os"
	//"os/exec"
)


const (
	UP = 0
	DOWN = 1
)


type node struct{
	value Source.ButtonMessage
	next *node
}

type linkedList struct{
	head *node
	last *node
	length int
}

var allExternalQueues  = make(map[int] [][]int)
var allElevatorsInfo = make(map[int] Source.Elevator)


var queue = linkedList{nil,nil,0}
var qList [] int

/*
func init(){
	//var queue = linkedList{nil,nil,0} 
	//fetchMyQueue()
	//allQueues[elevatorID] = queue
}*/


func Queue(addOrderChannel chan Source.ButtonMessage, removeOrderChannel chan int, nextOrderChannel chan int, checkOrdersChannel chan int, orderInEmptyQueue chan int, orderRemovedChannel chan int, newElevInfoChannel chan Source.Elevator, fromUdpToQueue chan Source.Messagee){//, findBestElevator chan Source.ButtonMessage ){
	direction := -1
	currentFloor := 0
	elevatorID := 2
	//queueInit()
	
	for{
		/*
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
		PrintQueue()
		*/
		select{
			case newOrder := <- addOrderChannel:
				go addOrder(elevatorID, newOrder, currentFloor, direction, orderInEmptyQueue)
			
			
				
				

			case <- removeOrderChannel: 
				//println("QUEUE: REMOVEORDER")
				go removeOrder(orderRemovedChannel)
			
				
			case floor := <- checkOrdersChannel:
				
				//println("QUEUE: CHECKORDER DIRECTION:", direction)
				currentFloor = floor
				nextOrderedFloor := checkOrders(1)
				nextOrderChannel <- nextOrderedFloor
				direction = nextOrderedFloor - currentFloor
				if(direction > 0){
					direction = UP
				}else{
					direction = DOWN
				}

			case newUpdate := <- fromUdpToQueue:
				if(!newUpdate.FromMaster){				
					if(newUpdate.NewOrder && newUpdate.MessageTo == elevatorID){
						addOrderChannel <- newUpdate.Button
						//go recieveExternalQueue(newUpdate.MessageTo, newUpdate.Button)
					} else if( newUpdate.CompletedOrder && newUpdate.MessageTo != elevatorID){
						go recieveExternalQueue(newUpdate.MessageTo, newUpdate.Button)
					} else if (newUpdate.UpdatedElevInfo){
						go updateElevInfo(newUpdate.ElevInfo)
					} // FORTSETT MED SLAVE
					
						
				
			//case <- findBestElevatorChannel:

		}
	}
}

func checkOrders(elevatorID int) int {
	if queue.head == nil {
		return -1	
	}else {
		return queue.head.value.Floor
	}
}

func addOrder(elevatorID int, order Source.ButtonMessage, currentFloor int, movingDirection int, orderInEmptyQueue chan int) {

	var newOrder = node{order, nil}
	
	if (queue.length == 0) {
		queue.head = &newOrder
		queue.last = &newOrder
		queue.length = 1
		orderInEmptyQueue <- 1
		return
	} else if (queue.length == 1) {
		if equalOrders(queue.head.value, order) {
			return
		} else {
			queue.length++
			if equalOrders(compareOrders(queue.head.value, order, currentFloor, movingDirection), order) {
				newOrder.next = queue.last
				queue.head = &newOrder
			} else {
				queue.head.next = &newOrder
				queue.last = &newOrder
			}
		return
		}
	} else {

		var nodePointer *node = queue.head
		if equalOrders(nodePointer.value, order) {
			return
		} else if equalOrders(compareOrders((*nodePointer).value, order, currentFloor, movingDirection), order) {
			newOrder.next = queue.head
			queue.head = &newOrder
			queue.length++
			return
		}
		for i:=0; i < queue.length-1; i++ {
			 if equalOrders((*nodePointer).next.value, order) {
			 	return
			 } else {
			 	if equalOrders(compareOrders((*nodePointer).next.value, order, currentFloor, movingDirection), order) {
					newOrder.next = (*nodePointer).next
					(*nodePointer).next = &newOrder
					queue.length++
					return
				} else {
					nodePointer = (*nodePointer).next
				}
			 }
		}
		queue.last.next = &newOrder
		queue.last = &newOrder
		queue.length++
	}
}

func compareOrders(oldOrder Source.ButtonMessage, newOrder Source.ButtonMessage, currentFloor int, direction int) Source.ButtonMessage {
	
	if newOrder.Button == Source.BUTTON_COMMAND {
		if newOrder.Floor < currentFloor {
			//direction DOWN
			if oldOrder.Floor >  newOrder.Floor {
				return oldOrder
			} else if oldOrder.Floor < newOrder.Floor {
				return newOrder
			} 
		} else if newOrder.Floor > currentFloor {
			//direction UP
			if oldOrder.Floor <  newOrder.Floor {
				return oldOrder
			} else if oldOrder.Floor > newOrder.Floor {
				return newOrder
			} else if newOrder.Floor == currentFloor {
				return oldOrder
			}	
		}
	} else if newOrder.Button == Source.BUTTON_CALL_DOWN {
		if direction == UP {
			if (oldOrder.Button == Source.BUTTON_CALL_DOWN && oldOrder.Floor < newOrder.Floor){
				return newOrder
			} else if (oldOrder.Button != Source.BUTTON_CALL_DOWN && oldOrder.Floor < currentFloor) {
				return newOrder
			} else {
				return oldOrder
			}
		} else if direction == DOWN {
			if (oldOrder.Floor> newOrder.Floor) {
				return oldOrder
			} else {
				return newOrder
			}
		}
	} else if newOrder.Button== Source.BUTTON_CALL_UP {
		if direction == DOWN {
			if (oldOrder.Button== Source.BUTTON_CALL_UP && oldOrder.Floor > newOrder.Floor) {
				println("if")
				return newOrder
			}  else if (oldOrder.Button != Source.BUTTON_CALL_UP && oldOrder.Floor > currentFloor) {
				println("elseif")				
				return newOrder
			} else {
				println("else")
				return oldOrder
			}
		} else if direction == UP {
			if (oldOrder.Floor< newOrder.Floor) {
				return oldOrder
			} else {
				return newOrder
			}
		}
	}
	return oldOrder
}

func equalOrders(oldOrder Source.ButtonMessage, newOrder Source.ButtonMessage) bool {
	return (oldOrder.Floor == newOrder.Floor && oldOrder.Button== newOrder.Button)
}

func removeOrder(orderRemovedChannel chan int) {
	//nodePointer := queue.head

	for {
		if (queue.length > 1) {
			if (queue.head.value.Floor == queue.head.next.value.Floor) {
				queue.head = queue.head.next
				queue.length--
			} else {
				queue.head = queue.head.next
				queue.length--
				orderRemovedChannel <- 1
				return
			}
		} else {
			queue.head = nil
			queue.length = 0
			orderRemovedChannel <- 1
			return
		}
	}
}

func clearAllOrders(){
	queue.head = nil
	queue.last = nil
	queue.length = 0
}

func PrintQueue() {
	//println("KØ: ", queue.length)
	if queue.length == 0 {
		return
	}
	println("Element 1:\nEtasje: ", queue.head.value.Floor, "\tKnapp: ", queue.head.value.Button,"\n")
	var newOrder *node
	newOrder = queue.head.next
	for i:=1 ; i < queue.length-1; i++ {
		println("Element", i+1,":\nEtasje: ", newOrder.value.Floor, "\tKnapp: ", newOrder.value.Button,"\n")
		newOrder = newOrder.next
	}
}
/*
func fetchMyQueue() {
	q := FileHandler.Read(&NumOfElevs, &NumOfFloors)
	
	//queue.length = q[0], NEI!!
	clearAllOrders()
		
	for j:=0; j < len(q); j+=2 {
		ord := Source.ButtonMessage{q[j],q[j+1]}
		addOrder(elevatorID , ord, currentFloor , movingDirection)
	}
}
*/

func recieveExternalQueue(elevatorID int, button Source.ButtonMessage) {
		
	allExternalQueues[elevatorID][button.Floor][button.Button] = button.Value
	
}

func updateElevInfo(newElevInfo Source.Elevator){
	allElevatorsInfo[newElevInfo.ID] = newElevInfo
}
/*
func findBestElevator() int {
	best = 1
	min := findElevatorCost(1, newOrder)
	
	for elev := 2; elev <= NumOfElevs; elev++ {
		itsCost := findElevatorCost(elev, newOrder)
		
		if min > itsCost {
			best = elev
		}
	} 
	return best
}
*/
/*
func findElevatorCost(elevatorID int, newOrder Source.ButtonMessage) int {
	//for my ID:
	//		find which position in queue will new order get
	// Omtrent lik addOrder, bare returnerer pos, uten å legge til i kø
	return position
}
*/
 /*
func saveAndSendQueue() {
	
	qList = append([]int(nil), queue.head.value.Floor)
	qList = append(qList, queue.head.value.Button)
	
	var newOrder *node
	newOrder = queue.head.next
	for i:=1 ; i < queue.length; i++ {
		qList = append(qList, newOrder.value.Floor)
		qList = append(qList, newOrder.value.Button)
		newOrder = newOrder.next
	}
	
	FileHandler.Write(NumElevs, NumOfFloors, qList)
	//UDP.sendQueue()
	
}
*/
