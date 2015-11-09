package sequencenumber

import (
	"fmt"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/orm"
	"time"
)

type SequenceStatus int

const (
	SequenceStatus_Used      = 1
	SequenceStatus_Reserved  = 10
	SequenceStatus_Available = 999
)

var packageName = "sequencenumber"
var objSequence = "Sequence"
var objUsedSequence = "UsedSequence"

type Sequence struct {
	orm.ModelBase
	Id      string `bson:_id,json:_id`
	Title   string
	LastNo  int
	LogUsed bool
	Format  string
}

type UsedSequence struct {
	orm.ModelBase
	Id         int `bson:_id,json:_id`
	SequenceId string
	No         int
	Used       time.Time
	Status     SequenceStatus
}

var Ctx *orm.DataContext

func Get(id string, init bool) (*Sequence, error) {
	if Ctx == nil {
		return nil, errorlib.Error(packageName, objSequence, "Get", "Context not yet initialized")
	}

	s := new(Sequence)
	b, e := Ctx.GetById(s, id)
	if e != nil {
		e = errorlib.Error(packageName, objSequence, "Get", e.Error())
	} else {
		if init {
			s.Id = id
		} else {
			e = errorlib.Error(packageName, objSequence, "Get", "Record not found")
		}
	}
	return s, e
}

func (s *Sequence) Claim() (int, error) {
	var e error
	if Ctx == nil {
		return 0, errorlib.Error(packageName, objSequence, "Claim", "Context not yet initialized")
	}
	ret := s.LastNo + 1
	s.LastNo = ret
	e = Ctx.Save(s)
	return ret, e
}

func (s *Sequence) ClaimString() (string, error) {
	i, e := s.Claim()
	if e != nil {
		return "", e
	}
	return fmt.Sprintf(s.Format, i), nil
}

func (s *Sequence) Save() error {
	if Ctx == nil {
		return errorlib.Error(packageName, objSequence, "Save", "Context not yet initialized")
	}

	e := Ctx.Save(s)
	if e != nil {
		return errorlib.Error(packageName, objSequence, "Save", e.Error())
	}
	return nil
}
