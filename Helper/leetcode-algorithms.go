package main

import (
	"encoding/json"
	"log"
	"strconv"
)

// databases 保存API信息
type databases struct {
	Name     string          `json:"category_slug"`
	User     string          `json:"user_name"`
	ACEasy   int             `json:"ac_easy"`
	ACMedium int             `json:"ac_medium"`
	ACHard   int             `json:"ac_hard"`
	AC       int             `json:"num_solved"`
	Problems []problemStatus `json:"stat_status_pairs"`
}

type problemStatus struct {
	Status     string `json:"status"`
	State      `json:"stat"`
	IsFavor    bool `json:"is_favor"`
	IsPaid     bool `json:"paid_only"`
	Difficulty `json:"difficulty"`
}

// State 保存单个问题的解答状态
type State struct {
	ACs       int    `json:"total_acs"`
	Title     string `json:"question__title"`
	IsNew     bool   `json:"is_new_question"`
	Submitted int    `json:"total_submitted"`
	ID        string `json:"frontend_question_id"`
	TitleSlug string `json:"question__title_slug"`
}

// Difficulty 问题的难度
type Difficulty struct {
	Level int `json:"level"`
}

func getDatabases() *databases {
	URL := "https://leetcode-cn.com/api/problems/Database/"

	raw := getRaw(URL)

	res := new(databases)
	if err := json.Unmarshal(raw, res); err != nil {
		log.Panicf("无法把json转换成Category: %s\n", err.Error())
	}

	// 如果，没有登录的话，也能获取数据，但是用户名，就不是本人
	if res.User != getConfig().Username {
		log.Printf("res.User = %s\n", res.User)
		log.Fatal("没有获取到本人的数据")
	}

	problems := make([]problemStatus, 0)
	for _, problem := range res.Problems {
		_, err := strconv.Atoi(problem.State.ID)
		if err == nil {
			problems = append(problems, problem)
		}
	}
	res.Problems = problems
	return res
}
