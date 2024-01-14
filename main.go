package main

import "github.com/SyedDevop/gitpuller/cliapp"

// _ "github.com/SyedDevop/gitpuller/cliapp"

func main() {
	cliapp.CliAppInit()
	// req, err := http.NewRequest("GET", "https://api.github.com/repos/SyedDevop/linux-setup/contents", nil)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer resp.Body.Close()
	// respBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(respBytes))
}
