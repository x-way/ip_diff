package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"inet.af/netaddr"
)

func readFile(name string) []netaddr.IPPrefix {
	var prefixes []netaddr.IPPrefix
	var f *os.File
	if name == "-" {
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.Contains(line, "/") {
			if strings.Contains(line, ":") {
				line = line + "/128"
			} else {
				line = line + "/32"
			}
		}
		prefix, err := netaddr.ParseIPPrefix(line)
		if err != nil {
			log.Fatal(err)
		}
		prefixes = append(prefixes, prefix)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return prefixes
}

func main() {
	file2 := "-"
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Usage: ip_diff <file_a> [<file_b>]")
		return
	}
	if len(os.Args) > 2 {
		file2 = os.Args[2]
	}
	var aBuilder netaddr.IPSetBuilder
	var bBuilder netaddr.IPSetBuilder
	aPrefixes := readFile(os.Args[1])
	bPrefixes := readFile(file2)
	for _, prefix := range aPrefixes {
		aBuilder.AddPrefix(prefix)
	}
	aSet := aBuilder.IPSet()
	for _, prefix := range bPrefixes {
		bBuilder.AddPrefix(prefix)
	}
	bSet := bBuilder.IPSet()
	fmt.Println("ip_diff " + os.Args[1] + " " + file2)
	fmt.Println("--- " + os.Args[1])
	fmt.Println("+++ " + file2)

	removedBuilder := bBuilder.Clone()
	removedBuilder.Complement()
	removedBuilder.Intersect(aSet)
	for _, prefix := range removedBuilder.IPSet().Prefixes() {
		fmt.Println("-" + prefix.String())
	}

	addedBuilder := aBuilder.Clone()
	addedBuilder.Complement()
	addedBuilder.Intersect(bSet)
	for _, prefix := range addedBuilder.IPSet().Prefixes() {
		fmt.Println("+" + prefix.String())
	}
}
