package infra

func IfThenElse[T any](p bool, then, els T) T {
	if p {
		return then
	} else {
		return els
	}
}
