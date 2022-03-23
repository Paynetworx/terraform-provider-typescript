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
	additional_files {
		filename = "test.txt"
		content = "hello world"
	}
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
						r.TestCheckResourceAttr(dsn,"output_sha","cc4994ff823cf9b6a00aee89bcf08362ea01c138"),
					),
				},
			},
			CheckDestroy: func(*terraform.State) error {
				return nil
			},
		})
	}
}
