package main

import (
	"flag"
	"fmt"
	"github.com/fluturenet/dht"
	"github.com/fluturenet/ed25519"
	"log"
	"os"
)

var (
	mode    = flag.String("m", "", "command mode one of: get put createkey")
	target  = flag.String("target", "", "target for get mode")
	value   = flag.String("value", "", "value for put mode")
	mutable = flag.Bool("mut", false, "store as mutable")
	seq     = flag.Uint64("seq", 0, "sequential number for mutable")
	keyfile = flag.String("keyfile", "", "file to store or read privatekey")
	s       *dht.Server
)

func main() {
	fmt.Println("DHT army knife.")
	flag.Parse()

	switch *mode {
	case "get":
		if *target == "" {
			gentlyDie()
		}
		ok, targetHash := stringTo20Byte(target)
		if !ok {
			gentlyDie()
		}
		startServer()
		ad, _ := s.ArbitraryData(targetHash, nil)
	values:
		for {
			select {
			case value, ok := <-ad.Value:
				if !ok {
					break values
				}
				fmt.Printf("Received StorageItem from %v with value: %s seq:%d\n", value.NodeInfo.Addr, value.StorageItem.V, value.StorageItem.Seq)
				if value.StorageItem.Check() == nil {
					//s.AddStorageItem(value.StorageItem)
				}
			}
		}
		i, _ := s.GetStorageItem(targetHash)
		fmt.Printf("%s\n", i.V)
	case "put":
		if *value == "" {
			gentlyDie()
		}
		startServer()
		item := dht.StorageItem{V: *value}
		if *mutable {
			if *keyfile == "" {
				fmt.Println("Generating Key")
				keyPair, _ := ed25519.GenerateKey(nil)
				item.PrivateKey = keyPair.PrivateKey()
			} else {
				file, err := os.Open(*keyfile)
				if err != nil {
					fmt.Println(err)
					return
				}
				var buf [64]byte
				n, err := file.Read(buf[:])
				if err != nil {
					fmt.Println(err)
					return
				}
				if n != 64 {
					fmt.Printf("corrupted keyfile, read %d\n", n)
					return
				}
				item.PrivateKey = buf[:]
				file.Close()
			}
			item.Seq = *seq
		}
		item.Calc()
		s.AddStorageItem(item)
		log.Printf("storing: %s at %x seq:%d\n", *value, item.Target, item.Seq)
		ad, err := s.ArbitraryData(item.Target, nil)
		if err != nil {
			log.Fatal(err)
		}
	bput:
		for {
			select {
			case _, ok := <-ad.Value:
				if !ok {
					break bput
				}
				fmt.Printf(".")
			}
		}

	case "createkey":
		fmt.Println("Generating Key")
		keyPair, _ := ed25519.GenerateKey(nil)
		file, err := os.Create(*keyfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		file.Write(keyPair.PrivateKey())
		file.Close()

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
			return false, r
		}
	default:
		return false, r
	}

	copy(r[:], *t)
	return true, r
}

func startServer() {
	var err error
	s, err = dht.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("dht server on %s, ID is %x\n", s.Addr(), s.ID())
}
