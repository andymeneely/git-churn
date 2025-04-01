package helper

import (
	"fmt"
	"log"
	"os"
)

// INFO exported
var INFO *log.Logger

// ERROR exported
var ERROR *log.Logger

func init() {
	fmt.Println("Creating log file")
	generalLog, err := os.OpenFile("git-churn.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	INFO = log.New(generalLog, "INFO:\t", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	ERROR = log.New(generalLog, "ERROR:\t", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}
