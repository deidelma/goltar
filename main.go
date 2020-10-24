package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	jobs "github.com/deidelma/goltar/process"
)

func main() {
	var jobPath string

	flag.StringVar(&jobPath, "job", "jobs.toml", "Path to the job file")
	flag.Parse()

	log.Printf("The job path is \"%s\".\n", jobPath)
	// path := fmt.Sprintf("%s/%s", "/home/david/Projects/goltar", jobPath)
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Unable to find the current working directory!")
	}
	path := path.Join(cwd, jobPath)
	job, err := jobs.ReadJobFile(path)
	if err != nil {
		log.Fatalf("Read error:%v\n", err)
	}
	fmt.Printf("Read job named %s\n", job.Name)

	for _, search := range job.Searches {
		startDate := search.Years[0]
		fmt.Printf("The first search term is: %s\n", search.Ands[0])
		fmt.Printf("The year is %d\n", startDate)
	}
}
