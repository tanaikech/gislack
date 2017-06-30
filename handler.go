// Package main (handler.go) :
// Handler for gistslack
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

// const :
const (
	appname    = "gislack"
	cfgFile    = "gislack.cfg"
	cfgpathenv = "GISLACK_CFG_PATH"

	gisturl         = "https://api.github.com/gists"
	gistauthcode    = "https://github.com/login/oauth/authorize?"
	gistaccesstoken = "https://github.com/login/oauth/access_token"
	gistchktoken    = "https://api.github.com/shibumori/whatever?"

	slackurl         = "https://slack.com/api/"
	slackauthcode    = "https://slack.com/oauth/authorize?"
	slackaccesstoken = "https://slack.com/api/oauth.access?"
	slackchkat       = "https://slack.com/api/auth.test?"
)

// initVal : Initial values
type initVal struct {
	pstart  time.Time
	workdir string
}

// gist : Commands for gist
func gist(c *cli.Context) error {
	g := getAugs(c).getCfg().initGistContainer()
	if c.Bool("list") || c.Bool("listasjson") {
		g.gistList()
		return nil
	}
	if len(c.String("get")) > 0 && len(c.String("gethistory")) == 0 {
		g.gistGet().disp()
		return nil
	}
	if len(c.String("gethistory")) > 0 && len(c.String("get")) == 0 {
		g.gistGet().disphis()
		return nil
	}
	if len(c.String("getversion")) > 0 && len(c.String("gethistory")) == 0 && len(c.String("get")) == 0 {
		g.gistGet().disp()
		return nil
	}
	if len(c.String("title")) > 0 && len(c.String("files")) > 0 &&
		(len(c.String("updateoverwrite")) == 0 && len(c.String("updateadd")) == 0) {
		g.defGistContainer().gistSubmit()
		if c.Bool("simpleresult") {
			g.simpleDisp()
		} else {
			g.disp()
		}
		return nil
	}
	if (len(c.String("updateoverwrite")) > 0 || len(c.String("updateadd")) > 0) &&
		(len(c.String("title")) > 0 || len(c.String("files")) > 0) {
		g.gistUpdate(g.defGistContainer().gistMakeUpdate()).disp()
		return nil
	}
	if len(c.String("delete")) > 0 {
		g.gistDel()
		return nil
	}
	if c.Bool("deleteall") {
		g.gistDeleteAll()
		return nil
	}
	fmt.Printf("Usage is `%s gist --help'\n", appname)
	return nil
}

// slack : Commands for slack
func slack(c *cli.Context) error {
	s := getAugs(c).getCfg().initSlackContainer()
	if c.Bool("channellist") {
		s.slackGetChannels().slackDispChannel()
		return nil
	}
	if c.Bool("filelist") || c.Bool("filelistasjson") {
		s.slackGetFileList().slackOutFilelist()
		return nil
	}
	if len(c.String("getfile")) > 0 {
		s.slackGetFile()
		return nil
	}
	if c.Bool("channelhistory") {
		s.slackGetChannelHistory().slackDispChannelHistory()
		return nil
	}
	if (len(c.String("file")) > 0 || len(c.String("content")) > 0) && len(c.String("channel")) > 0 {
		s.slackSubmit().disp()
		return nil
	}
	if len(c.String("deletefile")) > 0 {
		s.slackDeleteFile()
		return nil
	}
	if c.Bool("deletefiles") {
		s.slackGetFileList().slackDeleteAllFiles()
		return nil
	}
	if len(c.String("deletehistory")) > 0 {
		s.slackGetChannelHistory().slackDeleteHistory()
		return nil
	}
	if c.Int("deletehistories") > 0 {
		s.slackGetChannelHistory().slackDeleteChannelAllHistory()
		return nil
	}
	fmt.Printf("Usage is `%s slack --help'\n", appname)
	return nil
}

