// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"time"

	"github.com/harness/gitness/cache"
	"github.com/harness/gitness/gitrpc/internal/gitea"
	"github.com/harness/gitness/gitrpc/internal/service"
	"github.com/harness/gitness/gitrpc/internal/types"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideServer,
	ProvideHTTPServer,
	ProvideGITAdapter,
	ProvideGoGitRepoProvider,
	ProvideLastCommitCache,
)

func ProvideGoGitRepoProvider() *gitea.GoGitRepoProvider {
	const objectCacheSize = 16 << 20 // 16MiB
	return gitea.NewGoGitRepoProvider(objectCacheSize, 15*time.Minute)
}

func ProvideLastCommitCache(
	config Config,
	redisClient redis.UniversalClient,
	repoProvider *gitea.GoGitRepoProvider,
) cache.Cache[gitea.CommitEntryKey, *types.Commit] {
	cacheDuration := time.Duration(config.LastCommitCache.DurationSeconds) * time.Second

	if config.LastCommitCache.Mode == ModeNone || cacheDuration < time.Second {
		return gitea.NoLastCommitCache(repoProvider)
	}

	if config.LastCommitCache.Mode == ModeRedis && redisClient != nil {
		return gitea.NewRedisLastCommitCache(redisClient, cacheDuration, repoProvider)
	}

	return gitea.NewInMemoryLastCommitCache(cacheDuration, repoProvider)
}

func ProvideGITAdapter(
	repoProvider *gitea.GoGitRepoProvider,
	lastCommitCache cache.Cache[gitea.CommitEntryKey, *types.Commit],
) (service.GitAdapter, error) {
	return gitea.New(repoProvider, lastCommitCache)
}

func ProvideServer(config Config, adapter service.GitAdapter) (*GRPCServer, error) {
	return NewServer(config, adapter)
}

func ProvideHTTPServer(config Config) (*HTTPServer, error) {
	return NewHTTPServer(config)
}
