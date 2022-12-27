package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xuliangTang/athena/athena"
	"testing"
)

type bar struct {
	Id   int
	Name string
}

func Test_SlicePage(t *testing.T) {
	var bars []*bar
	for i := 0; i < 23; i++ {
		bars = append(bars, &bar{Id: i + 1, Name: fmt.Sprintf("test-%d", i+1)})
	}

	ibars := make([]any, len(bars))
	for i, bar := range bars {
		ibars[i] = bar
	}

	is := assert.New(t)

	page1 := athena.NewPage(1, 10)
	s1, e1 := page1.SlicePage(ibars)
	is.Equal(s1, 0)
	is.Equal(e1, int64(10))
	is.Equal(page1.TotalPage, 3)
	is.Equal(page1.TotalSize, int64(23))

	page2 := athena.NewPage(3, 10)
	s2, e2 := page2.SlicePage(ibars)
	is.Equal(s2, 20)
	is.Equal(e2, int64(23))
	is.Equal(page2.TotalPage, 3)
	is.Equal(page2.TotalSize, int64(23))

	page3 := athena.NewPage(4, 10)
	s3, e3 := page3.SlicePage(ibars)
	is.Equal(s3, 0)
	is.Equal(e3, int64(0))
	is.Equal(page3.TotalPage, 3)
	is.Equal(page3.TotalSize, int64(23))
}
