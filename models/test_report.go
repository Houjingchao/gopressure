package models

// 压力测试报告
type TestReport struct {
	// id
	Id int `db:"id" json:"id"`

	// 测试名称
	Name string `db:"name" json:"name"`

	// 请求总数
	Count int `db:"count" json:"count"`

	// 成功总数
	SuccessCount int `db:"success_count" json:"successCount"`

	// 成功率
	SuccessRate int `db:"success_rate" json:"successRate"`

	// 平均耗时ms
	AverageCost int `db:"average_cost" json:"averageCost"`

	// min_cost
	MinCost int `db:"min_cost" json:"minCost"`

	// max_cost
	MaxCost int `db:"max_cost" json:"maxCost"`

	// qps
	Qps int `db:"qps" json:"qps"`

	// create_time
	CreateTime time.Time `db:"create_time" json:"createTime"`
}
