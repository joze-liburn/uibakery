package cryptlex

import (
	"errors"
	"fmt"
	"time"
)

type (
	OptCmpString struct {
		cmp     string
		operand string
	}

	OptCmpTime struct {
		cmp     string
		operand time.Time
	}

	CryptlexOpts struct {
		page      int          // The page number.
		limit     int          // The number of records per page. Must be a number between 1 and 100.
		sort      string       // The sort string to sort the returned records e.g. "-createdAt"
		name      OptCmpString // Name of the reseller
		email     OptCmpString // Notification email address of the reseller.
		search    OptCmpString // Search string.
		id        OptCmpString // ID of the resource.
		createdAt OptCmpTime   // Date of creation
		updatedAt OptCmpTime   // Date of last update
	}

	CallOptions func(*CryptlexOpts) error
)

var (
	errCallOption = errors.New("call option error")
)

func (cs *OptCmpString) String() string {
	return fmt.Sprintf("%s%%20%s", cs.cmp, cs.operand)
}

func (ct *OptCmpTime) String() string {
	return fmt.Sprintf("%s%%20%s", ct.cmp, ct.operand.Format("2006-01-02T15:04:05Z"))
}

func WithPage(page int) CallOptions {
	return func(co *CryptlexOpts) error {
		if page <= 0 {
			return fmt.Errorf("%w: invalid page %d (1..2147483647)", errCallOption, page)
		}
		co.page = page
		return nil
	}
}

func WithPageSize(limit int) CallOptions {
	return func(co *CryptlexOpts) error {
		if limit <= 0 || limit > 100 {
			return fmt.Errorf("%w: invalid page size %d (1..100)", errCallOption, limit)
		}
		co.limit = limit
		return nil
	}
}

func WithSort(on string) CallOptions {
	return func(co *CryptlexOpts) error {
		co.sort = on
		return nil
	}
}

func WithId(cmp string, id string) CallOptions {
	return func(co *CryptlexOpts) error {
		if len(cmp) == 0 || len(id) == 0 {
			return fmt.Errorf("%w: invalid operation %q or target id %q", errCallOption, cmp, id)
		}
		co.id = OptCmpString{cmp, id}
		return nil
	}
}

func WithName(cmp string, name string) CallOptions {
	return func(co *CryptlexOpts) error {
		if len(cmp) == 0 || len(name) == 0 {
			return fmt.Errorf("%w: invalid operation %q or target name %q", errCallOption, cmp, name)
		}
		co.name = OptCmpString{cmp, name}
		return nil
	}
}

func WithEmail(cmp string, email string) CallOptions {
	return func(co *CryptlexOpts) error {
		if len(cmp) == 0 || len(email) == 0 {
			return fmt.Errorf("%w: invalid operation %q or target email %q", errCallOption, cmp, email)
		}
		co.email = OptCmpString{cmp, email}
		return nil
	}
}

func WithSearch(cmp string, term string) CallOptions {
	return func(co *CryptlexOpts) error {
		if len(cmp) == 0 || len(term) == 0 {
			return fmt.Errorf("%w: invalid operation %q or search string %q", errCallOption, cmp, term)
		}
		co.search = OptCmpString{cmp, term}
		return nil
	}
}

func WithCreatedAt(cmp string, created time.Time) CallOptions {
	return func(co *CryptlexOpts) error {
		co.createdAt = OptCmpTime{cmp, created}
		return nil
	}
}

func WithUpdatedAt(cmp string, updated time.Time) CallOptions {
	return func(co *CryptlexOpts) error {
		co.updatedAt = OptCmpTime{cmp, updated}
		return nil
	}
}

func (co *CryptlexOpts) ToSearchParam() string {
	ret := ""
	sep := "?"
	if co.page > 0 {
		ret = fmt.Sprintf("%s%spage%%20%d", ret, sep, co.page)
		sep = "&"
	}
	if co.limit > 0 {
		ret = fmt.Sprintf("%s%slimit%%20%d", ret, sep, co.limit)
		sep = "&"
	}
	if len(co.sort) > 0 {
		ret = fmt.Sprintf("%s%ssort%%20%s", ret, sep, co.sort)
		sep = "&"
	}
	if len(co.name.cmp) > 0 {
		ret = fmt.Sprintf("%s%sname%%20%s", ret, sep, co.name.String())
		sep = "&"
	}
	if len(co.email.cmp) > 0 {
		ret = fmt.Sprintf("%s%semail%%20%s", ret, sep, co.email.String())
		sep = "&"
	}
	if len(co.search.cmp) > 0 {
		ret = fmt.Sprintf("%s%ssearch%%20%s", ret, sep, co.search.String())
		sep = "&"
	}
	if len(co.id.cmp) > 0 {
		ret = fmt.Sprintf("%s%sid%%20%s", ret, sep, co.id.String())
		sep = "&"
	}
	if len(co.createdAt.cmp) > 0 {
		ret = fmt.Sprintf("%s%screatedAt%%20%s", ret, sep, co.createdAt.String())
		sep = "&"
	}
	if len(co.updatedAt.cmp) > 0 {
		ret = fmt.Sprintf("%s%supdatedAt%%20%s", ret, sep, co.updatedAt.String())
		sep = "&"
	}
	return ret
}
