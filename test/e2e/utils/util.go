/*
Copyright 2019 The Kubernetes Authors.

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

package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
)

const (
	// KubebuilderBinName define the name of the kubebuilder binary to be used in the tests
	KubebuilderBinName = "kubebuilder"
)

// RandomSuffix returns a 4-letter string.
func RandomSuffix() (string, error) {
	source := []rune("abcdefghijklmnopqrstuvwxyz")
	res := make([]rune, 4)
	for i := range res {
		bi := new(big.Int)
		r, err := rand.Int(rand.Reader, bi.SetInt64(int64(len(source))))
		if err != nil {
			return "", err
		}
		res[i] = source[r.Int64()]
	}
	return string(res), nil
}

// GetNonEmptyLines converts given command output string into individual objects
// according to line breakers, and ignores the empty elements in it.
func GetNonEmptyLines(output string) []string {
	var res []string
	elements := strings.Split(output, "\n")
	for _, element := range elements {
		if element != "" {
			res = append(res, element)
		}
	}

	return res
}

// InsertCode searches target content in the file and insert `toInsert` after the target.
func InsertCode(filename, target, code string) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	idx := strings.Index(string(contents), target)
	out := string(contents[:idx+len(target)]) + code + string(contents[idx+len(target):])
	// false positive
	// nolint:gosec
	return ioutil.WriteFile(filename, []byte(out), 0644)
}

// UncommentCode searches for target in the file and remove the comment prefix
// of the target content. The target content may span multiple lines.
func UncommentCode(filename, target, prefix string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	strContent := string(content)

	idx := strings.Index(strContent, target)
	if idx < 0 {
		return nil
	}

	out := new(bytes.Buffer)
	_, err = out.Write(content[:idx])
	if err != nil {
		return err
	}

	strs := strings.Split(target, "\n")
	for _, str := range strs {
		_, err := out.WriteString(strings.TrimPrefix(str, prefix) + "\n")
		if err != nil {
			return err
		}
	}

	_, err = out.Write(content[idx+len(target):])
	if err != nil {
		return err
	}
	// false positive
	// nolint:gosec
	return ioutil.WriteFile(filename, out.Bytes(), 0644)
}

// ImplementWebhooks will mock an webhook data
func ImplementWebhooks(filename string) error {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	str := string(bs)

	str, err = EnsureExistAndReplace(
		str,
		"import (",
		`import (
	"errors"`)
	if err != nil {
		return err
	}

	// implement defaulting webhook logic
	str, err = EnsureExistAndReplace(
		str,
		"// TODO(user): fill in your defaulting logic.",
		`if r.Spec.Count == 0 {
		r.Spec.Count = 5
	}`)
	if err != nil {
		return err
	}

	// implement validation webhook logic
	str, err = EnsureExistAndReplace(
		str,
		"// TODO(user): fill in your validation logic upon object creation.",
		`if r.Spec.Count < 0 {
		return errors.New(".spec.count must >= 0")
	}`)
	if err != nil {
		return err
	}
	str, err = EnsureExistAndReplace(
		str,
		"// TODO(user): fill in your validation logic upon object update.",
		`if r.Spec.Count < 0 {
		return errors.New(".spec.count must >= 0")
	}`)
	if err != nil {
		return err
	}
	// false positive
	// nolint:gosec
	return ioutil.WriteFile(filename, []byte(str), 0644)
}

// EnsureExistAndReplace check if the content exists and then do the replace
func EnsureExistAndReplace(input, match, replace string) (string, error) {
	if !strings.Contains(input, match) {
		return "", fmt.Errorf("can't find %q", match)
	}
	return strings.Replace(input, match, replace, -1), nil
}
