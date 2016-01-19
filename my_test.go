package sequencenumber

import (
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm/v1"
	"testing"
)

func prepareOrm() (*orm.DataContext, error) {
	conn, e := dbox.NewConnection("mongo",
		&dbox.ConnectionInfo{"localhost:27123", "ectest", "", "", nil})
	e = conn.Connect()
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
	s, e := Get("General", true)
	if e != nil {
		t.Fatal(e)
	}
	s.Save()
}

func TestClaim(t *testing.T) {
	ctx, e := prepareOrm()
	if e != nil {
		t.Error(e.Error())
	}
	defer ctx.Close()

	Ctx = ctx
	s, e := Get("General", false)
	if e != nil {
		t.Error(e.Error())
	}
	i := s.LastNo + 1
	claimed, e := s.Claim()
	if e != nil {
		t.Error(e.Error())
	}
	if i != claimed {
		t.Errorf("Error, want %d got %d", i, claimed)
	}
}
