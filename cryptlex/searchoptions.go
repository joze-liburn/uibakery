package cryptlex

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

type (
	OptCmpString struct {
		Cmp     string
		Operand string
	}

	OptCmpTime struct {
		Cmp     string
		Operand time.Time
	}

	CryptlexOpts struct {
		Page      int          // The page number.
		Limit     int          // The number of records per page. Must be a number between 1 and 100.
		Sort      string       // The sort string to sort the returned records e.g. "-createdAt"
		Name      OptCmpString // Name of the reseller
		Email     OptCmpString // Notification email address of the reseller.
		Search    string       // Search string.
		Id        OptCmpString // ID of the resource.
		CreatedAt OptCmpTime   // Date of creation
		UpdatedAt OptCmpTime   // Date of last update
	}

	CallOptions func(*CryptlexOpts) error
)

var (
	errCallOption      = errors.New("call option error")
	OperatorsId        = []string{"eq", "ne", "in", "nin"}
	OperatorsName      = []string{"eq", "ne", "sw", "ew", "cn", "nc", "in", "nin"}
	OperatorsEmail     = []string{"eq", "ne", "sw", "ew", "cn", "nc", "in", "nin"}
	OperatorsCreatedAt = []string{"eq", "ne", "gt", "ge", "lt", "le"}
	OperatorsUpdatedAt = []string{"eq", "ne", "gt", "ge", "lt", "le"}
	TimeFormats        = []string{"2006-01-02", "2006-01-02 15:04", "2006-01-02 15:04:05"}
)

func (cs *OptCmpString) String() string {
	return fmt.Sprintf("%s%%20%s", cs.Cmp, cs.Operand)
}

func (ct *OptCmpTime) String() string {
	return fmt.Sprintf("%s%%20%s", ct.Cmp, ct.Operand.Format("2006-01-02T15:04:05Z"))
}

func parseTime(layouts []string, val string) (time.Time, error) {
	for _, layout := range layouts {
		if t, e := time.Parse(layout, val); e == nil {
			return t, nil
		}
	}
	return time.Time{}, &time.ParseError{Value: val}
}

func cmdlineToOperatorOperand(cmd string) (string, string) {
	op, target, _ := strings.Cut(cmd, " ")
	return op, target
}

func withFromFlag(val string, supported []string, fnc func(string, string) CallOptions) (CallOptions, error) {
	op, tg := cmdlineToOperatorOperand(val)
	if i := slices.Index(supported, op); i < 0 {
		return nil, fmt.Errorf("%w: operator (%s) must be one of %s", errCallOption, op, strings.Join(supported, ", "))
	}
	return fnc(op, tg), nil
}

func withFromTimeFlag(val string, supported []string, fnc func(string, time.Time) CallOptions) (CallOptions, error) {
	op, tg := cmdlineToOperatorOperand(val)
	if i := slices.Index(supported, op); i < 0 {
		return nil, fmt.Errorf("%w: operator (%s) must be one of %s", errCallOption, op, strings.Join(supported, ", "))
	}
	if tm, err := parseTime(TimeFormats, tg); err != nil {
		return nil, fmt.Errorf("%w: unsupported time format (%q)", errCallOption, tg)
	} else {
		return fnc(op, tm), nil
	}
}

func WithPage(page int) CallOptions {
	return func(co *CryptlexOpts) error {
		if page <= 0 {
			return fmt.Errorf("%w: invalid page %d (1..2147483647)", errCallOption, page)
		}
		co.Page = page
		return nil
	}
}

func WithPageSize(limit int) CallOptions {
	return func(co *CryptlexOpts) error {
		if limit <= 0 || limit > 100 {
			return fmt.Errorf("%w: invalid page size %d (1..100)", errCallOption, limit)
		}
		co.Limit = limit
		return nil
	}
}

func WithSort(on string) CallOptions {
	return func(co *CryptlexOpts) error {
		co.Sort = on
		return nil
	}
}

