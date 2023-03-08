package main

import "github.com/DeltaLaboratory/dynamicdns/application"

func main() {
	service := application.DDNSService{}
	service.Run()
}
