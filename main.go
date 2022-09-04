package main

import (
	"fmt"
	"regexp"

	"github.com/amaury95/regex-cache/cache"
	"github.com/hashicorp/golang-lru/simplelru"
)

func main() {
	/* define regular expressions */
	re := regexp.MustCompile(`(fan){2}`)

	/* initiate regex cache */
	lru, err := simplelru.NewLRU(256, nil)
	if err != nil {
		panic(err)
	}
	r := cache.NewRegexCache(lru)

	/* populate cache */
	r.DefineSysRegexp(re, "fansea")

	/* check cache */
	fmt.Printf("\"testing\": %v\n", "fan")
	fmt.Println(r.Get("fan"))
	fmt.Println(r.Get("fanfan"))
	fmt.Println(r.Get("fanfanfan"))

	/* update cache */
	re2 := regexp.MustCompile(`(sea){2}`)
	r.UpdateSysRegexp(re, re2)

	/* it should fail now */
	fmt.Printf("\"testing\": %v\n", "fan")
	fmt.Println(r.Get("fan"))
	fmt.Println(r.Get("fanfan"))
	fmt.Println(r.Get("fanfanfan"))

	/* it should succeed now */
	fmt.Printf("\"testing\": %v\n", "sea")
	fmt.Println(r.Get("sea"))
	fmt.Println(r.Get("seasea"))
	fmt.Println(r.Get("seaseasea"))

	/* it should update the value */
	r.Add("seaseasea", "seafan")
	fmt.Println(r.Get("seasea"))
}
