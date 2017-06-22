// Package main (gistslack.go) :
// This file is included all commands and options.
package main

import (
	"os"

	"github.com/urfave/cli"
)

// main : main method
func main() {
	app := cli.NewApp()
	app.Name = appname
	app.Author = "tanaike [ https://github.com/tanaikech/gislack ] "
	app.Email = "tanaike@hotmail.com"
	app.Usage = "Submit files to Gist, Slack and both."
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:        "gist",
			Aliases:     []string{"g"},
			Usage:       "Submits files to gist.",
			Description: "In this mode, an access token is required for both gist and slack.",
			Action:      gist,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "title, t",
					Usage: "Value is submission title.",
				},
				cli.StringFlag{
					Name:  "files, f",
					Usage: "Value is submit files. You can set several files.",
				},
				cli.StringFlag{
					Name:  "filenames, fn",
					Usage: "Value is file names. If you want to use different names for from submitting files, please use this.",
				},
				cli.BoolFlag{
					Name:  "public, p",
					Usage: "Submitting as a public. Default is non public.",
				},
				cli.BoolFlag{
					Name:  "list, l",
					Usage: "Display list of gists.",
				},
				cli.BoolFlag{
					Name:  "listasjson, lj",
					Usage: "Display list of gists as JSON.",
				},
				cli.StringFlag{
					Name:  "get, g",
					Usage: "Get a single gist. You can check ID by list command.",
				},
				cli.StringFlag{
					Name:  "gethistory, gh",
					Usage: "Value is gist ID. Get history of gist ID. You can check ID by list command.",
				},
				cli.StringFlag{
					Name:  "getversion, gv",
					Usage: "Value is URL for a version. Get version of gist ID.",
				},
				cli.StringFlag{
					Name:  "updateoverwrite, uo",
					Usage: "Value is gist ID. File is overwritten. You can check ID by list command.",
				},
				cli.StringFlag{
					Name:  "updateadd, ua",
					Usage: "Value is gist ID. File is added. You can check ID by list command.",
				},
				cli.StringFlag{
					Name:  "delete, d",
					Usage: "Value is gist ID. You can check ID by list command.",
				},
				cli.BoolFlag{
					Name:  "deleteall",
					Usage: "Warning : Delete all files.",
				},
				cli.BoolFlag{
					Name:  "anonymous",
					Usage: "Warning : Submit files as anonymous.",
				},
				cli.BoolFlag{
					Name:  "jsonparser, j",
					Usage: "Displays results by JSON parser.",
				},
				cli.StringFlag{
					Name:  "cfgdirectory, cfgdir",
					Usage: "Value is path of directory with gislack.cfg.",
				},
			},
		},
		{
			Name:        "slack",
			Aliases:     []string{"s"},
			Usage:       "Submits files to slack.",
			Description: "In this mode, an access token is required for both gist and slack.",
			Action:      slack,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "Value is a file you want to submit.",
				},
				cli.StringFlag{
					Name:  "channel, ch",
					Usage: "Value is a submission channel. Channel name or channel ID.",
				},
				cli.StringFlag{
					Name:  "content, co",
					Usage: "Value is content strings you want to submit. File is prioritized.",
				},
				cli.StringFlag{
					Name:  "title, ti",
					Usage: "Value is a submission title.",
				},
				cli.StringFlag{
					Name:  "filetype, ft",
					Usage: "Value is file type.",
				},
				cli.StringFlag{
					Name:  "initialcomment, ic",
					Usage: "Value is initial comment for submission.",
				},
				cli.BoolFlag{
					Name:  "channellist, cl",
					Usage: "Display channel list.",
				},
				cli.BoolFlag{
					Name:  "filelist, fl",
					Usage: "Display file list.",
				},
				cli.BoolFlag{
					Name:  "filelistasjson, fj",
					Usage: "Display file list as JSON.",
				},
				cli.StringFlag{
					Name:  "getfile, gf",
					Usage: "Value is file ID. You can check ID by filelist command.",
				},
				cli.StringFlag{
					Name:  "user, u",
					Usage: "Value is a submitted user ID. This is used to retrieve file list.",
				},
				cli.BoolFlag{
					Name:  "channelhistory, hi",
					Usage: "Display history list for a channel.",
				},
				cli.StringFlag{
					Name:  "deletefile, df",
					Usage: "Value is a file ID you want to delete.",
				},
				cli.BoolFlag{
					Name:  "deletefiles, dfs",
					Usage: "Delete files. This is bool.",
				},
				cli.StringFlag{
					Name:  "deletehistory, dh",
					Usage: "Value is a history ID you want to delete.",
				},
				cli.IntFlag{
					Name:  "deletehistories, dhs",
					Usage: "Value is number of histories you want to delete.",
					Value: 0,
				},
				cli.BoolFlag{
					Name:  "jsonparser, j",
					Usage: "Displays results by JSON parser.",
				},
				cli.StringFlag{
					Name:  "cfgdirectory, cfgdir",
					Usage: "Value is path of directory with gislack.cfg.",
				},
			},
		},
		{
			Name:        "doublesubmit",
			Aliases:     []string{"d"},
			Usage:       "Submits files to both gist and slack.",
			Description: "In this mode, an access token is required for both gist and slack.",
			Action:      doubleSubmit,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "title, t",
					Usage: "Value is submission title for both.",
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "Value is a file for both. In this mode, it can upload only one file.",
				},
				cli.StringFlag{
					Name:  "filename, fn",
					Usage: "Value is file name. If you want to use different name for from submitting file, please use this.",
				},
				cli.BoolFlag{
					Name:  "public, p",
					Usage: "Gist : Submitting as a public.",
				},
				cli.StringFlag{
					Name:  "updateoverwrite, uo",
					Usage: "Value is gist ID. File is overwritten. You can check ID by list command.",
				},
				cli.StringFlag{
					Name:  "updateadd, ua",
					Usage: "Value is gist ID. File is added. You can check ID by list command.",
				},
				cli.StringFlag{
					Name:  "filetype, ft",
					Usage: "Slack : Value is file type.",
				},
				cli.StringFlag{
					Name:  "channel, ch",
					Usage: "Slack : Value is a submission channel.",
				},
				cli.StringFlag{
					Name:  "initialcomment, ic",
					Usage: "Slack : Value is initial comment.",
				},
				cli.BoolFlag{
					Name:  "simpleresult, s",
					Usage: "Displays simple results.",
				},
				cli.BoolFlag{
					Name:  "jsonparser, j",
					Usage: "Displays results by JSON parser.",
				},
				cli.StringFlag{
					Name:  "cfgdirectory, cfgdir",
					Usage: "Value is path of directory with gislack.cfg.",
				},
			},
		},
		{
			Name:        "auth",
			Aliases:     []string{"a"},
			Usage:       "Retrieves access tokens for gist and slack.",
			Description: "In this mode, client ID and client secret are required for gist and slack.",
			Action:      getaccesstopen,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "gistclientid, gi",
					Usage: "Client ID for gist.",
				},
				cli.StringFlag{
					Name:  "gistclientsecret, gs",
					Usage: "Client secret for gist.",
				},
				cli.StringFlag{
					Name:  "slackclientid, si",
					Usage: "Client ID for slack.",
				},
				cli.StringFlag{
					Name:  "slackclientsecret, ss",
					Usage: "Client secret for slack.",
				},
				cli.BoolFlag{
					Name:  "chkgisttoken, cgt",
					Usage: "Check access token for gist.",
				},
				cli.StringFlag{
					Name:  "cfgdirectory, cfgdir",
					Usage: "Value is path of directory with gislack.cfg.",
				},
				cli.IntFlag{
					Name:  "port, p",
					Usage: "Value is port number of HTML server for redirect.",
					Value: 8080,
				},
			},
		},
		{
			Name:        "json",
			Aliases:     []string{"j"},
			Usage:       "Submission control using JSON data. Please check document at GitHub about how to use this command.",
			Description: "In this mode, client ID and client secret are required for gist and slack.",
			Action:      useJSON,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json",
					Usage: "Value is JSON data.",
				},
			},
		},
	}
	app.CommandNotFound = commandNotFound
	app.Run(os.Args)
}
