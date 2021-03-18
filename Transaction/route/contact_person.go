package route

import (
	"LiveServer/Base/jwt"
	"LiveServer/Transaction/db_connect"
	"LiveServer/common"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type ContactPersons struct {
	x *xorm.Engine
}

// 联系人太多时,需要分页
func (cp *ContactPersons) getContactPersons(c *gin.Context) {
	var queryContact common.QueryContact
	var errorCode = common.ERROR_CODE_NONE
	err := c.BindJSON(&queryContact)
	var userInfos []common.AccountDevice
	defer func() {
		var message string
		if err == nil {
			message = "get contact persons sucess"
		} else {
			message = err.Error()
		}

		c.JSON(http.StatusOK, gin.H{
			"errorCode":     errorCode,
			"errorCodeName": common.ERROR_CODE_name[errorCode],
			"message":       message,
			"contacts":      userInfos,
		})
	}()

	if err != nil {
		errorCode = common.ERROR_CODE_DATA_INVALID
	} else {
		_, err = jwt.VerifyToken(queryContact.Token)
		if err == nil {
			x := db_connect.MysqlEngine
			contactPersonTableName := x.TableName(&common.ContactPerson{})
			accountDeviceName := x.TableName(&common.AccountDevice{})
			strSelect := fmt.Sprintf("%s.phone,%s.device_id", accountDeviceName, accountDeviceName)
			strJoin := fmt.Sprintf("%s.contact_phone = %s.phone", contactPersonTableName, accountDeviceName)
			strWhere := fmt.Sprintf("%s.phone=%s", contactPersonTableName, queryContact.Phone)
			err = x.Table(contactPersonTableName).Select(strSelect).Join("INNER", accountDeviceName, strJoin).Where(strWhere).Find(&userInfos)
			if err != nil {
				log.Printf("getContactPersons fail:%v", err)
				errorCode = common.ERROR_CODE_INTERNAL_ERROR
			}
		} else {
			errorCode = common.ERROR_CODE_TOKEN_INVALID
		}
	}

	// if errorCode == common.ERROR_CODE_NONE {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"errorCode":     errorCode,
	// 		"errorCodeName": common.ERROR_CODE_name[errorCode],
	// 		"contacts":      userInfos,
	// 	})
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"errorMessage":  err.Error(),
	// 		"errorCode":     errorCode,
	// 		"errorCodeName": common.ERROR_CODE_name[errorCode],
	// 	})
	// }
}

func (cp *ContactPersons) Route(e *gin.Engine) {
	e.POST("/getContactPersons", cp.getContactPersons)
}

func (cp *ContactPersons) Init() {
	x := db_connect.MysqlEngine
	session := x.NewSession()
	defer session.Close()
	session.Begin()
	if isTableExit, _ := x.IsTableExist(&common.ContactPerson{}); !isTableExit {
		err := session.Sync2(&common.ContactPerson{})
		if err != nil {
			log.Fatalf("failed to create table ContactPerson: %v", err)
			return
		}

		fkIDName := "fk_contact_persons_phone"
		tableName := x.TableName(&common.ContactPerson{})
		sqlStr := fmt.Sprintf(`alter table %s add constraint %s foreign key (phone)
		references account_data(phone) ON DELETE CASCADE ON UPDATE CASCADE`, tableName, fkIDName)
		_, err = session.Exec(sqlStr)
		if err != nil {
			log.Fatalf("failed to create foreign key %s: %v", fkIDName, err)
			return
		}

		fkContactIDName := "fk_contact_phone"
		sqlStr = fmt.Sprintf(`alter table %s add constraint %s foreign key (contact_phone)
		references account_data(phone) ON DELETE CASCADE ON UPDATE CASCADE`, tableName, fkContactIDName)
		_, err = session.Exec(sqlStr)
		if err != nil {
			log.Fatalf("failed to create foreign key %s: %v", fkContactIDName, err)
			return
		}

		err = session.Commit()
		if err != nil {
			log.Fatalf("failed to table %s commit error: %v", tableName, err)
		}
	}
}
