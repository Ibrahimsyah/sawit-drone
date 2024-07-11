package main

func main() {
	util := NewUtil()
	app := NewApp(util)
	app.Start()
}
