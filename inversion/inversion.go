package inversion

type Inversion interface {
	Count([]int64) uint64
}

var InvDivCon Inversion = &invDivCon{}

type invDivCon struct{}

func (idc *invDivCon) Count(a []int64) uint64 {
	if len(a) <= 1 {
		return 0
	}
	n := len(a) >> 1
	left := a[:n]
	right := a[n:]
	x := idc.Count(left)
	y := idc.Count(right)
	z, b := idc.mergeCount(left, right)
	copy(a, b)
	return x + y + z
}

func (idc *invDivCon) mergeCount(li, ri []int64) (r uint64, s []int64) {
	for i, j := 0, 0; i < len(li) || j < len(ri); {
		if i >= len(li) {
			s = append(s, ri[j:]...)
			break
		}
		if j >= len(ri) {
			s = append(s, li[i:]...)
			break
		}

		switch {
		case li[i] < ri[j]:
			s = append(s, li[i])
			i++
		case li[i] == ri[j]:
			s = append(s, li[i])
			i++
			s = append(s, ri[j])
			j++
			r += uint64(len(li) - i)
		default: // li[i] > ri [j]
			s = append(s, ri[j])
			j++
			r += uint64(len(li) - i)
		}
	}
	return
}
