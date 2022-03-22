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
						r.TestCheckResourceAttr(dsn,"output_sha","9b7d55491b87d06bc0412a420b244f1bce9c929e"),
					),
				},
			},
			CheckDestroy: func(*terraform.State) error {
				return nil
			},
		})
	}
}
