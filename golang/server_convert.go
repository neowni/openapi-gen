package golang

import (
	c "columba-livia/content"
)

func serverConvert() (render render) {
	return func() c.C {
		file.importMap["io"] = ""
		file.importMap["net/http"] = ""
		file.importMap["github.com/gin-gonic/gin"] = ""

		return c.C(serverConvertContent).TrimSpace()
	}
}

const serverConvertContent = `
var ErrorHandler func(ctx *gin.Context, err error)
var ErrorLogger func(err error)

//                                                                              转换函数

var convert *_convert

type _convert struct{}

func (*_convert) BindReqJSON(ctx *gin.Context, req any) (abort bool) {
	err := ctx.ShouldBind(req)
	if err == nil {
		return false
	}

	if ErrorLogger != nil {
		ErrorLogger(err)
	}

	ginErr := ctx.AbortWithError(http.StatusBadRequest, err)
	if ginErr != nil {
		if ErrorLogger != nil {
			ErrorLogger(err)
		}
	}

	return true
}

func (*_convert) BindReqText(ctx *gin.Context, req *string) (abort bool) {
	content, err := io.ReadAll(ctx.Request.Body)
	if err == nil {
		*req = string(content)
		return false
	}

	if ErrorLogger != nil {
		ErrorLogger(err)
	}

	ginErr := ctx.AbortWithError(http.StatusBadRequest, err)
	if ginErr != nil {
		if ErrorLogger != nil {
			ErrorLogger(err)
		}
	}

	return true
}

func (*_convert) BindURI(ctx *gin.Context, uri any) (abort bool) {
	err := ctx.ShouldBindUri(uri)
	if err == nil {
		return false
	}

	if ErrorLogger != nil {
		ErrorLogger(err)
	}

	ginErr := ctx.AbortWithError(http.StatusBadRequest, err)
	if ginErr != nil {
		if ErrorLogger != nil {
			ErrorLogger(err)
		}
	}

	return true
}

func (*_convert) BindQry(ctx *gin.Context, qry any) (abort bool) {
	err := ctx.ShouldBindQuery(qry)
	if err == nil {
		return false
	}

	if ErrorLogger != nil {
		ErrorLogger(err)
	}

	ginErr := ctx.AbortWithError(http.StatusBadRequest, err)
	if ginErr != nil {
		if ErrorLogger != nil {
			ErrorLogger(err)
		}
	}

	return true
}

func (*_convert) ResponseJSON(
	ctx *gin.Context,
	code int,
	rsp any,
) {
	ctx.JSON(code, rsp)
}

func (*_convert) ResponseText(ctx *gin.Context, code int, rsp *string) {
	ctx.String(code, *rsp)
}

func (*_convert) ResponseEmpty(ctx *gin.Context, code int, rsp any) {
	ctx.Status(code)
}

func (*_convert) ResponseError(ctx *gin.Context, err error) (abort bool) {
	if err == nil {
		return false
	}

	if ErrorHandler == nil {
		ginErr := ctx.AbortWithError(http.StatusInternalServerError, err)
		if ginErr != nil {
			if ErrorLogger != nil {
				ErrorLogger(err)
			}
		}
	} else {
		ErrorHandler(ctx, err)
	}

	return true
}
`
