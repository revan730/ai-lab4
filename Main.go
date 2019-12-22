package main

import "fmt"

func main() {
	generalScore := 0.0
	const iterCount = 100
	for i := 0; i < iterCount; i++ {
		env := NewEnvironment()
		env.InitField(1, 1)
		const usePrint = true
		env.StartLoop(usePrint)
		generalScore += float64(env.GetScore())
	}
	fmt.Printf("Average score: %f, after %d iterations\n", generalScore / iterCount, iterCount)
}
