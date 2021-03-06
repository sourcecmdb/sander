// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"context"
	"net/url"

	"sander/db"
	"sander/logger"
	"sander/model"

	"github.com/polaris1119/goutils"
)

type TopicNodeLogic struct{}

var DefaultNode = TopicNodeLogic{}

func (self TopicNodeLogic) FindOne(nid int) *model.TopicNode {
	topicNode := &model.TopicNode{}
	_, err := db.MasterDB.Id(nid).Get(topicNode)
	if err != nil {
		logger.Error("TopicNodeLogic FindOne,nid:%+v,error:%+v", nid, err)
	}

	return topicNode
}

func (self TopicNodeLogic) FindByEname(ename string) *model.TopicNode {
	topicNode := &model.TopicNode{}
	_, err := db.MasterDB.Where("ename=?", ename).Get(topicNode)
	if err != nil {
		logger.Error("TopicNodeLogic FindByEname ename:%+v, error:", ename, err)
	}

	return topicNode
}

func (self TopicNodeLogic) FindByNids(nids []int) map[int]*model.TopicNode {
	nodeList := make(map[int]*model.TopicNode, 0)
	err := db.MasterDB.In("nid", nids).Find(&nodeList)
	if err != nil {
		logger.Error("TopicNodeLogic FindByNids nids:%+v, error:", nids, err)
	}

	return nodeList
}

func (self TopicNodeLogic) FindByParent(pid, num int) []*model.TopicNode {
	nodeList := make([]*model.TopicNode, 0)
	err := db.MasterDB.Where("parent=?", pid).Limit(num).Find(&nodeList)
	if err != nil {
		logger.Error("TopicNodeLogic FindByParent parent:%+v, error:", pid, err)
	}

	return nodeList
}

func (self TopicNodeLogic) FindAll(ctx context.Context) []*model.TopicNode {
	nodeList := make([]*model.TopicNode, 0)
	err := db.MasterDB.Asc("seq").Find(&nodeList)
	if err != nil {
		logger.Error("TopicNodeLogic FindAll error:%+v", err)
	}

	return nodeList
}

func (self TopicNodeLogic) Modify(ctx context.Context, form url.Values) error {

	node := &model.TopicNode{}
	err := schemaDecoder.Decode(node, form)
	if err != nil {
		logger.Error("TopicNodeLogic Modify decode error:", err)
		return err
	}

	nid := goutils.MustInt(form.Get("nid"))
	if nid == 0 {
		// 新增
		_, err = db.MasterDB.Insert(node)
		if err != nil {
			logger.Error("TopicNodeLogic Modify insert error:", err)
		}
		return err
	}

	change := make(map[string]interface{})

	fields := []string{"parent", "logo", "name", "ename", "intro", "seq", "show_index"}
	for _, field := range fields {
		change[field] = form.Get(field)
	}

	_, err = db.MasterDB.Table(new(model.TopicNode)).Id(nid).Update(change)
	if err != nil {
		logger.Error("TopicNodeLogic Modify update error:", err)
	}
	return err
}

func (self TopicNodeLogic) ModifySeq(ctx context.Context, nid, seq int) error {
	_, err := db.MasterDB.Table(new(model.TopicNode)).Id(nid).Update(map[string]interface{}{"seq": seq})
	return err
}

func (self TopicNodeLogic) FindParallelTree(ctx context.Context) []*model.TopicNode {
	nodeList := make([]*model.TopicNode, 0)
	err := db.MasterDB.Asc("parent").Asc("seq").Find(&nodeList)
	if err != nil {
		logger.Error("TopicNodeLogic FindTreeList error:%+v", err)

		return nil
	}

	showNodeList := make([]*model.TopicNode, 0, len(nodeList))
	self.tileNodes(&showNodeList, nodeList, 0, 1, 3, 0)

	return showNodeList
}

func (self TopicNodeLogic) tileNodes(showNodeList *[]*model.TopicNode, nodeList []*model.TopicNode, parentId, curLevel, showLevel, pos int) {
	for num := len(nodeList); pos < num; pos++ {
		node := nodeList[pos]

		if node.Parent == parentId {
			*showNodeList = append(*showNodeList, node)

			if node.Level == 0 {
				node.Level = curLevel
			}

			if curLevel <= showLevel {
				self.tileNodes(showNodeList, nodeList, node.Nid, curLevel+1, showLevel, pos+1)
			}
		}

		if node.Parent > parentId {
			break
		}
	}
}
