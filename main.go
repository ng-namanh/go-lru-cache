package main

import (
	"fmt"
	"log"
	"os"

	"go-lru-cache/cache"
)

func main() {
	c, err := cache.NewCache(2)
	if err != nil {
		log.Fatal(err)
	}

	c.Put("1", 1)
	c.Put("2", 2)
	if v, ok := c.Get("1"); ok {
		fmt.Println("Get 1 =", v)
	}
	c.Put("3", 3)

	if v, ok := c.Get("2"); !ok {
		fmt.Println("Get 2: miss (evicted) — ok =", ok)
	} else {
		fmt.Fprintln(os.Stderr, "unexpected hit for 2:", v)
		os.Exit(1)
	}
	if v, ok := c.Get("1"); ok {
		fmt.Println("Get 1 =", v)
	}
	if v, ok := c.Get("3"); ok {
		fmt.Println("Get 3 =", v)
	}
	fmt.Println("Len =", c.Len())
}
