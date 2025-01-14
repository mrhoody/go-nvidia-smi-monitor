package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main(){
	var testData gpuUsageData
	testData.timeStamp = "2021-09-01 12:00:00"
	testData.currentGPUMemoryUtilisation = "50%"
	testData.GPUCapacityUtilisation = "50%"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// get environment variables
	configureNumericEnvVariable("SCAN_INTERVAL", 5)
	configureNumericEnvVariable("SCAN_DURATION", 5)


}

type gpuUsageData struct {
	timeStamp string
	currentGPUMemoryUtilisation string
	GPUCapacityUtilisation string
}

func configureNumericEnvVariable(envVariableName string, fallBackValue int) int {

	// check if environment variable exists
	if os.Getenv(envVariableName) == "" {
		log.Println("Environment variable ", envVariableName, " not set - using fallback value: ", fallBackValue)
		return fallBackValue
	}

	if s, err := strconv.Atoi(os.Getenv(envVariableName)); err == nil {
		log.Println(envVariableName, " = ", s)
		return s
	} else {
		log.Println("Error converting ",envVariableName, " to int - using fallback value: ", fallBackValue)
		return fallBackValue
	}
}

func captureTerminalOutputFromCommand(command string) string {
	out, err := exec.Command(command).Output()
	log.Default().Println("Command: ", command)
	if err != nil {
        log.Panicln(err)
		return string(err.Error())
    }
	return string(out)
}

func appendToCSV (filename string, data gpuUsageData) {
	// check if file ends with .csv
	if !strings.HasSuffix(filename, ".csv") {
		log.Panic("Filename must end with .csv")
	}

	// check if file already exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// file does not exist
		log.Println("File does not exist, creating a new file with filename: ", filename)
		_, err := os.Create(filename); if err != nil {
			log.Panic("Error creating file: ", err)
		}
	}

	// append data to file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend) ; if err != nil {
		log.Panic("Error opening file: ", err)
	}

	// prepare buffer to write to file
	var b strings.Builder
	b.WriteString(data.timeStamp)
	b.WriteString(",")
	b.WriteString(data.currentGPUMemoryUtilisation)
	b.WriteString(",")
	b.WriteString(data.GPUCapacityUtilisation)
	b.WriteString("\n")

	// write to file

	_, err = f.WriteString(b.String()); if err != nil {
		log.Panic("Error writing to file: ", err)
	}
	log.Println("GPU usage data logged.")
}