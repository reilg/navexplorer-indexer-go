package explorer

type Vins []Vin

func (vins *Vins) Empty() bool {
	return len(*vins) == 0
}

func (vins *Vins) First() Vin {
	return (*vins)[0]
}

func (vins *Vins) HasAddress(address string) bool {
	for _, vin := range *vins {
		if vin.Address == address {
			return true
		}
	}

	return false
}

func (vins *Vins) GetAmount() uint64 {
	var amount uint64 = 0
	for _, i := range *vins {
		amount += i.ValueSat
	}
	return amount
}

func (vins *Vins) GetAmountByAddress(address string) (value float64, valuesat uint64) {
	for _, i := range *vins {
		if i.Address == address {
			value += i.Value
			valuesat += i.ValueSat
		}
	}

	return
}
