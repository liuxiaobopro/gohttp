package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
	logrus.Infof("Hello World! %s", time.Now())
}

func main() {
	http.HandleFunc("/", greet)
	if err := http.ListenAndServe(":9000", nil); err != nil {
		panic(err)
	}
}
