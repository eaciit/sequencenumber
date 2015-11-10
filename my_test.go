package sequencenumber

import (
	"github.com/eaciit/database/mongodb"
	"github.com/eaciit/orm"
	"testing"
)

func prepareOrm() (*orm.DataContext, error) {
	conn := mongodb.NewConnection("localhost:27123", "", "", "ectest")
	e := conn.Connect()
	if e != nil {
		return nil, e
	}

	ctx := orm.New(conn)
	return ctx, nil
}

func TestCreate(t *testing.T) {
	ctx, e := prepareOrm()
	if e != nil {
		t.Error(e.Error())
	}
	defer ctx.Close()

	Ctx = ctx
	s, _ := Get("General", true)
	s.Save()
}

func TestClaim(t *testing.T) {
	ctx, e := prepareOrm()
	if e != nil {
		t.Error(e.Error())
	}
	defer ctx.Close()

	Ctx = ctx
	s, e := Get("General", true)
	if e != nil {
		t.Error(e.Error())
	}
	s.UseLog = true
	i := s.LastNo + 1
	claimed, e := s.Claim()
	if e != nil {
		t.Error(e.Error())
	}
	if i != claimed {
		t.Errorf("Error, want %d got %d", i, claimed)
	}
}
