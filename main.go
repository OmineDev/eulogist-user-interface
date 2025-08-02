package main

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/function"
	"github.com/OmineDev/eulogist-user-interface/server"
)

func main() {
	test()
}

func test() {
	s := server.NewServer()

	fmt.Println(s.RunServer("127.0.0.1:19132"))
	fmt.Println(s.WaitConnect())

	conn := s.MinecraftConn()
	interact := server.NewInteract(conn)
	clientFunction, err := function.NewFunction(interact)
	if err != nil {
		panic(err)
	}

	// clientFunction.RegisterOrLogin()
	clientFunction.RequestUserInfo()
	clientFunction.ShowPanel()
}
