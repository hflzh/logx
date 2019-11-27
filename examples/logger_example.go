// An example that shows how to configure a logger and write log messages at
// different levels using the logx package.
package main

import (
	"fmt"
	"os"

	"github.com/hflzh/logx"
)

func main() {
	var lgr *logx.Logger

	fmt.Println("Logging level: Debug ")
	lgr = logx.New(os.Stdout, logx.Debug, false)
	lgr.Fine("logx: fine messages are only visible when the logging level is less than or equals to Fine.")
	lgr.Debug("logx: this is a debug message.")
	lgr.Info("logx: hello world!")
	lgr.Error("logx: the server is not responding, please check logs for more details.")

	fmt.Println("Logging level: Info ")
	lgr = logx.New(os.Stderr, logx.Info, true)
	lgr.Fine("logx: fine messages are only visible when the logging level is less than or equals to Fine.")
	lgr.Fine("logx: debug messages are only visible when the logging level is less than or equals to Debug.")
	lgr.Info("logx: hello world!")
	lgr.Error("logx: the server is not responding, please check logs for more details.")
}
