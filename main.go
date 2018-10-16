package unicorn

import (
	"flag"
	"fmt"
	"github.com/Waitfantasy/unicorn/conf"
	"log"
)



var filename string

func main()  {
	var (
		err error
		cfg *conf.Conf
	)
	flag.StringVar(&filename, "config", "", "")
	flag.Parse()

	if cfg, err = conf.InitConf(filename); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
}