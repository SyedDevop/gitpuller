package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/SyedDevop/gitpuller/log"
)

var LogPath = "network_test.log"

func main() {
	var reqTimes int
	flag.IntVar(&reqTimes, "rt", 10, "How many times should the request be sent")
	flag.Parse()

	logF := log.LogFile(&LogPath)
	defer logF.Close()

	tookChan := make(chan time.Duration, reqTimes)
	wg := sync.WaitGroup{}
	fmt.Fprint(logF, "Starting the test with: "+fmt.Sprint(reqTimes)+" times\n")
	wg.Add(reqTimes)
	for i := range reqTimes {
		go func(logFile *os.File, count int) {
			defer wg.Done()
			start := time.Now()
			res := req()
			if res.StatusCode == http.StatusOK {
				logString := fmt.Sprintf("(%d):Request took: %0.3d ms\n", count, time.Since(start).Milliseconds())
				fmt.Fprint(logF, logString)
				tookChan <- time.Since(start)
			} else {
				return
			}
		}(logF, i)
	}

	wg.Wait()
	close(tookChan)
	var avgTime int64
	for t := range tookChan {
		avgTime += t.Milliseconds()
	}

	fmt.Println("Done")
	fmt.Fprintf(logF, "The average time of the requests was: %0.3d ms for %d requests\n\n", avgTime/int64(reqTimes), reqTimes)
}

func req() *http.Response {
	resp, err := http.Get("https://api.github.com")
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()
	return resp
}
