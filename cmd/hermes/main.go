// Package main contains the entry point for the Hermes application.
// It bootstraps the application and starts its execution.
package main

import "Hermes/internal/app"

func main() {

	// main is the program entry point. It initializes the application using app.Boot
	// and starts it by calling Run, which blocks until a shutdown signal is received.
	app.Boot().Run()

}
