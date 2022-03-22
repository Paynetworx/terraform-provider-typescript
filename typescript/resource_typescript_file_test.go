package typescript

import (
	"fmt"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const templateTypescriptFileConfig = `
resource "typescript_file" "handler" {
	source = "%s"
	target = "es2020"
}

`


func TestTemplateDirRendering(t *testing.T) {
	var cases = []struct {
		file string
	}{
		{
			file: `_testdata/example_handler.ts`,
		},
	}
	dsn:="typescript_file.handler"
	for _, tt := range cases {
		// Run test case.
		r.UnitTest(t, r.TestCase{
			Providers: testProviders,
			Steps: []r.TestStep{
				{
					Config: fmt.Sprintf(templateTypescriptFileConfig, tt.file),
					Check: r.ComposeTestCheckFunc(
						r.TestCheckResourceAttr(dsn,"output_files.#","4"),
						r.TestCheckResourceAttr(dsn,"output_files.0.filename","example_handler.js"),
						r.TestCheckResourceAttr(dsn,"output_files.0.content","import * as lib from './lib';\nexport function handler() {\n    console.log(lib);\n}\n"),
					),
				},
			},
			CheckDestroy: func(*terraform.State) error {
				return nil
			},
		})
	}
}
