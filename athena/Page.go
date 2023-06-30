package athena

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"strconv"
)

// DbQuery 查询对象
type DbQuery struct {
	Db       *gorm.DB    `json:"-"`
	Debug    bool        `json:"-"`
	Order    string      `json:"-"`
	Fields   []string    `json:"-"`
	Preloads []*Preload  `json:"-"`
	Joins    []*Join     `json:"-"`
	Where    *Conditions `json:"-"`
}

func NewDbQuery() *DbQuery {
	return &DbQuery{}
}

// SetDb 设置db对象
func (this *DbQuery) SetDb(db *gorm.DB) *DbQuery {
	this.Db = db
	return this
}

// SetDb 设置Debug
func (this *DbQuery) SetDebug() *DbQuery {
	this.Debug = true
	return this
}

// SetWhere 设置查询条件
func (this *DbQuery) SetWhere(where *Conditions) *DbQuery {
	this.Where = where
	return this
}

// SetOrder 设置排序
func (this *DbQuery) SetOrder(order string) *DbQuery {
	this.Order = order
	return this
}

// AddPreloads 设置预加载
func (this *DbQuery) AddPreloads(preloads ...*Preload) *DbQuery {
	this.Preloads = append(this.Preloads, preloads...)
	return this
}

// AddJoins 设置关联
func (this *DbQuery) AddJoins(joins ...*Join) *DbQuery {
	this.Joins = append(this.Joins, joins...)
	return this
}

// AddFields 设置查询字段
func (this *DbQuery) AddFields(fields ...string) *DbQuery {
	this.Fields = append(this.Fields, fields...)
	return this
}

// SetCountBuildOpts 设置查询选项
func (this *DbQuery) SetBuildOpts(build *gorm.DB) *gorm.DB {
	if this.Debug {
		build = build.Debug()
	}

	if len(this.Fields) > 0 {
		build = build.Select(this.Fields)
	}

	if this.Order != "" {
		build = build.Order(this.Order)
	}

	if len(this.Preloads) > 0 {
		for _, p := range this.Preloads {
			build = build.Preload(p.Query, p.Args...)
		}
	}

	if len(this.Joins) > 0 {
		for _, j := range this.Joins {
			build = build.Joins(j.Query, j.Args...)
		}
	}

	if this.Where != nil {
		build = build.Where(this.Where.Query, this.Where.Args...)
	}

	return build
}

// SetCountBuildOpts 设置统计数量的查询选项
func (this *DbQuery) SetCountBuildOpts(countBuild *gorm.DB) *gorm.DB {
	if this.Debug {
		countBuild = countBuild.Debug()
	}

	if len(this.Joins) > 0 {
		for _, j := range this.Joins {
			countBuild = countBuild.Joins(j.Query, j.Args...)
		}
	}

	if this.Where != nil {
		countBuild = countBuild.Where(this.Where.Query, this.Where.Args...)
	}

	return countBuild
}

// Page 快速分页对象
type Page struct {
	*DbQuery
	CurrentPage int   `json:"current_page"`     // 当前页
	PerPage     int   `json:"per_page"`         // 每页条数
	TotalSize   int64 `json:"total_size"`       // 总条数
	TotalPage   int   `json:"total_page"`       // 总页数
	Extend      any   `json:"extend,omitempty"` // 扩展字段
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

	return &Page{CurrentPage: currentPage, PerPage: perPage, DbQuery: &DbQuery{Order: getOrder}}
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
	build := this.Db.Limit(this.PerPage).Offset(this.GetOffset())
	buildCount := this.Db.Model(items)

	// 设置query参数
	build = this.SetBuildOpts(build)
	buildCount = this.SetCountBuildOpts(buildCount)

	// 查询
	tx := build.Find(items)

	// 统计总数
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
