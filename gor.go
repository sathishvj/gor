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
	listAll := flag.Bool("l", false, "to list all entries")
	help := flag.Bool("h", false, "to display this usage listing")
	curr := flag.Bool("c", false, "to show go file for current directory")
	rem := flag.Bool("r", false, "to remove the entry corresponding to current directory")
	del := flag.Bool("x", false, "to delete the current .gorrc")
	build := flag.Bool("b", false, "to use build tool (default is run)")
	test := flag.Bool("t", false, "to use test tool (default is run)")
	doc := flag.Bool("d", false, "to use test tool (default is run)")
	flag.Parse()
	anyFlags := *listAll || *help || *curr || *rem || *del

	if (*build && *test) || (*test && *doc) || (*build && *doc) {
		fmt.Println("Error: You can have only one of build/test/doc.\n")
		usage()
		return
	}
	tool := "run"
	if *build {
		tool = "build"
	}
	if *test {
		tool = "test"
	}
	if *doc {
		tool = "doc"
	}

	if *help {
		usage()
		return
	}

	u, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user.", err.Error())
		return
	}
	//fmt.Println("HomeDir is:", u.HomeDir)

	gorrc := u.HomeDir + string(os.PathSeparator) + ".gorrc"
	gorExists := true
	_, err = os.Stat(gorrc)
	if err != nil {
		if os.IsNotExist(err) {
			gorExists = false
		} else {
			fmt.Println("Error reading file", gorrc, ":", err)
			return
		}
	}

	if *del {
		if !gorExists {
			return
		}
		err = os.Remove(gorrc)
		if err != nil {
			fmt.Println("Error removing file", gorrc, ":", err)
		} else {
			fmt.Println("Removed", gorrc)
		}
		return
	}

	var gorList string
	if gorExists {
		b, err := ioutil.ReadFile(gorrc)
		if err != nil {
			fmt.Println("Error reading file", gorrc, ":", err)
		}
		gorList = string(b)
		//fmt.Println("gorfile contains::", gorList)
	}

	m := make(map[string]map[string]string)
	lines := strings.Split(gorList, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, "=")
		if _, exists := m[parts[0]]; exists {
			m[parts[0]][parts[1]] = parts[2]
		} else {
			tmpM := make(map[string]string)
			tmpM[parts[1]] = parts[2] //run = file
			m[parts[0]] = tmpM        //dir=run=file
		}
	}

	if *listAll {
		if !gorExists {
			return
		}
		for k, v := range m {
			for k1, v1 := range v {
				fmt.Println(k, "=", k1, "=", v1)
			}
		}
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory.", err.Error())
		return
	}
	//fmt.Println("pwd is:", pwd)

	if *rem {
		delete(m[pwd], tool)
		write(m, gorrc)
		return
	}

	if *curr {
		tmpM, exists := m[pwd]
		if !exists {
			return
		}
		for k, v := range tmpM {
			fmt.Println(pwd, "=", k, "=", v)
		}
		return
	}

	// if any flags are on, then it is just a options setting request
	if anyFlags {
		return
	}

	var runArgs string
	var ok bool
	runArgs, ok = m[pwd][tool]
	if len(flag.Args()) == 0 && !ok {
		fmt.Println("Error! No arguments given and no previous entry in gorfile.\n")
		usage()
		return
	} else if len(flag.Args()) > 0 {
		runArgs = strings.Join(flag.Args(), " ")
	}

	//fmt.Println("Runargs is:", runArgs)

	if !ok || len(flag.Args()) > 0 { //this entry didn't exist before
		if _, exists := m[pwd]; exists {
			m[pwd][tool] = runArgs
		} else {
			tmpM := make(map[string]string)
			tmpM[tool] = runArgs //run = file
			m[pwd] = tmpM        //dir=run=file
		}
		err = write(m, gorrc)
		if err != nil {
			fmt.Println("Err! Unable to write to ", gorrc, ":", err)
			return
		}
	}

	fmt.Println("go", tool, runArgs)
	cmd := exec.Command("go", tool, runArgs)
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
			//fmt.Println("Err!", s)
			fmt.Println(s)
		}
	}
	//fmt.Println("Finished")

	err = cmd.Wait()
	if err != nil {
		// Typically a compile error. Don't print any message.
		//fmt.Println("Waiting error", err)
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

func usage() {
	pgm := os.Args[0]
	fmt.Printf("%s will re-run/build/test/doc with previous arguments in the current directory.  Usage:\n", pgm)
	flag.PrintDefaults()
	fmt.Printf("\n%s hello.go\n%s\n", pgm, pgm)
	return
}

func write(m map[string]map[string]string, f string) error {
	var s string
	for k, v := range m {
		for k1, v1 := range v {
			s = s + k + "=" + k1 + "=" + v1 + "\n"
		}
	}
	err := ioutil.WriteFile(f, []byte(s), 0644)
	if err != nil {
		return err
	}
	return nil
}
