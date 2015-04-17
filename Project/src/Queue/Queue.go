package Queue

import (
	"driver"
)


const (
	UP = 0
	DOWN = 1
)

type node struct{
	value driver.ButtonMessage
	next *node
}

type linkedList struct{
	head *node
	last *node
	length int
}

//var allQueues map[int]linkedList

var queue = linkedList{nil,nil,0}
var qList [] int
/*
func init(){
	//var queue = linkedList{nil,nil,0} 
	//fetchMyQueue()
	//allQueues[elevatorID] = queue
}*/


func Queue(addOrderChannel chan driver.ButtonMessage, removeOrderChannel chan int, nextOrderChannel chan int, checkOrdersChannel chan int, orderInEmptyQueueChannel chan int){//, findBestElevator chan driver.ButtonMessage ){
	direction := 0
	currentFloor := 0
	//queueInit()
	
	for{
		//PrintQueue()
		select{
			case newOrder := <- addOrderChannel:
				
				addOrder(1, newOrder, currentFloor, direction, orderInEmptyQueueChannel)

			case <- removeOrderChannel: 
				println("QUEUE: REMOVEORDER")
				removeOrder()
				
				
				
			case floor := <- checkOrdersChannel:
				println("QUEUE: CHECKORDER DIRECTION:", direction)
				currentFloor = floor
				nextOrderChannel <- checkOrders(1)
				
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

func addOrder(elevatorID int, order  driver.ButtonMessage, currentFloor int, movingDirection int, orderInEmptyQueueChannel chan int) {
	var noden = node{order, nil}
	
	if (queue.length == 0) {
		queue.head = &noden
		queue.last = &noden
		queue.length = 1
		orderInEmptyQueueChannel <- 1
	} else if (queue.length == 1) {
		if equalOrders(queue.head.value, order) {
			return
		} else {
			queue.length = 2
			if equalOrders(compareOrders(queue.head.value, order, currentFloor, movingDirection), order) {
				noden.next = queue.last
				queue.head = &noden
			} else {
				queue.head.next = &noden
				queue.last = &noden
			}
		}
	} else {
		var nodePointer = &node{queue.head.value, queue.head.next}
		if equalOrders(nodePointer.value, order) {
			return
		} else if equalOrders(compareOrders(nodePointer.value, order, currentFloor, movingDirection), order) {
			noden.next = nodePointer
			queue.head = &noden
			queue.length++
			return
		}
		for i:=0; i < queue.length-1; i++ {
			 if equalOrders(nodePointer.next.value, order) {
			 	return
			 } else {
			 	if equalOrders(compareOrders(nodePointer.next.value, order, currentFloor, movingDirection), order) {
					noden.next = nodePointer.next
					nodePointer.next = &noden
					queue.length++
					return
				} else {
					nodePointer = nodePointer.next
				}
			 }
		}
		queue.last.next = &noden
		queue.last = &noden
		queue.length++
	}
}

func compareOrders(oldOrder driver.ButtonMessage, newOrder driver.ButtonMessage, currentFloor int, direction int) driver.ButtonMessage {
	if newOrder.Button== driver.BUTTON_COMMAND {
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
	} else if newOrder.Button== driver.BUTTON_CALL_DOWN {
		if direction == UP {
			if (oldOrder.Button== driver.BUTTON_CALL_DOWN && oldOrder.Floor < newOrder.Floor){
				return newOrder
			} else if (oldOrder.Button!= driver.BUTTON_CALL_DOWN && oldOrder.Floor < currentFloor) {
				return newOrder
			} else {
				return oldOrder
			}
		} else if direction == DOWN {
			if (oldOrder.Button> newOrder.Button) {
				return oldOrder
			} else {
				return newOrder
			}
		}
	} else if newOrder.Button== driver.BUTTON_CALL_UP {
		if direction == DOWN {
			if (oldOrder.Button== driver.BUTTON_CALL_UP && oldOrder.Floor > newOrder.Floor) {
				return newOrder
			}  else if (oldOrder.Button!= driver.BUTTON_CALL_UP && oldOrder.Floor > currentFloor) {
				return newOrder
			} else {
				return oldOrder
			}
		} else if direction == UP {
			if (oldOrder.Button< newOrder.Button) {
				return oldOrder
			} else {
				return newOrder
			}
		}
	}
	return oldOrder
}

func equalOrders(oldOrder driver.ButtonMessage, newOrder driver.ButtonMessage) bool {
	return (oldOrder.Floor == newOrder.Floor && oldOrder.Button== newOrder.Button)
}

func removeOrder() {
	nodePointer := queue.head
	
	for {
		if (nodePointer.next != nil) {
			if (nodePointer.value.Floor == nodePointer.next.value.Floor) {
			nodePointer = queue.head.next
			queue.head = nodePointer
			queue.length--
			}
		} else {
			nodePointer = queue.head.next
			queue.head = nodePointer
			queue.length--
			break
		}
	}
}

func clearAllOrders(){
	queue.head = nil
	queue.last = nil
	queue.length = 0
}

func PrintQueue() {
	if queue.length == 0 {
		return
	}
	println("Element 1:\nEtasje: ", queue.head.value.Floor, "\tKnapp: ", queue.head.value.Button,"\n")
	var noden *node
	noden = queue.head.next
	for i:=1 ; i < queue.length; i++ {
		println("Element", i+1,":\nEtasje: ", noden.value.Floor, "\tKnapp: ", noden.value.Button,"\n")
		noden = noden.next
	}
}
/*
func fetchMyQueue() {
	q := FileHandler.Read(&NumOfElevs, &NumOfFloors)
	
	//queue.length = q[0], NEI!!
	clearAllOrders()
		
	for j:=0; j < len(q); j+=2 {
		ord := driver.ButtonMessage{q[j],q[j+1]}
		addOrder(elevatorID , ord, currentFloor , movingDirection)
	}
}
*/
/*
func recieveExternalQueue() {
	// newQ := noko---
	// its direction
	// its currFloor
	// its ID
	// "Lag" ny kø
	// allQueues[ID] = nyKø 
	// kun ha ordre i allQueues?
}
*/
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
func findElevatorCost(elevatorID int, newOrder driver.ButtonMessage) int {
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
	
	var noden *node
	noden = queue.head.next
	for i:=1 ; i < queue.length; i++ {
		qList = append(qList, noden.value.Floor)
		qList = append(qList, noden.value.Button)
		noden = noden.next
	}
	
	FileHandler.Write(NumElevs, NumOfFloors, qList)
	//UDP.sendQueue()
	
}
*/
