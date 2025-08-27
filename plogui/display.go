package main

import "fmt"

func maxWeight(data []weightedData) int64 {
	var maximum int64 = 0
	for _, item := range data {
		if item.weight > maximum {
			maximum = item.weight
		}
	}
	return maximum
}

func displayData(title string, data []weightedData) {
	fmt.Println(title)
	// 	_ := maxWeight(data)

	for _, item := range data {
		fmt.Printf("%d\t%s", item.weight, item.data)
	}
}