func WithId(cmp string, id string) CallOptions {
	return func(co *CryptlexOpts) error {
		if len(cmp) == 0 || len(id) == 0 {
			return fmt.Errorf("%w: invalid operation %q or target id %q", errCallOption, cmp, id)
		}
		co.Id = OptCmpString{cmp, id}
		return nil
	}
}

func WithIdFlag(fs string) (CallOptions, error) {
	return withFromFlag(fs, OperatorsId, WithId)
}

func WithName(cmp string, name string) CallOptions {
	return func(co *CryptlexOpts) error {
		if len(cmp) == 0 || len(name) == 0 {
			return fmt.Errorf("%w: invalid operation %q or target name %q", errCallOption, cmp, name)
		}
		co.Name = OptCmpString{cmp, name}
		return nil
	}
}

func WithNameFlag(fs string) (CallOptions, error) {
	return withFromFlag(fs, OperatorsName, WithName)
}

func WithEmail(cmp string, email string) CallOptions {
	return func(co *CryptlexOpts) error {
		if len(cmp) == 0 || len(email) == 0 {
			return fmt.Errorf("%w: invalid operation %q or target email %q", errCallOption, cmp, email)
		}
		co.Email = OptCmpString{cmp, email}
		return nil
	}
}

func WithEmailFlag(fs string) (CallOptions, error) {
	return withFromFlag(fs, OperatorsEmail, WithEmail)
}

func WithSearch(what string) CallOptions {
	return func(co *CryptlexOpts) error {
		co.Search = what
		return nil
	}
}

func WithCreatedAt(cmp string, created time.Time) CallOptions {
	return func(co *CryptlexOpts) error {
		co.CreatedAt = OptCmpTime{cmp, created}
		return nil
	}
}

func WithCreatedAtFlag(what string) (CallOptions, error) {
	return withFromTimeFlag(what, OperatorsCreatedAt, WithCreatedAt)
}

func WithUpdatedAt(cmp string, updated time.Time) CallOptions {
	return func(co *CryptlexOpts) error {
		co.UpdatedAt = OptCmpTime{cmp, updated}
		return nil
	}
}

func WithUpdatedAtFlag(what string) (CallOptions, error) {
	return withFromTimeFlag(what, OperatorsUpdatedAt, WithUpdatedAt)
}

func (co *CryptlexOpts) Apply(opts ...CallOptions) error {
	for _, op := range opts {
		if err := op(co); err != nil {
			return err
		}
	}
	return nil
}

func (co *CryptlexOpts) ToSearchParam() string {
	ret := ""
	sep := "?"
	if co.Page > 0 {
		ret = fmt.Sprintf("%s%spage=%d", ret, sep, co.Page)
		sep = "&"
	}
	if co.Limit > 0 {
		ret = fmt.Sprintf("%s%slimit=%d", ret, sep, co.Limit)
		sep = "&"
	}
	if len(co.Sort) > 0 {
		ret = fmt.Sprintf("%s%ssort=%s", ret, sep, co.Sort)
		sep = "&"
	}
	if len(co.Name.Cmp) > 0 {
		ret = fmt.Sprintf("%s%sname=%s", ret, sep, co.Name.String())
		sep = "&"
	}
	if len(co.Email.Cmp) > 0 {
		ret = fmt.Sprintf("%s%semail=%s", ret, sep, co.Email.String())
		sep = "&"
	}
	if len(co.Search) > 0 {
		ret = fmt.Sprintf("%s%s%s", ret, sep, co.Search)
		sep = "&"
	}
	if len(co.Id.Cmp) > 0 {
		ret = fmt.Sprintf("%s%sid=%s", ret, sep, co.Id.String())
		sep = "&"
	}
	if len(co.CreatedAt.Cmp) > 0 {
		ret = fmt.Sprintf("%s%screatedAt=%s", ret, sep, co.CreatedAt.String())
		sep = "&"
	}
	if len(co.UpdatedAt.Cmp) > 0 {
		ret = fmt.Sprintf("%s%supdatedAt=%s", ret, sep, co.UpdatedAt.String())
		sep = "&"
	}
	return ret
}
