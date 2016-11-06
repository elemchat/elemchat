package attr

type Value int
type Type int

const (
	ILLEGAL Type = iota
	BLOOD
)

type Diff map[Type]Value
type Attr struct {
	Blood Value
}

func (after Attr) Diff(before Attr) Diff {
	diff := make(Diff)
	if diffBlood := after.Blood - before.Blood; diffBlood != 0 {
		diff[BLOOD] = diffBlood
	}

	return diff
}

func (attr *Attr) Effect(diff Diff) {
	for t, v := range diff {
		switch t {
		case BLOOD:
			attr.Blood += v
		}
	}
}
