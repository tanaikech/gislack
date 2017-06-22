// Package main (materials_double.go) :
// Materials for doubleSubmit.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tanaikech/gislack/utl"
)

// doubleParam : Parameter for double submissions
type doubleParam struct {
	GistS       *utl.RequestParams
	SlackS      *utl.RequestParams
	Pstart      time.Time
	JSONControl *jsonControl
}

// doubleResults : Results from double submissions
type doubleResults struct {
	Gist    gistGetList   `json:"gist_response"`
	Slack   slackFileList `json:"slack_response"`
	TotalEt float64       `json:"TotalElapsedTime,omitempty"`
}

// doubleSubmitInit : Initialize doubleParam
func (p *iniparamsContainer) doubleSubmitInit(g, s *utl.RequestParams) *doubleParam {
	return &doubleParam{
		GistS:       g,
		SlackS:      s,
		Pstart:      p.pstart,
		JSONControl: p.jsonControl,
	}
}

// doubleSubmitting : Do double submissions under parallel process
func (d *doubleParam) doubleSubmitting() map[string]interface{} {
	var wg sync.WaitGroup
	submit := make(chan *utl.RequestParams, 2)
	workers := 2
	var res [][]byte
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, submit chan *utl.RequestParams) {
			defer wg.Done()
			for {
				p, done := <-submit
				if !done {
					return
				}
				body, err := p.FetchAPI()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %+v, %s\n", string(body), err)
					return
				}
				res = append(res, body)
			}
		}(&wg, submit)
	}
	submit <- d.GistS
	submit <- d.SlackS
	close(submit)
	wg.Wait()
	return map[string]interface{}{
		"r1": res,
		"r2": d.Pstart,
	}
}

// gistSubmit : Request to Gist
func (g *gistContainer) gistSubmitReq() *utl.RequestParams {
	var err error
	g.GistPayload.Description = g.jsonControl.Options["title"].(string)
	g.GistPayload.Public = g.jsonControl.Options["public"].(bool)
	if err != nil {
		panic(err)
	}
	g.GistPayload.Files = func(f string) map[string]interface{} {
		obj := map[string]interface{}{}
		var fpath string
		if filepath.Dir(f) == "." {
			fpath = filepath.Join(g.workdir, f)
		} else {
			fpath = f
		}
		data, err := ioutil.ReadFile(fpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		value := map[string]interface{}{
			"content": string(data),
		}
		obj[filepath.Base(f)] = value
		return obj
	}(g.jsonControl.Options["file"].(string))
	payload, _ := json.Marshal(g.GistPayload)
	r := &utl.RequestParams{
		Method:      "POST",
		APIURL:      gisturl,
		Data:        bytes.NewBuffer(payload),
		Contenttype: "application/x-www-form-urlencoded",
		Accesstoken: g.Accesstoken,
		Dtime:       10,
	}
	return r
}

// slackSubmit : Request to Slack
func (s *slackContainer) slackSubmitReq() *utl.RequestParams {
	s.slackParams.SlackPayload.Title = s.jsonControl.Options["title"].(string)
	s.slackParams.SlackPayload.Filetype = s.jsonControl.Options["filetype"].(string)
	s.slackParams.SlackPayload.Channels = s.slackGetChannels().slackChannelNameToID()
	s.slackParams.SlackPayload.InitialComment = s.jsonControl.Options["initialcomment"].(string)
	s.slackParams.SlackPayload.Filename = s.jsonControl.Options["file"].(string)
	var file string
	if filepath.Dir(s.slackParams.SlackPayload.Filename) == "." {
		file = filepath.Join(s.workdir, s.slackParams.SlackPayload.Filename)
	} else {
		file = s.slackParams.SlackPayload.Filename
	}
	p := url.Values{}
	p.Set("token", s.Token)
	p.Set("channels", s.slackParams.SlackPayload.Channels)
	p.Set("title", s.slackParams.SlackPayload.Title)
	p.Set("filetype", s.slackParams.SlackPayload.Filetype)
	p.Set("initial_comment", s.slackParams.SlackPayload.InitialComment)
	p.Set("filename", filepath.Base(file))
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part := make(textproto.MIMEHeader)
	data, err := w.CreatePart(part)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		os.Exit(1)
	}
	fs, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		os.Exit(1)
	}
	defer fs.Close()
	data, err = w.CreateFormFile("file", file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		os.Exit(1)
	}
	if _, err = io.Copy(data, fs); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		os.Exit(1)
	}
	w.Close()
	r := &utl.RequestParams{
		Method:      "POST",
		APIURL:      slackurl + "files.upload?" + p.Encode(),
		Data:        &b,
		Contenttype: w.FormDataContentType(),
		Accesstoken: s.Token,
		Dtime:       10,
	}
	return r
}

// doubleSubmittingDisp : Display results
func (p *iniparamsContainer) doubleSubmittingDisp(ar map[string]interface{}) {
	res, _ := ar["r1"].([][]byte)
	d := &doubleResults{}
	if json.Unmarshal(res[0], &d.Slack); d.Slack.OK {
		json.Unmarshal(res[1], &d.Gist)
	} else {
		json.Unmarshal(res[0], &d.Gist)
		json.Unmarshal(res[1], &d.Slack)
	}
	d.Gist.CreatedAt = d.Gist.CreatedAt.In(time.Local)
	d.Gist.UpdatedAt = d.Gist.UpdatedAt.In(time.Local)
	d.Slack.File.CreatedTime = time.Unix(d.Slack.File.Created, 0)
	if p.jsonControl.Options["simpleresult"].(bool) {
		d.simpleResult()
		return
	}
	pstime, _ := ar["r2"].(time.Time)
	d.TotalEt = math.Trunc(time.Now().Sub(pstime).Seconds()*1000) / 1000
	var result []byte
	if p.jsonControl.Options["jsonparser"].(bool) {
		result, _ = json.MarshalIndent(d, "", "  ")
	} else {
		result, _ = json.Marshal(d)
	}
	fmt.Println(string(result))
	return
}

// simpleResult : Display simple results
func (d *doubleResults) simpleResult() {
	fmt.Printf(
		"{\"gist_created_at\": \"%s\", \"gist_id\": \"%s\", \"slack_created_at\": \"%s\", \"slack_id\": \"%s\"}",
		d.Gist.CreatedAt.Format("20060102_15:04:05"),
		d.Gist.ID,
		d.Slack.File.CreatedTime.Format("20060102_15:04:05"),
		d.Slack.File.ID,
	)
	return
}
