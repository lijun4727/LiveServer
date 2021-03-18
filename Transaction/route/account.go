package route

import (
	"LiveServer/Base/jwt"
	"LiveServer/Transaction/db_connect"
	"LiveServer/common"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	ziu_mysql "github.com/ziutek/mymysql/mysql"
)

type AccountRoute struct {
	x *xorm.Engine
}

func (ar *AccountRoute) login(c *gin.Context) {
	var loginData common.AccountDevice
	var errorCode = common.ERROR_CODE_NONE
	var err error
	var token string
	defer func() {
		var message string
		if err == nil {
			message = "login sucess"
		} else {
			message = err.Error()
		}

		c.JSON(http.StatusOK, gin.H{
			"errorCode":     errorCode,
			"errorCodeName": common.ERROR_CODE_name[errorCode],
			"token":         token,
			"message":       message,
		})
	}()

	err = c.BindJSON(&loginData)
	if err != nil {
		errorCode = common.ERROR_CODE_DATA_INVALID
	} else {
		if isNum, _ := regexp.MatchString("^[0-9]*$", loginData.Phone); !isNum || len(loginData.Phone) != 11 {
			errorCode = common.ERROR_CODE_PHONE_INVALID
			return
		}
		if len(loginData.DeviceId) == 0 {
			errorCode = common.ERROR_CODE_DEVICE_INVALID
			return
		}

		has, err := ar.x.Table(&common.AccountData{}).Where("phone=?", loginData.Phone).Exist()
		if err != nil {
			errorCode = common.ERROR_CODE_INTERNAL_ERROR
			return
		}
		if !has {
			_, err = ar.x.Insert(common.AccountData{
				Phone: loginData.Phone,
			})
			if err != nil {
				var mySQLErr *mysql.MySQLError
				if errors.As(err, &mySQLErr) && mySQLErr.Number != ziu_mysql.ER_DUP_ENTRY {
					errorCode = common.ERROR_CODE_INTERNAL_ERROR
					return
				}
			}
		}

		has, err = ar.x.Table(&common.AccountDevice{}).Where("phone=?", loginData.Phone).And("device_id=?", loginData.DeviceId).Exist()
		if err != nil {
			errorCode = common.ERROR_CODE_INTERNAL_ERROR
			return
		}
		if !has {
			_, err = ar.x.Insert(loginData)
			if err != nil {
				var mySQLErr *mysql.MySQLError
				if errors.As(err, &mySQLErr) && mySQLErr.Number != ziu_mysql.ER_DUP_ENTRY {
					errorCode = common.ERROR_CODE_INTERNAL_ERROR
					return
				}
			}
		}

		token, err = jwt.CreateToken(loginData.Phone, loginData.DeviceId, common.TokenExpireTime)
		if err != nil {
			errorCode = common.ERROR_CODE_INTERNAL_ERROR
		}
	}
}

func (ar *AccountRoute) Route(e *gin.Engine) {
	r := e.Group("/account")
	r.POST("/login", ar.login)
}

func (ar *AccountRoute) Init() {
	ar.x = db_connect.MysqlEngine
	session := ar.x.NewSession()
	defer session.Close()
	session.Begin()
	if isTableExit, _ := session.IsTableExist(&common.AccountData{}); !isTableExit {
		err := session.Sync2(&common.AccountData{})
		if err != nil {
			log.Fatalf("failed to create table AccountData: %v", err)
		}
	}

	if isTableExit, _ := session.IsTableExist(&common.AccountDevice{}); !isTableExit {
		err := session.Sync2(&common.AccountDevice{})
		if err != nil {
			log.Fatalf("failed to create table AccountData: %v", err)
		}

		fkPhoneName := "fk_account_device_phone"
		tableName := ar.x.TableName(&common.AccountDevice{})
		sqlStr := fmt.Sprintf(`alter table %s add constraint %s foreign key (phone)
		references account_data(phone) ON DELETE CASCADE ON UPDATE CASCADE`, tableName, fkPhoneName)
		_, err = session.Exec(sqlStr)
		if err != nil {
			log.Fatalf("failed to create foreign key %s: %v", fkPhoneName, err)
			return
		}
	}
}
