// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"time"

	"math/rand"
	"sander/cmd"
	"sander/config"

	"github.com/polaris1119/keyword"
	"github.com/polaris1119/logger"
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	logger.Init(config.ROOT+"/log", config.ConfigFile.MustValue("global", "log_level", "DEBUG"), "crawl")
	go keyword.Extractor.Init(keyword.DefaultProps, true, config.ROOT+"/data/programming.txt,"+config.ROOT+"/data/dictionary.txt")

	server.CrawlServer()

	select {}
}
