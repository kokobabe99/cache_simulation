package main

// 缓存块结构
type Block struct {
	SeqNumber int  // seqNumber
	Age       int  // age
	Valid     bool // 有效位
}

type Set struct {
	Blocks []Block
}

// 缓存配置结构
type CacheConfig struct {
	BlockNumber int // 缓存大小
	WayNum      int // N-way set associative
	SetNum      int // 组数
	LineSize    int // 行大小
}

// cache ,4way = 4block/per
type Cache struct {
	config      CacheConfig
	Sets        []Set
	AccessCount int
	HitCount    int
	MissCount   int
}

func NewCache(way, blockNumber, lineSize int) *Cache {

	var (
		SetNum = blockNumber / way
		cache  = &Cache{
			config: CacheConfig{
				WayNum:      way,
				SetNum:      SetNum,
				LineSize:    lineSize,
				BlockNumber: blockNumber,
			},
			Sets: make([]Set, SetNum),
		}
	)

	// 初始化每个Set
	for i := range cache.Sets {
		cache.Sets[i] = Set{
			Blocks: make([]Block, way),
		}
	}

	return cache
}

// 根据序列号获取对应的SET索引
func (c *Cache) GetSetIndex(seqNumber int) int {
	return seqNumber % (c.config.SetNum)
}

func (c *Cache) Access(seqNumber int) (bool, int) {

	var (
		setIndex = c.GetSetIndex(seqNumber)
		set      = &c.Sets[setIndex]
	)

	c.AccessCount++

	for i := 0; i < c.config.WayNum; i++ {
		if set.Blocks[i].Valid && set.Blocks[i].SeqNumber == seqNumber {
			c.HitCount++
			set.Blocks[i].Age = c.AccessCount
			return true, setIndex
		}
	}

	c.MissCount++

	replaceIndexOfYoungest := 0

	for i := 0; i < c.config.WayNum; i++ {
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
		// C = 1ns
		// M = 1 + 16*10 + 1 = 162ns
		cacheAccessTime := 1.0
		//missPenalty := 162.0 // 1 + 16*10 + 1
		missPenalty := c.config.LineSize*10 + 1 + 1

		stats["avgAccessTime"] = hitRate*cacheAccessTime + missRate*float64(missPenalty)
		stats["totalAccessTime"] = float64(c.AccessCount) * stats["avgAccessTime"]
	}

	return stats
}
