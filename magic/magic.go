package magic

type Magic int

const (
	ILLEGAL Magic = iota
	FIRE_BALL
	WATER_WALL
	GRASS_GROWTH

	// length of magic
	magicLength
)

func (m Magic) String() string {
	switch m {
	case FIRE_BALL:
		return "fireBall"
	case WATER_WALL:
		return "waterWall"
	case GRASS_GROWTH:
		return "grassGrowth"
	}
	return "illegal"
}

func ToMagic(m string) Magic {
	switch m {
	case "fireBall":
		return FIRE_BALL
	case "waterWall":
		return WATER_WALL
	case "grassGrowth":
		return GRASS_GROWTH
	}
	return ILLEGAL
}

func (m Magic) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m *Magic) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		*m = ILLEGAL
		return nil
	}

	data = data[1 : len(data)-1]
	*m = ToMagic(string(data))
	return nil
}