// doubleSubmit : Submit a file to both Gist and Slack.
func doubleSubmit(c *cli.Context) error {
	if len(c.String("title")) > 0 &&
		len(c.String("file")) > 0 &&
		len(c.String("channel")) > 0 &&
		len(c.String("updateoverwrite")) == 0 &&
		len(c.String("updateadd")) == 0 {
		p := getAugs(c).getCfg()
		res := p.doubleSubmitInit(
			p.initGistContainer().gistSubmitReq(),
			p.initSlackContainer().slackSubmitReq(),
		).doubleSubmitting()
		p.doubleSubmittingDisp(res)
		return nil
	}
	if len(c.String("channel")) > 0 &&
		(len(c.String("title")) > 0 ||
			len(c.String("file")) > 0) &&
		(len(c.String("updateoverwrite")) > 0 ||
			len(c.String("updateadd")) > 0) {
		p := getAugs(c).getCfg()
		res := p.doubleSubmitInit(
			p.initGistContainer().defGistContainer().gistMakeUpdate(),
			p.initSlackContainer().slackSubmitReq(),
		).doubleSubmitting()
		p.doubleSubmittingDisp(res)
		return nil
	}
	fmt.Printf("Usage is `%s doublesubmit --help'\n", appname)
	return nil
}

// getaccesstopen : Rerieves access token from gist and slack.
func getaccesstopen(c *cli.Context) error {
	if len(c.String("gistclientid")) > 0 && len(c.String("gistclientsecret")) > 0 {
		getAugs(c).authInit().getGistAccesstoken().makecfgfile()
		return nil
	}
	if len(c.String("slackclientid")) > 0 && len(c.String("slackclientsecret")) > 0 {
		getAugs(c).authInit().getSlackAccesstoken().makecfgfile()
		return nil
	}
	if c.Bool("chkgisttoken") {
		getAugs(c).authInit().getGistChkToken()
		return nil
	}
	fmt.Printf("Usage is `%s auth --help'\n", appname)
	return nil
}

// useJSON : For contorolling JSON data
func useJSON(c *cli.Context) error {
	i := getAugs(c)
	switch i.jsonControl.Command {
	case "gist":
		j := i.getCfg().keyChk()
		g := j.initGistContainer()
		switch {
		case j.chkArgs("list").(bool) || j.chkArgs("listasjson").(bool):
			g.gistList()
		case j.chkArgs("get").(string) != "" && j.chkArgs("gethistory").(string) == "":
			g.gistGet().disp()
		case j.chkArgs("gethistory").(string) != "" && j.chkArgs("get").(string) == "":
			g.gistGet().disphis()
		case (j.chkArgs("updateoverwrite").(string) == "" || j.chkArgs("updateadd").(string) == "") &&
			(j.chkArgs("title").(string) != "" && j.chkArgs("files").(string) != ""):
			g.defGistContainer().gistSubmit().disp()
		case (j.chkArgs("updateoverwrite").(string) != "" || j.chkArgs("updateadd").(string) != "") &&
			(j.chkArgs("title").(string) != "" || j.chkArgs("files").(string) != ""):
			g.gistUpdate(g.defGistContainer().gistMakeUpdate()).disp()
		case j.chkArgs("getversion").(string) != "" && j.chkArgs("gethistory").(string) == "" && j.chkArgs("get").(string) == "":
			g.gistGet().disp()
		case j.chkArgs("delete").(string) != "":
			g.gistDel()
		case j.chkArgs("deleteall").(bool):
			g.gistDeleteAll()
		default:
			fmt.Println("no parameters")
		}
	case "slack":
		j := i.getCfg().keyChk()
		s := j.initSlackContainer()
		switch {
		case j.chkArgs("channellist").(bool):
			s.slackGetChannels().slackDispChannel()
		case j.chkArgs("filelist").(bool) || j.chkArgs("filelistasjson").(bool):
			s.slackGetFileList().slackOutFilelist()
		case j.chkArgs("getfile").(string) != "":
			s.slackGetFile()
		case j.chkArgs("channelhistory").(bool):
			s.slackGetChannelHistory().slackDispChannelHistory()
		case (j.chkArgs("file").(string) != "" || j.chkArgs("content").(string) != "") &&
			j.chkArgs("channel").(string) != "":
			s.slackSubmit().disp()
		case j.chkArgs("deletefile").(string) != "":
			s.slackDeleteFile()
		case j.chkArgs("deletefiles").(bool):
			s.slackGetFileList().slackDeleteAllFiles()
		case j.chkArgs("deletehistory").(string) != "":
			s.slackGetChannelHistory().slackDeleteHistory()
		case j.chkArgs("deletehistories").(int) > 0:
			s.slackGetChannelHistory().slackDeleteChannelAllHistory()
		}
	case "doublesubmit":
		j := i.getCfg().keyChk()
		switch {
		case j.chkArgs("title").(string) != "" &&
			j.chkArgs("file").(string) != "" &&
			j.chkArgs("channel").(string) != "" &&
			j.chkArgs("updateoverwrite").(string) == "" &&
			j.chkArgs("updateadd").(string) == "":
			res := j.doubleSubmitInit(
				j.initGistContainer().gistSubmitReq(),
				j.initSlackContainer().slackSubmitReq(),
			).doubleSubmitting()
			j.doubleSubmittingDisp(res)
		case j.chkArgs("channel").(string) != "" &&
			(j.chkArgs("title").(string) != "" ||
				j.chkArgs("file").(string) != "") &&
			(j.chkArgs("updateoverwrite").(string) != "" ||
				j.chkArgs("updateadd").(string) != ""):
			res := j.doubleSubmitInit(
				j.initGistContainer().defGistContainer().gistMakeUpdate(),
				j.initSlackContainer().slackSubmitReq(),
			).doubleSubmitting()
			j.doubleSubmittingDisp(res)
		}
	case "auth":
		a := getAugs(c).keyChk().authInit()
		switch {
		case a.chkArgs("gistclientid").(string) != "" && a.chkArgs("gistclientsecret").(string) != "" && a.chkArgs("gistcode").(string) == "":
			a.showCodeURLGist()
		case a.chkArgs("slackclientid").(string) != "" && a.chkArgs("slackclientsecret").(string) != "" && a.chkArgs("slackcode").(string) == "":
			a.showCodeURLSlack()
		case a.chkArgs("gistcode").(string) != "":
			a.getGistAccesstokenJSON().makecfgfile()
		case a.chkArgs("slackcode").(string) != "":
			a.getSlackAccesstokenJSON().makecfgfile()
		}
	case "appcheck":
		switch {
		case getAugs(c).chkArgs("appcheck").(bool):
			fmt.Println("ok")
		}
	}
	return nil
}

