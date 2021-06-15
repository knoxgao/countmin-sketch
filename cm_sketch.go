package countmin_sketch

import (
	"hash/fnv"
	"math"
	"math/rand"
)

type cmsRow struct {
	n    int
	A    []int32
	seed int64
}

func (c *cmsRow) Incr(v int64) {
	v = v ^ c.seed // random
	index := v % int64(c.n)
	c.A[index]++
}

func (c *cmsRow) Value(v int64) int32 {
	v = v ^ c.seed // random
	index := v % int64(c.n)
	return c.A[index]
}

func newCmsRow(n int) *cmsRow {
	seed := rand.Int63()
	return &cmsRow{
		n:    n,
		A:    make([]int32, n),
		seed: seed,
	}
}

type cms struct {
	m, n int // m rows, n columns
	rows []*cmsRow
}

func newCMS(w, d int) *cms {
	rows := make([]*cmsRow, d)
	for i := range rows {
		rows[i] = newCmsRow(w)
	}
	return &cms{
		m:    w,
		n:    d,
		rows: rows,
	}
}

func (c *cms) Incr(key string) {
	value := hash([]byte(key))
	for _, row := range c.rows {
		row.Incr(value)
	}
}

func (c *cms) Count(key string) int32 {
	value := hash([]byte(key))
	var ans int32 = math.MaxInt32
	for _, row := range c.rows {
		ans = min(ans, row.Value(value))
	}
	return ans
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func hash(key []byte) int64 {
	tmp := fnv.New64()
	tmp.Write(key) // nolint(errcheck)
	return int64(tmp.Sum64())
}
