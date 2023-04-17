
const tokenKindToCSSName: Record<MetaRevisionKind | MetaRegexpKind | MetaPredicateKind | MetaStructuralKind, string> = {
  Separator: 'separator',
  IncludeGlobMarker: 'include-glob-marker',
  ExcludeGlobMarker: 'exclude-glob-marker',
  CommitHash: 'commit-hash',
  Label: 'label',
  ReferencePath: 'reference-path',
  Wildcard: 'wildcard',
  Assertion: 'assertion',
  Alternative: 'alternative',
  Delimited: 'delimited',
  EscapedCharacter: 'escaped-character',
  CharacterSet: 'character-set',
  CharacterClass: 'character-class',
  CharacterClassRange: 'character-class-range',
  CharacterClassRangeHyphen: 'character-class-range-hyphen',
  CharacterClassMember: 'character-class-member',
  LazyQuantifier: 'lazy-quantifier',
  RangeQuantifier: 'range-quantifier',
  NameAccess: 'name-access',
  Dot: 'dot',
  Parenthesis: 'parenthesis',
  Hole: 'hole',
  RegexpHole: 'regexp-hole',
  Variable: 'variable',
  RegexpSeparator: 'regexp-separator',
}

export const toCSSClassName = (token: DecoratedToken): string => {
  switch (token.type) {
      case 'field':
          return 'search-filter-keyword'
      case 'keyword':
      case 'openingParen':
      case 'closingParen':
      case 'metaRepoRevisionSeparator':
      case 'metaContextPrefix':
          return 'search-keyword'
      case 'metaFilterSeparator':
          return 'search-filter-separator'
      case 'metaPath':
          return 'search-path-separator'

      case 'metaRevision': {
          return `search-revision-${tokenKindToCSSName[token.kind]}`
      }

      case 'metaRegexp': {
          return `search-regexp-meta-${tokenKindToCSSName[token.kind]}`
      }

      case 'metaPredicate': {
          return `search-predicate-${tokenKindToCSSName[token.kind]}`
      }

      case 'metaStructural': {
          return `search-structural-${tokenKindToCSSName[token.kind]}`
      }

      default:
          return 'search-query-text'
  }
}

export const createLiteral = (value: string, range: CharacterRange, quoted = false): Literal => ({
  type: 'literal',
  value,
  range,
  quoted,
})

const decorateContext = (token: Literal): DecoratedToken[] => {
  if (!token.value.startsWith('@')) {
      return [token]
  }

  const { start, end } = token.range
  return [
      { type: 'metaContextPrefix', range: { start, end: start + 1 }, value: '@' },
      createLiteral(token.value.slice(1), { start: start + 1, end }),
  ]
}


export enum PatternKind {
  Literal = 1,
  Regexp,
  Structural,
}

export enum MetaRegexpKind {
  Assertion = 'Assertion', // like ^ or \b
  Alternative = 'Alternative', // like |
  Delimited = 'Delimited', // like ( or )
  EscapedCharacter = 'EscapedCharacter', // like \(
  CharacterSet = 'CharacterSet', // like \s
  CharacterClass = 'CharacterClass', // like [a-z]
  CharacterClassRange = 'CharacterClassRange', // the a-z part in [a-z]
  CharacterClassRangeHyphen = 'CharacterClassRangeHyphen', // the - part in [a-z]
  CharacterClassMember = 'CharacterClassMember', // a character inside a charcter class like [abcd]
  LazyQuantifier = 'LazyQuantifier', // the ? after a range quantifier
  RangeQuantifier = 'RangeQuantifier', // like +
}

