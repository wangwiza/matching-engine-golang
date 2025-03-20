package submission

type SortedList struct {
	elements []*Order
	isAsc    bool
	counter  int
}

// NewSortedList creates a new SortedList
// isAsc determines if the list should be sorted in ascending (true) or descending (false) order
func NewSortedList(isAsc bool) *SortedList {
	return &SortedList{
		elements: make([]*Order, 0),
		isAsc:    isAsc,
		counter:  0,
	}
}

// Push adds a new Order to the list in the correct sorted position
func (sl *SortedList) Push(order *Order) {
	sl.counter++

	// For empty list, just append
	if len(sl.elements) == 0 {
		sl.elements = append(sl.elements, order)
		return
	}

	// Find the insertion position using insertion sort approach
	insertPos := 0

	if sl.isAsc {
		// Ascending order: lower prices at front, higher prices at back
		for i := len(sl.elements) - 1; i >= 0; i-- {
			current := sl.elements[i]

			// If new order's price is greater than or equal to current element's price
			if order.Price >= current.Price {
				// If prices are equal, newer insertion goes after existing elements
				if order.Price == current.Price {
					insertPos = i + 1
				} else {
					insertPos = i + 1
				}
				break
			}

			// If we've checked all elements and price is lower than all,
			// insert at the beginning
			if i == 0 {
				insertPos = 0
			}
		}
	} else {
		// Descending order: higher prices at front, lower prices at back
		for i := len(sl.elements) - 1; i >= 0; i-- {
			current := sl.elements[i]

			// If new order's price is less than or equal to current element's price
			if order.Price <= current.Price {
				// If prices are equal, newer insertion goes after existing elements
				if order.Price == current.Price {
					insertPos = i + 1
				} else {
					insertPos = i + 1
				}
				break
			}

			// If we've checked all elements and price is higher than all,
			// insert at the beginning
			if i == 0 {
				insertPos = 0
			}
		}
	}

	// Insert the element at the determined position
	if insertPos == len(sl.elements) {
		sl.elements = append(sl.elements, order)
	} else {
		sl.elements = append(sl.elements[:insertPos+1], sl.elements[insertPos:]...)
		sl.elements[insertPos] = order
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
			sl.elements = append(sl.elements[:i], sl.elements[i+1:]...)
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
