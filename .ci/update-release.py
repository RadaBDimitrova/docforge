#!/usr/bin/env python3

# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

import pathlib
import os

import ccc.github

VERSION_FILE_NAME='VERSION'

repo_owner_and_name = util.check_env('SOURCE_GITHUB_REPO_OWNER_AND_NAME')
repo_dir = util.check_env('MAIN_REPO_DIR')
output_dir = util.check_env('BINARY')

repo_owner, repo_name = repo_owner_and_name.split('/')

repo_path = pathlib.Path(repo_dir).resolve()
output_path = pathlib.Path(output_dir).resolve()
version_file_path = repo_path / VERSION_FILE_NAME

version_file_contents = version_file_path.read_text()

github_api = ccc.github.github_api(repo_url=f'github.com/{repo_owner_and_name}')
repository = github_api.repository(repo_owner, repo_path)


gh_release = repository.release_from_tag(version_file_contents)

for dir, dirs, files in os.walk(os.path.join(output_path, "bin", "rel")):
    for binName in files:
        dir_path = pathlib.Path(dir).resolve()
        binFilePath = dir_path / binName
        gh_release.upload_asset(
            content_type='application/octet-stream',
            name=f'{binName}',
            asset=binFilePath.open(mode='rb'),
        )
