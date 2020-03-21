package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/mikioh/ipaddr"
)

func readFile(name string) ([]ipaddr.Prefix, []ipaddr.Prefix) {
	var prefixesv6 []ipaddr.Prefix
	var prefixesv4 []ipaddr.Prefix
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
		_, ipNet, err := net.ParseCIDR(line)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(line, ":") {
			prefixesv6 = append(prefixesv6, *(ipaddr.NewPrefix(ipNet)))
		} else {
			prefixesv4 = append(prefixesv4, *(ipaddr.NewPrefix(ipNet)))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return prefixesv6, prefixesv4
}

func splitPrefix(pfx ipaddr.Prefix) (ipaddr.Prefix, ipaddr.Prefix) {
	splits := pfx.Subnets(1)
	return splits[0], splits[1]
}

func excludePrefix(a ipaddr.Prefix, b ipaddr.Prefix) ([]ipaddr.Prefix, []ipaddr.Prefix) {
	if a.Equal(&b) {
		return []ipaddr.Prefix{}, []ipaddr.Prefix{}
	}
	if !a.Contains(&b) {
		return []ipaddr.Prefix{a}, []ipaddr.Prefix{}
	}
	var low, up []ipaddr.Prefix
	s1, s2 := splitPrefix(a)
	for !s1.Equal(&b) && !s2.Equal(&b) {
		if s1.Contains(&b) {
			up = append([]ipaddr.Prefix{s2}, up...)
			s1, s2 = splitPrefix(s1)
		} else {
			low = append(low, s1)
			s1, s2 = splitPrefix(s2)
		}
	}
	if s1.Equal(&b) {
		up = append([]ipaddr.Prefix{s2}, up...)
	} else {
		low = append(low, s1)
	}
	return low, up
}

func subtractPrefixes(add []ipaddr.Prefix, sub []ipaddr.Prefix) []ipaddr.Prefix {
	var res []ipaddr.Prefix
	for len(add) > 0 && len(sub) > 0 {
		ap := add[0]
		sp := sub[0]
		if sp.Equal(&ap) || sp.Contains(&ap) {
			add = add[1:]
			continue
		}
		if ap.Contains(&sp) {
			add = add[1:]
			sub = sub[1:]
			low, up := excludePrefix(ap, sp)
			res = append(res, low...)
			add = append(up, add...)
			continue
		}
		if ipaddr.Compare(&ap, &sp) < 0 {
			add = add[1:]
			res = append(res, ap)
			continue
		}
		sub = sub[1:]
	}
	res = append(res, add...)
	return res
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
	aPfx6, aPfx4 := readFile(os.Args[1])
	bPfx6, bPfx4 := readFile(file2)
	fmt.Println("ip_diff " + os.Args[1] + " " + file2)
	fmt.Println("--- " + os.Args[1])
	fmt.Println("+++ " + file2)
	for _, prefix := range subtractPrefixes(ipaddr.Aggregate(aPfx6), ipaddr.Aggregate(bPfx6)) {
		fmt.Println("-" + prefix.String())
	}
	for _, prefix := range subtractPrefixes(ipaddr.Aggregate(bPfx6), ipaddr.Aggregate(aPfx6)) {
		fmt.Println("+" + prefix.String())
	}
	for _, prefix := range subtractPrefixes(ipaddr.Aggregate(aPfx4), ipaddr.Aggregate(bPfx4)) {
		fmt.Println("-" + prefix.String())
	}
	for _, prefix := range subtractPrefixes(ipaddr.Aggregate(bPfx4), ipaddr.Aggregate(aPfx4)) {
		fmt.Println("+" + prefix.String())
	}
}
