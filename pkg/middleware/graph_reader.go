package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type bodyGraphQL struct {
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
	Query         string                 `json:"query"`
}

func (rs *Resource) GraphQueryReader(whitelistFunctions []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fix unexpected result due to graphQL playground is currently not working properly
			if r.Method == http.MethodGet {
				// passes the captured values to the context
				ctx := r.Context()
				ctx = context.WithValue(ctx, RequestFunctionNameKey, "-")
				ctx = context.WithValue(ctx, PublicFunctionKey, true)

				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)

				return
			}

			buf, err := io.ReadAll(r.Body)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			rdr := io.NopCloser(bytes.NewBuffer(buf))
			r.Body = io.NopCloser(bytes.NewBuffer(buf))

			var body bodyGraphQL

			err = json.Unmarshal(func(reader io.Reader) []byte {
				buf := new(bytes.Buffer)
				buf.ReadFrom(reader)
				return buf.Bytes()
			}(rdr), &body)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// TODO: to make better approach
			// searches if it contains the whitelisted functions
			publicFunction := false // sets default value
			for _, funcName := range whitelistFunctions {
				if strings.Contains(body.Query, fmt.Sprintf("%s(", funcName)) {
					// skip auth validation when querying this function
					// it allows unauthenticated users to use this function
					publicFunction = true
				}
			}

			// captures indexes of the `{` symbol
			firstCurlyBracket := strings.IndexAny(body.Query, "{")

			// captures indexes of the `(` symbol
			firstParenthesesIdx := strings.Index(body.Query, "(")

			// re-capture the query
			subStrQuery := strings.TrimSpace(subStr(body.Query, firstCurlyBracket+1, len(body.Query)))
			firstParenthesesIdx = strings.Index(subStrQuery, "(")
			functionName := subStr(subStrQuery, 0, firstParenthesesIdx)

			// passes the captured values to the context
			ctx := r.Context()
			ctx = context.WithValue(ctx, RequestFunctionNameKey, functionName)
			ctx = context.WithValue(ctx, PublicFunctionKey, publicFunction)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

			return
		})
	}
}

// NOTE: this isn't multi-Unicode-codepoint aware, like specifying skintone or
//
//	gender of an emoji: https://unicode.org/emoji/charts/full-emoji-modifiers.html
func subStr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
