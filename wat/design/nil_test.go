package design_test

import "fmt"

type Explodes interface {
	Bang()
	Boom()
}

type Bomb struct{}

func (*Bomb) Bang() {}
func (Bomb) Boom()  {}

func ExampleInterfaceAndNilValue() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	var bomb *Bomb = nil
	var explodes Explodes = bomb
	if explodes != nil {
		explodes.Bang()
		explodes.Boom()
	}

	// Output:
	// value method github.com/kamilsk/go-research/wat/design_test.Bomb.Boom called using nil *Bomb pointer
}
