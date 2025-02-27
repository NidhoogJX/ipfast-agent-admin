package bloomFilter

import (
	"github.com/bits-and-blooms/bloom/v3"
)

type ChannelIPFilter struct {
	filter *bloom.BloomFilter
}

func NewChannelIPFilter(n uint, p float64) *ChannelIPFilter {
	return &ChannelIPFilter{
		filter: bloom.NewWithEstimates(n, p),
	}
}

func (c *ChannelIPFilter) Add(id string) {
	c.filter.Add([]byte(id))
}

func (c *ChannelIPFilter) Test(id string) bool {
	return c.filter.Test([]byte(id))
}
