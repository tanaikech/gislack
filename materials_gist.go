// Package main (materials_gist.go) :
// Materials for gist.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/tanaikech/gislack/utl"
)

// gistParams : Parameters for Gist
type gistParams struct {
	Accesstoken string
	GistPayload gistPayload
	GistGetList []gistGetList
}

// gistContainer : Container included parameters
type gistContainer struct {
	*initVal
	*gistParams
	*jsonControl
}

// gistGetList : Struct for get list to Gist
type gistGetList struct {
	ID          string                 `json:"id,omitempty"`
	HTMLURL     string                 `json:"html_url,omitempty"`
	Files       map[string]interface{} `json:"files,omitempty"`
	Public      bool                   `json:"public,omitempty"`
	CreatedAt   time.Time              `json:"created_at,omitempty"`
	UpdatedAt   time.Time              `json:"updated_at,omitempty"`
	Description string                 `json:"description,omitempty"`
	History     []struct {
		CommittedAt time.Time `json:"committed_at,omitempty"`
		URL         string    `json:"url,omitempty"`
		Version     string    `json:"version,omitempty"`
	} `json:"history,omitempty"`
	Owner struct {
		Login string `json:"login,omitempty"`
	} `json:"owner,omitempty"`
}

// gistPayload : Payload for requesting to Gist
type gistPayload struct {
	Description string                 `json:"description,omitempty"`
	Public      bool                   `json:"public"`
	Files       map[string]interface{} `json:"files,omitempty"`
}

// gistfiles : Files for submitting and updating files to Gist
type gistfiles struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
	Language string `json:"language"`
	RawURL   string `json:"raw_url"`
	Size     int    `json:"size"`
}

// initGistContainer : Initialize parameters for Gist
func (i *iniparamsContainer) initGistContainer() *gistContainer {
	g := &gistContainer{
		&initVal{
			pstart: i.authParams.pstart,
		},
		&gistParams{
			Accesstoken: i.authParams.GislackCfg.Gist.GistAccesstoken.Accesstoken,
		},
		&jsonControl{},
	}
	g.initVal.workdir = i.authParams.WorkDir
	g.jsonControl = i.jsonControl
	if len(g.Accesstoken) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Access token of GitHub is NOT found. Please retrieve Client ID and Client Secret from GitHub, and run 'gislack auth -gi clientid -gs clientsecret'.\n")
		os.Exit(1)
	}
	if g.jsonControl.Command == "doublesubmit" {
		g.jsonControl.Options["files"] = g.jsonControl.Options["file"]
		g.jsonControl.Options["filenames"] = g.jsonControl.Options["filename"]
	}
	return g
}

// defGistContainer : Initialize a container for Gist
func (g *gistContainer) defGistContainer() *gistContainer {
	g.GistPayload.Description = g.jsonControl.Options["title"].(string)
	g.GistPayload.Public = g.jsonControl.Options["public"].(bool)
	if files := g.jsonControl.Options["files"].(string); len(files) > 0 {
		filenames := g.jsonControl.Options["filenames"].(string)
		filesAr := strings.Split(files, ",")
		filenamesAr := strings.Split(filenames, ",")
		if (len(filesAr) == len(filenamesAr) || len(filesAr) < len(filenamesAr)) && len(filenames) > 0 {
			g.GistPayload.Files = func(files, filenames []string) map[string]interface{} {
				obj := map[string]interface{}{}
				for i, file := range files {
					fns := strings.TrimSpace(filenames[i])
					e := strings.TrimSpace(file)
					var fpath string
					if filepath.Dir(e) == "." {
						fpath = filepath.Join(g.workdir, e)
					} else {
						fpath = e
					}
					data, err := ioutil.ReadFile(fpath)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error: %s\n", err)
						os.Exit(1)
					}
					value := map[string]interface{}{
						"content":  string(data),
						"filename": fns,
					}
					obj[fns] = value
				}
				return obj
			}(filesAr, filenamesAr)
		}
		if len(filesAr) > len(filenamesAr) || len(filenames) == 0 {
			g.GistPayload.Files = func(files []string) map[string]interface{} {
				obj := map[string]interface{}{}
				for _, file := range files {
					e := strings.TrimSpace(file)
					var fpath string
					if filepath.Dir(e) == "." {
						fpath = filepath.Join(g.workdir, e)
					} else {
						fpath = e
					}
					data, err := ioutil.ReadFile(fpath)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error: %s\n", err)
						os.Exit(1)
					}
					value := map[string]interface{}{
						"content": string(data),
					}
					obj[filepath.Base(e)] = value
				}
				return obj
			}(filesAr)
		}
	}
	return g
}

