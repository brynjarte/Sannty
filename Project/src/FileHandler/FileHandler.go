package FileHandler

import ( 
	"os"
	"fmt"
	"time"
)

func Read(NumOfElevs *int, NumOfFloors *int) [] int{
	
	fd, err := os.Open("backup.txt")
	if err != nil {
		*NumOfElevs = 3
		*NumOfFloors = 4
		fmt.Println(err)
		ErrorLog(err)
		return nil
	}
	
	defer func() {
		if err := fd.Close(); err != nil {
			panic(err)
		}
    	}()

	var buf int

	_, err = fmt.Fscanf(fd, "%d\n", &buf)
	if err != nil {
		*NumOfElevs = 3
	} else{
		*NumOfElevs = buf
	}

	_, err = fmt.Fscanf(fd, "%d\n", &buf)
	if err != nil {
		*NumOfFloors = 4
	} else{
		*NumOfFloors = buf
	}

	qLength := 0	

	_, err = fmt.Fscanf(fd, "%d\n", &buf)
	if err != nil {
		qLength = 0
	} else{
		qLength = buf
	}

	QueueList := make([]int, 0)
	
	for i:=0; i < qLength; i++{
		
		_, err = fmt.Fscanf(fd, "%d", &buf)
		QueueList  = append(QueueList, buf)
		
		_, err = fmt.Fscanf(fd, "%d\n", &buf)
		QueueList = append(QueueList, buf)
	}
		
	return QueueList
}

func Write(numOfElevs int, numOfFloors int, length int, queue []int) {
    fo, err := os.Create("backup.txt")
    if err != nil {
        panic(err)
    }

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
        if _, err := fmt.Fprintf(fo,"%d\t",queue[i]); err != nil {
            panic(err)
        }
        if _, err := fmt.Fprintf(fo,"%d\n",queue[i+1]); err != nil {
            panic(err)
        }
    }
}

func ErrorLog(occuredError error) {

	timestamp := time.Now().Format("02-01-2006 15:04:05 CET")
	
	fe, err := os.OpenFile("errlog.txt", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
	        panic(err)
    	}
    	
	defer func() {
		if err := fe.Close(); err != nil {
			panic(err)
		}
 	}()

	if _, err = fmt.Fprintf(fe,"Error occured at time %s \nError type:\t", timestamp); err != nil {
		panic(err)
    }
    	
    if _, err = fmt.Fprintf(fe,"%s\n\n",occuredError); err != nil {
		panic(err)
    }
}
