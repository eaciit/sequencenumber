package sequencenumber

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/orm"
	"time"
)

type NumberStatus int

const (
	NumberStatus_Used      = 1
	NumberStatus_Reserved  = 10
	NumberStatus_Available = 999
)

var packageName = "sequencenumber"
var objSequence = "Sequence"
var objUsedSequence = "UsedSequence"

type Sequence struct {
	orm.ModelBase
	Id     string `bson:"_id"`
	Title  string
	LastNo int
	UseLog bool
	Format string
}

type UsedSequence struct {
	orm.ModelBase
	Id         string `bson:"_id"`
	SequenceId string
	No         int
	Used       time.Time
	Status     NumberStatus
}

func NewSequence(id string) *Sequence {
	s := new(Sequence)
	s.Id = id
	return s
}

func NewUsedSequence(sequenceid string, no int, status NumberStatus) *UsedSequence {
	us := new(UsedSequence)
	us.SequenceId = sequenceid
	us.No = no
	us.Status = status
	return us
}

func (s *Sequence) RecordId() interface{} {
	return s.Id
}

func (u *UsedSequence) PrepareId() interface{} {
	u.Id = fmt.Sprintf("%s_%d", u.SequenceId, u.No)
	return u.Id
}

func (s *Sequence) TableName() string {
	return "sequences"
}

func (u *UsedSequence) TableName() string {
	return "usedsequences"
}

func (u *UsedSequence) PreSave() error {
	if u.Status == NumberStatus_Used {
		u.Used = time.Now()
	}
	return nil
}

var Ctx *orm.DataContext

func Get(id string, init bool) (*Sequence, error) {
	if Ctx == nil {
		return nil, errorlib.Error(packageName, objSequence, "Get", "Context not yet initialized")
	}

	s := new(Sequence)
	_, e := Ctx.GetById(s, id)
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

func (s *Sequence) ChangeNumberStatus(n int, status NumberStatus) error {
	if Ctx == nil {
		return errorlib.Error(packageName, objSequence, "ChangeNumberStatus", "Context not yet initialized")
	}
	used := NewUsedSequence(s.Id, n, status)
	e := Ctx.Save(used)
	if e != nil {
		return errorlib.Error(packageName, objSequence, "ChangeNumberStatus", e.Error())
	}
	return nil
}

func (s *Sequence) Claim() (int, error) {
	var e error
	if Ctx == nil {
		return 0, errorlib.Error(packageName, objSequence, "Claim", "Context not yet initialized")
	}
	var latestNo int
	latest, e := Get(s.Id, true)
	if e != nil {
		return 0, errorlib.Error(packageName, objSequence, "Claim",
			"Unable to get latest number - "+e.Error())
	}
	latestNo = latest.LastNo + 1

	if s.UseLog {
		used := new(UsedSequence)
		c := Ctx.Connection.Query().From(s.TableName()).OrderBy("no").Where(
			base.Eq("sequenceid", s.Id), base.Eq("status", NumberStatus_Available)).Cursor(nil)
		found, e := c.FetchClose(used)
		if e != nil {
			return 0, errorlib.Error(packageName, objSequence, "Claim",
				"Unable to get latest available - "+e.Error())
		}
		if found && used.No < latestNo {
			latestNo = used.No
		}
	}

	s.LastNo = latestNo
	ret := s.LastNo
	e = Ctx.Save(s)

	if e == nil && s.UseLog {
		e = s.ChangeNumberStatus(ret, NumberStatus_Used)
	}
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
