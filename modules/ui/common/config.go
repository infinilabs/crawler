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

package common

type AuthConfig struct {
	Enabled           bool     `config:"enabled"`
	OAuthProvider     string   `config:"oauth_provider"`
	oauthAuthorizeUrl string   `config:"oauth_authorize_url"`
	oauthTokenUrl     string   `config:"oauth_token_url"`
	oauthRedirectUrl  string   `config:"oauth_redirect_url"`
	AuthorizedAdmins  []string `config:"authorized_admin"`
	ClientSecret      string   `config:"client_secret"`
	ClientID          string   `config:"client_id"`
}

type UIConfig struct {
	AuthConfig AuthConfig `config:"auth"`
}
