package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

//	"math"
//	"os/exec"
//	"time"
)

var f = flag.String("f", "test.dat", "Path of the file specifying the Problem.")
var out = flag.String("o", "data.lp", "Path of data file")
var model = flag.String("model", "model.lp", "path to model file")
var ver = flag.Bool("ver", false, "Show version info.")
var solve = flag.Bool("solve", false, "Pass problem to clasp and solve.")
var dbg = flag.Bool("d", false, "Print debug information.")
var dbgfile = flag.String("df", "", "File to print debug information.")

//var ttimeout = flag.Int("timeout", 10, "Timeout in seconds.")

var digitRegexp = regexp.MustCompile("([0-9]+ )*[0-9]+")

var dbgoutput *os.File

func main() {
	flag.Parse()

	if *dbgfile != "" {
		var err error
		dbgoutput, err = os.Create(*dbgfile)
		if err != nil {
			panic(err)
		}
		defer dbgoutput.Close()
	}

	debug("Running Debug Mode...")

	if *ver {
		fmt.Println(`Task Scheduling Formating Tool, Version Tag 0.1a
Copyright (C) NICTA and Valentin Mayer-Eichberger
License GPLv2+: GNU GPL version 2 or later <http://gnu.org/licenses/gpl.html>
There is NO WARRANTY, to the extent permitted by law.`)
		return
	}

	tasks, workers := parse(*f)

	debug("Tasks", tasks)
	debug("Workers", workers)

	if *solve {
		//solveDirect(tasks, workers)
	}
}

func debug(arg ...interface{}) {
	if *dbg {
		if *dbgfile == "" {
			fmt.Print("dbg: ")
			for _, s := range arg {
				fmt.Print(s, " ")
			}
			fmt.Println()
		} else {
			ss := "dbg: "
			for _, s := range arg {
				ss += fmt.Sprintf("%v", s) + " "
			}
			ss += "\n"

			if _, err := dbgoutput.Write([]byte(ss)); err != nil {
				panic(err)
			}
		}
	}
}

type Task struct {
	id    int
	start int
	end   int
}

type Worker struct {
	id    int
	skill []int
}

func parse(filename string) (tasks []Task, workers []Worker) {

	input, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Please specifiy correct path to instance. File does not exist: ", filename)
		panic(err)
	}

	output, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	lines := strings.Split(string(input), "\n")

	state := 0
	id := 0

	for _, l := range lines {

		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}

		numbers := strings.Fields(l)
		//		if digitregexp.matchstring(numbers[0]) {
		switch state {
		case 0:
			{
				debug("Is of Type 1", numbers[2] == "1")
				state = 1
			}
		case 1:
			{
				debug("Number of Jobs", numbers[2])
				n, _ := strconv.Atoi(numbers[2])
				tasks = make([]Task, n)
				state = 2
			}
		case 2:
			start, _ := strconv.Atoi(numbers[0])
			end, _ := strconv.Atoi(numbers[1])
			tasks[id] = Task{id, start, end}
			id++
			if id == len(tasks) {
				state = 3
			}
		case 3:
			{
				debug("Number of Workers", numbers[2])
				n, _ := strconv.Atoi(numbers[2])
				workers = make([]Worker, n)
				id = 0
				state = 4
			}
		case 4:
			{
				n, _ := strconv.Atoi(strings.TrimRight(numbers[0], ":"))
				workers[id].id = id
				workers[id].skill = make([]int, n)

				for i, _ := range workers[id].skill {
					x, _ := strconv.Atoi(numbers[i+1])
					workers[id].skill[i] = x
				}
				id++
			}
		}
	}

	for _, w := range tasks {
		s := "task(" + strconv.Itoa(w.id) + "," + strconv.Itoa(w.start) + "," + strconv.Itoa(w.end) + ").\n"
		if _, err := output.Write([]byte(s)); err != nil {
			panic(err)
		}
	}

	for _, w := range workers {
		for _, t := range w.skill {
			s := "worker2task(" + strconv.Itoa(w.id) + "," + strconv.Itoa(t) + ").\n"
			if _, err := output.Write([]byte(s)); err != nil {
				panic(err)
			}
		}
	}

	return
}

