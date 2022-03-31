package webserver

import (
	"bitbucket.org/digi-sense/gg-core"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func WriteContentType(ctx *fiber.Ctx, bytes []byte, contentType string) error {
	if len(contentType) == 0 {
		contentType = "application/octet-stream"
	}
	ctx.Response().Header.Set("Content-Type", contentType)
	ctx.Response().SetBody(bytes)
	return nil
}

func WriteDownload(ctx *fiber.Ctx, bytes []byte, contentType, filename string) error {
	if len(contentType) == 0 {
		contentType = "application/octet-stream"
	}
	ctx.Response().Header.Set("Content-Type", contentType)
	ctx.Response().Header.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", filename))
	ctx.Response().SetBody(bytes)
	return nil
}

func WriteHtml(ctx *fiber.Ctx, html string) error {
	ctx.Response().Header.Set("Content-Type", "text/html")
	ctx.Response().SetBody([]byte(html))
	return nil
}

func WriteImage(ctx *fiber.Ctx, bytes []byte) error {
	ctx.Response().Header.Set("Content-Type", "image/jpeg")
	ctx.Response().Header.Set("Content-Length", strconv.Itoa(len(bytes)))
	ctx.Response().SetBody(bytes)
	return nil
}

func WriteResponse(ctx *fiber.Ctx, data interface{}, err error) error {
	m := map[string]interface{}{}
	if nil != err {
		errMessage := err.Error()
		m["error"] = errMessage
		if errMessage == HttpUnauthorizedError.Error() || errMessage == HttpInvalidCredentialsError.Error() {
			ctx.Response().SetStatusCode(401)
			if nil == data {
				data = map[string]interface{}{"code": 401, "message": HttpUnauthorizedError.Error()}
			}
		} else if errMessage == AccessTokenExpiredError.Error() {
			ctx.Response().SetStatusCode(403)
			if nil == data {
				data = map[string]interface{}{"code": 403, "message": AccessTokenExpiredError.Error()}
			}
		} else if errMessage == RefreshTokenExpiredError.Error() {
			ctx.Response().SetStatusCode(401)
			if nil == data {
				data = map[string]interface{}{"code": 401, "message": RefreshTokenExpiredError.Error()}
			}
		} else if errMessage == AccessTokenInvalidError.Error() {
			ctx.Response().SetStatusCode(401)
			if nil == data {
				data = map[string]interface{}{"code": 401, "message": AccessTokenInvalidError.Error()}
			}
		} else {
			ctx.Response().SetStatusCode(500)
			if nil == data {
				data = map[string]interface{}{"code": 500, "message": err.Error()}
			}
		}
	}
	if nil != data {
		m["response"] = toJson(data)
	} else {
		m["response"] = "EMPTY_RESPONSE"
	}
	ctx.Response().Header.Set("Content-Type", "text/json")
	ctx.Response().SetBody([]byte(gg.JSON.Stringify(m)))

	return nil
}

func GetAuthTokens(ctx *fiber.Ctx) []string {
	response := make([]string, 0)
	t := getApplicationToken(ctx)
	if len(t) > 0 {
		response = append(response, t)
	}
	t = getAuthenticationToken(ctx)
	if len(t) > 0 {
		response = append(response, t)
	}
	return response
}

func BodyMap(ctx *fiber.Ctx, flatData bool) map[string]interface{} {
	var m map[string]interface{}
	body := ctx.Body()
	_ = gg.JSON.Read(body, &m)

	if data, b := m["data"].(map[string]interface{}); b && flatData {
		delete(m, "data")
		for k, v := range data {
			m[k] = v
		}
	}

	return m
}

func Params(ctx *fiber.Ctx, flatData bool, selectOnly ...string) map[string]interface{} {
	// get all body params
	response := BodyMap(ctx, flatData)
	if nil == response {
		response = make(map[string]interface{})
	}

	// try add form params
	if form, err := ctx.MultipartForm(); nil == err && nil != form && nil != form.Value {
		for k, v := range form.Value {
			if _, b := response[k]; !b {
				if len(v) == 1 {
					response[k] = v[0]
				} else {
					response[k] = v
				}
			}
		}
	}

	// route
	routeParams := ctx.Route().Params
	if len(routeParams) > 0 {
		for _, rp := range routeParams {
			response[rp] = ctx.Params(rp)
		}
	}

	// url query
	if path := ctx.OriginalURL(); len(path) > 0 {
		uri, err := url.Parse(path)
		if nil == err {
			query := uri.Query()
			if nil != query && len(query) > 0 {
				for k, v := range query {
					if len(v) == 1 {
						response[k] = v[0]
					} else {
						response[k] = v
					}
				}
			}
		}
	}

	// evaluate filter
	if len(selectOnly) > 0 {
		m := map[string]interface{}{}
		for _, name := range selectOnly {
			if v, b := response[name]; b {
				m[name] = v
			}
		}
		return m
	}

	return response
}

func AssertParams(ctx *fiber.Ctx, names []string) (map[string]interface{}, error) {
	params := Params(ctx, true)
	if len(names) > 0 {
		missing := make([]string, 0)
		// get missing parameters
		for _, name := range names {
			if value, b := params[name]; !b || len(gg.Convert.ToString(value)) == 0 {
				missing = append(missing, name)
			}
		}
		if len(missing) > 0 {
			return nil, errors.New(fmt.Sprintf("missing_params:%v", strings.Join(missing, ",")))
		}
	}
	return params, nil
}

func Upload(ctx *fiber.Ctx, root string, sizeLimit int64) ([]string, error) {
	// get form
	form, err := ctx.MultipartForm()
	if nil != err {
		return nil, err // not a form request
	}

	root = gg.Paths.Absolute(root) // absolute

	response := make([]string, 0)
	// loop on files
	for _, files := range form.File {
		// Loop through files:
		for _, file := range files {
			size := file.Size
			if sizeLimit > 0 && size > sizeLimit {
				continue
			}
			filename := file.Filename
			ext := gg.Paths.Extension(filename)
			name := gg.Paths.FileName(filename, false)
			dir := gg.Paths.Dir(filename)
			filename = gg.Paths.Concat(dir, name+"_"+gg.Coding.MD5(gg.Rnd.Uuid())+ext)
			path := gg.Paths.DatePath(root, filename, 3, true)
			err = ctx.SaveFile(file, path)
			if nil != err {
				return response, err
			}
			response = append(response, strings.Replace(path, root, ".", 1)) // relative path
		}
	}
	return response, nil
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func toJson(data interface{}) interface{} {
	if nil != data {
		return gg.JSON.Parse(gg.Convert.ToString(data))
	}
	return nil
}

func getAuthToken(auth *Authorization) (string, error) {
	value := auth.Value
	switch auth.Type {
	case "base", "basic":
		// BASE AUTH
		if strings.Index(value, ":") > -1 {
			return value, nil
		}
		data, err := gg.Coding.DecodeBase64(value)
		if nil != err {
			return "", err
		}
		return string(data), nil
	case "bearer":
		// BEARER AUTH
		return value, nil
	}
	return value, nil
}

func getAuthenticationToken(ctx *fiber.Ctx) string {
	data := getAuthentication(ctx)
	accessToken := gg.Reflect.GetString(data, "token")
	if len(accessToken) == 0 {
		// look into params
		params := Params(ctx, true)
		accessToken = gg.Reflect.GetString(params, "access_token")
	}
	return accessToken
}

func getApplicationToken(ctx *fiber.Ctx) string {
	// look into params
	params := Params(ctx, true)
	token := gg.Reflect.GetString(params, "token")
	if len(token) == 0 {
		token = gg.Reflect.GetString(params, "app_token")
		if len(token) == 0 {
			token = gg.Reflect.GetString(params, "application_token")
			if len(token) == 0 {
				token = gg.Reflect.GetString(params, "auth_token")
				if len(token) == 0 {
					token = gg.Reflect.GetString(params, "app-token")
					if len(token) == 0 {
						token = gg.Reflect.GetString(params, "auth-token")
					}
				}
			}
		}
	}
	return token
}

func getAuthentication(ctx *fiber.Ctx) map[string]string {
	response := map[string]string{}

	v := ctx.Get("Authorization", "")
	tokens := gg.Strings.Split(v, " ")
	response["authorization"] = v

	if len(tokens) == 2 {
		mode := tokens[0]
		response["mode"] = mode
		switch mode {
		case "Bearer":
			response["token"] = tokens[1]
		case "Basic":
			response["username"] = ""
			response["password"] = ""
			data, _ := gg.Coding.DecodeBase64(tokens[1])
			response["token"] = string(data)
			if len(data) > 0 {
				t := strings.Split(string(data), ":")
				if len(t) == 2 {
					response["username"] = t[0]
					response["password"] = t[1]
				}
			}
		}
	}

	return response
}
