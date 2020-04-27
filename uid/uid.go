package uid

import (
	"encoding/json"
	"gitlab.livedev.shika2019.com/go/common/logging"
	"strconv"
)

type UserId struct {
	Low  uint64
	High uint64
}

func NewUserId(s string) UserId {
	if len(s) > 32 {
		logging.Log.Errorf("invalid userId[%s]", s)
		return UserId{}
	}

	var newOjb UserId
	var err error
	if len(s) < 16 {
		newOjb.Low, err = strconv.ParseUint(s, 16, 64)
		if err != nil {
			logging.Log.Errorf("invalid userId [%s]", s)
			return UserId{}
		}
	} else {
		newOjb.Low, err = strconv.ParseUint(s[len(s)-16:], 16, 64)
		if err != nil {
			logging.Log.Errorf("invalid userId [%s]", s)
			return UserId{}
		}
		newOjb.High, err = strconv.ParseUint(s[:len(s)-16], 16, 64)
		if err != nil {
			logging.Log.Errorf("invalid userId [%s]", s)
			return UserId{}
		}
	}
	return newOjb
}

func (uid UserId) IsNull() bool {
	return uid.High == 0 && uid.Low == 0
}

func (uid UserId) Equal(b UserId) bool {
	return uid.High == b.High && uid.Low == b.Low
}

func (uid UserId) ToString() string {
	if uid.High > 0 {
		return strconv.FormatUint(uid.High, 16) + strconv.FormatUint(uid.Low, 16)
	} else {
		return strconv.FormatUint(uid.Low, 16)
	}
}

func (uid UserId) MarshalJSON() ([]byte, error) {
	return json.Marshal(uid.ToString())
}

func (uid *UserId) UnmarshalJSON(b []byte) (err error) {
	s := ""
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}

	*uid = NewUserId(s)
	return
}
