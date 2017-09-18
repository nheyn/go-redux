package main

import "github.com/nheyn/go-redux/store"

// Selector
type countInfo struct {
  value int
  isPositive bool
}

func (c *countInfo) SelectFrom(s *store.State) {
  counterData := (*s)["COUNTER_STATE"].(counter)

  c.value = int(counterData)
  c.isPositive = counterData > 0
}

func (c *countInfo) isNegitive() bool {
  return c.value < 0
}
