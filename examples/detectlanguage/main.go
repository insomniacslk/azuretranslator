package main

import (
	"flag"
	"fmt"
	//"github.com/insomniacslk/azuretranslator"
	azuretranslator "../.."
	"os"
	"sync"
)

var apiKey = flag.String("apikey", "", "Azure translator API key")

func main() {
	flag.Parse()
	if *apiKey == "" {
		fmt.Println("Must specify API key")
		os.Exit(1)
	}
	c, err := azuretranslator.NewTranslatorClient(*apiKey)
	if err != nil {
		panic(err)
	}
	phrases := []string{
		"the pen is on the table",
		"je suis un chien",
		"no tiengo dinero",
		"ciao amico",
	}
	var wg sync.WaitGroup
	for _, phrase := range phrases {
		wg.Add(1)
		go func(phrase string) {
			defer wg.Done()
			trans, err := c.Translate(phrase, "en", "")
			if err != nil {
				panic(err)
			}
			fmt.Printf("'%v' -> '%v'\n", phrase, trans)
		}(phrase)
	}
	wg.Wait()
}
