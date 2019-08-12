package random

type DicePrng interface {
	RollDice(int, int, int) int
}

type Dice struct {
	dnum, dval, dmod int
}

func (d *Dice) roll(prng DicePrng) int {
	return prng.RollDice(d.dnum, d.dval, d.dmod)
}