const mapRegexpMeta = (pattern: Pattern): DecoratedToken[] | undefined => {
  const tokens: DecoratedToken[] = []

  if (pattern.delimited) {
      tokens.push({
          type: 'metaRegexp',
          range: { start: pattern.range.start, end: pattern.range.start + 1 },
          value: '/',
          kind: MetaRegexpKind.Delimited,
      })
      tokens.push({
          type: 'metaRegexp',
          range: { start: pattern.range.end - 1, end: pattern.range.end },
          value: '/',
          kind: MetaRegexpKind.Delimited,
      })
  }

  const offset = pattern.delimited ? pattern.range.start + 1 : pattern.range.start

  try {
      const ast = new RegExpParser().parsePattern(pattern.value)
      visitRegExpAST(ast, {
          onAlternativeEnter(node: Alternative) {
              // regexpp doesn't tell us where a '|' operator is. We infer it by visiting any
              // pattern of an Alternative node, and for a '|' directly after it. Based on
              // regexpp's implementation, we know this is a true '|' operator, and _not_ an
              // escaped \| or part of a character class like [abcd|].
              if (pattern.value[node.end] && pattern.value[node.end] === '|') {
                  tokens.push({
                      type: 'metaRegexp',
                      range: { start: offset + node.end, end: offset + node.end + 1 },
                      value: '|',
                      kind: MetaRegexpKind.Alternative,
                  })
              }
          },
          onAssertionEnter(node: Assertion) {
              tokens.push({
                  type: 'metaRegexp',
                  range: { start: offset + node.start, end: offset + node.end },
                  value: node.raw,
                  kind: MetaRegexpKind.Assertion,
              })
          },
          onGroupEnter(node: Group) {
              // Push the leading '('
              tokens.push({
                  type: 'metaRegexp',
                  range: { start: offset + node.start, end: offset + node.start + 1 },
                  value: '(',
                  kind: MetaRegexpKind.Delimited,
              })
              // Push the trailing ')'
              tokens.push({
                  type: 'metaRegexp',
                  range: { start: offset + node.end - 1, end: offset + node.end },
                  value: ')',
                  kind: MetaRegexpKind.Delimited,
              })
          },
          onCapturingGroupEnter(node: CapturingGroup) {
              // Push the leading '('
              tokens.push({
                  type: 'metaRegexp',
                  range: { start: offset + node.start, end: offset + node.start + 1 },
                  groupRange: { start: offset + node.start, end: offset + node.end },
                  value: '(',
                  kind: MetaRegexpKind.Delimited,
              })
              // Push the trailing ')'
              tokens.push({
                  type: 'metaRegexp',
                  range: { start: offset + node.end - 1, end: offset + node.end },
                  groupRange: { start: offset + node.start, end: offset + node.end },
                  value: ')',
                  kind: MetaRegexpKind.Delimited,
              })
          },
          onCharacterSetEnter(node: CharacterSet) {
              tokens.push({
                  type: 'metaRegexp',
                  range: { start: offset + node.start, end: offset + node.end },
                  value: node.raw,
                  kind: MetaRegexpKind.CharacterSet,
              })
          },
          onCharacterClassEnter(node: CharacterClass) {
              const negatedOffset = node.negate ? 1 : 0
              // Push the leading '['
              tokens.push({
                  type: 'metaRegexp',
                  range: { start: offset + node.start, end: offset + node.start + 1 + negatedOffset },
                  groupRange: { start: offset + node.start, end: offset + node.end },
                  value: node.negate ? '[^' : '[',
                  kind: MetaRegexpKind.CharacterClass,
              })
              // Push the trailing ']'
              tokens.push({
                  type: 'metaRegexp',
                  range: { start: offset + node.end - 1, end: offset + node.end },
                  groupRange: { start: offset + node.start, end: offset + node.end },
                  value: ']',
                  kind: MetaRegexpKind.CharacterClass,
              })
          },
          onCharacterClassRangeEnter(node: CharacterClassRange) {
              // Push the min and max characters of the range to associate them with the
              // same groupRange for hovers.
              tokens.push(
                  {
                      type: 'metaRegexp',
                      range: { start: offset + node.min.start, end: offset + node.min.end },
                      groupRange: { start: offset + node.start, end: offset + node.end },
                      value: node.raw,
                      kind: MetaRegexpKind.CharacterClassRange,
                  },
                  // Highlight the '-' in [a-z]. Take care to use node.min.end, because we
                  // don't want to highlight the first '-' in [--z], nor an escaped '-' with a
                  // two-character offset as in [\--z].
                  {
                      type: 'metaRegexp',
                      range: { start: offset + node.min.end, end: offset + node.min.end + 1 },
                      groupRange: { start: offset + node.start, end: offset + node.end },
                      value: node.raw,
                      kind: MetaRegexpKind.CharacterClassRangeHyphen,
                  },
                  {
                      type: 'metaRegexp',
                      range: { start: offset + node.max.start, end: offset + node.max.end },
                      groupRange: { start: offset + node.start, end: offset + node.end },
                      value: node.raw,
                      kind: MetaRegexpKind.CharacterClassRange,
                  }
              )
          },
          onQuantifierEnter(node: Quantifier) {
              // the lazy quantifier ? adds one
              const lazyQuantifierOffset = node.greedy ? 0 : 1
              if (!node.greedy) {
                  tokens.push({
                      type: 'metaRegexp',
                      range: { start: offset + node.end - 1, end: offset + node.end },
                      value: '?',
                      kind: MetaRegexpKind.LazyQuantifier,
                  })
              }

              const quantifier = node.raw[node.raw.length - lazyQuantifierOffset - 1]
              if (quantifier === '+' || quantifier === '*' || quantifier === '?') {
                  tokens.push({
                      type: 'metaRegexp',
                      range: {
                          start: offset + node.end - 1 - lazyQuantifierOffset,
                          end: offset + node.end - lazyQuantifierOffset,
                      },
                      value: quantifier,
                      kind: MetaRegexpKind.RangeQuantifier,
                  })
              } else {
                  // regexpp provides no easy way to tell whether the quantifier is a range '{number, number}',
                  // nor the offsets of this range.
                  // At this point we know it is none of +, *, or ?, so it is a ranged quantifier.
                  // We need to then find the opening brace of {number, number}, and go backwards from the end
                  // of this quantifier to avoid dealing with other leading braces that are not part of it.
                  let openBrace = node.end - 1 - lazyQuantifierOffset
                  while (pattern.value[openBrace] && pattern.value[openBrace] !== '{') {
                      openBrace = openBrace - 1
                  }
                  tokens.push({
                      type: 'metaRegexp',
                      range: { start: offset + openBrace, end: offset + node.end - lazyQuantifierOffset },
                      value: pattern.value.slice(openBrace, node.end - lazyQuantifierOffset),
                      kind: MetaRegexpKind.RangeQuantifier,
                  })
              }
          },
          onCharacterEnter(node: Character) {
              if (node.end - node.start > 1 && node.raw.startsWith('\\')) {
                  // This is an escaped value like `\.`, `\u0065`, `\x65`.
                  // If this escaped value is part of a range, like [a-\\],
                  // set the group range to associate it with hovers.
                  const groupRange =
                      node.parent.type === 'CharacterClassRange'
                          ? { start: offset + node.parent.start, end: offset + node.parent.end }
                          : undefined
                  tokens.push({
                      type: 'metaRegexp',
                      range: { start: offset + node.start, end: offset + node.end },
                      groupRange,
                      value: node.raw,
                      kind: MetaRegexpKind.EscapedCharacter,
                  })
                  return
              }
              if (node.parent.type === 'CharacterClassRange') {
                  return // This unescaped character is handled by onCharacterClassRangeEnter.
              }
              if (node.parent.type === 'CharacterClass') {
                  // This character is inside a character class like [abcd] and is contextually special for hover tooltips.
                  tokens.push({
                      type: 'metaRegexp',
                      range: { start: offset + node.start, end: offset + node.end },
                      value: node.raw,
                      kind: MetaRegexpKind.CharacterClassMember,
                  })
                  return
              }
              tokens.push({
                  type: 'pattern',
                  range: { start: offset + node.start, end: offset + node.end },
                  value: node.raw,
                  kind: PatternKind.Regexp,
              })
          },
      })
  } catch {
      return undefined
  }
  // The AST is not necessarily traversed in increasing range. We need
  // to sort by increasing range because the ordering is significant to Monaco.
  tokens.sort((left, right) => left.range.start - right.range.start)
  return coalescePatterns(tokens)
}

