package test

// 测试入口

import (
	"fmt"
	"testing"
)

// PrintLn 测试打印
func PrintLn(m string, err error, a ...interface{}) {
	s := fmt.Sprintf("---------------------------%s-----------------------", m)
	fmt.Println("err: ", err)
	fmt.Println("")
	fmt.Println(a...)
	fmt.Println(s)
}

func TestRun(t *testing.T) {
	// InitMatrix()

	// SetPointIndex()

	// InitCostTest()

	// LightUpTest()

	// TurnOffTest()

	// RemovePointTest()
}