//func computeValue(warehouses []Warehouse, customers []Customer, assignment []int) (cost float64) {
//
//	added := make([]bool, len(warehouses))
//
//	for c, w := range assignment {
//
//		cost += customers[c].costs[w]
//
//		if !added[w] {
//			cost += warehouses[w].setup
//
//			added[w] = true
//		}
//	}
//	return
//}
//
//func solveDirect(warehouses []Warehouse, customers []Customer) {
//
//	result := make(chan Result)
//	optimum := make(chan bool)
//	timeout := make(chan bool, 1)
//	go func() {
//		time.Sleep(time.Duration(*ttimeout) * time.Second)
//		timeout <- true
//	}()
//
//	go solveProblem(result, optimum)
//
//	assignment := make([]int, len(customers))
//
//	stop := false
//	current := 0.0
//	optimal := 0
//
//	answer := 0
//
//	for !stop {
//		select {
//		case r := <-result:
//
//			if answer < r.i && parseResult(r.s, assignment) {
//				current = computeValue(warehouses, customers, assignment)
//			}
//		case b := <-optimum:
//			if b {
//				optimal = 1
//			}
//			stop = true
//		case <-timeout:
//			debug("recieved timeout from solver")
//			stop = true
//		}
//	}
//
//	close(result)
//	close(optimum)
//
//	fmt.Printf("%v %v\n", current, optimal)
//	for _, x := range assignment {
//		fmt.Printf("%v ", x)
//	}
//	fmt.Printf("\n")
//}
//
//func round(in float64) int {
//	if in-math.Floor(in) > 0.5 {
//		return int(math.Floor(in) + 1)
//	} else {
//		return int(math.Floor(in))
//	}
//}
//
//func parseResult(s string, assignment []int) bool {
//	ss := strings.Split(string(s), " ")
//
//	ok := len(assignment) == len(ss)
//
//	if ok {
//		for _, x := range ss {
//
//			if strings.HasPrefix(x, "assign") {
//				numbers := digitRegexp.FindAllString(x, -1)
//
//				if 2 == len(numbers) {
//
//					customer, b1 := strconv.Atoi(numbers[0])
//					warehouse, b2 := strconv.Atoi(numbers[1])
//					if b1 != nil || b2 != nil {
//						panic("bad conversion of numbers in result")
//					}
//					assignment[customer] = warehouse
//				} else {
//					ok = false
//					break
//				}
//			}
//		}
//	}
//
//	return ok
//}

//func solveProblem(result chan<- Result, optimum chan<- bool) {
//
//	gringo := exec.Command("gringo", *out, *model)
//	clasp := exec.Command("clasp", "--configuration=chatty", "--stat", "-t 4", "--time-limit", strconv.Itoa(*ttimeout))
//	clasp.Stdin, _ = gringo.StdoutPipe()
//
//	mw := NewClaspFilter(result, optimum)
//	clasp.Stdout = &mw
//
//	_ = clasp.Start()
//	_ = gringo.Run()
//	_ = clasp.Wait()
//
//}
//
//type Result struct {
//	s string
//	i int
//}
//
//type ClaspFilter struct {
//	result  chan<- Result
//	optimum chan<- bool // true: optimium, false: satisfiable
//	answer  int         // nth answer from clasp
//	backup  string      // keep string if last string no new line
//}
//
//func NewClaspFilter(result chan<- Result, optimum chan<- bool) (cf ClaspFilter) {
//	cf.result = result
//	cf.optimum = optimum
//	return
//}
//
//func (cf *ClaspFilter) Write(p []byte) (n int, err error) {
//
//	cf.backup += string(p)
//
//	lines := strings.Split(cf.backup, "\n")
//
//	if lines[len(lines)-1] == "" {
//		cf.backup = ""
//
//		for _, s := range lines {
//
//			if strings.Contains(s, "assign") {
//				cf.answer++
//				cf.result <- Result{s, cf.answer}
//			} else {
//				debug("clasp:", s)
//			}
//
//			if strings.Contains(s, "OPTIMUM FOUND") {
//				cf.optimum <- true
//			}
//			if strings.Contains(s, "SATISFIABLE") {
//				cf.optimum <- false
//			}
//		}
//	}
//
//	return len(p), nil
//}
