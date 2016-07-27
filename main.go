package main

import (
	"go/parser"
	"go/token"
	"log"
	"strings"

	xi "github.com/rkusa/xi-peer"
)

var done chan bool
var peer *xi.Peer

func main() {
	done = make(chan bool)
	peer = xi.New()
	peer.Handle("ping", handlePing)
	peer.Handle("ping_from_editor", handlePingFromEditor)

	<-done
	log.Println("Closing ...")
}

func handlePing(params interface{}) {
	log.Println("ping received")
}

func handlePingFromEditor(params interface{}) {
	log.Println("ping_from_editor received")
	retrieveAllLines(10)
}

func retrieveAllLines(concurrency int) {
	var n float64 = 0
	if err := peer.RequestSync("n_lines", nil, &n); err != nil {
		log.Fatal(err)
	}

	log.Println("n_lines", n)

	// start := time.Now()

	response := make(chan *xi.Call, concurrency)
	lines := make([]string, int(n))
	remaining := len(lines)
	receiving := 0

	for remaining > 0 || receiving > 0 {
		if receiving >= concurrency || (remaining == 0 && receiving > 0) {
			// wait for a response to arrive, before making a new request
			call := <-response
			if call.Error != nil {
				log.Fatal(call.Error)
			}
			receiving--
		}

		if remaining == 0 {
			continue
		}

		lnr := remaining - 1
		peer.Request("get_line", map[string]int{"line": lnr}, &lines[lnr], response)

		remaining--
		receiving++
	}

	// elapsed := time.Since(start)
	// alert := fmt.Sprintf("Retrieving all %d lines took %s (with %d concurrent requests)", len(lines), elapsed, concurrency)
	// log.Println(alert)

	// if err := peer.CallSync("alert", map[string]string{"msg": alert}, nil); err != nil {
	// 	log.Fatal(err)
	// }

	processLine(strings.Join(lines, ""))

	done <- true

	// if err := p.lint(); err != nil {
	// 	return err
	// }
}

type Span struct {
	Start token.Pos
	End   token.Pos
	Color int64
}

func processLine(src string) {
	// log.Println("SRC", src)

	// expr, err := parser.ParseExpr(line)
	// log.Println(expr, err)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		log.Println("Error", err)
		return
	}

	// Print the AST.
	// ast.Fprint(os.Stderr, fset, f, ast.NotNilFilter)
	// return

	spans := make(chan *Span, 10)
	done := make(chan bool)

	go func() {
		lnr := 0
		lsp := []map[string]interface{}{}

		for s := range spans {
			start := fset.Position(s.Start)
			end := fset.Position(s.End)
			params := map[string]interface{}{
				"start": start.Column - 1,
				"end":   end.Column - 1,
				"fg":    s.Color,
				"font":  0,
			}

			if start.Line-1 == lnr {
				lsp = append(lsp, params)
				continue
			} else {
				if len(lsp) > 0 {
					err := peer.Notify("set_line_fg_spans", map[string]interface{}{
						"line": lnr, "spans": lsp,
					})
					if err != nil {
						log.Fatal(err)
					}
				}

				lnr = start.Line - 1
				lsp = []map[string]interface{}{params}
			}
		}

		// TODO: remove redundant code
		if len(lsp) > 0 {
			err := peer.Notify("set_line_fg_spans", map[string]interface{}{
				"line": lnr, "spans": lsp,
			})
			if err != nil {
				log.Fatal(err)
			}
		}

		done <- true
	}()

	walk(f, spans)

	close(spans)
	<-done // wait until all spans are send
}
