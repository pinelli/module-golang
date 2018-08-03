package main

func main(){
	server := CreateServer("localhost:9090")
	server.Run()
	select{}
}
