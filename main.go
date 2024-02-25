/*
Copyright © 2024 Syed Uzair Ahmed <syeds.devops007@gmail.com>
*/
package main

import (
	"github.com/SyedDevop/gitpuller/cmd"
	_ "github.com/joho/godotenv/autoload"
)

// var (
// 	cpuprofile = flag.String("cpuprofile", "cpu.prof", "write cpu profile to `file`")
// 	memprofile = flag.String("memprofile", "mem.prof", "write memory profile to `file`")
// )

func main() {
	// flag.Parse()
	// if *cpuprofile != "" {
	// 	// fmt.Println("Hello")
	// 	f, err := os.Create("./profiles/" + *cpuprofile)
	// 	if err != nil {
	// 		log.Fatal("could not create CPU profile: ", err)
	// 	}
	// 	defer f.Close() // error handling omitted for example
	// 	if err := pprof.StartCPUProfile(f); err != nil {
	// 		log.Fatal("could not start CPU profile: ", err)
	// 	}
	// 	// fmt.Println(f.Name())
	// 	defer pprof.StopCPUProfile()
	// }

	cmd.Execute()

	// if *memprofile != "" {
	// 	f, err := os.Create("./profiles/" + *memprofile)
	// 	if err != nil {
	// 		log.Fatal("could not create memory profile: ", err)
	// 	}
	// 	defer f.Close() // error handling omitted for example
	// 	runtime.GC()    // get up-to-date statistics
	// 	if err := pprof.WriteHeapProfile(f); err != nil {
	// 		log.Fatal("could not write memory profile: ", err)
	// 	}
	// }
}
