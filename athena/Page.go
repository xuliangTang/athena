package athena

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"strconv"
)

type Page struct {
	Db          *gorm.DB    `json:"-"`
	CurrentPage int         `json:"current_page"`     // 当前页
	PerPage     int         `json:"per_page"`         // 每页条数
	TotalSize   int64       `json:"total_size"`       // 总条数
	TotalPage   int         `json:"total_page"`       // 总页数
	Extend      any         `json:"extend,omitempty"` // 扩展字段
	Order       string      `json:"-"`
	Fields      []string    `json:"-"`
	Preloads    []*Preload  `json:"-"`
	Where       *Conditions `json:"-"`
}

func NewPage(currentPage int, perPage int) *Page {
	return &Page{CurrentPage: currentPage, PerPage: perPage}
}

// NewPageWithCtx 通过 query 参数创建
func NewPageWithCtx(ctx *gin.Context) *Page {
	getCurrentPage := ctx.Query("page")
	currentPage, err := strconv.Atoi(getCurrentPage)
	if err != nil {
		currentPage = 1
	}

	getPerPage := ctx.Query("per_page")
	perPage, err := strconv.Atoi(getPerPage)
	if err != nil {
		perPage = 20
	}

	getOrder := ctx.Query("order")
	if getOrder == "" {
		getOrder = "id DESC"
	}

	return &Page{CurrentPage: currentPage, PerPage: perPage, Order: getOrder, Fields: []string{"*"}}
}

// IsValid 是否有效
func (this *Page) IsValid() bool {
	return this.CurrentPage > 0 && this.PerPage > 0
}

// GetOffset 获取 Offset
func (this *Page) GetOffset() int {
	return (this.CurrentPage - 1) * this.PerPage
}

// SetTotal 设置总条数和总页数
func (this *Page) SetTotal(totalSize int64) {
	this.TotalSize = totalSize
	this.TotalPage = int(math.Ceil(float64(this.TotalSize) / float64(this.PerPage)))
}

// SelectList 列表查询
func (this *Page) SelectList(items any) *gorm.DB {
	build := this.Db.Select(this.Fields).Limit(this.PerPage).Offset(this.GetOffset()).Order(this.Order)
	buildCount := this.Db.Model(items)

	if len(this.Preloads) > 0 {
		for _, p := range this.Preloads {
			build = build.Preload(p.Query, p.Args...)
		}
	}

	if this.Where != nil {
		build = build.Where(this.Where.Query, this.Where.Args...)
		buildCount = buildCount.Where(this.Where.Query, this.Where.Args...)
	}

	tx := build.Find(items)

	var totalSize int64
	buildCount.Count(&totalSize)
	this.SetTotal(totalSize)

	return tx
}

// SlicePage 切片分页
func (this *Page) SlicePage(list []any) (start int, end int64) {
	this.TotalSize = int64(len(list))

	this.TotalPage = int(math.Ceil(float64(this.TotalSize) / float64(this.PerPage)))
	if this.CurrentPage > this.TotalPage {
		// this.CurrentPage = this.TotalPage
		return 0, 0
	}

	start = (this.CurrentPage - 1) * this.PerPage
	end = int64(start + this.PerPage)
	if end > this.TotalSize {
		end = this.TotalSize
	}

	//ret = list[start:end]
	return start, end
}

// Collection 分页集合
type Collection struct {
	Items any   `json:"items"`
	Page  *Page `json:"page"`
}

func NewCollection(items any, page *Page) *Collection {
	return &Collection{Items: items, Page: page}
}
