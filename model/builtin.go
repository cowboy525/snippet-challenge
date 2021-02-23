package model

import (
	"time"
)

// NewBool create new bool pointer
func NewBool(b bool) *bool { return &b }

// NewInt create new int pointer
func NewInt(n int) *int { return &n }

// NewInt64 create new int64 pointer
func NewInt64(n int64) *int64 { return &n }

// NewUint64 create new uint64 pointer
func NewUint64(u uint64) *uint64 { return &u }

// NewString create new string pointer
func NewString(s string) *string { return &s }

// NewTime create new time.Time pointer
func NewTime(t time.Time) *time.Time { return &t }
