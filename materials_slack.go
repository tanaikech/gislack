// Package main (materials_slack.go) :
// Materials for slack.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/tanaikech/gislack/utl"
)

// slackParams : Parameters for Slack
type slackParams struct {
	Token          string
	Channel        string
	SlackPayload   slackPayload
	SlackDelFile   slackDelFile
	SlackFile      slackFile
	SlackFilesList slackFilesList
	SlackFileList  slackFileList
	ChannelHistory channelHistory
	ChannelList    channelList
}

// slackContainer : Container included parameters
type slackContainer struct {
	*initVal
	*slackParams
	*jsonControl
}

// channelList : Channel list
type channelList struct {
	Channels []channelar `json:"channels"`
}

// channelar :
type channelar struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Creator string `json:"creator"`
}

// slackFilesList : File list for Slack
type slackFilesList struct {
	Files  []slackfiles `json:"files"`
	Paging struct {
		Count int `json:"count"`
		Total int `json:"total"`
		Page  int `json:"page"`
		Pages int `json:"pages"`
	} `json:"paging"`
}

// slackfiles : Struct for a file
type slackfiles struct {
	ID          string    `json:"id"`
	Created     int64     `json:"created"`
	CreatedTime time.Time `json:"createdtime"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Filetype    string    `json:"filetype"`
	User        string    `json:"user"`
	Channels    []string  `json:"channels"`
}

// channelHistory : Channel histories
type channelHistory struct {
	Latest   string `json:"latest"`
	Messages []struct {
		Type     string `json:"type"`
		User     string `json:"user"`
		Username string `json:"username"`
		Text     string `json:"text"`
		Ts       string `json:"ts"`
	} `json:"messages"`
	HasMore bool `json:"has_more"`
}

// slackPayload : Payload for requesting to Slack
type slackPayload struct {
	Content        string `json:"content,omitempty"`
	Filename       string `json:"filename,omitempty"`
	Filetype       string `json:"filetype,omitempty"`
	Title          string `json:"title,omitempty"`
	InitialComment string `json:"initial_comment,omitempty"`
	Channels       string `json:"channels,omitempty"`
}

// slackInputJSON : Struct for submitting using JSON data
type slackInputJSON struct {
	File           string `json:"file,omitempty"`
	Content        string `json:"content,omitempty"`
	Filetype       string `json:"filetype,omitempty"`
	Filename       string `json:"filename,omitempty"`
	Title          string `json:"title,omitempty"`
	InitialComment string `json:"initial_comment,omitempty"`
	Channels       string `json:"channels,omitempty"`
}

// slackFileList : File list
type slackFileList struct {
	OK   bool `json:"ok,omitempty"`
	File struct {
		ID          string    `json:"id,omitempty"`
		Created     int64     `json:"created,omitempty"`
		CreatedTime time.Time `json:"createdtime,omitempty"`
		Name        string    `json:"name,omitempty"`
		Title       string    `json:"title,omitempty"`
		Filetype    string    `json:"filetype,omitempty"`
		User        string    `json:"user,omitempty"`
		Channels    []string  `json:"channels,omitempty"`
	} `json:"file,omitempty"`
}

// slackDelFile : Struct for deleting files
type slackDelFile struct {
	File string `json:"file"`
}

// slackError : Errors from Slack
type slackError struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// slackFile : File of Slack
type slackFile struct {
	OK   bool `json:"ok"`
	File struct {
		Name        string    `json:"name,omitempty"`
		MimeType    string    `json:"mimetype,omitempty"`
		Title       string    `json:"title,omitempty"`
		Created     int64     `json:"created,omitempty"`
		CreatedTime time.Time `json:"createdtime,omitempty"`
	} `json:"file,omitempty"`
	Content string `json:"content,omitempty"`
}

// slackChkAt : For checking slack access token
type slackChkAt struct {
	OK   bool   `json:"ok"`
	URL  string `json:"url,omitempty"`
	Team string `json:"team,omitempty"`
	User string `json:"user,omitempty"`
}

// initSlackContainer : Initialize parameters for Slack
func (i *iniparamsContainer) initSlackContainer() *slackContainer {
	s := &slackContainer{
		&initVal{
			pstart: i.authParams.pstart,
		},
		&slackParams{
			Token: i.authParams.GislackCfg.Slack.SlackAccesstoken.Accesstoken,
		},
		&jsonControl{},
	}
	s.initVal.workdir = i.authParams.WorkDir
	s.jsonControl = i.jsonControl
	if len(s.Token) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Access token of Slack is NOT found. Please retrieve Client ID and Client Secret from Slack, and run 'gislack auth -si clientid -ss clientsecret'.\n")
		os.Exit(1)
	}
	return s
}

// slackGetChannels : Retrieve channel list
func (s *slackContainer) slackGetChannels() *slackContainer {
	p := url.Values{}
	p.Set("token", s.slackParams.Token)
	r := &utl.RequestParams{
		Method:      "GET",
		APIURL:      slackurl + "channels.list?" + p.Encode(),
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	json.Unmarshal(body, &s.ChannelList)
	return s
}

// slackDispChannel : Display retrieved channale list
func (s *slackContainer) slackDispChannel() {
	if len(s.ChannelList.Channels) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No channels.\n")
		os.Exit(1)
	} else {
		buffer := &bytes.Buffer{}
		w := new(tabwriter.Writer)
		w.Init(buffer, 0, 4, 1, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\t%s\n", "# channelname", "# channalID", "# creator")
		for _, e := range s.ChannelList.Channels {
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				e.Name,
				e.ID,
				e.Creator,
			)
		}
		w.Flush()
		fmt.Printf("%s", buffer)
	}
	return
}

// slackGetFileList : Retrieve file list
func (s *slackContainer) slackGetFileList() *slackContainer {
	p := url.Values{}
	p.Set("token", s.slackParams.Token)
	p.Set("channel", s.jsonControl.Options["channel"].(string))
	p.Set("user", s.jsonControl.Options["user"].(string))
	p.Set("count", strconv.FormatInt(100, 10))
	p.Set("page", strconv.Itoa(s.slackParams.SlackFilesList.Paging.Page))
	r := &utl.RequestParams{
		Method:      "POST",
		APIURL:      slackurl + "files.list?" + p.Encode(),
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	var fl slackFilesList
	json.Unmarshal(body, &fl)
	s.slackParams.SlackFilesList.Files = append(s.slackParams.SlackFilesList.Files, fl.Files...)
	if fl.Paging.Page <= fl.Paging.Pages {
		s.slackParams.SlackFilesList.Paging.Page = fl.Paging.Page + 1
		s.slackGetFileList()
	}
	return s
}

// slackGetFiles : Display file list
func (s *slackContainer) slackOutFilelist() {
	if s.jsonControl.Options["filelistasjson"].(bool) {
		s.slackDispFilesJSON()
	} else {
		s.slackDispFiles()
	}
}

// slackDispFilesJSON : Display file list as JSON
func (s *slackContainer) slackDispFilesJSON() {
	ar := s.slackParams.SlackFilesList.Files
	if len(ar) > 0 {
		for i := range s.slackParams.SlackFilesList.Files {
			s.slackParams.SlackFilesList.Files[i].CreatedTime = time.Unix(s.slackParams.SlackFilesList.Files[i].Created, 0)
		}
		result, _ := json.Marshal(s.slackParams.SlackFilesList)
		fmt.Println(string(result))
	} else {
		fmt.Println("No files.")
	}
	return
}

// slackDispFiles : Display file list
func (s *slackContainer) slackDispFiles() {
	ar := s.slackParams.SlackFilesList.Files
	if len(ar) > 0 {
		buffer := &bytes.Buffer{}
		w := new(tabwriter.Writer)
		w.Init(buffer, 0, 4, 1, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", "# Title", "# Created time", "# fileID", "# channel", "# user", "# fileType")
		for i := len(ar) - 1; i >= 0; i-- {
			ar[i].CreatedTime = time.Unix(ar[i].Created, 0)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				ar[i].Title,
				ar[i].CreatedTime.Format("20060102_15:04:05"),
				ar[i].ID,
				ar[i].Channels,
				ar[i].User,
				ar[i].Filetype,
			)
		}
		w.Flush()
		fmt.Printf("%s", buffer)
		fmt.Printf("\n Total : %d", len(ar))
	} else {
		fmt.Println("No files.")
	}
	return
}

// slackGetFile : Get file from file ID.
func (s *slackContainer) slackGetFile() {
	p := url.Values{}
	p.Set("token", s.slackParams.Token)
	p.Set("file", s.jsonControl.Options["getfile"].(string))
	r := &utl.RequestParams{
		Method:      "POST",
		APIURL:      slackurl + "files.info?" + p.Encode(),
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	json.Unmarshal(body, &s.slackParams.SlackFile)
	s.slackParams.SlackFile.File.CreatedTime = time.Unix(s.slackParams.SlackFile.File.Created, 0)

	if s.jsonControl.Options["usejsoncontrol"].(bool) {
		result, _ := json.Marshal(s.slackParams.SlackFile)
		fmt.Println(string(result))
		return
	}
	outfile := strings.TrimSpace(s.slackParams.SlackFile.File.Name)
	if _, err := os.Stat(outfile); err == nil {
		fmt.Fprintf(os.Stderr, "Error: %s already exists. Content was not saved to a file.\n", outfile)
	} else {
		ioutil.WriteFile(filepath.Join(s.workdir, outfile), []byte(s.slackParams.SlackFile.Content), 0777)
		s.slackParams.SlackFile.Content = fmt.Sprintf("Content was saved to a file (%s).", outfile)
	}
	result, _ := json.MarshalIndent(s.slackParams.SlackFile, "", "  ")
	fmt.Println(string(result))
	return
}

// slackChannelNameToID : Convert from name to ID for Slack channel
func (s *slackContainer) slackChannelNameToID() string {
	for _, e := range s.ChannelList.Channels {
		if s.jsonControl.Options["channel"].(string) == e.Name {
			return e.ID
		}
	}
	return fmt.Sprintf("No channel ID for %s.", s.jsonControl.Options["channel"].(string))
}

// slackGetChannelHistory : Retrieve channel histories
func (s *slackContainer) slackGetChannelHistory() *slackContainer {
	if len(s.jsonControl.Options["channel"].(string)) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Please input channel name using '-ch'.\n")
		os.Exit(1)
	}
	p := url.Values{}
	p.Set("token", s.slackParams.Token)
	p.Set("channel", s.slackGetChannels().slackChannelNameToID())
	p.Set("latest", s.ChannelHistory.Latest)
	p.Set("count", strconv.FormatInt(1000, 10))
	r := &utl.RequestParams{
		Method:      "POST",
		APIURL:      slackurl + "channels.history?" + p.Encode(),
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	var ch channelHistory
	json.Unmarshal(body, &ch)
	s.ChannelHistory.Messages = append(s.ChannelHistory.Messages, ch.Messages...)
	if ch.HasMore {
		s.ChannelHistory.Latest = ch.Messages[len(ch.Messages)-1].Ts
		s.slackGetChannelHistory()
	}
	return s
}

// slackDispChannelHistory : Display channel histories
func (s *slackContainer) slackDispChannelHistory() {
	ar := s.ChannelHistory.Messages
	if len(ar) > 0 {
		buffer := &bytes.Buffer{}
		w := new(tabwriter.Writer)
		w.Init(buffer, 0, 4, 1, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "# Created date", "# text", "# user", "# ts(historyID)")
		for i := len(ar) - 1; i >= 0; i-- {
			ut, err := strconv.ParseFloat(ar[i].Ts, 64)
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				time.Unix(int64(ut), 0).Format("20060102_15:04:05"),
				ar[i].Text,
				func(user, username string) string {
					if len(username) > 0 {
						return username
					}
					if len(user) > 0 {
						return user
					}
					return ""
				}(ar[i].User, ar[i].Username),
				ar[i].Ts,
			)
		}
		w.Flush()
		fmt.Printf("%s", buffer)
		fmt.Printf("\n Total : %d", len(ar))
	} else {
		fmt.Println("No history.")
	}
}

// slackSubmit : Submit a file
func (s *slackContainer) slackSubmit() *slackContainer {
	if len(s.jsonControl.Options["file"].(string)) > 0 {
		s.slackParams.SlackPayload.Filename = s.jsonControl.Options["file"].(string)
	}
	if len(s.jsonControl.Options["file"].(string)) == 0 {
		s.slackParams.SlackPayload.Content = s.jsonControl.Options["content"].(string)
	}
	s.slackParams.SlackPayload.Title = s.jsonControl.Options["title"].(string)
	s.slackParams.SlackPayload.Filetype = s.jsonControl.Options["filetype"].(string)
	s.slackParams.SlackPayload.Channels = s.slackGetChannels().slackChannelNameToID()
	s.slackParams.SlackPayload.InitialComment = s.jsonControl.Options["initialcomment"].(string)
	var file string
	if filepath.Dir(s.slackParams.SlackPayload.Filename) == "." {
		file = filepath.Join(s.workdir, s.slackParams.SlackPayload.Filename)
	} else {
		file = s.slackParams.SlackPayload.Filename
	}
	r := &utl.RequestParams{}
	if len(file) > 0 {
		p := url.Values{}
		p.Set("token", s.slackParams.Token)
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
		r = &utl.RequestParams{
			Method:      "POST",
			APIURL:      slackurl + "files.upload?" + p.Encode(),
			Data:        &b,
			Contenttype: w.FormDataContentType(),
			Dtime:       10,
		}
	} else if len(file) == 0 {
		p := url.Values{}
		p.Set("token", s.slackParams.Token)
		p.Set("channels", s.slackParams.SlackPayload.Channels)
		p.Set("title", s.slackParams.SlackPayload.Title)
		p.Set("filetype", s.slackParams.SlackPayload.Filetype)
		p.Set("initial_comment", s.slackParams.SlackPayload.InitialComment)
		p.Set("content", s.slackParams.SlackPayload.Content)
		r = &utl.RequestParams{
			Method:      "POST",
			APIURL:      slackurl + "files.upload?" + p.Encode(),
			Data:        nil,
			Contenttype: "application/x-www-form-urlencoded",
			Dtime:       10,
		}
	}
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		os.Exit(1)
	}
	json.Unmarshal(body, &s.slackParams.SlackFileList)
	s.slackParams.SlackFileList.File.CreatedTime = time.Unix(s.slackParams.SlackFileList.File.Created, 0)
	return s
}

// slackDeleteFile : Delete a file
func (s *slackContainer) slackDeleteFile() {
	p := url.Values{}
	p.Set("token", s.slackParams.Token)
	p.Set("file", s.jsonControl.Options["deletefile"].(string))
	r := &utl.RequestParams{
		Method:      "POST",
		APIURL:      slackurl + "files.delete?" + p.Encode(),
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	var rs map[string]interface{}
	json.Unmarshal(body, &rs)
	if !rs["ok"].(bool) {
		fmt.Println("Error: File has already been deleted.")
	} else {
		fmt.Println("Done.")
	}
	return
}

// slackDeleteAllFiles : Delete all files.
func (s *slackContainer) slackDeleteAllFiles() {
	p := url.Values{}
	p.Set("token", s.slackParams.Token)
	r := &utl.RequestParams{
		Method:      "POST",
		APIURL:      slackchkat + p.Encode(),
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	var sc slackChkAt
	if json.Unmarshal(body, &sc); err != nil || !sc.OK {
		fmt.Fprintf(os.Stderr, "Error: %s\n", string(body))
		os.Exit(1)
	}
	var input string
	fmt.Printf("Here is a team '%s' on Slack. You are %s.\n[WARNING] Will you delete all files here? [y or n] ... ", sc.Team, sc.User)
	if _, err := fmt.Scan(&input); err != nil {
		log.Fatalf("Error: %v.\n", err)
	}
	if input == "y" {
		ar := s.slackParams.SlackFilesList.Files
		count := len(ar)
		if count > 0 {
			bar := pb.StartNew(count)
			for i := 0; i < count; i++ {
				bar.Increment()
				p := url.Values{}
				p.Set("token", s.slackParams.Token)
				p.Set("file", ar[i].ID)
				r := &utl.RequestParams{
					Method:      "POST",
					APIURL:      slackurl + "files.delete?" + p.Encode(),
					Data:        nil,
					Contenttype: "application/x-www-form-urlencoded",
					Dtime:       10,
				}
				body, err := r.FetchAPI()
				var se slackError
				if json.Unmarshal(body, &se); err != nil || !se.OK {
					fmt.Fprintf(os.Stderr, "\nError: [ %s ] Overuse of API, or owner of this channel may not be you.\n", se.Error)
					os.Exit(1)
				}
			}
			bar.FinishPrint("Done.")
		} else {
			fmt.Println("No files.")
		}
	} else {
		fmt.Printf("Team %s's files on Slack were not deleted.", sc.Team)
	}
	return
}

// slackDeleteHistory : Delete a history
func (s *slackContainer) slackDeleteHistory() {
	p := url.Values{}
	p.Set("token", s.Token)
	p.Set("ts", s.jsonControl.Options["deletehistory"].(string))
	p.Set("channel", s.slackGetChannels().slackChannelNameToID())
	r := &utl.RequestParams{
		Method:      "POST",
		APIURL:      slackurl + "chat.delete?" + p.Encode(),
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       10,
	}
	body, err := r.FetchAPI()
	var se slackError
	if json.Unmarshal(body, &se); err != nil || !se.OK {
		fmt.Fprintf(os.Stderr, "Error: %s\n", se.Error)
		os.Exit(1)
	} else {
		fmt.Println("Done.")
	}
}

// slackDeleteChannelAllHistory : Delete histories. But it can delete just 50 histories at once.
func (s *slackContainer) slackDeleteChannelAllHistory() {
	fmt.Printf("# 50 histories can be deleted for one time, because of the limitation.\n")
	fmt.Printf("# If you want to delete a lot of histories, please run several times.\n")
	fmt.Printf("# Histories are deleted in order of the old date.\n")
	ar := s.ChannelHistory.Messages
	count := len(ar)
	if count > 0 {
		var numdel int
		if s.jsonControl.Options["deletehistories"].(int) < 50 {
			numdel = s.jsonControl.Options["deletehistories"].(int)
		} else {
			numdel = 50
		}
		if count < numdel {
			numdel = count
		}
		var j int
		bar := pb.StartNew(numdel)
		for i := count - numdel; i < count; i++ {
			j = (count - 1) - (i - count + numdel)
			bar.Increment()
			p := url.Values{}
			p.Set("token", s.slackParams.Token)
			p.Set("ts", ar[j].Ts)
			p.Set("channel", s.slackGetChannels().slackChannelNameToID())
			r := &utl.RequestParams{
				Method:      "POST",
				APIURL:      slackurl + "chat.delete?" + p.Encode(),
				Data:        nil,
				Contenttype: "application/x-www-form-urlencoded",
				Dtime:       10,
			}
			body, err := r.FetchAPI()
			var se slackError
			if json.Unmarshal(body, &se); err != nil || !se.OK {
				fmt.Fprintf(os.Stderr, "\nError: [ %s, %s ] Overuse of API, or owner of this channel may not be you.\n", se.Error, err)
				os.Exit(1)
			}
		}
		bar.FinishPrint("Done.")
	} else {
		fmt.Printf("No history in channel %s.\n", s.jsonControl.Options["channel"].(string))
	}
	return
}
