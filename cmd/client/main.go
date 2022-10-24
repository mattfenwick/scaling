package main

import (
	"github.com/mattfenwick/scaling/pkg/client"
	"github.com/mattfenwick/scaling/pkg/utils"
)

func main() {
	utils.DoOrDie(client.RunSmallBatchOfRequests("localhost", 8765))
}
