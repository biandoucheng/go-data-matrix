package test

import "sort"

// SortInt64 int64 冒泡排序
// asc = true 正序
// asc = false 倒叙
func SortInt64(its []int64, asc bool) {
	sort.Slice(its, func(i, j int) bool {
		if asc {
			return its[i] < its[j]
		}
		return its[i] > its[j]
	})
}