// gistList : Retrieve file list
func (g *gistContainer) gistList() {
	r := &utl.RequestParams{
		Method:      "GET",
		APIURL:      gisturl,
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Accesstoken: g.Accesstoken,
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(body, &g.GistGetList)
	if len(g.GistGetList) > 0 {
		if g.jsonControl.Options["list"].(bool) && !g.jsonControl.Options["listasjson"].(bool) {
			buffer := &bytes.Buffer{}
			w := new(tabwriter.Writer)
			w.Init(buffer, 0, 4, 1, ' ', 0)
			for i, e := range g.GistGetList {
				g.GistGetList[i].CreatedAt = e.CreatedAt.In(time.Local)
				g.GistGetList[i].UpdatedAt = e.UpdatedAt.In(time.Local)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					e.Description,
					g.GistGetList[i].UpdatedAt.Format("20060102_15:04:05"),
					func(p bool) string {
						if p {
							return "Public"
						}
						return "Secret"
					}(e.Public),
					e.ID,
				)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "# Description", "# Updated time", "# Public", "# id")
			w.Flush()
			s := bufio.NewScanner(buffer)
			var ar []string
			for s.Scan() {
				ar = append(ar, s.Text())
			}
			for i := len(ar) - 1; i >= 0; i-- {
				fmt.Println(ar[i])
			}
		}
		if g.jsonControl.Options["listasjson"].(bool) {
			for i, e := range g.GistGetList {
				g.GistGetList[i].CreatedAt = e.CreatedAt.In(time.Local)
				g.GistGetList[i].UpdatedAt = e.UpdatedAt.In(time.Local)
			}
			listjson, _ := json.MarshalIndent(g.GistGetList, "", "  ")
			fmt.Println(string(listjson))
		}
	} else {
		fmt.Println("No gists.")
	}
	return
}

// gistGet : Retrieve a gist from ID
func (g *gistContainer) gistGet() *gistContainer {
	var gistid string
	if g.jsonControl.Options["get"].(string) != "" {
		gistid = g.jsonControl.Options["get"].(string)
	}
	if g.jsonControl.Options["gethistory"].(string) != "" {
		gistid = g.jsonControl.Options["gethistory"].(string)
	}
	if g.jsonControl.Options["getversion"].(string) != "" {
		gistid = strings.Replace(g.jsonControl.Options["getversion"].(string), gisturl+"/", "", 1)
	}
	if gistid != "" {
		g.gistGetMain(gistid)
	} else {
		fmt.Fprintf(os.Stderr, "Error: Gist ID was not found.\n")
		os.Exit(1)
	}
	if g.jsonControl.Options["usejsoncontrol"].(bool) ||
		g.jsonControl.Options["gethistory"].(string) != "" {
		return g
	}
	for i, e := range g.GistGetList[0].Files {
		file, _ := json.Marshal(e)
		obj := map[string]interface{}{}
		json.Unmarshal(file, &obj)
		outfile := strings.TrimSpace(obj["filename"].(string))
		if _, err := os.Stat(outfile); err == nil {
			fmt.Fprintf(os.Stderr, "Error: %s already exists. Content was not saved to a file.\n", outfile)
		} else {
			ioutil.WriteFile(filepath.Join(g.workdir, outfile), []byte(obj["content"].(string)), 0777)
			obj["content"] = fmt.Sprintf("Content was saved to a file (%s).", outfile)
		}
		g.GistGetList[0].Files[i] = obj
	}
	return g
}

// gistGetMain : Main method for retrieving a gist from ID
func (g *gistContainer) gistGetMain(id string) *gistContainer {
	r := &utl.RequestParams{
		Method:      "GET",
		APIURL:      gisturl + "/" + id,
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Accesstoken: g.Accesstoken,
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: gist ID '%s' was not found.\n%v\n", id, err)
		os.Exit(1)
	}
	var gg gistGetList
	json.Unmarshal(body, &gg)
	gg.CreatedAt = gg.CreatedAt.In(time.Local)
	gg.UpdatedAt = gg.UpdatedAt.In(time.Local)
	for i, e := range gg.History {
		gg.History[i].CommittedAt = e.CommittedAt.In(time.Local)
	}
	g.GistGetList = append(g.GistGetList, gg)
	return g
}

// gistUpdate : Update Gist
func (g *gistContainer) gistUpdate(r *utl.RequestParams) *gistContainer {
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n%s\n", err, string(body))
		os.Exit(1)
	}
	var p gistGetList
	json.Unmarshal(body, &p)
	g.GistGetList = append(g.GistGetList, p)
	return g
}