// chkArgs : Check args
func (i *iniparamsContainer) chkArgs(key string) interface{} {
	value, ok := i.jsonControl.Options[key]
	if ok {
		return value
	}
	return nil
}

// chkExist : For controlling by JSON
func chkExist(data interface{}) bool {
	var res bool
	switch data.(type) {
	case string:
		if len(data.(string)) == 0 {
			res = false
		} else {
			res = true
		}
	case int:
		if data.(int) == 0 {
			res = false
		} else {
			res = true
		}
	}
	return res
}

// disp : Display results for Gist
func (g *gistContainer) disp() error {
	var result []byte
	if g.jsonControl.Options["jsonparser"].(bool) {
		result, _ = json.MarshalIndent(g.GistGetList, "", "  ")
	} else {
		result, _ = json.Marshal(g.GistGetList)
	}
	fmt.Println(string(result))
	return nil
}

// disphis : Display history of a gist
func (g *gistContainer) disphis() error {
	var result []byte
	if g.jsonControl.Options["jsonparser"].(bool) {
		result, _ = json.MarshalIndent(g.GistGetList[0].History, "", "  ")
	} else {
		result, _ = json.Marshal(g.GistGetList[0].History)
	}
	fmt.Println(string(result))
	return nil
}

// disp : Display results for Slack
func (s *slackContainer) disp() error {
	if !s.jsonControl.Options["simpleresult"].(bool) {
		var result []byte
		if s.jsonControl.Options["jsonparser"].(bool) {
			result, _ = json.MarshalIndent(s.slackParams.SlackFileList, "", "  ")
		} else {
			result, _ = json.Marshal(s.slackParams.SlackFileList)
		}
		fmt.Println(string(result))
	}
	return nil
}

// commandNotFound :
func commandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "'%s' is not a %s command. Check '%s --help' or '%s -h'.", command, c.App.Name, c.App.Name, c.App.Name)
	os.Exit(2)
}
