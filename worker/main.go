package main

func main() {
	_, err := config_from_env()
	if err != nil {
		panic(err)
	}
}
