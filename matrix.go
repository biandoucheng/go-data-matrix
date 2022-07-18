package godatamatrix

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

// 数据矩阵模型

// DataMatrix 数据矩阵
// 注意事项:
// 1、点阵中 1 代表符合,0 代表不符合
// 2、默认点阵为不符合 即: 默认点阵为 0的集合
// 3、数据所在点阵单元列索引 ridx 从0开始,但是数据点在二进制中的位置是从低到高从 1 开始
// 4、如果是新增索引的同时新增数据,需要预先将数据点阵通过 AddPoint 添加进去
// 5、矩阵名称只是个备注性文案
// 6、采用waitGroup 配合 chan 实现多线程处理
type DataMatrix struct {
	sync.RWMutex

	name string // 数据矩阵名称

	count    int      // 数据量(数据移除时这个不会改变,会随着数据添加而递增)
	defPoint []uint64 // 默认点阵,用于初始化一个索引点阵

	index       map[int64][2]int // 数据ID => 数据点阵坐标
	indexNumber uint64           // 最大点阵索引

	matrix map[uint64][]uint64 // 数据索引 => 数据点阵
	remove map[int64][2]int    // 被删除的数据点

	HandleFunc func(interface{}) error // 自定义方法

	createTime time.Time // 创建时间
	startMs    int64     // 开始检测时间,microseconds
	endMs      int64     // 结束检测时间,microseconds
}

// Init 初始化数据矩阵
// uint64 来表示一个点阵单元 64位
// uint64 的二进制位是从低到高利用的,即先遍历的ID其点位越高
// ids 传进来时候会去重一次,建议传入前先去重
func (d *DataMatrix) Init(name string, ids []int64) {
	defer d.Unlock()
	d.Lock()

	// 设置名称
	d.name = name

	if len(ids) == 0 {
		d.index = map[int64][2]int{}
		d.matrix = map[uint64][]uint64{}
		d.remove = map[int64][2]int{}
		return
	}

	// 数据去重
	unik := map[int64]bool{}
	n_ids := []int64{}
	for _, id := range ids {
		_, has := unik[id]
		if !has {
			n_ids = append(n_ids, id)
			unik[id] = true
		}
	}
	ids = n_ids
	d.count = len(ids)

	// 开始初始化
	bnum := len(ids) / 64
	numless := len(ids) % 64
	hasles := numless > 0
	if hasles {
		bnum += 1
	}

	d.defPoint = make([]uint64, bnum)
	d.index = map[int64][2]int{}
	d.matrix = map[uint64][]uint64{}
	d.remove = map[int64][2]int{}

	for idx, id := range ids {
		ridx := idx / 64
		res := idx % 64

		point := 64 - res
		if (ridx+1) >= bnum && hasles {
			point = numless - res
		}

		d.index[id] = [2]int{ridx, point}
	}

	d.createTime = time.Now()
}

// GetName 获取数据矩阵名称
func (d *DataMatrix) GetName() string {
	return d.name
}

// CreateTime 获取创建时间
func (d *DataMatrix) CreateTime() time.Time {
	return d.createTime
}

// GetIndex 获取一个新的点阵索引
func (d *DataMatrix) GetIndex() uint64 {
	defer d.Unlock()
	d.Lock()

	if d.indexNumber == math.MaxUint64 {
		return d.indexNumber
	}

	d.indexNumber += 1
	return d.indexNumber
}

// RemoveIndex 移除索引
func (d *DataMatrix) RemoveIndex(idx uint64) {
	d.Lock()

	delete(d.matrix, idx)

	d.Unlock()
}

// RemoveIndex 移除索引
func (d *DataMatrix) RemoveIndexs(idxs []uint64) {
	d.Lock()

	for _, idx := range idxs {
		delete(d.matrix, idx)
	}

	d.Unlock()
}

// CurrentIndex 获取当前最大索引值
func (d *DataMatrix) CurrentIndexNum() uint64 {
	return d.indexNumber
}

// LightUpPoint 点亮数据点阵单元
func (d *DataMatrix) LightUpPoint(id int64, idx uint64) uint64 {
	defer d.Unlock()
	d.Lock()

	// 传0则默认生成一个索引ID
	if idx == 0 {
		if d.indexNumber != math.MaxUint64 {
			d.indexNumber += 1
		}

		idx = d.indexNumber
	}

	// 初始化索引数据点阵
	points, has := d.matrix[idx]
	if !has {
		points = append([]uint64{}, d.defPoint...)
	}

	// 点亮对应位置的数据位
	coord, ex := d.index[id]
	if !ex {
		return idx
	}

	factor := uint64(1) << (coord[1] - 1)
	points[coord[0]] |= factor

	if !has {
		d.matrix[idx] = points
	}

	return idx
}

// lightUpPoint 点亮数据点阵单元 内部使用
func (d *DataMatrix) lightUpPoint(id int64, idx uint64) {
	// 初始化索引数据点阵
	points, has := d.matrix[idx]
	if !has {
		return
	}

	// 点亮对应位置的数据位
	coord, ex := d.index[id]
	if !ex {
		return
	}

	factor := uint64(1) << (coord[1] - 1)
	points[coord[0]] |= factor

	if !has {
		d.matrix[idx] = points
	}
}

