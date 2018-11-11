package dao

import (
	xe "be/common/error"
	"be/common/log"
	"be/mysql"
)

type ProjectDao struct {
}

func (d *ProjectDao) CreateProject(fullName string, sourceCodeIp string) (int64, error) {
	tx := mysql.DB.GetTx()
	sql := "INSERT INTO PROJECT (fullName, status, sourceCodeIp, config) VALUES(?, '项目已创建', ?, '')"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.WithFields(log.Fields{
			"sql": sql,
			"err": err.Error(),
		}).Error("prepare错误")
		tx.Rollback()
		return -1, err
	}
	result, err := stmt.Exec(fullName, sourceCodeIp)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("query错误")
		stmt.Close()
		tx.Rollback()
		return -1, err
	}

	projectId, err := result.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("LastInsertId错误")
		stmt.Close()
		tx.Rollback()
		return -1, err
	}
	stmt.Close()
	tx.Commit()
	return projectId, nil
}

func (d *ProjectDao) BindUserAndProject(projectId int64, username string) error {
	tx := mysql.DB.GetTx()
	sql := "INSERT INTO USER_PROJECT (username, projectId) VALUES(?, ?)"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.WithFields(log.Fields{
			"sql": sql,
			"err": err.Error(),
		}).Error("prepare错误")
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(username, projectId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("query错误")
		stmt.Close()
		tx.Rollback()
		return err
	}
	stmt.Close()
	tx.Commit()
	return nil
}

func (d *ProjectDao) GetProjectById(projectId int64) (id int64, fullName string, status string, sourceCodeIp string, config string, err error) {
	var cnt int64 = -1
	cnt, err = mysql.DB.SingleRowQuery("SELECT id, fullName, status, sourceCodeIp, config FROM PROJECT WHERE id=?", []interface{}{projectId}, &id, &fullName, &status, &sourceCodeIp, &config)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("SingleRowQuery错误")
		return 0, "", "", "", "", err
	} else {
		if cnt == 0 {
			return 0, "", "", "", "", xe.New("记录不存在")
		} else {
			return id, fullName, status, sourceCodeIp, config, nil
		}
	}
}

func (d *ProjectDao) UpdateStatus(projectId int64, status string) error {
	tx := mysql.DB.GetTx()
	sql := "UPDATE PROJECT SET status=? WHERE id=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.WithFields(log.Fields{
			"sql": sql,
			"err": err.Error(),
		}).Error("prepare错误")
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(status, projectId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("query错误")
		stmt.Close()
		tx.Rollback()
		return err
	}
	stmt.Close()
	tx.Commit()
	return nil
}
