/*
Package main (doc.go) :
This is a CLI tool to submit files to both Gist and Slack.

When I submitted a script to Slack, I had saved the script to Gist as a backup.
I have done manually this. Recently, I wished this process could be automatically run.
So I created this tool.

# Features of "gislack" are as follows.

1. Submits files to both Gist and Slack, simultaneously.

2. Submits, updates and deletes files for Gist.

3. Submits and deletes files for Slack.

4. Retrieves access token from client ID and client secret for Gist and Slack.

You can see the release page https://github.com/tanaikech/gislack/releases

More information is https://github.com/tanaikech/gislack

# APIs
Gist api document https://developer.github.com/v3/gists/

Slack api document https://api.slack.com/methods

---------------------------------------------------------------

# How to use gislack
At first, please retrieve client ID and client secret from GitHub and Slack.
Please set `http://localhost:8080` for each redirect_uri.
If you cannot port 8080, you can change it.


GitHub -> https://developer.github.com/apps/building-integrations/setting-up-and-registering-oauth-apps/about-authorization-options-for-oauth-apps/
Slack -> https://api.slack.com/apps

After retrieved them, please run as follows.

A1. Retrieves access token for Gist
$ gislack auth -gi {client_id of GitHub} -gs {client_secret of GitHub} -p {port(default is 8080.)}

A2. Retrieves access token for Slack
$ gislack auth -si {client_id of Slack} -ss {client_secret of Slack} -p {port(default is 8080.)}

By this, you can see `gislack.cfg` file on your current working directory.

In the case of submitting a file to both Gist and Slack,

$ gislack d -f {file} -t {title} -p -ft {file type for Slack(e.g. javascript)} -ch {channel for Slack(e.g. general)} -ic {initial comment for Slack}

-p : This is a boolean. If you want to submit as a public for Gist, please use this.
-ic : You can give initial comment using this.

In the case of submitting a file to Gist,

$ gislack g -f {file} -t {title} -p

In the case of submitting a file to Slack,

$ gislack s -f {file} -ti {title} -ch {channel} -ft {file type} -ic {initial comment}

gislack has more commands, please check it here. ( https://github.com/tanaikech/gislack )

---------------------------------------------------------------
*/
package main
