package spatial

type ConvertFunc func(Point) Point

type Projectable interface {
	Project(ConvertFunc)
	Copy() Projectable
}
