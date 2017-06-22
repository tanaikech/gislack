// Package main (auth.go) :
// Get access token
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tanaikech/getcode"
	"github.com/tanaikech/gislack/utl"
)

// authContainer : Authorization container
type authContainer struct {
	AuthURL string
	Scopes  []string
	Port    int
}

// gistAccesstoken : Access token for gist
type gistAccesstoken struct {
	Accesstoken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

// slackAccesstoken : Access token for Slack
type slackAccesstoken struct {
	Ok          string `json:"ok,omitempty"`
	Accesstoken string `json:"access_token,omitempty"`
	Scope       string `json:"scope,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	TeamName    string `json:"team_name,omitempty"`
	TeamID      string `json:"team_id,omitempty"`
	Error       string `json:"error,omitempty"`
}

// gislackCfg : Data of CFG file
type gislackCfg struct {
	Gist struct {
		ClientID        string `json:"client_id,omitempty"`
		ClientSecret    string `json:"client_secret,omitempty"`
		GistAccesstoken gistAccesstoken
	} `json:"gist,omitempty"`
	Slack struct {
		ClientID         string `json:"client_id,omitempty"`
		ClientSecret     string `json:"client_secret,omitempty"`
		SlackAccesstoken slackAccesstoken
	} `json:"slack,omitempty"`
}

// authParams : Parameters for authorization process
type authParams struct {
	WorkDir    string
	CfgDir     string
	pstart     time.Time
	GislackCfg gislackCfg
}

// iniparamsContainer : Initial parameters
type iniparamsContainer struct {
	*authParams
	*jsonControl
}

// authErrGist : Error messages for authorization
type authErrGist struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// authInit : Initialize authorization process
func (i *iniparamsContainer) authInit() *iniparamsContainer {
	if cfgdata, err := ioutil.ReadFile(filepath.Join(i.CfgDir, cfgFile)); err == nil {
		err = json.Unmarshal(cfgdata, &i.authParams.GislackCfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Format error of '%s'. ", cfgFile)
			os.Exit(1)
		}
	}
	if len(i.jsonControl.Options["gistclientid"].(string)) > 0 && len(i.jsonControl.Options["gistclientsecret"].(string)) > 0 {
		i.authParams.GislackCfg.Gist.ClientID = i.jsonControl.Options["gistclientid"].(string)
		i.authParams.GislackCfg.Gist.ClientSecret = i.jsonControl.Options["gistclientsecret"].(string)
	}
	if len(i.jsonControl.Options["slackclientid"].(string)) > 0 && len(i.jsonControl.Options["slackclientsecret"].(string)) > 0 {
		i.authParams.GislackCfg.Slack.ClientID = i.jsonControl.Options["slackclientid"].(string)
		i.authParams.GislackCfg.Slack.ClientSecret = i.jsonControl.Options["slackclientsecret"].(string)
	}
	return i
}

// makecfgfile : Make a configuration file.
func (i *iniparamsContainer) makecfgfile() {
	btok, _ := json.MarshalIndent(i.GislackCfg, "", "\t")
	ioutil.WriteFile(filepath.Join(i.CfgDir, cfgFile), btok, 0777)
	fmt.Println("Done.")
}

// getGistAccesstoken : Get access token for using gist APIs.
func (i *iniparamsContainer) getGistAccesstoken() *iniparamsContainer {
	a := &authContainer{
		AuthURL: gistauthcode,
		Scopes:  []string{"gist", "repo"},
		Port:    i.jsonControl.Options["port"].(int),
	}
	codepara := url.Values{}
	codepara.Set("client_id", i.GislackCfg.Gist.ClientID)
	codepara.Set("scope", strings.Join(a.Scopes, " "))
	a.AuthURL = a.AuthURL + codepara.Encode()
	code := getcode.Init(a.AuthURL, a.Port, 30, true, false).Do()
	return i.getGistAccesstokenDo(code)
}

// getGistAccesstokenDo : Retrieve access token for gist
func (i *iniparamsContainer) getGistAccesstokenDo(code string) *iniparamsContainer {
	tokenparams := url.Values{}
	tokenparams.Set("client_id", i.GislackCfg.Gist.ClientID)
	tokenparams.Set("client_secret", i.GislackCfg.Gist.ClientSecret)
	tokenparams.Set("code", code)
	r := &utl.RequestParams{
		Method:       "POST",
		APIURL:       gistaccesstoken,
		Data:         strings.NewReader(tokenparams.Encode()),
		AcceptHeader: "application/json",
		Contenttype:  "application/x-www-form-urlencoded",
		Dtime:        10,
	}
	body, err := r.FetchAPI()
	var errmsg authErrGist
	if json.Unmarshal(body, &errmsg); err != nil || errmsg.Error != "" {
		fmt.Fprintf(os.Stderr, "Error: [ %v ] - Code is wrong. ",
			func(a error, b string) interface{} {
				if a != nil {
					return a
				}
				if b != "" {
					return b
				}
				return nil
			}(err, errmsg.ErrorDescription))
		os.Exit(1)
	}
	json.Unmarshal(body, &i.GislackCfg.Gist.GistAccesstoken)
	return i
}

// getSlackAccesstoken : Get access token for using Slack APIs.
func (i *iniparamsContainer) getSlackAccesstoken() *iniparamsContainer {
	a := &authContainer{
		AuthURL: slackauthcode,
		Scopes:  []string{"channels:history", "channels:read", "chat:write:user", "files:read", "files:write:user"},
		Port:    i.jsonControl.Options["port"].(int),
	}
	codepara := url.Values{}
	codepara.Set("client_id", i.GislackCfg.Slack.ClientID)
	codepara.Set("scope", strings.Join(a.Scopes, " "))
	a.AuthURL = a.AuthURL + codepara.Encode()
	code := getcode.Init(a.AuthURL, a.Port, 30, true, false).Do()
	return i.getSlackAccesstokenDo(code)
}

// getSlackAccesstokenDo : Retrieve access token for slack
func (i *iniparamsContainer) getSlackAccesstokenDo(code string) *iniparamsContainer {
	tokenparams := url.Values{}
	tokenparams.Set("client_id", i.GislackCfg.Slack.ClientID)
	tokenparams.Set("client_secret", i.GislackCfg.Slack.ClientSecret)
	tokenparams.Set("code", code)
	r := &utl.RequestParams{
		Method:      "GET",
		APIURL:      slackaccesstoken + tokenparams.Encode(),
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	json.Unmarshal(body, &i.GislackCfg.Slack.SlackAccesstoken)
	if err != nil || i.GislackCfg.Slack.SlackAccesstoken.Error != "" {
		fmt.Fprintf(os.Stderr, "Error: [ %v ] - Code is wrong. ",
			func(a error, b string) interface{} {
				if a != nil {
					return a
				}
				if b != "" {
					return b
				}
				return nil
			}(err, i.GislackCfg.Slack.SlackAccesstoken.Error))
		os.Exit(1)
	}
	return i
}

// showCodeURLGist : Show URL for retrieving authorization code for gist. This is for controlling by JSON.
func (i *iniparamsContainer) showCodeURLGist() {
	a := &authContainer{
		AuthURL: gistauthcode,
		Scopes:  []string{"gist", "repo"},
		Port:    i.jsonControl.Options["port"].(int),
	}
	codepara := url.Values{}
	codepara.Set("client_id", i.GislackCfg.Gist.ClientID)
	codepara.Set("scope", strings.Join(a.Scopes, " "))
	fmt.Printf("%s", a.AuthURL+codepara.Encode())
}

// getGistAccesstokenJSON : Retrieve access token of gist using JSON data
func (i *iniparamsContainer) getGistAccesstokenJSON() *iniparamsContainer {
	return i.getGistAccesstokenDo(i.jsonControl.Options["gistcode"].(string))
}

// showCodeURLSlack : Show URL for retrieving authorization code for slack. This is for controlling by JSON.
func (i *iniparamsContainer) showCodeURLSlack() {
	a := &authContainer{
		AuthURL: slackauthcode,
		Scopes:  []string{"channels:history", "channels:read", "chat:write:user", "files:read", "files:write:user"},
		Port:    i.jsonControl.Options["port"].(int),
	}
	codepara := url.Values{}
	codepara.Set("client_id", i.GislackCfg.Slack.ClientID)
	codepara.Set("scope", strings.Join(a.Scopes, " "))
	fmt.Printf("%s", a.AuthURL+codepara.Encode())
}

// getSlackAccesstokenJSON : Retrieve access token of slack using JSON data
func (i *iniparamsContainer) getSlackAccesstokenJSON() *iniparamsContainer {
	return i.getSlackAccesstokenDo(i.jsonControl.Options["slackcode"].(string))
}

// getGistChkToken : Check the condition of github access token
func (i *iniparamsContainer) getGistChkToken() {
	para := url.Values{}
	para.Set("client_id", i.GislackCfg.Gist.ClientID)
	para.Set("client_secret", i.GislackCfg.Gist.ClientSecret)
	r := &utl.RequestParams{
		Method:      "GET",
		APIURL:      gistchktoken + para.Encode(),
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       10,
	}
	res, err := r.FetchAPIres()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: [ %v ] - Code is wrong. ", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	headobj := map[string]interface{}{}
	headobj["MaxLimit"] = res.Header.Get("X-Ratelimit-Limit")
	headobj["Remaining"] = res.Header.Get("X-RateLimit-Remaining")
	resetreq0, _ := strconv.ParseInt(res.Header.Get("X-RateLimit-Reset"), 10, 64)
	resetreq1 := time.Unix(resetreq0, 0).Format("20060102_15:04:05")
	headobj["ResetTime"] = resetreq1
	result, _ := json.MarshalIndent(headobj, "", "  ")
	fmt.Println(string(result))
	return
}
