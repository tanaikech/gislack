// Package main (jsoncontrol.go) :
// Control gistslack by JSON data
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/urfave/cli"
)

// jsonControl : Struct for controlling JSON
type jsonControl struct {
	Command string                 `json:"command"` // gist, slack, doublesubmit, auth
	Options map[string]interface{} `json:"options"`
}

// getAugs : Get augs and initialize them
func getAugs(c *cli.Context) *iniparamsContainer {
	j := &jsonControl{}
	j.Command = c.Command.Names()[0]
	if j.Command == "json" {
		err := json.Unmarshal([]byte(c.String("json")), &j)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: JSON Format error. Please confirm inputted JSON again.\n")
			os.Exit(1)
		}
	} else {
		obj := map[string]interface{}{}
		for _, e := range c.FlagNames() {
			switch reflect.TypeOf(c.Generic(e)).String() {
			case "*flag.stringValue":
				obj[e] = c.String(e)
			case "*flag.intValue":
				obj[e] = c.Int(e)
			case "*flag.boolValue":
				obj[e] = c.Bool(e)
			}
		}
		j.Options = obj
	}
	i := &iniparamsContainer{
		&authParams{},
		&jsonControl{},
	}
	i.jsonControl = j

	workdir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	i.WorkDir = workdir
	var cfgdir string
	cfgdir = os.Getenv(cfgpathenv)
	if st, _ := i.jsonControl.Options["cfgdirectory"].(string); len(st) > 0 {
		cfgdir = st
	}
	if cfgdir == "" {
		cfgdir = workdir
	}
	i.CfgDir = cfgdir
	return i
}

// getCfg : Get data from a CFG file
func (i *iniparamsContainer) getCfg() *iniparamsContainer {
	p := &authParams{}
	p.pstart = time.Now()
	if cfgdata, err := ioutil.ReadFile(filepath.Join(i.CfgDir, cfgFile)); err == nil {
		err = json.Unmarshal(cfgdata, &p.GislackCfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Format error of '%s'. ", cfgFile)
			os.Exit(1)
		}
	} else {
		if i.jsonControl.Command != "auth" {
			fmt.Printf("Error: %s.cfg is not found. Please authorization for gist and/or slack you want to use. Please access token by executing '%s auth'.\n", appname, appname)
			os.Exit(1)
		}
	}
	i.authParams = p
	i.jsonControl.Options["usejsoncontrol"] = false
	return i
}

// keyChk : Check keys for controlling JSON
func (i *iniparamsContainer) keyChk() *iniparamsContainer {
	boolkeys := []string{
		"jsonparser",
		"public",
		"list",
		"deleteall",
		"anonymous",
		"listasjson",
		"channellist",
		"filelist",
		"channelhistory",
		"deletefiles",
		"simpleresult",
		"chkgisttoken",
		"filelistasjson",
		"appcheck",
	}
	for _, key := range boolkeys {
		if i.chkArgs(key) == nil {
			i.jsonControl.Options[key] = false
		}
	}
	stringkeys := []string{
		"title",
		"cfgdirectory",
		"files",
		"get",
		"updateoverwrite",
		"updateadd",
		"delete",
		"filenames",
		"file",
		"channel",
		"content",
		"filetype",
		"initialcomment",
		"user",
		"deletefile",
		"deletehistory",
		"gistclientid",
		"gistclientsecret",
		"slackclientid",
		"slackclientsecret",
		"gistcode",
		"slackcode",
		"getfile",
		"workdir",
		"getversion",
		"gethistory",
	}
	for _, key := range stringkeys {
		if i.chkArgs(key) == nil {
			i.jsonControl.Options[key] = ""
		}
	}
	if i.chkArgs("deletehistories") == nil {
		i.jsonControl.Options["deletehistories"] = 0
	} else {
		i.jsonControl.Options["deletehistories"] = int(i.jsonControl.Options["deletehistories"].(float64))
	}
	if i.chkArgs("port") == nil {
		i.jsonControl.Options["port"] = 8080
	} else {
		i.jsonControl.Options["port"] = int(i.jsonControl.Options["port"].(float64))
	}
	i.jsonControl.Options["usejsoncontrol"] = true
	return i
}
