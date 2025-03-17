package main

import "math/rand"

// 生成顺序序列：0到2n-1，重复4次
func GenerateSequential(n int) []int {
	sequence := make([]int, 0)
	// 重复4次
	for i := 0; i < 4; i++ {
		for j := 0; j < 2*n; j++ {
			sequence = append(sequence, j)
		}
	}
	return sequence
}

// 生成随机序列：包含4n个主存块
func GenerateRandom(n int) []int {
	sequence := make([]int, 4*n)
	for i := range sequence {
		sequence[i] = rand.Intn(4 * n)
	}
	return sequence
}

// 生成中间重复序列
func GenerateMidRepeat(n int) []int {
	sequence := make([]int, 0)

	// 一次完整序列
	baseSeq := make([]int, 0)

	// 开始：0
	baseSeq = append(baseSeq, 0)

	// 中间重复两次：1到n-1
	for k := 0; k < 2; k++ {
		for i := 1; i < n; i++ {
			baseSeq = append(baseSeq, i)
		}
	}

	// 继续到2n
	for i := n; i < 2*n; i++ {
		baseSeq = append(baseSeq, i)
	}

	// 重复4次
	for i := 0; i < 4; i++ {
		sequence = append(sequence, baseSeq...)
	}

	return sequence
}

// GenerateSameSetSequence
func GenerateSameSetSequence(setIndex int) []int {
	// 生成 10 个映射到同一个 set 的序列
	sequence := make([]int, 4)

	// 由于 setIndex = seqNumber % 8
	// 所以 seqNumber = 8k + setIndex 的数都会映射到同一个 set
	for i := 0; i < 4; i++ {
		sequence[i] = i*4 + setIndex
	}

	return sequence
}
