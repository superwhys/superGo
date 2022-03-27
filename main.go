package main

import (
	"github.com/superwhys/superGo/superFlags"
	lg "github.com/superwhys/superGo/superLog"
)

var (
	serviceName = superFlags.String("serviceName", "test_service", "the name of this service!")
)

func main() {
	superFlags.Parse()
	lg.Info(serviceName())
	//superLog.Fatalf("user define error:", errors.New("this is error"))
}
