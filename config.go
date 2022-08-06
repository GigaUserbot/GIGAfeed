package main

import "os"

// first the token will be read from the environment variable
// if environment variable is empty, will use the constant value
var BOT_TOKEN = func() string {
	if os.Getenv("BOT_TOKEN") == "" {
		// constant bot token value
		return "3838833:efuhuefhuefhuefheu"
	}
	return os.Getenv("BOT_TOKEN")
}()

// first the token will be read from the environment variable
// if environment variable is empty, will use the constant value
var PORT = func() string {
	if os.Getenv("PORT") == "" {
		// constant bot token value
		return "8080"
	}
	return os.Getenv("PORT")
}()
