package magic

type Magic int

const (
	ILLEGAL Magic = iota
	FIRE_BALL
	WATER_WALL
	GRASS_GROWTH
	size
)

func Size() int {
	return int(size)
}
