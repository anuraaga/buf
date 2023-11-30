// Copyright 2020-2023 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bufisk

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bufbuild/buf/private/pkg/app/appcmd"
	"github.com/bufbuild/buf/private/pkg/app/appflag"
	"github.com/spf13/pflag"
)

const (
	useBufVersionEnvKey         = "USE_BUF_VERSION"
	bufVersionFileName          = ".bufversion"
	bufVersionLatestGithubValue = "latest"
)

// Main is the entrypoint to the buf CLI.
func Main(name string) {
	appcmd.Main(context.Background(), newCommand(name))
}

func newCommand(name string) *appcmd.Command {
	builder := appflag.NewBuilder(
		name,
		appflag.BuilderWithTimeout(120*time.Second),
		appflag.BuilderWithTracing(),
	)
	flags := newFlags()
	return &appcmd.Command{
		Use: name,
		Run: builder.NewRunFunc(
			func(ctx context.Context, container appflag.Container) error {
				return run(ctx, container, flags)
			},
		),
		BindFlags:           flags.Bind,
		BindPersistentFlags: builder.BindRoot,
	}
}

type flags struct{}

func newFlags() *flags {
	return &flags{}
}

func (f *flags) Bind(flagSet *pflag.FlagSet) {}

func run(
	ctx context.Context,
	container appflag.Container,
	flags *flags,
) error {
	bufVersion, err := getBufVersion()
	if err != nil {
		return err
	}
	bufFilePath := filepath.Join(container.CacheDirPath())
	_ = bufVersion
	return nil
}

func getBufVersion() (string, error) {
	if useBufVersionEnvValue := os.Getenv(useBufVersionEnvKey); useBufVersionEnvValue != "" {
		return useBufVersionEnvValue, nil
	}
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	curDirPath := pwd
	for {
		data, err := os.ReadFile(filepath.Join(curDirPath, bufVersionFileName))
		// Ignore all errors, not just fs.ErrNotExist - we don't want this program to fail
		// on bad permissions. We could choose to stop on the first bad permissions error.
		if err == nil {
			return strings.TrimSpace(string(data)), nil
		}
		if curDirPath == string(os.PathSeparator) {
			break
		}
		curDirPath = filepath.Dir(curDirPath)
	}
	return "", fmt.Errorf("%s not set and no %s file found", useBufVersionEnvKey, bufVersionFileName)
}
