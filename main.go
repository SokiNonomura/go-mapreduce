package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// KeyValue mapping
type KeyValue struct {
	Key   string
	Value int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s: see usage comments in file\n", os.Args[0])
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var wg sync.WaitGroup
	var reduce map[string]int = make(map[string]int)

	reducer := make(chan []KeyValue)
	scanner := bufio.NewScanner(f)

	go func() {
		Reduce(reduce, reducer, &wg)
	}()
	for scanner.Scan() {
		t := scanner.Text()
		wg.Add(1)
		go func(t string) {
			reducer <- Map(t)
		}(t)
	}
	wg.Wait()
	close(reducer)
	fmt.Println(reduce)
}

// Map Mapping
func Map(t string) []KeyValue {
	kvs := []KeyValue{}
	slice := strings.Split(t, "")
	len := len(slice)
	for i := 0; i < len; i++ {
		kvs = append(kvs, KeyValue{slice[i], 1})
	}
	return kvs
}

// Reduce reduce
func Reduce(reduce map[string]int, reducer <-chan []KeyValue, wg *sync.WaitGroup) {
	for kvs := range reducer {
		func() {
			defer wg.Done()
			for _, v := range kvs {
				reduce[v.Key]++
			}
		}()
	}
}
