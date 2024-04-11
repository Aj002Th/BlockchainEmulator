//go:build wireinject
// +build wireinject

package boot

// Copyright 2018 The Wire Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The build tag makes sure the stub is not built in the final build.

// 从Wire官方教程抄来
// 第一句是必须的。保证不会被编译的时候包含进来。而是留给wire的板子。

import "github.com/google/wire"

func InitializeApp() (App, error) {
	wire.Build(NewApp, ParseAndBuildArg)
	return App{}, nil
}
