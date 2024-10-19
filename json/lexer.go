package json

import (
	"errors"
	"fmt"
)

type tokenType int

const (
	tStr      tokenType = iota // "hello"
	tNum                       // 1.23
	tTrue                      // true
	tFalse                     // false
	tNull                      // null
	tLBracket                  // [
	tRBracket                  // ]
	tLBrace                    // {
	tRBrace                    // }
	tColon                     // :
	tComma                     // ,
)

type token struct {
	tType tokenType
	value any
}

func newStringToken(value string) token {
	return token{tType: tStr, value: value}
}

func newNumberToken(value float64) token {
	return token{tType: tNum, value: value}
}

func newTrueToken() token {
	return token{tType: tTrue}
}

func newFalseToken() token {
	return token{tType: tFalse}
}

func newNullToken() token {
	return token{tType: tNull}
}

func newLBracketToken() token {
	return token{tType: tLBracket}
}

func newRBracketToken() token {
	return token{tType: tRBracket}
}

func newLBraceToken() token {
	return token{tType: tLBrace}
}

func newRBraceToken() token {
	return token{tType: tRBrace}
}

func newColonToken() token {
	return token{tType: tColon}
}

func newCommaToken() token {
	return token{tType: tComma}
}

func tokenize(src []byte) ([]token, error) {
	res := []token{}
	cur := 0

	for cur < len(src) {
		switch {
		case src[cur] == '"':
			token, newCur, err := extractStringToken(src, cur)
			if err != nil {
				return nil, err
			}
			res = append(res, token)
			cur = newCur
		case '0' <= src[cur] && src[cur] <= '9':
			token, newCur, err := extractNumberToken(src, cur)
			if err != nil {
				return nil, err
			}
			res = append(res, token)
			cur = newCur
		case src[cur] == 'n':
			token, newCur, err := extractKeywordToken(src, cur, "null", newNullToken)
			if err != nil {
				return nil, err
			}
			res = append(res, token)
			cur = newCur
		case src[cur] == 't':
			token, newCur, err := extractKeywordToken(src, cur, "true", newTrueToken)
			if err != nil {
				return nil, err
			}
			res = append(res, token)
			cur = newCur
		case src[cur] == 'f':
			token, newCur, err := extractKeywordToken(src, cur, "false", newFalseToken)
			if err != nil {
				return nil, err
			}
			res = append(res, token)
			cur = newCur
		case src[cur] == '[':
			res = append(res, newLBracketToken())
			cur++
		case src[cur] == ']':
			res = append(res, newRBracketToken())
			cur++
		case src[cur] == '{':
			res = append(res, newLBraceToken())
			cur++
		case src[cur] == '}':
			res = append(res, newRBraceToken())
			cur++
		case src[cur] == ':':
			res = append(res, newColonToken())
			cur++
		case src[cur] == ',':
			res = append(res, newCommaToken())
			cur++
		default:
			cur++
		}
	}

	return res, nil
}

func extractStringToken(src []byte, cur int) (token, int, error) {
	cur++
	start := cur

	for cur < len(src) && src[cur] != '"' {
		cur++
	}

	if cur == len(src) {
		return token{}, 0, errors.New("unterminated string, reached EOF")
	}

	str := string(src[start:cur])
	// skip last "
	cur++

	return newStringToken(str), cur, nil
}

func extractNumberToken(src []byte, cur int) (token, int, error) {
	start := cur
	hasDot := false

	for cur < len(src) && (src[cur] >= '0' && src[cur] <= '9' || src[cur] == '.') {
		if src[cur] == '.' {
			if hasDot {
				return token{}, 0, errors.New("invalid number format")
			}
			hasDot = true
		}
		cur++
	}

	if start == cur || src[cur-1] == '.' {
		return token{}, 0, errors.New("invalid number format")
	}

	numStr := string(src[start:cur])
	var num float64
	_, err := fmt.Sscanf(numStr, "%f", &num)
	if err != nil {
		return token{}, 0, err
	}

	return newNumberToken(num), cur, nil
}

func extractKeywordToken(src []byte, cur int, keyword string, tokenFunc func() token) (token, int, error) {
	if cur+len(keyword) > len(src) {
		return token{}, 0, errors.New("unterminated token, reached EOF")
	}
	maybeKeyword := string(src[cur : cur+len(keyword)])
	if maybeKeyword != keyword {
		return token{}, 0, fmt.Errorf("invalid token: %v", maybeKeyword)
	}

	return tokenFunc(), cur + len(keyword), nil
}
