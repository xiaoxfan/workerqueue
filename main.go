/*
@Time : 2020/4/9 4:43 PM
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"workerqueue/workerqueue"
)

func main() {
	workerQueue := workerqueue.New(10)
	workerQueue.Start()
	// for i in {1..4096}; do curl localhost:8000/submit-work -d name=$USER -d delay=$(expr $i % 11)s; done
	http.HandleFunc("/submit-work", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		delay, _ := time.ParseDuration(r.FormValue("delay"))
		work := func() {
			time.Sleep(delay)
			fmt.Printf("after delay %f seconds,Hello, %s!\n", delay.Seconds(), name)
		}
		workerQueue.SubmitWork(work)
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}