const mapRegexpMetaSucceed = (token: Pattern): DecoratedToken[] => mapRegexpMeta(token) || [token]

export const decorate = (token: Token): DecoratedToken[] => {
  const decorated: DecoratedToken[] = []
  switch (token.type) {
      case 'pattern':
          switch (token.kind) {
              case PatternKind.Regexp:
                  decorated.push(...mapRegexpMetaSucceed(token))
                  break
              case PatternKind.Structural:
                  decorated.push(...mapStructuralMeta(token))
                  break
              case PatternKind.Literal:
                  decorated.push(token)
                  break
          }
          break
      case 'filter': {
          decorated.push({
              type: 'field',
              range: { start: token.field.range.start, end: token.field.range.end },
              value: token.field.value,
          })
          decorated.push({
              type: 'metaFilterSeparator',
              range: { start: token.field.range.end, end: token.field.range.end + 1 },
              value: ':',
          })
          const predicate = scanPredicate(token.field.value, token.value?.value || '')
          if (predicate && token.value) {
              decorated.push(...decoratePredicate(predicate, token.value.range))
              break
          }
          if (
              token.value &&
              token.field.value.toLowerCase().match(/^-?(repo|r)$/i) &&
              !token.value.quoted &&
              specifiesRevision(token.value.value)
          ) {
              decorated.push(...decorateRepoRevision(token.value))
          } else if (token.value && token.field.value.toLowerCase().match(/rev|revision/i) && !token.value.quoted) {
              decorated.push(...mapRevisionMeta(createLiteral(token.value.value, token.value.range)))
          } else if (token.value && !token.value.quoted && hasRegexpValue(token.field.value)) {
              // Highlight fields with regexp values.
              if (hasPathLikeValue(token.field.value) && token.value?.type === 'literal') {
                  decorated.push(...mapPathMetaForRegexp(token.value))
              } else {
                  decorated.push(...mapRegexpMetaSucceed(toPattern(token.value)))
              }
          } else if (token.field.value === 'context' && token.value && !token.value.quoted) {
              decorated.push(...decorateContext(token.value))
          } else if (token.field.value === 'select' && token.value && !token.value.quoted) {
              decorated.push(...decorateSelector(token.value))
          } else if (token.value) {
              decorated.push(token.value)
          }
          break
      }
      default:
          decorated.push(token)
  }
  return decorated
}

