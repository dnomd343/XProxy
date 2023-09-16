package main

import (
	"XProxy/next/process"
	"time"
)

func main() {

	//cmd := exec.Command("ls", "-al")
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//err := cmd.Start()
	//fmt.Println(err)
	//err = cmd.Wait()
	//fmt.Println(err)

	p := process.NewDaemon("demo", []string{"sleep", "3"}, map[string]string{"DEBUG": "TRUE"})
	p.Start()

	time.Sleep(10 * time.Second)

}
