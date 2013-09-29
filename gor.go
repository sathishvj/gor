package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
	log.Println("HomeDir is:", u.HomeDir)

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
		fmt.Println("gorfile contains::", gorList)
	}

	list := strings.Split(gorList, ":")

	if *listAll {
		for i := 0; i < len(list) && len(list) > 0; i++ {
			fmt.Println(i+1, ":", list[i])
		}
		return
	}

	files := flag.Args()
	fmt.Println("Command line files are:", files)
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
	fmt.Println("Runfile is:", runFile)

	err = ioutil.WriteFile(gorFile, []byte(gorList), 0644)

	fmt.Println("Going to run:", runFile)

	cmd := exec.Command("go", "run", runFile)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(out.String())
}
