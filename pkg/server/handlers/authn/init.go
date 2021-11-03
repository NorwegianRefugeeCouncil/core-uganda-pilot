package authn

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(&time.Time{})
}
