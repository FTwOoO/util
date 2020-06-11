package logging

import (
	"testing"
)

type People struct {
	Name     string `json:"name_title"`
	Age      int    `json:"age_size, ex"`
	Nickname string `json:"nickname,omitempty"`
}

func TestStructToMap(t *testing.T) {
	student := People{"jqw", 18, ""}
	data := structToMap(student)
	t.Log(data)

	data = structToMap(&student)
	t.Log(data)

	student.Nickname = "hahaha"
	data = structToMap(student)
	t.Log(data)

	var cases []interface{} = []interface{}{
		student,
		&student,
		map[string]interface{}{"a": 1},
		map[string]string{"a": "2"},
	}

	for _, C := range cases {
		if _, ok := C.(map[string]interface{}); ok {
			t.Logf("%v is map[string]interface{}", C)

		}
	}

}
