package search

type itemInfoSet struct {
	Data           []itemInfo
	shouldCopyNext bool
}

func newItemInfoSet() *itemInfoSet {
	return &itemInfoSet{nil, false}
}

func (s *itemInfoSet) IsEmpty() bool {
	return s.Data != nil && len(s.Data) == 0
}

func (s *itemInfoSet) Intersect(other []itemInfo) {
	if s.IsEmpty() {
		return
	}
	if other == nil {
		s.Data = []itemInfo{}
		return
	}
	if s.Data == nil {
		s.Data = other
		s.shouldCopyNext = true
	} else {
		cur := s.Data
		if s.shouldCopyNext {
			s.Data = make([]itemInfo, min(len(cur), len(other)))
		}
		di := 0
		ci := 0
		oi := 0
		for ci < len(cur) && oi < len(other) {
			if cur[ci].index == other[oi].index {
				index := cur[ci].index
				for oi < len(other) && other[oi].index == index && cur[ci].position >= other[oi].position {
					oi++
				}
				if oi < len(other) && other[oi].index == index && cur[ci].position < other[oi].position {
					s.Data[di] = other[oi]
					di++
				}
				for ci < len(cur) && cur[ci].index == index {
					ci++
				}
				for oi < len(other) && other[oi].index == index {
					oi++
				}
			} else if cur[ci].index < other[oi].index {
				ci++
			} else {
				oi++
			}
		}
		s.Data = s.Data[:di]
		s.shouldCopyNext = false
	}
}

func (s *itemInfoSet) Intersection(other []itemInfo) *itemInfoSet {
	result := newItemInfoSet()
	if s.Data != nil {
		result.Intersect(s.Data)
	}
	result.Intersect(other)
	return result
}
