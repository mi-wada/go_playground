# json

## Syntax of JSON

Syntax of the JSON, written in [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form).

```ebnf
json   = object | array | string | number | bool | null
object = “{“ ((string ”:” json) (“,” string ”:” json)*)? “}”
array  = “[“ (json (“,” json)*)? “]”
string = “\”” char* “\””
char   = All UTF-8 Encoded Unicode Characters
number = "-"? digit digit* (“.” digit digit*)?
digit  = “0”..”9”
bool   = “true” | “false”
null   = “null”
```
