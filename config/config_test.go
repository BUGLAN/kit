package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewKitConfig(t *testing.T) {
	c := NewKitConfig()
	assert.NotNil(t, c)
}