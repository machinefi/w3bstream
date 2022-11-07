package id

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/base/types/snowflake_id"
)

func NewSFIDGenerator() (SFIDGenerator, error) {
	return NewSFIDGeneratorWithWorkerID(snowflake_id.WorkerIDFromLocalIP())
}

func NewSFIDGeneratorWithWorkerID(wid uint32) (SFIDGenerator, error) {
	sf, err := sff.NewSnowflake(wid)
	if err != nil {
		return nil, err
	}
	return &SFIDGeneratorImpl{sf}, nil
}

func MustNewSFIDGenerator() SFIDGenerator {
	g, err := NewSFIDGenerator()
	if err != nil {
		panic(err)
	}
	return g
}

func MustNewSFIDGeneratorWithWorkerID(wid uint32) SFIDGenerator {
	g, err := NewSFIDGeneratorWithWorkerID(wid)
	if err != nil {
		panic(err)
	}
	return g
}

type SFIDGeneratorImpl struct{ *snowflake_id.Snowflake }

func (sfg *SFIDGeneratorImpl) MustGenSFID() types.SFID {
	id, err := sfg.ID()
	if err != nil {
		panic(err)
	}
	return types.SFID(id)
}

func (sfg *SFIDGeneratorImpl) MustGenSFIDs(n int) types.SFIDs {
	var ids = make(types.SFIDs, 0, n)
	for i := 0; i < n; i++ {
		id, err := sfg.ID()
		if err != nil {
			panic(err)
		}
		ids = append(ids, types.SFID(id))
	}
	return ids
}
