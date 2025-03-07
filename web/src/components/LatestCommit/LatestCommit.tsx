/*
 * Copyright 2023 Harness, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React from 'react'
import { Container, Layout, FlexExpander, Text, Avatar } from '@harnessio/uicore'
import { Color, FontVariation } from '@harnessio/design-system'
import { Link } from 'react-router-dom'
import { Render } from 'react-jsx-match'
import ReactTimeago from 'react-timeago'
import cx from 'classnames'
import type { TypesCommit } from 'services/code'
import { CommitActions } from 'components/CommitActions/CommitActions'
import { useAppContext } from 'AppContext'
import { formatBytes, formatDate } from 'utils/Utils'
import type { GitInfoProps } from 'utils/GitUtils'
import { PipeSeparator } from 'components/PipeSeparator/PipeSeparator'
import css from './LatestCommit.module.scss'

interface LatestCommitProps extends Pick<GitInfoProps, 'repoMetadata'> {
  latestCommit?: TypesCommit
  standaloneStyle?: boolean
  size?: number
}

export function LatestCommitForFolder({ repoMetadata, latestCommit, standaloneStyle }: LatestCommitProps) {
  const { routes } = useAppContext()
  const commitURL = routes.toCODECommit({
    repoPath: repoMetadata.path as string,
    commitRef: latestCommit?.sha as string
  })

  return (
    <Render when={latestCommit}>
      <Container>
        <Layout.Horizontal spacing="small" className={cx(css.latestCommit, { [css.standalone]: standaloneStyle })}>
          <Avatar hoverCard={false} size="small" name={latestCommit?.author?.identity?.name || ''} />
          <Text font={{ variation: FontVariation.SMALL_BOLD }}>
            {latestCommit?.author?.identity?.name || latestCommit?.author?.identity?.email}
          </Text>
          <Link to={commitURL}>
            <Text className={css.commitLink} lineClamp={1}>
              {latestCommit?.title}
            </Text>
          </Link>
          <FlexExpander />
          <CommitActions sha={latestCommit?.sha as string} href={commitURL} enableCopy />
          <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_400} className={css.time}>
            <ReactTimeago date={latestCommit?.author?.when as string} />
          </Text>
        </Layout.Horizontal>
      </Container>
    </Render>
  )
}

export function LatestCommitForFile({ repoMetadata, latestCommit, standaloneStyle, size }: LatestCommitProps) {
  const { routes } = useAppContext()
  const commitURL = routes.toCODECommit({
    repoPath: repoMetadata.path as string,
    commitRef: latestCommit?.sha as string
  })

  return (
    <Render when={latestCommit}>
      <Container>
        <Layout.Horizontal
          spacing="medium"
          className={cx(css.latestCommit, css.forFile, { [css.standalone]: standaloneStyle })}>
          <Avatar hoverCard={false} size="small" name={latestCommit?.author?.identity?.name || ''} />
          <Text font={{ variation: FontVariation.SMALL_BOLD }}>
            {latestCommit?.author?.identity?.name || latestCommit?.author?.identity?.email}
          </Text>
          <PipeSeparator height={9} />

          <Text lineClamp={1} tooltipProps={{ portalClassName: css.popover }}>
            <Link to={commitURL} className={css.commitLink}>
              {latestCommit?.title}
            </Link>
          </Text>
          <PipeSeparator height={9} />
          <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_400}>
            {formatDate(latestCommit?.author?.when as string)}
          </Text>
          {size && size > 0 && (
            <>
              <PipeSeparator height={9} />
              <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_400}>
                {formatBytes(size)}
              </Text>
            </>
          )}

          <FlexExpander />
          <CommitActions sha={latestCommit?.sha as string} href={commitURL} enableCopy />
        </Layout.Horizontal>
      </Container>
    </Render>
  )
}
