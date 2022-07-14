package test

import "strconv"

// 测试数据

// Butterfly 蝴蝶
type Butterfly struct {
	Species string // 品种
	Life    uint16 // 寿命
	Habitat string // 栖息地
}

var (
	// 数据源定义
	ButterflyMap map[int64]Butterfly
)

// DataGererate 生产指定数据量的数据
func DataGererate(count int) {
	ButterflyMap = map[int64]Butterfly{}

	id := int64(0)
	for i := 0; i < count; i++ {
		id += 1
		ButterflyMap[id] = Butterfly{
			Species: "Btf" + strconv.Itoa(i+1),
			Life:    5,
			Habitat: "Asia",
		}
	}
}

// GetDataIds 获取数据源ID
func GetDataIds() []int64 {
	ids := make([]int64, len(ButterflyMap))

	idx := 0
	for id, _ := range ButterflyMap {
		ids[idx] = id
		idx += 1
	}

	SortInt64(ids, true)

	return ids
}
