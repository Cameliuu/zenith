package utils

type Vector3 struct {
	X float32
	Y float32
	Z float32
}
type Vector2 struct {
	X, Y float32
}

func DistanceBetweenTwoPoints(a, b Vector2) float32 {
	//WITHOUT SQRT CAUSE IT'S FASTERRR
	return float32(((b.X - a.X) * (b.X - a.X)) + ((b.Y - a.Y) * (b.Y - a.Y)))
}

var modelToTeams = map[string]uint{
	"gign":     2,
	"gsg":      2,
	"sas":      2,
	"urban":    2,
	"terror":   1,
	"leet":     1,
	"arctic":   1,
	"guerilla": 1,
}

func IsEntityInTheSameTeam(model string, playerTeam int32) bool {
	team, _ := modelToTeams[model]
	return team == uint(playerTeam)
}
