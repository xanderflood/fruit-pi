package tools

import (
	uuid "github.com/nu7hatch/gouuid"
)

//UUIDer UUID generator
//go:generate counterfeiter . UUIDer
type UUIDer interface {
	UUID() string
}

//UUIDAgent generates random UUIDs
type UUIDAgent struct{}

func (*UUIDAgent) UUID() string {
	count := 0
	for {
		uuid, err := uuid.NewV4()
		if err == nil {
			return uuid.String()
		}
		if count >= 10 {
			panic("can't generate UUID")
		}
	}
}

//IsUUID checks whether the string is a valid v4UUID
func IsUUID(s string) bool {
	_, err := uuid.ParseHex(s)
	return err == nil
}
