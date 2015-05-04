package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/op/go-logging"
	"github.com/spikeekips/kasi"
)

func main() {
	yml, _ := ioutil.ReadFile(os.Args[1])

	kasi.SetLogging(logging.DEBUG)

	fmt.Println("hello! this is `kasi`.")
	kasi.Run(string(yml))
}
