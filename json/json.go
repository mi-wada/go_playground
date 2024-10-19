package json

import (
	"errors"
	"fmt"
)

func Decode(src []byte) (any, error) {
	tokens, err := tokenize(src)
	if err != nil {
		return nil, err
	}

	res, cur, err := decodeJson(tokens, 0)
	if err != nil {
		return nil, err
	}
	if cur != len(tokens) {
		return nil, errors.New("extra tokens after JSON data")
	}

	return res, nil
}

func decodeJson(tokens []token, cur int) (any, int, error) {
	switch tokens[cur].tType {
	case tStr, tNum:
		return tokens[cur].value, cur + 1, nil
	case tTrue:
		return true, cur + 1, nil
	case tFalse:
		return false, cur + 1, nil
	case tNull:
		return nil, cur + 1, nil
	case tLBracket:
		return decodeArray(tokens, cur)
	case tLBrace:
		return decodeObject(tokens, cur)
	default:
		return nil, 0, fmt.Errorf("invalid token: %v", tokens[cur])
	}
}

func decodeArray(tokens []token, cur int) ([]any, int, error) {
	cur++ // skip first [

	res := []any{}
	for cur < len(tokens) && tokens[cur].tType != tRBracket {
		json, newCur, err := decodeJson(tokens, cur)
		if err != nil {
			return nil, 0, err
		}

		res = append(res, json)

		cur = newCur
		if cur < len(tokens) && tokens[cur].tType == tComma {
			cur++ // skip comma
		}
	}

	if cur == len(tokens) {
		return nil, 0, errors.New("unexpected end of JSON input")
	}
	if tokens[cur-1].tType == tComma && tokens[cur].tType == tRBracket {
		return nil, 0, errors.New("trailing comma in array")
	}

	cur++ // skip last ]

	return res, cur, nil
}

func decodeObject(tokens []token, cur int) (map[string]any, int, error) {
	cur++ // skip first {

	res := make(map[string]any)
	for cur < len(tokens) && tokens[cur].tType != tRBrace {
		key := tokens[cur]
		if key.tType != tStr {
			return nil, 0, errors.New("expected string key in object")
		}
		cur++ // skip key

		colon := tokens[cur]
		if colon.tType != tColon {
			return nil, 0, errors.New("expected colon after key in object")
		}
		cur++ // skip colon

		value, newCur, err := decodeJson(tokens, cur)
		if err != nil {
			return nil, 0, err
		}

		res[key.value.(string)] = value

		cur = newCur
		if cur < len(tokens) && tokens[cur].tType == tComma {
			cur++ // skip comma
		}
	}

	if cur == len(tokens) {
		return nil, 0, errors.New("unexpected end of JSON input")
	}
	if tokens[cur-1].tType == tComma && tokens[cur].tType == tRBrace {
		return nil, 0, errors.New("trailing comma in object")
	}

	cur++ // skip last }

	return res, cur, nil
}
