package main

import (
	embedScript "DiskBurn/embedScript"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)


func runTest(disk string, arguments []string){
	var tempArgument []string

	tempArgument = append(tempArgument,"./burntest.sh")
	tempArgument = append(tempArgument, arguments...)
	tempArgument = append(tempArgument, disk)

	cmd := exec.Command("sudo", tempArgument...)

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	multiOut := io.MultiReader(stdout, stderr)
	scannerOut := bufio.NewScanner(multiOut)

	cmd.Start()
	cmd.Process.Release()
	for scannerOut.Scan() {
		m := scannerOut.Text()
		fmt.Println(m)
	}

}

func rootTest(){
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		log.Fatal(err)
	}
	if i == 0 {
		log.Println("Confirmed root user, running tests.")
	} else {
		log.Fatal("This program must be run as root! (sudo)")
	}
}

func main(){
	rootTest()

	scriptText, _ := embedScript.Fs.ReadFile("burntest.sh")
	tempFile, _ := os.OpenFile("burntest.sh", os.O_CREATE|os.O_RDWR, 0755)
	tempFile.Write(scriptText)
	tempFile.Close()

	if _, err := os.Stat("disks.txt"); errors.Is(err, os.ErrNotExist) {
		log.Println("Disk file not found, creating empty file.")
		os.OpenFile("disks.txt", os.O_CREATE|os.O_RDWR,0644)
		log.Println("disks.txt created")
		log.Fatal("Please add disks to disk file and run again")
	}

	diskList, _ := ioutil.ReadFile("disks.txt")

	diskString := strings.Split(string(diskList), "\n")
	arguments := os.Args[1:]
	for _, disk := range diskString{
		runTest(disk, arguments)
	}

}