package dao

import (
	xe "be/common/error"
	"be/common/log"
	"be/mysql"
	"be/options"
	"be/structs"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type AuthDao struct {
}

func (d *AuthDao) GenTokenByUsernameAndPassword(username string, password string) (string, error) {
	var id int64 = -1
	cnt, err := mysql.DB.SingleRowQuery("SELECT id FROM USER_AUTH_INFO WHERE username=? AND password=PASSWORD(PASSWORD(?))", []interface{}{username, fmt.Sprintf("%s-%s", password, options.Options.PasswordSalt)}, &id)
	if err != nil {
		log.Errorf("用户密码校验失败: %s %s", username, password)
		return "", err
	}
	if cnt == 0 {
		log.Errorf("用户密码校验失败: %s %s", username, password)
		return "", xe.New("用户密码校验失败")
	}

	token := fmt.Sprintf("%s", uuid.NewV4())
	err = mysql.DB.SimpleInsert("INSERT INTO TOKEN(token, username, expire_time) VALUES(?, ?, ?)", token, username, time.Now().Add(time.Duration(4)*time.Hour).Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Errorf("用户密码校验失败: %s %s", username, password)
		return "", err
	}

	return token, nil
}

func (d *AuthDao) GetUserInfoByToken(token string) (*structs.UserInfo, error) {
	result := &structs.UserInfo{}
	cnt, err := mysql.DB.SingleRowQuery("SELECT username FROM TOKEN WHERE token=? AND expire_time>NOW()", []interface{}{token}, &result.Username)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("GetUserInfoByToken DB错误")
		return nil, xe.DBError()
	}

	if cnt == 0 {
		return nil, xe.AuthError()
	}

	return result, nil
}
