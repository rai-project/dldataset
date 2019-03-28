package object_detection

func (s *StringIntLabelMap) Get(idx int) *StringIntLabelMapItem {
	items := s.GetItem()
	if idx < 0 {
		return nil
	}
	if idx >= len(items) {
		return nil
	}
	return items[idx]
}

func (s *StringIntLabelMap) GetName(idx int) string {
	item := s.Get(idx)
	return item.GetName()
}

func (s *StringIntLabelMap) GetId(idx int) int32 {
	item := s.Get(idx)
	return item.GetId()
}

func (s *StringIntLabelMap) GetDisplayName(idx int) string {
	item := s.Get(idx)
	return item.GetDisplayName()
}
