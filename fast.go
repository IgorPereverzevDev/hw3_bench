package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"io"
	"os"
	"strings"
	"sync"
)

type User struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}

var pool = sync.Pool{
	New: func() interface{} {
		return &User{}
	},
}

// FastSearch вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {

	file, _ := os.Open(filePath) // For read access.
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var seenBrowsers []string
	var unique bool
	var notSeen bool
	var isAndroid bool
	var isMSIE bool

	scanner := bufio.NewScanner(file)

	_, err := fmt.Fprintln(out, fmt.Sprintf("found users:"))
	if err != nil {
		return
	}

	var step int
	for scanner.Scan() {
		step = step + 1
		row := scanner.Bytes()

		if bytes.Contains(row, []byte("Android")) == false && bytes.Contains(row, []byte("MSIE")) == false {
			continue
		}

		user := pool.Get().(*User)
		err := easyjson.Unmarshal(row, user)
		if err != nil {
			return
		}

		isAndroid = false
		isMSIE = false

		for _, browser := range user.Browsers {
			notSeen = true
			unique = false

			if strings.Contains(browser, "Android") {
				isAndroid = true
				unique = true
			} else if strings.Contains(browser, "MSIE") {
				isMSIE = true
				unique = true
			}

			if unique == true {
				for _, element := range seenBrowsers {
					if element == browser {
						notSeen = false
					}
				}

				if notSeen {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}
		}

		pool.Put(user)
		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.Replace(user.Email, "@", " [at] ", -1)
		_, err = fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", step-1, user.Name, email))
		if err != nil {
			return
		}
	}

	_, err = fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
	if err != nil {
		return
	}
}

var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson97766e5aDecodeMainTestGoBench(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = in.String()
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = in.String()
		case "name":
			out.Name = in.String()
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson97766e5aEncodeMainTestGoBench(out *jwriter.Writer, in User) {
	out.RawByte('{')
	_ = true
	_ = true
	{
		const prefix string = ",\"browsers\":"
		out.RawString(prefix[1:])
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(v3)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(in.Email)
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(in.Name)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson97766e5aEncodeMainTestGoBench(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson97766e5aEncodeMainTestGoBench(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson97766e5aDecodeMainTestGoBench(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson97766e5aDecodeMainTestGoBench(l, v)
}
