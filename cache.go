package main

// 缓存块结构
type Block struct {
	SeqNumber int  // seqNumber
	Age       int  // age
	Valid     bool // 有效位
}

// SET结构，每个SET包含4个块
type Set struct {
	Blocks [4]Block
}

// cache ,4way = 4block/per
type Cache struct {
	Sets        [8]Set
	AccessCount int
	HitCount    int
	MissCount   int
}

func NewCache() *Cache {
	return &Cache{}
}

// 根据序列号获取对应的SET索引
func (c *Cache) GetSetIndex(seqNumber int) int {
	return seqNumber % 8
}

func (c *Cache) Access(seqNumber int) (bool, int) {

	var (
		setIndex = c.GetSetIndex(seqNumber)
		set      = &c.Sets[setIndex]
	)

	c.AccessCount++

	// 检查是否命中
	for i := 0; i < 4; i++ {
		if set.Blocks[i].Valid && set.Blocks[i].SeqNumber == seqNumber {
			c.HitCount++
			set.Blocks[i].Age = c.AccessCount
			return true, setIndex
		}
	}

	// 未命中
	c.MissCount++

	// 查找空块或替换最老的块
	replaceIndexOfYoungest := 0

	// 优先查找无效块
	for i := 0; i < 4; i++ {
		if !set.Blocks[i].Valid {
			replaceIndexOfYoungest = i
			break
		}
		if set.Blocks[i].Age < set.Blocks[replaceIndexOfYoungest].Age {
			replaceIndexOfYoungest = i
		}
	}

	// 替换块
	set.Blocks[replaceIndexOfYoungest].SeqNumber = seqNumber
	set.Blocks[replaceIndexOfYoungest].Age = c.AccessCount
	set.Blocks[replaceIndexOfYoungest].Valid = true

	return false, setIndex
}

func (c *Cache) GetStats() map[string]float64 {
	stats := make(map[string]float64)

	stats["accessCount"] = float64(c.AccessCount)
	stats["hitCount"] = float64(c.HitCount)
	stats["missCount"] = float64(c.MissCount)

	if c.AccessCount > 0 {
		hitRate := float64(c.HitCount) / float64(c.AccessCount)
		missRate := float64(c.MissCount) / float64(c.AccessCount)

		stats["hitRate"] = hitRate
		stats["missRate"] = missRate

		// Tavg = h*C + (1-h)*M
		// C = 1ns (缓存访问时间)
		// M = 1 + 16*10 + 1 = 162ns (未命中惩罚：检查缓存 + 加载16个字 + 写入缓存)
		cacheAccessTime := 1.0
		missPenalty := 162.0 // 1 + 16*10 + 1

		stats["avgAccessTime"] = hitRate*cacheAccessTime + missRate*missPenalty
		stats["totalAccessTime"] = float64(c.AccessCount) * stats["avgAccessTime"]
	}

	return stats
}
