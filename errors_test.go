package orm_test

import (
	"errors"
	"testing"

	"github.com/novakit/orm"
)

func TestErrorsCanBeUsedOutsideOrm(t *testing.T) {
	errs := []error{errors.New("First"), errors.New("Second")}

	gErrs := orm.Errors(errs)
	gErrs = gErrs.Add(errors.New("Third"))
	gErrs = gErrs.Add(gErrs)

	if gErrs.Error() != "First; Second; Third" {
		t.Fatalf("Gave wrong error, got %s", gErrs.Error())
	}
}
