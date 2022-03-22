package typescript

import (
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"os"
	"os/exec"
	"fmt"
	"path"
	"path/filepath"
	"io/fs"
	"strings"
	"math/rand"
)

func resourceTypescriptFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTypescriptCreate,
		Read:   resourceTypescriptRead,
		Delete: resourceTypescriptDelete,

		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeString,
				Description: "Absolut paths to source typescript file",
				Required:    true,
				ForceNew:    true,
			},
			"triggers":{
				Type:        schema.TypeMap,
				Description: "Map of arbitrary strings to trigger a rebuild",
				Optional:    true,
				ForceNew:    true,
			},
			"output_files": {
				Type:        schema.TypeList,
				Description: "file names and contents of generated javascript files",
				Computed:    true,
				Elem:        &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content":{
										Type: schema.TypeString,
										Computed: true,
									},
									"filename":{
										Type: schema.TypeString,
										Computed: true,
									},
								},
							 },
			},
			"target":{
				Type:		 schema.TypeString,
				Description: "javascript version to target",
				Required:	 true,
				ForceNew:	 true,
			},
		},
	}
}


func resourceTypescriptCreate(d *schema.ResourceData, meta interface{}) error {
	source := d.Get("source").(string)
	target := d.Get("target").(string)

	dir, err := ioutil.TempDir("","*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	working_dir := path.Dir(source)
	source_file := path.Base(source)
	cmd := exec.Command("npx","tsc", 
		"--target",target,
		"--moduleResolution","node",
		"--outdir",dir,
		"--listFiles",
		"--pretty","false",
		source_file)
	cmd.Dir = working_dir

	stdout, err:= cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(stdout))
		return err
	}

	output_files := []map[string]string{}
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filename := strings.Replace(path,fmt.Sprint(dir,"/"),"",1)
			out := map[string]string{}
			out["filename"]=filename
			contents, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			out["content"]=string(contents)
			output_files = append(output_files,out)
		}
		return nil
	})
	if err != nil {
		return err
	}
	
	if err:=d.Set("output_files",output_files); err != nil {
		return err	
	}
	d.SetId(fmt.Sprintf("%d", rand.Int()))
	return nil
}

func resourceTypescriptRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTypescriptDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")
	return nil
}
