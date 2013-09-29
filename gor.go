package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func main() {
	listAll := flag.Bool("l", false, "-l")
	flag.Parse()

	u, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user.", err.Error())
		return
	}
	//fmt.Println("HomeDir is:", u.HomeDir)

	gorFile := u.HomeDir + string(os.PathSeparator) + ".gorfile"
	_, err = os.Stat(gorFile)
	var gorList string
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		b, err := ioutil.ReadFile(gorFile)
		if err != nil {
			panic(err)
		}
		gorList = string(b)
		//fmt.Println("gorfile contains::", gorList)
	}

	list := strings.Split(gorList, ":")

	if *listAll {
		for i := 0; i < len(list) && len(list) > 0; i++ {
			fmt.Println(i+1, ":", list[i])
		}
		return
	}

	files := flag.Args()
	//fmt.Println("Command line files are:", files)
	if len(files) == 0 && len(gorList) == 0 {
		fmt.Println("Err! One go file needs to be given as there was nothing in the previous list either.")
		return
	}
	if len(files) > 1 {
		fmt.Println("Err! Specify only one file please.")
		return
	}

	var runFile string

	if len(files) != 0 {
		if len(gorList) > 0 {
			gorList += ":"
		}
		gorList += files[0]
		runFile = files[0]
	} else {
		runFile = list[len(list)-1]
	}
	//fmt.Println("Runfile is:", runFile)

	// for now, just write the last one.
	//err = ioutil.WriteFile(gorFile, []byte(gorList), 0644)
	err = ioutil.WriteFile(gorFile, []byte(runFile), 0644)

	fmt.Println("Going to run:", runFile, "...")

	cmd := exec.Command("go", "run", runFile)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error getting StdoutPipe")
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error getting StdErrPipe")
		panic(err)
	}

	outPipe := bufio.NewReader(stdout)
	errPipe := bufio.NewReader(stderr)
	outCh := make(chan string)
	errCh := make(chan string)

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting cmd")
		panic(err)
	}
	go getStdOutput(outCh, outPipe)
	go getStdOutput(errCh, errPipe)
outside:
	for {
		select {
		case s, ok := <-outCh:
			if !ok {
				break outside
			}
			//fmt.Println("From Out:", s)
			fmt.Println(s)
		case s, ok := <-errCh:
			if !ok {
				break outside
			}
			fmt.Println("Err!", s)
		}
	}
	//fmt.Println("Finished")

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Waiting error", err)
		return
	}
}

func getStdOutput(c chan string, p *bufio.Reader) {
	for {
		line, _, err := p.ReadLine()
		if err != nil {
			if err != io.EOF {
				fmt.Println("ReadLine error:", err)
				//return
			}
			break
		}
		//fmt.Println("Got line:", string(line))
		c <- string(line)
	}
	//fmt.Println("Closing c")
	close(c)
}
