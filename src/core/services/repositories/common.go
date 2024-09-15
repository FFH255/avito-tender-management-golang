package repositories

// Limit Максимальное кол-во возвращаемых объектов
type Limit uint

func NewLimit(n *int) Limit {
	if n == nil {
		return 5
	}
	if *n < 0 {
		return 5
	}
	return Limit(*n)
}

// Offset Кол-во объектов, которое должно быть пропущено с начала
type Offset uint

func NewOffset(n *int) Offset {
	if n == nil {
		return 0
	}
	if *n < 0 {
		return 0
	}
	return Offset(*n)
}
