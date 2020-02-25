package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// create our args type
type Args struct {
	Addr    string
	Workers int
	Wait    int
	Range   []int
}

// init some vars
var (
	open_ports []int
	addr       string
	workers    int
	wait       int
	_range_    []int
	err        error
)

// the help msg for usage
var help_msg string = `
usage: pscan ADDR [ARGS]

optional args:
--help      show this help msg and exit
--workers   how many workers to dispatch (max is 1000)
--wait      how long to wait in ms before we fail the port (default is 90)
--range     range of ports to scan (42-6666)

examples: (on windows command is pscan.exe)

$ pscan 192.168.1.87 (defaults applied are ports 0-65535 with 250 workers
                      waiting 90ms)

$ pscan 192.168.1.87 --workers 666 --wait 100

$ pscan 192.168.1.87 --workers 666 --wait 100 --range 40000-42000
`

func parse_args() Args {
	/*
		this function parses args via the os package and stores them in
		the Args type
	*/
	// set default args
	workers = 250
	wait = 90
	_range_ = []int{1, 65535}
	// if the total args is ess than the minimum, esit
	if len(os.Args) < 2 {
		fmt.Printf("expected an addr as an arg%s\n...exiting\n", help_msg)
		os.Exit(1)
	}
	// we expect the addr to be the first arg
	addr = os.Args[1]
	// iter over the args and store them. if an arg doesnt start with
	// "--", skip it.
	for i, a := range os.Args[1:] {
		if !strings.HasPrefix(a, "-") {
			continue
		} else if a == "--help" {
			fmt.Println(help_msg)
			os.Exit(0)
		} else if a == "--workers" {
			workers, err = strconv.Atoi(os.Args[i+2])
			if err != nil {
				fmt.Printf("expected an int, got: %v\n...exiting", os.Args[i+2])
				os.Exit(1)
			}
			if workers > 1000 {
				fmt.Printf("%v exceeds 1000 worker max, setting to 1000\n", workers)
				workers = 1000
			}
		} else if a == "--wait" {
			wait, err = strconv.Atoi(os.Args[i+2])
			if err != nil {
				fmt.Printf("expected an int, got: %v\n...exiting", os.Args[i+2])
				os.Exit(1)
			}
		} else if a == "--range" {
			if !strings.Contains(a, "-") {
				fmt.Println("expected range to be startINT-endINT\n...exiting")
				os.Exit(1)
			}
			range_args := strings.Split(os.Args[i+2], "-")
			start_int, err := strconv.Atoi(range_args[0])
			if err != nil {
				fmt.Printf("expected an int, got: %v\n...exiting", range_args[0])
				os.Exit(1)
			}
			end_int, err := strconv.Atoi(range_args[1])
			if err != nil {
				fmt.Printf("expected an int, got: %v\n...exiting", range_args[1])
				os.Exit(1)
			}
			_range_ = []int{start_int, end_int}
		} else {
			fmt.Printf("received an unexpected arg: %v\n%s\n...exiting", a, help_msg)
			os.Exit(1)
		}
	}
	// store it
	args := Args{
		addr,
		workers,
		wait,
		_range_,
	}
	return args
}

func port_worker(addr string, wait int, ports chan int, results chan int) {
	/*
		this function tries to connect to the addr on a given port
		according to the current worker pool
	*/
	// for each port in the pool
	for p := range ports {
		// concat it to the addr so we can connect
		address := fmt.Sprintf("%v:%d", addr, p)
		fmt.Printf("\rworking on %d", p)
		conn, err := net.DialTimeout("tcp", address, time.Duration(wait)*time.Millisecond)
		// if we cant connect, store port as 0
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		// if port is open, add it to the list
		results <- p
	}
}

func main() {
	/*
	   entry point for program
	*/
	// mark the start time
	start := time.Now().Unix()
	// parse the args
	args := parse_args()
	// give the user the details of the scan
	fmt.Printf("dispatching %d workers with a %d ms timeout\n", args.Workers, args.Wait)
	// create our worker pool
	ports := make(chan int, args.Workers)
	// init the results chan
	results := make(chan int)
	// try and connect
	for i := 0; i < cap(ports); i++ {
		go port_worker(args.Addr, args.Wait, ports, results)
	}
	// iter over the pool
	go func() {
		for i := args.Range[0]; i <= args.Range[1]; i++ {
			ports <- i
		}
	}()
	// iter over the reuslts list
	for i := args.Range[0]; i <= args.Range[1]; i++ {
		port := <-results
		if port != 0 {
			open_ports = append(open_ports, port)
		}
	}
	close(ports)
	close(results)
	fmt.Println()
	for _, port := range open_ports {
		fmt.Printf("%d is open\n", port)
	}
	fin := time.Now().Unix()
	fmt.Printf("took %v seconds\n", fin-start)
}