// LightUpPoints 批量点亮数据点阵单元
// 1、该方法使用时各个索引需要已经生产 即: idx > 0
func (d *DataMatrix) LightUpPoints(lmps map[int64][]uint64) {
	defer d.Unlock()
	d.Lock()

	for id, idxs := range lmps {
		for _, idx := range idxs {
			d.lightUpPoint(id, idx)
		}
	}
}

// TurnOffPoint 关闭数据点阵单元
func (d *DataMatrix) TurnOffPoint(id int64, idx uint64) {
	defer d.Unlock()
	d.Lock()

	// 初始化索引数据点阵
	points, has := d.matrix[idx]
	if !has {
		return
	}

	// 关闭对应位置的数据位
	coord, ex := d.index[id]
	if !ex {
		return
	}

	factor := uint64(1) << (coord[1] - 1)
	points[coord[0]] &= ^factor
}

// turnOffPoint 关闭数据点阵单元,内部方法
func (d *DataMatrix) turnOffPoint(coord [2]int) {
	factor := uint64(1) << (coord[1] - 1)

	for idx := range d.matrix {
		d.matrix[idx][coord[0]] &= ^factor
	}
}

// TurnOffPoints 批量关闭数据点阵单元
func (d *DataMatrix) TurnOffPoints(lmps map[int64][]uint64) {
	defer d.Unlock()
	d.Lock()

	for id, idxs := range lmps {
		for _, idx := range idxs {
			// 初始化索引数据点阵
			points, has := d.matrix[idx]
			if !has {
				continue
			}

			// 关闭对应位置的数据位
			coord, ex := d.index[id]
			if !ex {
				continue
			}

			factor := uint64(1) << (coord[1] - 1)
			points[coord[0]] &= ^factor
		}
	}
}

// TurnOffPointsByIds 根据数据ID全部滚逼数据阵点单元
func (d *DataMatrix) TurnOffPointsByIds(ids []int64) {
	defer d.Unlock()
	d.Lock()

	for _, id := range ids {
		coord, ex := d.index[id]
		if !ex {
			continue
		}

		factor := uint64(1) << (coord[1] - 1)

		for idx := range d.matrix {
			points := d.matrix[idx]
			points[coord[0]] &= ^factor
		}
	}
}

// RemovePoint 移除一个数据点
func (d *DataMatrix) RemovePoint(id int64) {
	defer d.Unlock()
	d.Lock()

	points, has := d.index[id]
	if !has {
		return
	}

	// 标记移除数据点
	d.remove[id] = points

	// 关闭数据点阵
	d.turnOffPoint(points)

	// 移除数据点
	delete(d.index, id)
}

// removePoint 内部用批量移除数据点
func (d *DataMatrix) removePoint(id int64) {
	points, has := d.index[id]
	if !has {
		return
	}

	// 标记移除数据点
	d.remove[id] = points

	// 关闭数据点阵
	d.turnOffPoint(points)

	// 移除数据点
	delete(d.index, id)
}

// RemovePoints 批量移除数据点
func (d *DataMatrix) RemovePoints(ids []int64) {
	defer d.Unlock()
	d.Lock()

	for _, id := range ids {
		d.removePoint(id)
	}
}

// AddPoint 添加一个数据阵点
func (d *DataMatrix) AddPoint(id int64) {
	defer d.Unlock()
	d.Lock()

	// 判断数据是否被初始化过
	if len(d.defPoint) == 0 {
		return
	}

	// 判断是否已存在
	if _, has := d.index[id]; has {
		return
	}

	// 先复用删除的点位
	if len(d.remove) > 0 {
		fmt.Println("exp: 复用", id)
		for del_id, coord := range d.remove {
			d.index[id] = coord
			delete(d.remove, del_id)
			break
		}

		return
	}

	// 没有可复用的再新增

	// 计算新的点位
	res := d.count % 64
	old_rdx := len(d.defPoint) - 1
	new_rdx := old_rdx
	expand := res == 0
	if expand {
		new_rdx += 1
		d.defPoint = append(d.defPoint, uint64(0))
	}

	// 追加
	for old_id, coord := range d.index {
		if expand {
			continue
		}

		if coord[0] == old_rdx {
			coord[1] += 1
			d.index[old_id] = coord
		}
	}
	d.index[id] = [2]int{new_rdx, 1}

	// 索引点阵位移
	for idx := range d.matrix {
		if !expand {
			d.matrix[idx][old_rdx] <<= 1
		} else {
			d.matrix[idx] = append(d.matrix[idx], uint64(0))
		}
	}

	d.count += 1
}

// GetPoint 获取点阵
func (d *DataMatrix) GetPoint(idx uint64) ([]uint64, bool) {
	defer d.RUnlock()
	d.RLock()

	points, has := d.matrix[idx]
	if has {
		return points, true
	}

	return d.defPoint, false
}

