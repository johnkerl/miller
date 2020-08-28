//  Copyright 2020 Marius Ackerman
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

/*
Package md extracts code sections of markdown files
*/
package md

import (
	"io/ioutil"
)

var ch rune

/*
GetSource returns code sections eclosed in triple backticks.
*/
func GetSource(mdfile string) (string, error) {
	inbuf, err := ioutil.ReadFile(mdfile)
	if err != nil {
		return "", err
	}
	input := []rune(string(inbuf))
	loadMd(input)
	return string(input), nil
}

func loadMd(input []rune) {
	i := 0
	text := true
	for i < len(input) {
		if i <= len(input)-3 && input[i] == '`' && input[i+1] == '`' && input[i+2] == '`' {
			text = !text
			for j := 0; j < 3; j++ {
				input[i+j] = ' '
			}
			i += 3
		}
		if i < len(input) {
			if text {
				if input[i] == '\n' {
					input[i] = '\n'
				} else {
					input[i] = ' '
				}
			}
			i += 1
		}
	}
}
