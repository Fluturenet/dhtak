package main

import (
	"fmt"
	"flag"
	"os"
	"log"
	"github.com/fluturenet/dht"
)

var (
	mode 	= flag.String("m","","command mode one of: get put")
	target 	= flag.String("target","","target for get mode")
	value  = flag.String("value","","value for put mode")
	s	  *dht.Server
)


func main(){
        fmt.Println("DHT army knife.")
	flag.Parse()

	switch *mode {
	case "get":
		if *target =="" {
			gentlyDie()
		}
		ok,targetHash := stringTo20Byte(target)
		if !ok {
			gentlyDie()
		}
		startServer()
		ad,_:=s.ArbitraryData(targetHash)
		values:
		for {
			select {
				case value,ok :=  <- ad.Value:
					if !ok {
                                        break values
	                        }
                                fmt.Printf("Received StorageItem from %v with value: %s\n", value.NodeInfo.Addr,value.StorageItem.V)
                                if value.StorageItem.Check() ==nil {
                                        s.AddStorageItem(value.StorageItem)
        	                }
			}
		}
	case "put":
		if *value ==""{
			gentlyDie()
		}
		startServer()
                item := dht.StorageItem{ V: *value, }
                item.Calc()
                s.AddStorageItem(item)
                log.Printf("storing: %s at %x\n",*value,item.Target)
                ad, err := s.ArbitraryData(item.Target)
                if err != nil {
                        log.Fatal(err)
                }
		bput:
                for {
                        select {
                                case _,ok :=  <- ad.Value:
                                        if !ok {
                                        break bput
                                	}
                                	fmt.Printf(".")
                        }
                }

	default:
		gentlyDie()
	}
}

func gentlyDie() {
	fmt.Println("\nUsage:")
	flag.PrintDefaults()
	fmt.Println("Example:")
	fmt.Println("\t-m get -target 7d81e1e5c5a5512935b1d19365e157bd97e849d5")
	fmt.Println("\t-m put -value \"Hello world!\"")
	os.Exit(0)
}

func stringTo20Byte(t *string) (bool, [20]byte) {
	var r [20]byte
        switch len(*t) {
        case 20:
        case 40:
                _, err := fmt.Sscanf(*t, "%x", t)
                if err != nil {
                        return false,r
                }
        default:
                return false,r
        }

        copy(r[:], *t)
	return true,r
}

func startServer() {
        var err error
        s, err = dht.NewServer(nil)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Printf("dht server on %s, ID is %x\n", s.Addr(), s.ID())
}
