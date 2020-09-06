package common

import (
	"github.com/google/uuid"
)

//------------------------------------------------------------------------------

// UUID creates a universal unique id
func UUID() string {
	return uuid.New().String()
}

//------------------------------------------------------------------------------
