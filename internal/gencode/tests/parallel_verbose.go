// Program cooltool - A cool program
package main

/*
 * cooltool.go
 * A cool program
 * By MysteryDev
 * Created in the past
 * Last Modified in the past
 */

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	/* ProgramStart notes when the program has started for printing the
	elapsed time when the program ends. */
	ProgramStart = time.Now()

	/* Verbosef wil be a no-op if -verbose isn't given. */
	Verbosef = log.Printf
)

// Task contains the information necessary to accomplish a task.
type Task struct{}

func main() {
	/* Command-line flags. */
	var (
		noSummary = flag.Bool(
			"no-summary",
			false,
			"Don't print a summary on exit",
		)
		verbOn = flag.Bool(
			"verbose",
			false,
			"Enable verbose logging",
		)
		nPar = flag.Uint(
			"parallel",
			10,
			"Parallel task execution `count`",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options]

A cool program

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Work out verbose logging. */
	if !*verbOn {
		Verbosef = func(string, ...any) {}
	}

	/* Start some task executors. */
	var (
		ch = make(chan Task)
		wg sync.WaitGroup
	)
	for i := uint(0); i < *nPar; i++ {
		wg.Add(1)
		go taskExecutor(ch, &wg)
	}

	/* Send the tasks to be executed. */
	tasks, err := getTasks()
	if nil != err {
		log.Fatalf("Error getting tasks: %s", err)
	}
	for _, task := range tasks {
		ch <- task
	}

	/* Wait for the executors to finish executing. */
	close(ch)
	wg.Wait()

	/* All done. */
	if !*noSummary {
		log.Printf(
			"Done in %s.",
			time.Since(ProgramStart).Round(time.Millisecond),
		)
	}
}

/* getTasks returns a list of tasks to execute. */
func getTasks() ([]Task, error) {
	return make([]Task, 0), nil
}

/* taskExecutor executes the tasks sent on ch. */
func taskExecutor(ch <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range ch {
		executeTask(t)
	}
}

/* executeTask executes a single task. */
func executeTask(t Task) {
	log.Printf("Executing a task")
}
