package common

// AccountData 账户所有相关数据
type AccountData struct {
	Id    int64
	Phone string `xorm:"unique notnull VARCHAR(11)"`
}

type AccountDevice struct {
	Phone    string `xorm:"pk VARCHAR(11)" json:"phone" binding:"required"`
	DeviceId string `xorm:"pk VARCHAR(50)" json:"device_id" binding:"required"`
}

// ContactPerson 联系人
type ContactPerson struct {
	Phone        string `xorm:"pk VARCHAR(11)"`
	ContactPhone int64  `xorm:"pk VARCHAR(11)"`
}

type QueryContact struct {
	Phone string `json:"phone,required"`
	Token string `json:"token,required"`
}

type ErrorCode int32

const (
	ERROR_CODE_NONE           ErrorCode = 0
	ERROR_CODE_INTERNAL_ERROR ErrorCode = 1
	ERROR_CODE_DATA_INVALID   ErrorCode = 2
	ERROR_CODE_PHONE_INVALID  ErrorCode = 3
	ERROR_CODE_TOKEN_INVALID  ErrorCode = 4
	ERROR_CODE_DEVICE_INVALID ErrorCode = 5
)

var ERROR_CODE_name = map[ErrorCode]string{
	ERROR_CODE_NONE:           "ERROR_CODE_NONE",
	ERROR_CODE_INTERNAL_ERROR: "ERROR_CODE_INTERNAL_ERROR",
	ERROR_CODE_PHONE_INVALID:  "ERROR_CODE_PHONE_INVALID",
	ERROR_CODE_TOKEN_INVALID:  "ERROR_CODE_TOKEN_INVALID",
}

type JsonH map[string]interface{}
