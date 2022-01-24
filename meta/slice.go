package meta

type Slice struct {
	Id    uint64
	Size  uint32 //固定大小，get objectKey
	Pos   uint32 //chunk offset
	Off   uint32 //Slice offset
	Len   uint32
	Left  *Slice
	Right *Slice
}

type CompactedSlice struct {
	Skipped   int
	SliceList []*Slice
	Id        uint64
	Off       uint32
	Size      uint32
}

func newSlice(sliceId uint64, size, pos, off, len uint32) *Slice {
	if len == 0 {
		return nil
	}
	s := &Slice{}
	s.Id = sliceId
	s.Size = size
	s.Pos = pos
	s.Off = off
	s.Len = len
	s.Left = nil
	s.Right = nil
	return s
}

func (s *Slice) cut(pos uint32) (left, right *Slice) {
	if s == nil {
		return nil, nil
	}
	if pos <= s.Pos {
		if s.Left == nil {
			s.Left = newSlice(0, s.Pos-pos, pos, pos, s.Pos-pos)
		}
		left, s.Left = s.Left.cut(pos)
		return left, s
	} else if pos < s.Pos+s.Len {
		l := pos - s.Pos
		right = newSlice(s.Id, s.Size, pos, s.Off+l, s.Len-l)
		right.Right = s.Right
		s.Len = l
		s.Right = nil
		return s, right
	} else {
		if s.Right == nil {
			s.Right = newSlice(0, pos-s.Pos-s.Len, s.Pos+s.Len, s.Pos+s.Len, pos-s.Pos-s.Len)
		}
		s.Right, right = s.Right.cut(pos)
		return s, right
	}
}

func (s *Slice) visit(f func(*Slice)) {
	if s == nil {
		return
	}
	s.Left.visit(f)
	right := s.Right
	f(s) // s could be freed
	right.visit(f)
}

// buildSlice reutrn a slice tree root node and ordered slice
func buildSlice(root *Slice, sliceList []*Slice) (*Slice, []*Slice) {
	for _, s := range sliceList {
		if root != nil {
			var right *Slice
			s.Left, right = root.cut(s.Pos)
			_, s.Right = right.cut(s.Pos + s.Len)
		}
		root = s
	}
	var pos uint32
	var slices []*Slice
	root.visit(func(s *Slice) {
		if s.Pos > pos {
			slices = append(slices, &Slice{Id: 0, Size: s.Pos - pos, Pos: pos, Off: pos, Len: s.Pos - pos}) //hole
			pos = s.Pos
		}
		slices = append(slices, &Slice{Id: s.Id, Size: s.Size, Pos: s.Pos, Off: s.Off, Len: s.Len})
		pos += s.Len
	})
	return root, slices
}
