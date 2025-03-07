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

package service

import (
	"errors"
	"io"

	"github.com/harness/gitness/gitrpc/internal/types"
	"github.com/harness/gitness/gitrpc/rpc"

	"github.com/rs/zerolog/log"
)

func (s RepositoryService) GetBlob(request *rpc.GetBlobRequest, stream rpc.RepositoryService_GetBlobServer) error {
	if err := validateGetBlobRequest(request); err != nil {
		return err
	}

	ctx := stream.Context()
	base := request.GetBase()
	repoPath := getFullPathForRepo(s.reposRoot, base.GetRepoUid())

	// TODO: do we need to validate request for nil?
	gitBlob, err := s.adapter.GetBlob(ctx, repoPath, request.GetSha(), request.GetSizeLimit())
	if err != nil {
		return processGitErrorf(err, "failed to get blob")
	}
	defer func() {
		dErr := gitBlob.Content.Close()
		if dErr != nil {
			log.Ctx(ctx).Warn().Err(err).Msgf("failed to close blob content reader.")
		}
	}()

	err = stream.Send(&rpc.GetBlobResponse{
		Data: &rpc.GetBlobResponse_Header{
			Header: &rpc.GetBlobResponseHeader{
				Sha:         request.GetSha(),
				Size:        gitBlob.Size,
				ContentSize: gitBlob.ContentSize,
			},
		},
	})
	if err != nil {
		return ErrInternalf("failed to send header: %w", err)
	}

	const bufferSize = 16384
	contentBuffer := make([]byte, bufferSize)
	for {
		if ctx.Err() != nil {
			return ErrCanceledf("the context got canceled while streaming the blob content: %w", ctx.Err())
		}

		n, rErr := gitBlob.Content.Read(contentBuffer)
		// according to io.Reader documentation, returned bytes should always be processed before the error
		// This is crucial to handle cases were n > 0 bytes are returned together with io.EOF
		if n > 0 {
			sErr := stream.Send(&rpc.GetBlobResponse{
				Data: &rpc.GetBlobResponse_Content{
					Content: contentBuffer[:n],
				},
			})
			if sErr != nil {
				return ErrInternalf("failed to send content of buffer (potential read error: %s): %w", rErr, sErr)
			}
		}
		if errors.Is(rErr, io.EOF) {
			break
		}
		if rErr != nil {
			return rErr
		}
	}

	return nil
}

func validateGetBlobRequest(r *rpc.GetBlobRequest) error {
	if r.GetBase() == nil {
		return types.ErrBaseCannotBeEmpty
	}
	if r.GetSha() == "" {
		return types.ErrEmptySHA
	}

	return nil
}
