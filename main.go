package unicorn

import (
	"flag"
	"github.com/Waitfantasy/unicorn/conf"
	"log"
)



var filename string

func main()  {
	var (
		err error
		cfg *conf.Config
	)
	flag.StringVar(&filename, "config", "", "")
	flag.Parse()

	if cfg, err = conf.InitConfig(filename); err != nil {
		log.Fatal(err)
	}
}