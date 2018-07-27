package main

func main() {
	cfg := Config{
		StaticFileDir: "static",
		MainURLHash:   "hehe",
		SvrPort:       8080,
	}
	svr := HttpSvr{
		config: cfg,
	}
	svr.StartServe(nil)
}