// gistUpdate : Update Gist using ID
func (g *gistContainer) gistMakeUpdate() *utl.RequestParams {
	var udURL string
	if len(g.jsonControl.Options["updateoverwrite"].(string)) > 0 && len(g.jsonControl.Options["updateadd"].(string)) == 0 {
		udURL = gisturl + "/" + g.jsonControl.Options["updateoverwrite"].(string)
		var ufiles string
		if len(g.jsonControl.Options["filenames"].(string)) > 0 {
			ufiles = g.jsonControl.Options["filenames"].(string)
		} else {
			ufiles = g.jsonControl.Options["files"].(string)
		}
		g.gistGetMain(g.jsonControl.Options["updateoverwrite"].(string))
		upfiles := strings.Split(ufiles, ",")
		for _, e1 := range g.gistParams.GistGetList {
			for il := range e1.Files {
				var f bool
				for _, e2 := range upfiles {
					if il == strings.TrimSpace(e2) {
						f = true
					}
				}
				if !f {
					g.GistPayload.Files[il] = "null"
				}
			}
		}
	}
	if len(g.jsonControl.Options["updateoverwrite"].(string)) == 0 && len(g.jsonControl.Options["updateadd"].(string)) > 0 {
		udURL = gisturl + "/" + g.jsonControl.Options["updateadd"].(string)
	}
	payload, _ := json.Marshal(g.GistPayload)
	payloadStr := strings.Replace(string(payload), "\"null\"", "null", -1)
	r := &utl.RequestParams{
		Method:      "PATCH",
		APIURL:      udURL,
		Data:        strings.NewReader(payloadStr),
		Contenttype: "application/x-www-form-urlencoded",
		Accesstoken: g.Accesstoken,
		Dtime:       10,
	}
	return r
}

// gistSubmit : Submit files to Gist
func (g *gistContainer) gistSubmit() *gistContainer {
	if len(g.gistParams.Accesstoken) == 0 && !g.jsonControl.Options["anonymous"].(bool) {
		fmt.Fprintf(os.Stderr, "Warning: You have no access token. If you want to submit as anonymous, please run again by using option '--anonymous'.\n")
		os.Exit(1)
	}
	if g.jsonControl.Options["anonymous"].(bool) {
		g.gistParams.Accesstoken = ""
	}
	payload, _ := json.Marshal(g.GistPayload)
	r := &utl.RequestParams{
		Method:      "POST",
		APIURL:      gisturl,
		Data:        bytes.NewBuffer(payload),
		Contenttype: "application/x-www-form-urlencoded",
		Accesstoken: g.Accesstoken,
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s, %s\n", err, string(body))
		os.Exit(1)
	}
	var p gistGetList
	json.Unmarshal(body, &p)
	if len(p.Owner.Login) == 0 {
		p.Owner.Login = "### This was submitted as anonymous. ###"
	}
	p.CreatedAt = p.CreatedAt.In(time.Local)
	p.UpdatedAt = p.UpdatedAt.In(time.Local)
	g.GistGetList = append(g.GistGetList, p)
	return g
}

// gistDel : Delete a gist
func (g *gistContainer) gistDel() {
	r := &utl.RequestParams{
		Method:      "DELETE",
		APIURL:      gisturl + "/" + g.jsonControl.Options["delete"].(string),
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Accesstoken: g.Accesstoken,
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n%s\n", string(body), err)
		os.Exit(1)
	}
	if len(body) == 0 {
		fmt.Println("Done.")
	} else {
		fmt.Println(string(body))
	}
	return
}

// gistDeleteAll : Delete all gists. When you use this command, please be careful.
func (g *gistContainer) gistDeleteAll() {
	r := &utl.RequestParams{
		Method:      "GET",
		APIURL:      gisturl,
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Accesstoken: g.Accesstoken,
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(body, &g.GistGetList)
	if len(g.GistGetList) > 0 {
		var input string
		fmt.Printf("These are %s's Gists.\n[WARNING] Will you delete all gists? [y or n] ... ", g.GistGetList[0].Owner.Login)
		if _, err := fmt.Scan(&input); err != nil {
			log.Fatalf("Error: %v.\n", err)
		}
		if input == "y" {
			bar := pb.StartNew(len(g.GistGetList))
			for _, e := range g.GistGetList {
				bar.Increment()
				r := &utl.RequestParams{
					Method:      "DELETE",
					APIURL:      gisturl + "/" + e.ID,
					Data:        nil,
					Contenttype: "application/x-www-form-urlencoded",
					Accesstoken: g.Accesstoken,
					Dtime:       10,
				}
				_, err := r.FetchAPI()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %s\n", err)
					os.Exit(1)
				}
			}
			bar.FinishPrint("Done.")
		} else {
			fmt.Printf("%s's Gists were not deleted.", g.GistGetList[0].Owner.Login)
		}
	} else {
		fmt.Println("No gists.")
	}
}
