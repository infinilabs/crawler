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


package env

func GetWelcomeMessage() string {
	s := ("  ________ ________ __________  _____   \n")
	s += (" /  _____/ \\_____  \\\\______   \\/  _  \\  \n")
	s += ("/   \\  ___  /   |   \\|     ___/  /_\\  \\ \n")
	s += ("\\    \\_\\  \\/    |    \\    |  /    |    \\\n")
	s += (" \\______  /\\_______  /____|  \\____|__  /\n")
	s += ("        \\/         \\/                \\/ \n")

	commitLog := ""
	if len(LastCommitLog) > 0 {
		commitLog = "\n///last commit: " + LastCommitLog + "///"
	}
	s += ("[gopa] " + VERSION + "" + commitLog + "\n")
	return (s)
}
