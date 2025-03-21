package submission

import "slices"

type SortedList struct {
	elements []*Order
	comp     func(*Order, *Order) bool
}

// NewSortedList creates a new SortedList
// isAsc determines if the list should be sorted in ascending (true) or descending (false) order
func NewSortedList(isAsc bool) *SortedList {
	comp := func(o, co *Order) bool {
		return o.Price <= co.Price
	}
	if isAsc {
		comp = func(o, co *Order) bool {
			return o.Price >= co.Price
		}
	}
	return &SortedList{
		elements: make([]*Order, 0),
		comp:     comp,
	}
}

// Push adds a new Order to the list in the correct sorted position
func (sl *SortedList) Push(order *Order) {
	if sl.IsEmpty() {
		sl.elements = append(sl.elements, order)
		return
	}

	insertPos := 0
	for i := len(sl.elements) - 1; i >= 0; i-- {
		current := sl.elements[i]
		if sl.comp(order, current) {
			insertPos = i + 1
			break
		}
	}

	sl.insertAtPos(order, insertPos)
}

func (sl *SortedList) insertAtPos(order *Order, pos int) {
	if pos == len(sl.elements) {
		sl.elements = append(sl.elements, order)
	} else {
		sl.elements = append(sl.elements[:pos+1], sl.elements[pos:]...)
		sl.elements[pos] = order
	}
}

func (sl *SortedList) Peek() *Order {
	if len(sl.elements) == 0 {
		return nil
	}
	return sl.elements[0]
}

func (sl *SortedList) Pop() *Order {
	if len(sl.elements) == 0 {
		return nil
	}

	order := sl.elements[0]
	sl.elements = sl.elements[1:]
	return order
}

func (sl *SortedList) HasID(id uint32) bool {
	for _, order := range sl.elements {
		if order.ID == id {
			return true
		}
	}
	return false
}

// Remove removes an Order with the given ID from the list
// Returns true if an element was removed, false otherwise
func (sl *SortedList) Remove(id uint32) bool {
	for i, order := range sl.elements {
		if order.ID == id {
			sl.elements = slices.Delete(sl.elements, i, i+1)
			return true
		}
	}
	return false
}

func (sl *SortedList) Size() int {
	return len(sl.elements)
}

func (sl *SortedList) IsEmpty() bool {
	return len(sl.elements) == 0
}
