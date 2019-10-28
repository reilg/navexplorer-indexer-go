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
		for _, a := range vin.Addresses {
			if a == address {
				return true
			}
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
		for _, a := range i.Addresses {
			if a == address {
				value += i.Value
				valuesat += i.ValueSat
			}
		}
	}

	return
}
