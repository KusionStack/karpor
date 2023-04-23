/*
 * Copyright The Karbour Authors.
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

export const basicSyntaxColumns = [
  [
    {
      title: 'Scope search to specific repos',
      queryExamples: [
        { id: 'org-repos', query: 'repo:sourcegraph/.*' },
        { id: 'single-repo', query: 'repo:facebook/react' },
      ],
    },
    // {
    //     title: 'Jump into code navigation',
    //     queryExamples: [
    //         { id: 'file-filter', query: 'file:README.md' },
    //         { id: 'type-symbol', query: 'type:symbol SymbolName' },
    //     ],
    // },
    // {
    //     title: 'Get logical',
    //     queryExamples: [
    //         { id: 'not-operator', query: 'lang:go NOT file:main.go' },
    //         { id: 'or-operator', query: 'lang:javascript OR lang:typescript' },
    //         { id: 'and-operator', query: 'hello AND world' },
    //     ],
    // },
  ],
  [
    {
      title: 'Get logical',
      queryExamples: [
        { id: 'not-operator', query: 'lang:go NOT file:main.go' },
        { id: 'or-operator', query: 'lang:javascript OR lang:typescript' },
        { id: 'and-operator', query: 'hello AND world' },
      ],
    },
    // {
    //     title: 'Find content or patterns',
    //     queryExamples: [
    //         { id: 'exact-matches', query: 'some exact error message', helperText: 'No quotes needed' },
    //         { id: 'regex-pattern', query: '/regex.*pattern/' },
    //     ],
    // },
    // {
    //     title: 'Explore code history',
    //     queryExamples: [
    //         { id: 'type-diff-author', query: 'type:diff author:torvalds' },
    //         { id: 'type-commit-message', query: 'type:commit some message' },
    //     ],
    // },
    // {
    //     title: 'Get advanced',
    //     queryExamples: [
    //         { id: 'repo-has-description', query: 'repo:has.description(scientific computing)' },
    //         { id: 'commit-search', query: 'repo:has.commit.after(june 25 2017)' },
    //     ],
    // },
  ],
]