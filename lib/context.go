package lib

import (
	"errors"
	"net/url"
	"strconv"
)

type Context struct {
	Params    map[string]string
	GetParams url.Values
}

func (c Context) ID() uint {
	id, err := strconv.ParseUint(c.Params["id"], 10, 32)

	if err != nil {
		return 0
	}

	return uint(id)
}

func (c Context) Limit() (uint, error) {
	limit, err := strconv.ParseInt(c.Params["limit"], 10, 32)

	if err != nil {
		return 0, errors.New("Error when parsing limit parameter")
	}

	if limit <= 0 {
		return 0, errors.New("Limit must bigger than 0")
	}

	return uint(limit), nil
}

func (c Context) Offset() (uint, error) {
	offset, err := strconv.ParseInt(c.Params["offset"], 10, 32)

	if err != nil {
		return 0, errors.New("Error when parsing offset parameter")
	}

	if offset < 0 {
		return 0, errors.New("Offset must bigger than or equal to 0")
	}

	return uint(offset), nil
}
