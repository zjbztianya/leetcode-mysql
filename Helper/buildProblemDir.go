package main

import (
	"fmt"
	"github.com/aQuaYi/GoKit"
	"log"
	"os"
	"runtime/debug"
	"strings"
)

func buildProblemDir(problemNum int) {
	log.Printf("~~ 开始生成第 %d 题的文件夹 ~~\n", problemNum)

	// 获取 LeetCode 的记录文件
	lc := newLeetCode()

	// 检查 problemNum 的合法性
	if problemNum >= len(lc.Problems) {
		log.Panicf("%d 超出题目范围，请核查题号。", problemNum)
	}
	if lc.Problems[problemNum].ID == 0 {
		log.Panicf("%d 号题不存，请核查题号。", problemNum)
	}
	if lc.Problems[problemNum].IsPaid == true && getConfig().IsPaid == false {
		log.Panicf("%d 号题需要付费。如果已经订阅，请注修改config.toml的Ispaid选项。", problemNum)
	}

	if lc.Problems[problemNum].HasNoMysqlOption {
		log.Panicf("%d 号题，没有提供 Go 解答选项。请核查后，修改 unavailable.json 中的记录。", problemNum)
	}

	// 需要创建答题文件夹
	build(lc.Problems[problemNum])

	log.Printf("~~ 第 %d 题的文件夹，已经生成 ~~\n", problemNum)
}

func buildMultiProblemDir(l int, r int) {
	// 获取 LeetCode 的记录文件
	lc := newLeetCode()
	if l < 1 || r >= len(lc.Problems) {
		log.Panicf("最小题号%d或者最大题号%d超出题目范围，请核查题号。", l, r)
	}
	for i := l; i <= r; i++ {
		//过滤非database类型题目和非go题目
		if lc.Problems[i].ID == 0 || lc.Problems[i].HasNoMysqlOption {
			continue
		}
		if lc.Problems[i].IsPaid == true && getConfig().IsPaid == false {
			continue
		}
		log.Printf("~~ 开始生成第 %d 题的文件夹 ~~\n", i)
		// 需要创建答题文件夹
		build(lc.Problems[i])

		log.Printf("~~ 第 %d 题的文件夹，已经生成 ~~\n", i)
	}
}

func build(p problem) {
	if GoKit.Exist(p.Dir()) {
		log.Panicf("第 %d 题的文件夹已经存在，请 **移除** %s 文件夹后，再尝试。", p.ID, p.Dir())
	}

	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			log.Println(err)
			log.Println("清理不必要的文件")
			os.RemoveAll(p.Dir())
		}
	}()

	// windows用户注释这两行
	//mask := syscall.Umask(0)
	//defer syscall.Umask(mask)

	// 创建目录
	err := os.MkdirAll(p.Dir(), 0755)
	if err != nil {
		log.Panicf("无法创建目录，%s ：%s", p.Dir(), err)
	}

	log.Printf("开始创建 %d %s 的文件夹...\n", p.ID, p.Title)

	content, fc, mysqlSchemas := getGraphql(p)
	if fc == "" {
		log.Panicf("查无mysql写法")
	}

	// 利用 chrome 打开题目页面
	// go func() {
	// 	cmd := exec.Command("google-chrome", p.link())
	// 	_, err = cmd.Output()
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// }()

	// fc := getFunction(p.link())

	//fcName, para, ans, _ := parseFunction(fc)

	creatSql(p, mysqlSchemas)

	//creatGoTest(p, fcName, para, ans)

	creatREADME(p, content)

	log.Printf("%d.%s 的文件夹，创建完毕。\n", p.ID, p.Title)
}

var typeMap = map[string]string{
	"int":     "0",
	"float64": "0",
	"string":  "\"\"",
	"bool":    "false",
}

func creatSql(p problem, mysqlSchemas []string) {
	filename := fmt.Sprintf("%s/%s.sql", p.Dir(), p.TitleSlug)
	var builder strings.Builder
	for _, schemas := range mysqlSchemas {
		builder.WriteString(schemas + ";\n")
	}
	write(filename, builder.String())

	//vscodeOpen(filename)
}

// 把 函数的参数 变成 tc 的参数
func getTcPara(para string) string {
	// 把 para 按行切分
	paras := strings.Split(para, "\n")

	// 把单个参数按空格，切分成参数名和参数类型
	temp := make([][]string, len(paras))
	for i := range paras {
		temp[i] = strings.Split(strings.TrimSpace(paras[i]), ` `)
	}

	// 在参数名称前添加 "tc." 并组合在一起
	res := ""
	for i := 0; i < len(temp); i++ {
		res += ", tc." + temp[i][0]
	}

	return res[2:]
}

func (p problem) packageName() string {
	return fmt.Sprintf("problem%04d", p.ID)
}
