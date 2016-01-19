package sequencenumber

import (
	"fmt"
	"github.com/eaciit/dbox"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/orm/v1"
	"strings"
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
	Id          string `json:"_id",bson:"_id"`
	Title       string
	LastNo      int
	ReuseNumber bool
	Format      string
}

type UsedSequence struct {
	orm.ModelBase
	Id         string `json:"_id",bson:"_id"`
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

func (s *Sequence) RecordID() interface{} {
	return s.Id
}

func (u *UsedSequence) PrepareID() interface{} {
	u.Id = fmt.Sprintf("%s_%d", u.SequenceId, u.No)
	return u.Id
}

func (u *UsedSequence) RecordID() interface{} {
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
		return nil, errorlib.Error(packageName, objSequence,
			"Get", "Context not yet initialized")
	}

	s := new(Sequence)
	e := Ctx.GetById(s, id)
	if e != nil && !strings.Contains(e.Error(), "Not found") {
		fmt.Printf("Error: %s Found: %v\n", e.Error(), strings.Contains(e.Error(), "Not found"))
		e = errorlib.Error(packageName, objSequence,
			"Get", e.Error())
	} else if e != nil {
		//fmt.Printf("Error: %s Found: %v\n", e.Error(), strings.Contains(e.Error(), "Not found"))
		if init {
			s.Id = id
			e = nil
		} else {
			e = errorlib.Error(packageName, objSequence,
				"Get", "Not found")
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
		return 0, errorlib.Error(packageName, objSequence,
			"Claim", "Context not yet initialized")
	}
	var latestNo int
	latest, e := Get(s.Id, true)
	if e != nil {
		return 0, errorlib.Error(packageName, objSequence, "Claim",
			"Unable to get latest number - "+e.Error())
	}
	latestNo = latest.LastNo + 1

	if s.ReuseNumber {
		used := new(UsedSequence)
		c, e := Ctx.Connection.NewQuery().From(s.TableName()).Order("no").Where(
			dbox.Eq("sequenceid", s.Id), dbox.Eq("status", NumberStatus_Available)).Cursor(nil)
		e = c.Fetch(used, 1, true)
		if e != nil {
			return 0, errorlib.Error(packageName, objSequence, "Claim",
				"Unable to get latest available - "+e.Error())
		}
		if used.No < latestNo {
			latestNo = used.No
		}
	}

	s.LastNo = latestNo
	ret := s.LastNo
	e = Ctx.Save(s)

	if e == nil && s.ReuseNumber {
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
