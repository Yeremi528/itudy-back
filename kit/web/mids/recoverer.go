package mids

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/Yeremi528/itudy-back/kit/logger"
)

// Panics recovers from a panic and converts it to an error, which allows
// it to be reported and handled in the Errors middleware.
func Recoverer(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			// Defer a function to recover from a panic after executing the next handler.
			defer func() {
				if rec := recover(); rec != nil {
					if rec == http.ErrAbortHandler {
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						panic(rec)
					}

					s := prettyStack{}
					out, _ := s.parse(debug.Stack(), rec)

					stack := strings.Split(string(out), "\n")

					log.Info(r.Context(), "Recovering from panic", "stack", stack)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(h)
	}
}

type prettyStack struct{}

func (s prettyStack) parse(debugStack []byte, rvr interface{}) ([]byte, error) {
	var err error
	buf := &bytes.Buffer{}

	// Remove color logic, just plain text
	buf.WriteString("Message  ->   ")
	fmt.Fprintf(buf, "%v\n", rvr)

	// process debug stack info
	stack := strings.Split(string(debugStack), "\n")
	lines := []string{}

	// locate panic line, as we may have nested panics
	for i := len(stack) - 1; i > 0; i-- {
		lines = append(lines, stack[i])
		if strings.HasPrefix(stack[i], "panic(") {
			lines = lines[0 : len(lines)-2] // remove boilerplate
			break
		}
	}

	// reverse the lines
	for i := len(lines)/2 - 1; i >= 0; i-- {
		opp := len(lines) - 1 - i
		lines[i], lines[opp] = lines[opp], lines[i]
	}

	// decorate each line
	for i, line := range lines {
		lines[i], err = s.decorateLine(line, i)
		if err != nil {
			return nil, err
		}
	}

	for _, l := range lines {
		fmt.Fprintf(buf, "%s", l)
	}
	return buf.Bytes(), nil
}

func (s prettyStack) decorateLine(line string, num int) (string, error) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "\t") || strings.Contains(line, ".go:") {
		return s.decorateSourceLine(line, num)
	}
	if strings.HasSuffix(line, ")") {
		return s.decorateFuncCallLine(line, num)
	}
	if strings.HasPrefix(line, "\t") {
		return strings.Replace(line, "\t", "      ", 1), nil
	}
	return fmt.Sprintf("...%s\n", line), nil
}

func (s prettyStack) decorateFuncCallLine(line string, num int) (string, error) {
	idx := strings.LastIndex(line, "(")
	if idx < 0 {
		return "", errors.New("not a func call line")
	}

	buf := &bytes.Buffer{}
	pkg := line[0:idx]
	method := ""

	if idx := strings.LastIndex(pkg, string(os.PathSeparator)); idx < 0 {
		if idx := strings.Index(pkg, "."); idx > 0 {
			method = pkg[idx:]
			pkg = pkg[0:idx]
		}
	} else {
		method = pkg[idx+1:]
		pkg = pkg[0 : idx+1]
		if idx := strings.Index(method, "."); idx > 0 {
			pkg += method[0:idx]
			method = method[idx:]
		}
	}

	// No colors, just plain text formatting
	if num == 0 {
		buf.WriteString("Function ->   ")
	} else {
		buf.WriteString("...")
	}
	fmt.Fprintf(buf, "%s%s\n", pkg, method)

	return buf.String(), nil
}

func (s prettyStack) decorateSourceLine(line string, num int) (string, error) {
	idx := strings.LastIndex(line, ".go:")
	if idx < 0 {
		return "", errors.New("not a source line")
	}

	buf := &bytes.Buffer{}
	path := line[0 : idx+3]
	lineno := line[idx+3:]

	idx = strings.LastIndex(path, string(os.PathSeparator))
	dir := path[0 : idx+1]
	file := path[idx+1:]

	idx = strings.Index(lineno, " ")
	if idx > 0 {
		lineno = lineno[0:idx]
	}

	// Plain text formatting for source lines
	if num == 1 {
		buf.WriteString("File     ->   ")
	} else {
		buf.WriteString("...")
	}
	fmt.Fprintf(buf, "%s%s%s\n", dir, file, lineno)

	return buf.String(), nil
}
