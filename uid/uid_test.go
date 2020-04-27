package uid

import (
	"encoding/json"
	"testing"
)

func TestUserId_All(t *testing.T) {
	expect := "5d887151857aba0006168f8f"

	a := NewUserId(expect)
	b := NewUserId(expect)
	c := NewUserId("9d887151857ae30006168f8f")
	if !a.Equal(b) {
		t.Fatalf("%v!=%v", a, b)
	}

	if c.Equal(b) {
		t.Fatalf("%v==%v", c, b)
	}

	if a.ToString() != expect {
		t.Fatalf("expect %v, got %v", expect, a.ToString())
	}
}

func TestUserId_ToString(t *testing.T) {
	expect := "30006168f8f"
	a := NewUserId(expect)
	if a.ToString() != expect {
		t.Fatalf("expect %v, got %v", expect, a.ToString())
	}
}

type TestA struct {
	Uid  UserId `json:"user_id"`
	Name string `json:"name"`
}

func TestUserId_MarshalJSON(t *testing.T) {
	a := TestA{
		Uid:  NewUserId("9d887151857ae30006168f8f"),
		Name: "Lilei",
	}

	body, _ := json.Marshal(a)
	t.Logf("%v", string(body))
}

func TestUserId_UnmarshalJSON(t *testing.T) {

	j := `{
        "user_id": "9d887151857ae30006168f8f",
        "name": "lilei"
    }`

	var obj TestA
	err := json.Unmarshal([]byte(j), &obj)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(obj)

	if obj.Name != "lilei" {
		t.Fatalf("expect %v, got %v", "lilei", obj.Name)
	}

	if obj.Uid.ToString() != "9d887151857ae30006168f8f" {
		t.Fatalf("expect %v, got %v", "9d887151857ae30006168f8f", obj.Uid.ToString())
	}

}

func TestUserId_MarshalJSON_and_UnmarshalJSON(t *testing.T) {

	a := TestA{
		Uid:  NewUserId("9d887151857ae30006168f8f"),
		Name: "Lilei",
	}

	var obj TestA
	body, _ := json.Marshal(a)
	err := json.Unmarshal(body, &obj)
	if err != nil {
		t.Fatal(err)
	}

	if !(obj.Uid.Equal(a.Uid) && obj.Name == a.Name) {
		t.Fatalf("expect %v, got %v", a, obj)
	}
}

func TestUserId_ToString2(t *testing.T) {
	t.Logf("%v", NewUserId("9d887151857ae30006168f8f"))
}
