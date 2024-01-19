// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package handles the versioning in the API for prerecorded endpoint
package version

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/google/go-querystring/query"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/pkg/common"
)

const (
	// ReadAPIVersion current supported version
	ReadAPIVersion string = "v1"

	// ReadPath is the current path for prerecorded transcription
	ReadPath string = "read"
)

/*
GetAnalyzeAPI is a function which controls the versioning of the live transcription API and provides
mechanism for:

- overriding the host endpoint
- overriding the version used
- overriding the endpoint path
- additional arguments to the query string/parameters

The return value is the complete URL endpoint to be used for the live transcription
*/
func GetAnalyzeAPI(ctx context.Context, host, version, path string, options interfaces.AnalyzeOptions, args ...interface{}) (string, error) {
	if path == "" {
		path = ReadPath
	}

	if host == "" {
		host = common.DefaultHost
	}
	if version == "" {
		version = ReadAPIVersion
	}

	r, err := regexp.Compile("^(v[0-9]+|%%s)/")
	if err != nil {
		klog.V(1).Infof("regexp.Compile err: %v\n", err)
		return "", err
	}

	match := r.MatchString(path)
	klog.V(3).Infof("match: %t\n", match)

	if match {
		// version = r.FindStringSubmatch(path)[0]
		path = r.ReplaceAllString(path, "")
	}

	q, err := query.Values(options)
	if err != nil {
		return "", err
	}

	if parameters, ok := ctx.Value(interfaces.ParametersContext{}).(map[string][]string); ok {
		for k, vs := range parameters {
			for _, v := range vs {
				q.Add(k, v)
			}
		}
	}

	fullpath := fmt.Sprintf("%%s/%s", path)
	completeFullpath := fmt.Sprintf(fullpath, append([]interface{}{version}, args...)...)
	u := url.URL{Scheme: "https", Host: host, Path: completeFullpath, RawQuery: q.Encode()}

	return u.String(), nil
}