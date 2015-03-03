// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fs

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	fusefs "bazil.org/fuse/fs"
	"github.com/jacobsa/gcloud/gcs"
	"github.com/jacobsa/gcsfuse/timeutil"
)

var fEnableDebug = flag.Bool(
	"fs.debug",
	false,
	"Write gcsfuse/fs debugging messages to stderr.")

type fileSystem struct {
	logger *log.Logger
	clock  timeutil.Clock
	bucket gcs.Bucket
}

func (fs *fileSystem) Root() (fusefs.Node, error) {
	d := newDir(fs.logger, fs.clock, fs.bucket, "")
	return d, nil
}

func getLogger() *log.Logger {
	var writer io.Writer = ioutil.Discard
	if *fEnableDebug {
		writer = os.Stderr
	}

	return log.New(writer, "gcsfuse/fs: ", log.LstdFlags)
}

// Create a fuse file system whose root directory is the root of the supplied
// bucket. The supplied clock will be used for cache invalidation; it is *not*
// used for file modification times.
func NewFuseFS(clock timeutil.Clock, bucket gcs.Bucket) (fusefs.FS, error) {
	fs := &fileSystem{
		logger: getLogger(),
		clock:  clock,
		bucket: bucket,
	}

	return fs, nil
}