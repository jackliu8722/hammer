package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/daoone/hammer/model"
	"github.com/daoone/hammer/util"
	"strconv"
)

func TaskList(unionId, status string ) gin.H {
	account := &model.Account{}
	_, err := model.GetEngine().Cols("b_account.*").
		Join("LEFT", "b_user", "b_account.uid=b_user.uid").
		Where("b_user.union_id=?", unionId).Get(account)
	if err != nil || account.Id == 0{
		return util.WxError("账户不存在")
	}

	stat, _ := strconv.Atoi(status)
	tasks := make([]model.Task, 0)

	_, err3 := model.GetEngine().Cols("b_task.*").
		Join("LEFT", "b_task_user_rel","b_task.id=b_task_user_rel.task_id").
		Where("b_task_user_rel.uid=?",account.Id).
		And("b_task_user_rel.status=?",stat).Get(&tasks)
	if err3 != nil {
		return util.WxError("服务内部错误")
	}

	res := make([]map[string]string, 0)
	for _, t := range tasks {
		a := make(map[string]string)
		a["tid"] = strconv.FormatInt(t.Id, 10)
		a["title"] = t.Title
		a["start_at"] = t.StartAt.Format("2006/01/02 15:04:05")
		a["end_at"] = t.EndAt.Format("2006/01/02 15:04:05")
		res = append(res, a)
	}
	return gin.H{"success":true, "tasks": res}
}
