package FileHandler

import ( 
	"os"
	"fmt"
	//"encoding/binary"
	//"bytes"
	//"Queue"
	//"container/list"
)

type Directions struct{
	UP, DOWN int 
}

func Read(NumOfElevs *int, NumOfFloors *int) [] int{
	
	fd, err := os.Open("maxmekker.txt")
	if err != nil {
		panic(err)
	}

	var buf int

	_, err = fmt.Fscanf(fd, "%d\n", &buf)

	*NumOfElevs = buf
	//_,_ = fmt.Fscanf(fd, "%s", &buf)

	_, err = fmt.Fscanf(fd, "%d\n", &buf)

	*NumOfFloors = buf

	_, err = fmt.Fscanf(fd, "%d\n", &buf)
	
	qLength := buf

	QueueList := make([]int, 0)
	
	for i:=0; i < qLength; i++{
		
		_, err = fmt.Fscanf(fd, "%d", &buf)
		QueueList  = append(QueueList, buf)
		
		_, err = fmt.Fscanf(fd, "%d\n", &buf)
		QueueList = append(QueueList, buf)
	}
	
	return QueueList
}
/*
func Write(numOfElevs int, numOfFloors int, queue []Directions,ID int) {
	// open output file
    fo, err := os.Create("output.txt")
    if err != nil {
        panic(err)
    }
    // close fo on exit and check for its returned error
    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()
	if _, err := fmt.Fprintf(fo,"%d\n",numOfElevs); err != nil {
            panic(err)
    }
            
    if _, err := fmt.Fprintf(fo,"%d\n",numOfFloors); err != nil {
            panic(err)
    }
    println(numOfFloors," Antall etager")
	for i:=0;i<numOfFloors;i++{
	        // write a chunk
        if _, err := fmt.Fprintf(fo,"%d\t",queue[i].UP); err != nil {
            panic(err)
        }
        if _, err := fmt.Fprintf(fo,"%d\n",queue[i].DOWN); err != nil {
            panic(err)
        }
	}
	//for {}
}
*/

func Write(numOfElevs int, numOfFloors int, length int, queue []int) {
	fo, err := os.Create("output.txt")
    if err != nil {
        panic(err)
    }
    // close fo on exit and check for its returned error
    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()
	if _, err := fmt.Fprintf(fo,"%d\n",numOfElevs); err != nil {
            panic(err)
    }
            
    if _, err := fmt.Fprintf(fo,"%d\n",numOfFloors); err != nil {
            panic(err)
    }
    
    if _, err := fmt.Fprintf(fo,"%d\n",length); err != nil {
            panic(err)
    }
    
	for i := 0;i < len(queue);i+=2{
	        // write a chunk
        if _, err := fmt.Fprintf(fo,"%d\t",queue[i]); err != nil {
            panic(err)
        }
        if _, err := fmt.Fprintf(fo,"%d\n",queue[i+1]); err != nil {
            panic(err)
        }
	}
}
