package helper

import (
	"fmt"
	"github.com/lain/athena/athena"
	"gorm.io/gorm"
	"time"
)

type DateBetween struct {
	Start *time.Time
	End   *time.Time
}

// GetBetweenDates 生成区间内每一天的日期
func (this *DateBetween) GetBetweenDates() (allDate []string) {
	allDate = make([]string, 0)
	sInt := this.Start.Unix()
	eInt := this.End.Unix()
	for {
		st := time.Unix(sInt, 0).Format(athena.DateFormat)
		if sInt > eInt {
			return
		}
		allDate = append(allDate, st)
		sInt += 86400
	}
}

type StatisticsDaily struct {
	Db         *gorm.DB `json:"-"`
	Between    *DateBetween
	Model      athena.Model
	Field      string
	GroupField string
	Where      *athena.Conditions
}

func NewStatisticsDaily(db *gorm.DB, between *DateBetween, model athena.Model, field string) *StatisticsDaily {
	return &StatisticsDaily{Db: db, Between: between, Model: model, Field: field, GroupField: field}
}

type StatisticsDailyResult struct {
	Datetime   *time.Time `json:"-"`
	Num        int        `json:"num"`
	DateFormat string     `json:"date"`
}

func (this *StatisticsDaily) SetGroupField(groupField string) {
	this.GroupField = groupField
}

func (this *StatisticsDaily) SetWhere(where *athena.Conditions) {
	this.Where = where
}

// Exec 统计区间内每天的增长数量
func (this *StatisticsDaily) Exec() []*StatisticsDailyResult {
	var retItems, getItems []*StatisticsDailyResult

	build := this.Db.Model(&this.Model)
	if this.Where != nil {
		build.Where(this.Where.Query, this.Where.Args...)
	}
	build.Where(fmt.Sprintf("%s BETWEEN ? AND ?", this.Field), this.Between.Start, this.Between.End).
		Group(this.GroupField).
		Select(fmt.Sprintf("%s AS datetime, COUNT(id) AS num", this.GroupField)).
		Scan(&getItems)

	for _, getItem := range getItems {
		getItem.DateFormat = getItem.Datetime.Format(athena.DateFormat)
	}

	mItems := funk.ToMap(getItems, "DateFormat").(map[string]*StatisticsDailyResult)

	allDate := this.Between.GetBetweenDates()
	for _, date := range allDate {
		if item, ok := mItems[date]; ok {
			retItems = append(retItems, item)
		} else {
			retItems = append(retItems, &StatisticsDailyResult{
				DateFormat: date,
				Num:        0,
			})
		}
	}

	return retItems
}
