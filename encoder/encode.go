package main

import (
	"flag"
	"fmt"
	"github.com/vale1410/bule/sat"
	"github.com/vale1410/bule/sorters"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var f = flag.String("f", "test.dat", "Path of the file specifying the Problem.")
var out = flag.String("o", "data.lp", "Path of data file")
var typeIntersect = flag.String("intersect", "simple", "Type of encoding of non-intersection constriant for tasks.")
var ver = flag.Bool("ver", false, "Show version info.")
var asp = flag.Bool("asp", false, "Output instance in Gringo ASP format.")
var nWorkers = flag.Int("max", 0, "Max Number of Workers")

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

	if *asp {
		printASP(tasks, workers)
	} else {
		printSAT(tasks, workers)
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
	id     int
	start  int
	end    int
	worker map[int]bool
}

type ByStart []Task
type ByEnd []Task

func (a ByStart) Len() int             { return len(a) }
func (a ByStart) Swap(i, j int)        { a[i], a[j] = a[j], a[i] }
func (a ByStart) Less(i, j int) bool { //return a[i].start < a[j].start }
	if a[i].start < a[j].start {
		return true
	} else if a[i].start == a[j].start {
		return a[i].end < a[j].end
	} else {
		return false
	}
}

func (a ByEnd) Len() int      { return len(a) }
func (a ByEnd) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByEnd) Less(i, j int) bool {
	if a[i].end < a[j].end {
		return true
	} else if a[i].end == a[j].end {
		return a[i].start < a[j].start
	} else {
		return false
	}
}

type Worker struct {
	id     int
	skills []int
}

func parse(filename string) (tasks []Task, workers []Worker) {

	input, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Please specifiy correct path to instance. File does not exist: ", filename)
		panic(err)
	}

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
			tasks[id] = Task{id, start, end, make(map[int]bool, 0)}
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
				workers[id].skills = make([]int, n)

				for i, _ := range workers[id].skills {
					x, _ := strconv.Atoi(numbers[i+1])
					workers[id].skills[i] = x
				}
				id++
			}
		}
	}

	for _, w := range workers {
		for _, s := range w.skills {
			tasks[s].worker[w.id] = true
		}
	}

	return
}

func printASP(tasks []Task, workers []Worker) {

	output, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	for _, w := range tasks {
		s := "task(" + strconv.Itoa(w.id) + "," + strconv.Itoa(w.start) + "," + strconv.Itoa(w.end) + ").\n"
		if _, err := output.Write([]byte(s)); err != nil {
			panic(err)
		}
	}

	for _, w := range workers {
		for _, t := range w.skills {
			s := "worker2task(" + strconv.Itoa(w.id) + "," + strconv.Itoa(t) + ").\n"
			if _, err := output.Write([]byte(s)); err != nil {
				panic(err)
			}
		}
	}
}

func printSAT(tasks []Task, workers []Worker) {

	pAssign := sat.Pred("assign")
	pWorks := sat.Pred("works")

	sat.SetUp(4, sorters.Pairwise)
	var clauses sat.ClauseSet

	// at least one: simple clause
	for _, t := range tasks {
		lits := make([]sat.Literal, len(t.worker))
		i := 0
		for wId, _ := range t.worker {
			lits[i] = sat.Literal{true, sat.Atom{pAssign, wId, t.id}}
			i++
		}

		clauses.AddClause("al1", lits...)

		clauses.AddClauseSet(sat.CreateCardinality("am1", lits, 1, sorters.AtMost))

	}

	// count number of employees

	for _, w := range workers {
		for _, tId := range w.skills {

			l1 := sat.Literal{false, sat.Atom{pAssign, w.id, tId}}
			l2 := sat.Literal{true, sat.Atom{pWorks, w.id, 0}}

			clauses.AddClause("wrk", l1, l2)
		}
	}

	lits := make([]sat.Literal, len(workers))
	for i, w := range workers {
		lits[i] = sat.Literal{true, sat.Atom{pWorks, w.id, 0}}
	}

	clauses.AddClauseSet(sat.CreateCardinality("cWo", lits, *nWorkers, sorters.AtMost))

	// intersections on the timeline: two ways to do it
	// 1) list all intersecting tasks
	// 2) find maximal cliques in the interval graph, and post for that

	for _, w := range workers {

		ts := make([]Task, len(w.skills))
		for i, s := range w.skills {
			ts[i] = tasks[s]
		}
		sort.Sort(ByStart(ts))

		switch *typeIntersect {
		case "simple":
			for i, t1 := range ts {
				for j := i + 1; j < len(ts); j++ {
					t2 := ts[j]
					if t2.start < t1.end {
						l1 := sat.Literal{false, sat.Atom{pAssign, w.id, t1.id}}
						l2 := sat.Literal{false, sat.Atom{pAssign, w.id, t2.id}}
						clauses.AddClause("isc1", l1, l2)
					}
				}
			}
		case "clique":
			// find the maximal cliques in the interval graph and pose AMO on them

			clique := make([]Task, 0)

			for _, t := range ts {

				sort.Sort(ByEnd(clique))

				//todo: use a priority queue, e.g. heap
				//first one is earliest end time

				if len(clique) > 0 && clique[0].end <= t.start {
					// max clique reached
					//output the maximal clique!

					if len(clique) > 1 {

						lits := make([]sat.Literal, len(clique))

						for i, c := range clique {
							lits[i] = sat.Literal{true, sat.Atom{pAssign, w.id, c.id}}
							fmt.Print(c.id, "(", c.start, ",", c.end, ") ")
						}
						fmt.Println()

						//fmt.Println("clique", w.id, lits)
						clauses.AddClauseSet(sat.CreateCardinality("cli", lits, 1, sorters.AtMost))
					}

					//start removing elements:
					for len(clique) > 0 && clique[0].end <= t.start {
						clique = clique[1:]
					}
				}
				clique = append(clique, t)
			}
			if len(clique) > 1 {

				lits := make([]sat.Literal, len(clique))

				for i, c := range clique {
					lits[i] = sat.Literal{true, sat.Atom{pAssign, w.id, c.id}}
				}

				//fmt.Println("clique", w.id, lits)
				clauses.AddClauseSet(sat.CreateCardinality("cli", lits, 1, sorters.AtMost))
			}

		default:
			panic("Type not implemented")
		}
	}

	g := sat.IdGenerator(len(clauses) * 7)
	g.GenerateIds(clauses)

	//g.Filename = strings.Split(*f, ".")[0] + ".cnf"
	//g.Filename = *out

	if *dbg {
		g.PrintDebug(clauses)
	} else {
		g.PrintClausesDIMACS(clauses)
	}
}
