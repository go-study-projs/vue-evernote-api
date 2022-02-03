package utils

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type AggregateM []bson.M

type AggregatePipe struct {
	agg AggregateM
}

// 用于调试时直接数据构建原始语句
func (pipe *AggregatePipe) String() string {
	bt, err := json.Marshal(pipe.agg)
	if err != nil {
		return ""
	}

	return string(bt)
}

func NewAggregatePipe() *AggregatePipe {
	return &AggregatePipe{
		agg: make(AggregateM, 0),
	}
}

// 根据 $or 组数量动态生成 $match 中的语句
// matchOr = []$or
func buildMatchOr(match bson.M, matchOr []interface{}) {
	if match != nil && len(matchOr) > 0 {
		if len(matchOr) > 1 {
			and := make([]interface{}, 0)
			for _, or := range matchOr {
				if or != nil {
					and = append(and, bson.M{"$or": or})
				}
			}
			if len(and) > 0 {
				match["$and"] = and
			}
		} else {
			if matchOr[0] != nil {
				match["$or"] = matchOr[0]
			}
		}
	}
}

func (pipe *AggregatePipe) Match(m bson.M, or ...interface{}) *AggregatePipe {
	if len(or) > 0 {
		if len(m) == 0 {
			m = make(bson.M)
		}

		buildMatchOr(m, or)
	}

	if len(m) > 0 {
		pipe.agg = append(pipe.agg, bson.M{
			"$match": m,
		})
	}

	return pipe
}

// Sort 使用 $group 要注意 sort 位置（要在 $group 之后）及字段名的变化（_id.date）
func (pipe *AggregatePipe) Sort(s bson.M) *AggregatePipe {
	if len(s) > 0 {
		pipe.agg = append(pipe.agg, bson.M{
			"$sort": s,
		})
	}

	return pipe
}

// SortOne
//
//使用 $group 要注意 sort 位置（要在 $group 之后）及字段名的变化（_id.date）
func (pipe *AggregatePipe) SortOne(field string, isDesc bool) *AggregatePipe {
	if len(field) > 0 {
		sort := 1 // ASC 升序
		if isDesc {
			sort = -1 // DESC 降序
		}

		pipe.agg = append(pipe.agg, bson.M{
			"$sort": bson.M{
				field: sort,
			},
		})
	}

	return pipe
}

// Lookup 联表
func (pipe *AggregatePipe) Lookup(lookups ...bson.M) *AggregatePipe {
	if len(lookups) > 0 {
		for _, lp := range lookups {
			if len(lp) > 0 {
				pipe.agg = append(pipe.agg, bson.M{
					"$lookup": lp,
				})
			}
		}
	}

	return pipe
}

func (pipe *AggregatePipe) LookupOne(from, alias, localField, foreignField string) *AggregatePipe {
	if from == "" || alias == "" || localField == "" || foreignField == "" {
		return pipe
	}

	lookup := bson.M{
		"from":         from,
		"as":           alias,
		"localField":   localField,
		"foreignField": foreignField,
	}

	pipe.agg = append(pipe.agg, bson.M{
		"$lookup": lookup,
	})

	return pipe
}

func (pipe *AggregatePipe) Unwind(unwinds ...string) *AggregatePipe {
	if len(unwinds) > 0 {
		for _, uw := range unwinds {
			if len(uw) > 0 {
				pipe.agg = append(pipe.agg, bson.M{
					"$unwind": uw,
				})
			}
		}
	}

	return pipe
}

func (pipe *AggregatePipe) Group(groups ...bson.M) *AggregatePipe {
	if len(groups) > 0 {
		for _, gp := range groups {
			if len(gp) > 0 {
				if _, exist := gp["_id"]; exist {
					pipe.agg = append(pipe.agg, bson.M{
						"$group": gp,
					})
				}
			}
		}
	}

	return pipe
}

// Project 控制输出字段
func (pipe *AggregatePipe) Project(p bson.M) *AggregatePipe {
	if len(p) > 0 {
		pipe.agg = append(pipe.agg, bson.M{
			"$project": p,
		})
	}

	return pipe
}

func (pipe *AggregatePipe) Skip(skip int64) *AggregatePipe {
	if skip > 0 {
		pipe.agg = append(pipe.agg, bson.M{
			"$skip": skip,
		})
	}

	return pipe
}

func (pipe *AggregatePipe) Limit(limit int64) *AggregatePipe {
	if limit > 0 {
		pipe.agg = append(pipe.agg, bson.M{
			"$limit": limit,
		})
	}

	return pipe
}

func (pipe *AggregatePipe) Custom(ms ...bson.M) *AggregatePipe {
	for _, m := range ms {
		if len(m) > 0 {
			pipe.agg = append(pipe.agg, m)
		}
	}

	return pipe
}

// CountM 获取用于统计数据量的 AggregateM
// reservePaginate[0] == true 保留语句中所有的 $skip、$limit
func (pipe *AggregatePipe) CountM(reservePaginate ...bool) AggregateM {
	var agg AggregateM

	// 默认去除语句中所有的 $skip、$limit
	if len(reservePaginate) == 0 || reservePaginate[0] == false {
		for _, a := range pipe.agg {
			if _, exist := a["$skip"]; exist {
				continue
			}
			if _, exist := a["$limit"]; exist {
				continue
			}

			agg = append(agg, a)
		}
	}

	agg = append(agg, bson.M{
		"$count": "total",
	})

	return agg
}

// QueryM 获取用于查询数据的 AggregateM
func (pipe *AggregatePipe) QueryM() AggregateM {
	return pipe.agg
}
