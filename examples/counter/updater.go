package main

import (
  "context"
  "errors"
  "github.com/nheyn/go-redux/store"
)

// Updater
type counter int

func (c counter) Update(_ context.Context, action interface{}) (store.Updater, error) {
  switch amount := action.(type) {
  case incurment:
    if amount < 0 {
      return nil, errors.New("Do not use negitive incurment actions, use decurment instead")
    }

    return c + counter(amount), nil
  case decurment:
    if amount < 0 {
      return nil, errors.New("Do not use negitive decurment actions, use incurment instead")
    }

    return c - counter(amount), nil
  default:
    return c, nil
  }
}

// Actions
type incurment int

type decurment int
