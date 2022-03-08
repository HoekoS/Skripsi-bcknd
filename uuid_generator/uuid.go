package uuid_generator

import (
	"github.com/google/uuid"
)

func UuidGenerate() uuid.UUID {
	uuidWithHyphen := uuid.New()
	return uuidWithHyphen
	// uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	// fmt.Println(uuid)
}
