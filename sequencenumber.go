package sequencenumber

import (
	"github.com/eaciit/orm"
	"time"
)

type SequenceStatus int

const (
	SequenceStatus_Used      = 1
	SequenceStatus_Reserved  = 10
	SequenceStatus_Available = 999
)

type Sequence struct {
	orm.ModelBase
	Id      string `bson:_id,json:_id`
	Title   string
	Next    int
	LogUsed bool
}

type UsedSequence struct {
	Id         int `bson:_id,json:_id`
	SequenceId string
	No         int
	Used       time.Time
	Status     SequenceStatus
}
