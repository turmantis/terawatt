package terraform

import (
	"reflect"
	"strings"
	"testing"
)

func Test_parseVersionFile(t *testing.T) {
	want := []string{"3.2.1-bar", "1.2.3"}
	got, err := parseVersionFile(strings.NewReader(sampleVersionHtml))
	if err != nil {
		t.Errorf("parseVersionFile() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseVersionFile() got = %v, want %v", got, want)
	}
}

const sampleVersionHtml = `
<!DOCTYPE html>
<html>
    <head>
        <title>Some Versions</title>
		<meta name=viewport content="width=device-width, initial-scale=1">
		<style type="text/css">
			html {}
			body {}
			footer {}
			a {}
			a:hover {}
			a:visited {}
			ul {}
			li {}
		</style>
    </head>
    <body>
        <ul>
            <li>
            	<a href="../">../</a>
            </li>
            <li>
            	<a href="/foo/3.2.1-bar/">foo_3.2.1-bar</a>
            </li>
            <li>
            	<a href="/foo/1.2.3/">foo_1.2.3</a>
            </li>
        </ul>
    </body>
</html>
`
