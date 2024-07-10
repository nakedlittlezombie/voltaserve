// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useMemo } from 'react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { leViewerPermission } from '@/client/api/permission'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type ViewerPDFProps = {
  file: File
}

const ViewerPDF = ({ file }: ViewerPDFProps) => {
  const url = useMemo(() => {
    if (file.snapshot?.preview && file.snapshot?.preview.extension) {
      return `/proxy/api/v2/files/${file.id}/preview${
        file.snapshot?.preview.extension
      }?${new URLSearchParams({
        access_token: getAccessTokenOrRedirect(),
      })}`
    }
  }, [file])

  if (!file.snapshot?.preview) {
    return null
  }

  return (
    <iframe className={cx('w-full', 'h-full')} src={url} title={file.name} />
  )
}

export default ViewerPDF
