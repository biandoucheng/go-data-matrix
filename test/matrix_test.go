package test

import (
	"fmt"

	"github.com/biandoucheng/go-data-matrix/model"
)

// 测试类别

// InitMatrix 初始化测试
func InitMatrix() {
	// 生产测试数据
	DataGererate(10)
	ids := GetDataIds()
	name := "蝴蝶品种统计"

	mtx := model.DataMatrix{}
	mtx.Init(name, ids)

	s := mtx.String()
	fmt.Println(s)

	PrintLn("InitMatrix", nil)
}

// SetPointIndex 设置索引
func SetPointIndex() {
	DataGererate(191)
	ids := GetDataIds()
	name := "蝴蝶品种统计"

	// 初始化数据矩阵
	mtx := model.DataMatrix{}
	mtx.Init(name, ids)

	// 开启一个数据点
	idx := mtx.GetIndex()
	mtx.LightUpPoint(5, idx)

	// 开启一个数据点
	idx = mtx.GetIndex()
	mtx.LightUpPoint(6, idx)

	// 关闭一个数据点
	mtx.TurnOffPoint(6, idx)

	// 移除数据
	mtx.RemovePoints([]int64{6})

	// 添加数据点
	mtx.AddPoint(192)
	mtx.AddPoint(193)
	mtx.AddPoint(194)

	// 集合运算
	idx = mtx.GetIndex()
	mtx.LightUpPoint(194, idx)

	points := mtx.GetOrPoint([]uint64{1, 3}, true)
	cids := mtx.GetIds(points)

	s := mtx.String()
	fmt.Println(s)
	fmt.Println("get ids >>", cids)
}
