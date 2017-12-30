/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package public

import (
	"github.com/emirpasic/gods/sets/hashset"
	api "github.com/infinitbyte/gopa/core/http"
	"github.com/infinitbyte/gopa/modules/ui/common"
	"golang.org/x/oauth2"
)

const (
	defaultLayout = "templates/layout.html"
	templateDir   = "templates/"

	defaultConfigFile = "config.json"

	githubAuthorizeUrl = "https://github.com/login/oauth/authorize"
	githubTokenUrl     = "https://github.com/login/oauth/access_token"
	redirectUrl        = ""
)

var (
	oauthCfg *oauth2.Config

	// scopes
	scopes = []string{"repo"}
)

func InitUI(cfg common.AuthConfig) {

	ui := PublicUI{}
	api.HandleUIMethod(api.GET, "/redirect/", ui.RedirectHandler)

	if !cfg.Enabled {
		return
	}

	oauthCfg = &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  githubAuthorizeUrl,
			TokenURL: githubTokenUrl,
		},
		RedirectURL: redirectUrl,
		Scopes:      scopes,
	}

	ui.Admins = hashset.New()
	for _, v := range cfg.AuthorizedAdmins {
		if v != "" {
			ui.Admins.Add(v)
		}
	}

	api.HandleUIMethod(api.GET, "/auth/github/", ui.AuthHandler)
	api.HandleUIMethod(api.GET, "/auth/callback/", ui.CallbackHandler)
	api.HandleUIMethod(api.GET, "/auth/login/", ui.Login)
	api.HandleUIMethod(api.GET, "/auth/success/", ui.LoginSuccess)
	api.HandleUIMethod(api.GET, "/auth/fail/", ui.LoginFail)

}
