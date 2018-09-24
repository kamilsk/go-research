package design_test

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name string
}

func (u User) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"name":%q}`, u.Name)), nil
}

type UserBio struct {
	Gender string
}

func (b UserBio) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"gender":%q}`, b.Gender)), nil
}

func ExampleCollision() {
	marshal := func(v interface{}) []byte {
		b, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return b
	}

	fmt.Printf(`expected: {"name":"John"}; obtained: %s`+"\n", marshal(User{Name: "John"}))
	fmt.Printf(`expected: {"gender":"male"}; obtained: %s`+"\n", marshal(UserBio{Gender: "male"}))
	fmt.Printf(`expected: {"name":"John","Age":30}; obtained: %s`+"\n", marshal(struct {
		User
		Age uint
	}{User{Name: "John"}, 30}))
	fmt.Printf(`expected: {"name":"John","gender":"male"}; obtained: %s`+"\n", marshal(struct {
		User
		UserBio
	}{User{Name: "John"}, UserBio{Gender: "male"}}))

	// Output:
	// expected: {"name":"John"}; obtained: {"name":"John"}
	// expected: {"gender":"male"}; obtained: {"gender":"male"}
	// expected: {"name":"John","Age":30}; obtained: {"name":"John"}
	// expected: {"name":"John","gender":"male"}; obtained: {"Name":"John","Gender":"male"}
}
