package typescript

import (
	"github.com/hashicorp/terraform/helper/schema"
	"bytes"
	"archive/zip"
	"io/ioutil"
	"os"
	"io"
	"os/exec"
	"fmt"
	"path"
	"path/filepath"
	"io/fs"
	"strings"
	"math/rand"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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
			"output_md5": {
				Type:        schema.TypeString,
				Description: "base64 encoded zip file of generated code",
				Computed:    true,
			},
			"output_sha": {
				Type:        schema.TypeString,
				Description: "base64 encoded zip file of generated code",
				Computed:    true,
			},
			"output_base64sha256": {
				Type:        schema.TypeString,
				Description: "base64 encoded zip file of generated code",
				Computed:    true,
			},
			"output_content_base64": {
				Type:        schema.TypeString,
				Description: "base64 encoded zip file of generated code",
				Computed:    true,
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
	out := bytes.Buffer{}	
	zip := zip.NewWriter(&out)

	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filename := strings.Replace(path,fmt.Sprint(dir,"/"),"",1)
			contents, err := os.Open(path)
			if err != nil {
				return err
			}
			writer, err := zip.Create(filename)	
			if err != nil {
				return err
			}
			io.Copy(writer, contents)
		}
		return nil
	})
	zip.Close()
	data:= out.Bytes()

	h := sha1.New()
	h.Write(data)
	sha1 := hex.EncodeToString(h.Sum(nil))
	d.Set("output_sha",sha1)

	h256 := sha256.New()
	h256.Write(data)

	shaSum := h256.Sum(nil)
	sha256base64 := base64.StdEncoding.EncodeToString(shaSum[:])
	d.Set("output_base64sha256",sha256base64)

	md5 := md5.New()
	md5.Write(data)
	md5Sum := hex.EncodeToString(md5.Sum(nil))
	d.Set("output_md5",md5Sum)

	out_base64 := base64.StdEncoding.EncodeToString(data)
	d.Set("output_content_base64",out_base64)

	if err != nil {
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
