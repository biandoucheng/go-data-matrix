package test

// 性能测试

import (
	"fmt"
	"strconv"

	godatamatrix "github.com/biandoucheng/go-data-matrix"
)

var (
	ButterflyMatrix godatamatrix.DataMatrix = godatamatrix.DataMatrix{}
	ButterflyData   map[int64]Butterfly     = map[int64]Butterfly{}
)

// GenerateDataSource 生产数据源
func GenerateDataSource(count int) {
	id := int64(0)
	for i := 0; i < count; i++ {
		id += 1
		ButterflyData[id] = Butterfly{
			Species: "Btf" + strconv.Itoa(i+1),
			Life:    5,
			Habitat: "Asia",
		}
	}
}

// GetDataSourceIds 获取数据源ID
func GetDataSourceIds(count int) []int64 {
	idl := count
	if count == 0 {
		idl = len(ButterflyData)
	}
	ids := make([]int64, idl)

	idx := 0
	for id := range ButterflyData {
		ids[idx] = id
		idx += 1

		if count > 0 && idx >= count {
			break
		}
	}

	return ids
}

// InitDataMatrix 初始化数据矩阵
func InitDataMatrix(num int) {
	GenerateDataSource(num)
	n := "蝴蝶品种统计"
	ids := GetDataSourceIds(0)
	ButterflyMatrix.Init(n, ids)
}

// InitCostTest 初始化耗时测试
func InitCostTest() {
	GenerateDataSource(100000)
	n := "蝴蝶品种统计"
	ids := GetDataSourceIds(0)

	ButterflyMatrix.StartMrcs()

	ButterflyMatrix.Init(n, ids)

	cost := ButterflyMatrix.FinishMrcs()

	fmt.Println("\r\n初始化耗时[10 0000]: ", cost, "微秒")
	fmt.Println()
}

// LightUpTest 点亮测试
func LightUpTest() {
	InitDataMatrix(100000)

	ButterflyMatrix.StartMrcs()

	for i := 0; i < 100000; i++ {
		id := int64(i + 1)
		ButterflyMatrix.LightUpPoint(id, 0)
	}

	cost := ButterflyMatrix.FinishMrcs()

	fmt.Println("\r\n点亮阵点耗时[10 0000]: ", cost, "微秒")
	fmt.Println()
}

// TurnOffTest 熄灭测试
func TurnOffTest() {
	InitDataMatrix(100000)

	for i := 0; i < 100000; i++ {
		idx := ButterflyMatrix.GetIndex()
		id := int64(i + 1)
		ButterflyMatrix.LightUpPoint(id, idx)
	}

	ButterflyMatrix.StartMrcs()

	for i := 0; i < 100000; i++ {
		id := int64(i + 1)
		ButterflyMatrix.TurnOffPoint(id, uint64(id))
	}

	cost := ButterflyMatrix.FinishMrcs()

	fmt.Println("\r\n熄灭阵点耗时[10 0000]: ", cost, "微秒")
	fmt.Println()
}

// RemovePointTest 移除数据测试
// 移除操作非常耗时
func RemovePointTest() {
	InitDataMatrix(100000)

	for i := 0; i < 100000; i++ {
		idx := ButterflyMatrix.GetIndex()
		id := int64(i + 1)
		ButterflyMatrix.LightUpPoint(id, idx)
	}

	ids := GetDataSourceIds(10)
	fmt.Println("id len >>>", len(ids))

	ButterflyMatrix.StartMrcs()

	SortInt64(ids, true)
	ButterflyMatrix.RemovePoints(ids)

	cost := ButterflyMatrix.FinishMrcs()

	fmt.Println("\r\n移除数据耗时[10 0000]: ", cost, "微秒")
	fmt.Println()
}
