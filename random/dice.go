package random

type DicePrng interface {
	RollDice(int, int, int) int
}

type Dice struct {
	dnum, dval, dmod int
}

func NewDice(dnum, dval, dmod int) *Dice {
	return &Dice{dnum: dnum, dval: dval, dmod: dmod}
}

func (d *Dice) Roll(prng DicePrng) int {
	return prng.RollDice(d.dnum, d.dval, d.dmod)
}