export function toDecoration(query: string, token: DecoratedToken): Decoration {
  const className = toCSSClassName(token)

  switch (token.type) {
      case 'keyword':
      case 'field':
      case 'metaPath':
      case 'metaRevision':
      case 'metaRegexp':
      case 'metaStructural':
          return {
              value: token.value,
              key: token.range.start + token.range.end,
              className,
          }
      case 'openingParen':
          return {
              value: '(',
              key: token.range.start + token.range.end,
              className,
          }
      case 'closingParen':
          return {
              value: ')',
              key: token.range.start + token.range.end,
              className,
          }

      case 'metaFilterSeparator':
          return {
              value: ':',
              key: token.range.start + token.range.end,
              className,
          }
      case 'metaRepoRevisionSeparator':
      case 'metaContextPrefix':
          return {
              value: '@',
              key: token.range.start + token.range.end,
              className,
          }

      case 'metaPredicate': {
          let value = ''
          switch (token.kind) {
              case 'NameAccess':
                  value = query.slice(token.range.start, token.range.end)
                  break
              case 'Dot':
                  value = '.'
                  break
              case 'Parenthesis':
                  value = query.slice(token.range.start, token.range.end)
                  break
          }
          return {
              value,
              key: token.range.start + token.range.end,
              className,
          }
      }
  }
  return {
      value: query.slice(token.range.start, token.range.end),
      key: token.range.start + token.range.end,
      className,
  }
}