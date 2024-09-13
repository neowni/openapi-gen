package golang

import (
	c "columba-livia/content"
)

func clientConvert() (render render) {
	return func() c.C {
		file.importMap["fmt"] = ""
		file.importMap["strconv"] = ""
		file.importMap["encoding/json"] = ""
		file.importMap["github.com/go-resty/resty/v2"] = ""

		return c.C(clientConvertContent).TrimSpace()
	}
}

const clientConvertContent = `
// 																				转换器

var convert *_convert

type _convert struct{}

func (c *_convert) formatBool(v bool) string {
	return strconv.FormatBool(v)
}

func (c *_convert) formatInt(v int) string {
	return strconv.Itoa(v)
}

func (c *_convert) formatFloat64(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

// 																				uri

func (c *_convert) uriBool(r *resty.Request, n string, v bool) {
	r.SetPathParam(n, c.formatBool(v))
}

func (c *_convert) uriInt(r *resty.Request, n string, v int) {
	r.SetPathParam(n, c.formatInt(v))
}

func (c *_convert) uriFloat64(r *resty.Request, n string, v float64) {
	r.SetPathParam(n, c.formatFloat64(v))
}

func (c *_convert) uriString(r *resty.Request, n string, v string) {
	r.SetPathParam(n, v)
}

// 																				qry

func (c *_convert) qryBoolR(r *resty.Request, n string, v bool) {
	r.SetQueryParam(n, c.formatBool(v))
}

func (c *_convert) qryIntR(r *resty.Request, n string, v int) {
	r.SetQueryParam(n, c.formatInt(v))
}

func (c *_convert) qryFloat64R(r *resty.Request, n string, v float64) {
	r.SetQueryParam(n, c.formatFloat64(v))
}

func (c *_convert) qryStringR(r *resty.Request, n string, v string) {
	r.SetQueryParam(n, v)
}

func (c *_convert) qryBool(r *resty.Request, n string, v *bool) {
	if v == nil {
		return
	}

	r.SetQueryParam(n, c.formatBool(*v))
}

func (c *_convert) qryInt(r *resty.Request, n string, v *int) {
	if v == nil {
		return
	}

	r.SetQueryParam(n, c.formatInt(*v))
}

func (c *_convert) qryFloat64(r *resty.Request, n string, v *float64) {
	if v == nil {
		return
	}

	r.SetQueryParam(n, c.formatFloat64(*v))
}

func (c *_convert) qryString(r *resty.Request, n string, v *string) {
	if v == nil {
		return
	}

	r.SetQueryParam(n, *v)
}

// 																				rsp

func (*_convert) ResponseJSON(resp *resty.Response, rsp any) (err error) {
	body := resp.Body()
	if len(body) == 0 {
		body = []byte("{}")
	}

	return json.Unmarshal(body, rsp)
}

func (*_convert) ResponseText(resp *resty.Response, rsp *string) (err error) {
	*rsp = string(resp.Body())
	return nil
}

func (*_convert) ResponseEmpty(resp *resty.Response, rsp any) (err error) {
	return nil
}

func (*_convert) ResponseError(resp *resty.Response) (err error) {
	return &ErrorStatus{
		Status: resp.StatusCode(),
	}
}

//																				异常

type ErrorStatus struct {
	Status int
}

func (e *ErrorStatus) Error() string {
	return fmt.Sprintf("%d", e.Status)
}
`
