package helper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// INFO exported
var INFO *log.Logger

// ERROR exported
var ERROR *log.Logger

func init() {
	fmt.Println("Checking/creating logs folder")
	if _, err := os.Stat("../logs"); os.IsNotExist(err) {
		fmt.Println("Creating logs folder")

		os.MkdirAll("../logs", 0777)
	}
	absPath, err := filepath.Abs("../logs")
	fmt.Println("Created logs folder at : ", absPath)
	if err != nil {
		fmt.Println("Error reading given path:", err)
	}

	fmt.Println("Creating log file")
	generalLog, err := os.OpenFile(absPath+"/git-churn-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	INFO = log.New(generalLog, "INFO:\t", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	ERROR = log.New(generalLog, "ERROR:\t", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}
