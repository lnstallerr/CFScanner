package main

import (
	configuration "CFScanner/configuration"
	"CFScanner/scanner"
	"CFScanner/utils"
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func run(timer time.Time) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     os.Args[0],
		Short:   codename,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(VersionStatement())

			if v2raypath != "" {
				configuration.BIN = v2raypath
			}

			if Vpn {
				utils.CreateDir(configuration.DIR)
			}

			utils.CreateDir(configuration.RESULTDIR)

			if err := configuration.CreateInterimResultsFile(configuration.InterimResultsPath, nTries); err != nil {
				fmt.Printf("Error creating interim results file: %v\n", err)
			}
			// number of threads for scanning
			threadsCount := threads

			// lists of ip for scanning process
			var LIST []string

			LIST = IOScanner(LIST)

			// Parsing and Validating IP LISTS
			bigIPList = utils.IPParser(LIST)

			// Shuffling IPList
			if shuffle {
				rand.Shuffle(len(bigIPList), func(i, j int) {
					bigIPList[i], bigIPList[j] = bigIPList[j], bigIPList[i]
				})
			}

			// Total number of IPS
			numberIPS := utils.TotalIps(bigIPList)

			if int(numberIPS) <= 1 {
				log.Fatal("Scanning Failed : No IP detected")
			}

			// Create Configuration file
			Config, worker, _ := configuration.CreateTestConfig(configPath, startProcessTimeout, doUploadTest,
				minDLSpeed, minULSpeed, maxDLTime, maxULTime,
				frontingTimeout, fronting, maxDLLatency, maxULLatency,
				nTries, Vpn, threads, shuffle)

			fmt.Printf("Starting to scan %v%d%v IPS.\n\n", utils.Colors.OKGREEN, numberIPS, utils.Colors.ENDC)
			// Begin scanning process
			scanner.Start(&Config, &worker, bigIPList, threadsCount)

			fmt.Println("Results Written in :", configuration.InterimResultsPath)
			fmt.Println("Sorted IPS Written in :", configuration.FinalResultsPathSorted)
			fmt.Println("Time Elapse :", time.Since(timer))
		},
	}
	return rootCmd
}

func IOScanner(LIST []string) []string {
	file, _ := utils.Exists(subnets)

	if file && subnets != "" {
		subnetFilePath := subnets
		subnetFile, err := os.Open(subnetFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer func(subnetFile *os.File) {
			err := subnetFile.Close()
			if err != nil {

			}
		}(subnetFile)

		newScanner := bufio.NewScanner(subnetFile)
		for newScanner.Scan() {
			LIST = append(LIST, strings.TrimSpace(newScanner.Text()))
		}
		if err := newScanner.Err(); err != nil {
			log.Fatal(err)
		}

	} else {
		// type conversion of string subnet to []string
		var subnet []string
		subnet = append(subnet, subnets)

		ips := utils.IPParser(subnet)

		LIST = append(LIST, ips...)

	}
	return LIST
}
