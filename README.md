Playing with the http stack in golang, mostly an educational thing

`go run tapi.go` and open the browser to localhost:7678


fix

- open browser
- embed HTML in Go code
- proper httptrace timing
- handle redirs
- allow repeated keys (key: [v1, v2, v3])
- send errors in JSON
- syntax highlight JSON?
- pretty print HTML
- tests using httptest
- gzip compression
- treat timeout properly
- [atob decoding utf8 problem](http://stackoverflow.com/questions/30106476/using-javascripts-atob-to-decode-base64-doesnt-properly-decode-utf-8-strings)