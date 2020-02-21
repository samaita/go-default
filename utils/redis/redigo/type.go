package redigo

type Z struct {
	Score  float64
	Member interface{}
}

type ZRangeByScore struct {
	Min, Max      string
	Offset, Count int64
}