// BitAndOperate 与运算
func (d *DataMatrix) BitAndOperate(idx1 []uint64, idx2 []uint64) ([]uint64, bool) {
	flg := false

	l1 := len(idx1)
	l2 := len(idx2)
	if l1 > l2 {
		l1 = l2
	}

	if l1 == 0 {
		return []uint64{}, false
	}

	res := make([]uint64, l1)
	for i := 0; i < l1; i++ {
		res[i] = idx1[i] & idx2[i]
		flg = flg || (res[i] > 0)
	}

	return res, flg
}

// BitOrOperate 或运算
func (d *DataMatrix) BitOrOperate(idx1 []uint64, idx2 []uint64) ([]uint64, bool) {
	flg := false

	l1 := len(idx1)
	l2 := len(idx2)
	if l1 > l2 {
		l1 = l2
	}

	if l1 == 0 {
		return []uint64{}, false
	}

	res := make([]uint64, l1)
	for i := 0; i < l1; i++ {
		res[i] = idx1[i] | idx2[i]
		flg = flg || (res[i] > 0)
	}

	return res, flg
}

// GetAndPoint 与运算集合
// empFull 空处理, true 遇到空的直接跳过,false 返回空数据
func (d *DataMatrix) GetAndPoint(idxs []uint64, empFull bool) []uint64 {
	defer d.RUnlock()
	d.RLock()

	points := []uint64{}
	assign := true

	for _, idx := range idxs {
		tmp, has := d.matrix[idx]

		// 空处理
		if !has {
			if empFull {
				continue
			}
			return []uint64{}
		}

		// 初始赋值
		if assign {
			points = append(points, tmp...)
			assign = false
			continue
		}

		// 按位与
		for j, val := range tmp {
			points[j] &= val
		}
	}

	return points
}

// GetOrPoint 或运算集合
// empFull 空处理, true 遇到空的直接跳过,false 返回空数据
func (d *DataMatrix) GetOrPoint(idxs []uint64, empFull bool) []uint64 {
	defer d.RUnlock()
	d.RLock()

	points := []uint64{}
	assign := true

	for _, idx := range idxs {
		tmp, has := d.matrix[idx]

		// 空处理
		if !has {
			if empFull {
				continue
			}
			return points
		}

		// 初始赋值
		if assign {
			points = append(points, tmp...)
			assign = false
			continue
		}

		// 按位或
		for j, val := range tmp {
			points[j] |= val
		}
	}

	return points
}

// IsPointHit 判断坐标是否命中
func (d *DataMatrix) isPointHit(id int64, coord [2]int, points []uint64) bool {
	tl := coord[0] + 1
	if len(points) < tl {
		return false
	}

	point := points[coord[0]]
	judge := uint64(1) << (coord[1] - 1)

	return (judge & point) > 0
}

// GetIds 获取数据ID
func (d *DataMatrix) GetIds(points []uint64) []int64 {
	if len(points) == 0 {
		return []int64{}
	}

	defer d.RUnlock()
	d.RLock()

	ids := []int64{}

	for id, coord := range d.index {
		if d.isPointHit(id, coord, points) {
			ids = append(ids, id)
		}
	}

	return ids
}

// Handler 自定义方法
func (d *DataMatrix) Handler(in interface{}) error {
	return d.HandleFunc(in)
}

// String 输出打印信息
func (d *DataMatrix) String() string {
	defer d.RUnlock()
	d.RLock()

	template := `
矩阵名称: %s
数据个数: %d
最大索引: %d
默认点阵:
%s
点阵索引: 
%s
移除数据:
%s
数据矩阵:
%s
	`

	name := d.GetName()
	count := d.count
	idxmax := d.indexNumber
	def_s := fmt.Sprintf("	%v", d.defPoint)

	idxs := []string{}
	for id, coord := range d.index {
		s := fmt.Sprintf("	ID: %d / Position: %v", id, coord)
		idxs = append(idxs, s)
	}
	idxs_s := strings.Join(idxs, "\r\n")

	moves := []string{}
	for id, coord := range d.remove {
		s := fmt.Sprintf("	ID: %d / Position: %v", id, coord)
		moves = append(moves, s)
	}
	moves_s := strings.Join(moves, "\r\n")

	matrixs := []string{}
	for idx, points := range d.matrix {
		s := fmt.Sprintf("	Index: %d / BitMap: %v", idx, points)
		matrixs = append(matrixs, s)
	}
	matrixs_s := strings.Join(matrixs, "\r\n")

	return fmt.Sprintf(template, name, count, idxmax, def_s, idxs_s, moves_s, matrixs_s)
}

// StartMs 开始耗时统计 微秒
func (d *DataMatrix) StartMrcs() {
	d.Lock()

	d.startMs = time.Now().UnixNano() / 1e3

	d.Unlock()
}

// GetMs 获取执行耗时 微秒
func (d *DataMatrix) FinishMrcs() int64 {
	defer d.Unlock()
	d.Lock()

	d.endMs = time.Now().UnixNano() / 1e3
	return d.endMs - d.startMs
}
