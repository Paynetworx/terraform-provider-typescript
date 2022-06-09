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
			"output_file":{
				Type:        schema.TypeString,
				Description: "path to outputed file",
				Computed:    true,
			},
			"target":{
				Type:		 schema.TypeString,
				Description: "javascript version to target",
				Required:	 true,
				ForceNew:	 true,
			},
			"es_module_interop":{
				Type:			schema.TypeBool,
				Description:	"turns on esModuleInterop in the compile step",
				Optional:		true,
				Default:		false,
				ForceNew:		true,
			},
			"additional_files":{
				Type:		 schema.TypeList,
				Description: "additional files to put into zip file",
				Optional:    true,
				ForceNew:	 true,
				Elem:	     &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"filename": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},	
								},
				},
			},
		},
	}
}


func resourceTypescriptCreate(d *schema.ResourceData, meta interface{}) error {
	source := d.Get("source").(string)
	target := d.Get("target").(string)
	esModuleInterop_arg := ""
	esModuleInterop := d.Get("es_module_interop").(bool)
	if(esModuleInterop){
		esModuleInterop_arg = "--esModuleInterop"
	}
	dir, err := ioutil.TempDir("","*")
	if err != nil {
		return err
	}
	output_file, err := ioutil.TempFile("","lambda.zip")
	d.Set("output_file",output_file.Name())
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
		"--module","commonjs",
		"--pretty","false",
		esModuleInterop_arg,
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
	if additional_files, ok := d.GetOk("additional_files"); ok {
		additional_files_list := additional_files.([]interface{})
		for _, additional_file := range  additional_files_list {
			src := additional_file.(map[string]interface{})
			writer, err := zip.Create(src["filename"].(string))
			if err != nil {
				return err
			}
			writer.Write([]byte(src["content"].(string)))	
		}
	}
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
	
	_, err = io.Copy(output_file,bytes.NewBuffer(data))
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
