package region

// NOTICE: Untestd!!!
import (
	"fmt"
)

// Inspired by limetext's region code, but add an direction for a region

// Region from point From to To.
// From must be less than To, real direction is decide by the Dir field
type Region struct {
	From, To int
	Dir      Direction
}

// New a Region
func NewRegion(from, to int) Region {
	dir := POSITIVE
	if from > to {
		from, to, dir = to, from, REVERSE
	}
	return Region{from, to, dir}
}

// String return region's real (from,to)
func (r Region) String() string {
	return fmt.Sprintf("(%d, %d)", r.From, r.To)
}

// RealFrom return region's real from
func (r Region) RealFrom() int {
	return MinByDir(r.From, r.To, r.Dir)
}

// RealTo return real's real to
func (r Region) RealTo() int {
	return MaxByDir(r.From, r.To, r.Dir)
}

// Begin return start of Region
func (r Region) Begin() int {
	return r.From
}

// End return end of Region
func (r Region) End() int {
	return r.To
}

// Contains returns whether the region contains the given point or not.
func (r Region) Contains(point int) bool {
	return point >= r.From && point <= r.To
}

// MidIn returns whether the point is in the reign and is't the begin and end
func (r Region) MidIn(point int) bool {
	return point > r.From && point < r.To
}

// Cover returns whether the region fully covers the argument region
func (r Region) Cover(r2 Region) bool {
	return r.From <= r2.From && r2.To <= r.To
}

// Empty returns whether or not the region is empty
func (r Region) Empty() bool {
	return r.From == r.To
}

// Returns the size of the region
func (r Region) Size() int {
	return r.To - r.From
}

// Combine returns a region covering both regions, dir is same as r
func (r Region) Combine(r2 Region) Region {
	return Region{Min(r.From, r2.From), Max(r.To, r2.To), r.Dir}
}

// Clip return the cliped against another region
// if r is inside r2 or r2 inside r, return r,
// else return r that remove intesect part
func (r Region) Clip(r2 Region) Region {
	var ret Region = r
	if r2.Cover(r) {
		return r
	}
	if r2.Contains(ret.From) {
		ret.From = r2.To
	} else if r2.Contains(ret.To) {
		ret.To = r2.From
	}
	return ret
}

// Cuts remove the intersect part
func (r Region) Cut(r2 Region) (ret []Region) {
	if r.MidIn(r2.From) {
		ret = append(ret, Region{r.From, r2.From, r.Dir})
	}
	if r.MidIn(r2.To) {
		ret = append(ret, Region{r2.To, r.To, r.Dir})
	}
	if len(ret) == 0 && r2.Size() > 0 && !r2.Cover(r) {
		ret = append(ret, r)
	}
	return
}

// Intersects check whether the two regions intersects
func (r Region) Intersects(r2 Region) bool {
	return r.Intersection(r2).Size() > 0
}

// Intersection returns the region that is the intersection of the two
func (r Region) Intersection(r2 Region) (ret Region) {
	from := Max(r.From, r2.From)
	to := Min(r.To, r2.To)
	if from < to {
		ret = Region{from, to, r.Dir}
	}
	return
}

// Adjust apply the change in the position with delta(the change) to the region if need
// if point is after changed position, it must change also
// else if point is before the changed position, if the delta do not affect the point, needn't change point
// else apply the change also
func (r *Region) Adjust(position, delta int) {
	if r.From >= position {
		r.From += delta
	} else if diff := position + delta - r.From; diff < 0 {
		r.From += diff
	}
	if r.To >= position {
		r.To += delta
	} else if diff := position + delta - r.To; diff < 0 {
		r.To += diff
	}
}
