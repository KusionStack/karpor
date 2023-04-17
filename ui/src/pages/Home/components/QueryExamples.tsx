import React, { useMemo } from "@alipay/bigfish/react";
import classNames from "classnames";
import { Badge, Button } from "antd";
import { toDecoration } from "@/utils/tools";
import styles from "./styles.less";


export interface QueryExamplesSection {
  title: string
  productStatus?: string
  queryExamples: {
    id: string
    query: string
    helperText?: string
    slug?: string
  }[]
}

interface QueryExamplesLayout {
  queryColumns: QueryExamplesSection[][]
  onQueryExampleClick: (id: string | undefined, query: string, slug: string | undefined) => void
}

export enum SearchPatternType {
  keyword = 'keyword',
  literal = 'literal',
  lucky = 'lucky',
  regexp = 'regexp',
  standard = 'standard',
  structural = 'structural',
}



export const toPatternResult =
  (scanner: Scanner<Literal>, kind: PatternKind, delimited = false): Scanner<Pattern> =>
    (input, start) => {
      const result = scanner(input, start)
      if (result.type === 'success') {
        return createPattern(result.term.value, result.term.range, kind, result.term.quoted)
      }
      return result
    }

const scanPattern = (kind: PatternKind): Scanner<Pattern> =>
  toPatternResult(oneOf<Literal>(scanBalancedLiteral, literal), kind)


const keyword = oneOf<Keyword>(keywordAnd, keywordOr, keywordNot)

const scanStandard = (query: string): ScanResult<Token[]> => {
  const tokenScanner = [
    keyword,
    filter,
    toPatternResult(quoted('/'), PatternKind.Regexp),
    scanPattern(PatternKind.Literal),
  ]
  const earlyPatternScanner = [
    toPatternResult(quoted('/'), PatternKind.Regexp),
    toPatternResult(scanBalancedLiteral, PatternKind.Literal),
  ]

  const scan = zeroOrMore(
    oneOf<Term>(
      whitespace,
      ...earlyPatternScanner.map(token => followedBy(token, whitespaceOrClosingParen)),
      openingParen,
      closingParen,
      ...tokenScanner.map(token => followedBy(token, whitespaceOrClosingParen))
    )
  )

  return scan(query, 0)
}

function detectPatternType(query: string): SearchPatternType | undefined {
  const result = scanStandard(query)
  const tokens =
    result.type === 'success'
      ? result.term.filter(
        token => !!(token.type === 'filter' && token.field.value.toLowerCase() === 'patterntype')
      )
      : undefined
  if (tokens && tokens.length > 0) {
    return (tokens[0] as Filter).value?.value as SearchPatternType
  }
  return undefined
}


export const scanSearchQuery = (
  query: string,
  interpretComments?: boolean,
  searchPatternType = SearchPatternType.literal
): ScanResult<Token[]> => {
  const patternType = detectPatternType(query) || searchPatternType
  let patternKind
  switch (patternType) {
    case SearchPatternType.standard:
    case SearchPatternType.lucky:
    case SearchPatternType.keyword:
      return scanStandard(query)
    case SearchPatternType.literal:
      patternKind = PatternKind.Literal
      break
    case SearchPatternType.regexp:
      patternKind = PatternKind.Regexp
      break
    case SearchPatternType.structural:
      patternKind = PatternKind.Structural
      break
  }
  const scanner = createScanner(patternKind, interpretComments)
  return scanner(query, 0)
}

export function decorateQuery(query: string, searchPatternType?: SearchPatternType): Decoration[] | null {
  const tokens = searchPatternType ? scanSearchQuery(query, false, searchPatternType) : scanSearchQuery(query)
  return tokens.type === 'success'
    ? tokens.term.flatMap(token => decorate(token).map(token => toDecoration(query, token)))
    : null
}

export const SyntaxHighlightedSearchQuery: React.FunctionComponent<
  React.PropsWithChildren<SyntaxHighlightedSearchQueryProps>
> = ({ query, searchPatternType, ...otherProps }) => {
  const tokens = useMemo(() => {
    const decorations = decorateQuery(query, searchPatternType)

    return decorations
      ? decorations.map(({ value, key, className }) => (
        <span className={className} key={key}>
          {value}
        </span>
      ))
      : [<React.Fragment key="0">{query}</React.Fragment>]
  }, [query, searchPatternType])

  return (
    <span {...otherProps} className={classNames('text-monospace search-query-link', otherProps.className)}>
      {tokens}
    </span>
  )
}

export const QueryExampleChip: React.FunctionComponent<QueryExampleChipProps> = ({
  id,
  query,
  helperText,
  slug,
  className,
  onClick,
}) => (
  <li className={classNames(className)} style={{ display: 'flex', alignItems: "center" }}>
    <Button type="button" className={styles.queryExampleChip} onClick={() => onClick(id, query, slug || '')}>
      <SyntaxHighlightedSearchQuery query={query} searchPatternType={SearchPatternType.standard} />
    </Button>
    {helperText && (
      <span className="text-muted ml-2">
        <small>{helperText}</small>
      </span>
    )}
  </li>
)

export const ExamplesSection: React.FunctionComponent<ExamplesSection> = ({
  title,
  productStatus,
  queryExamples,
  onQueryExampleClick,
}) => (
  <div className={styles.queryExamplesSection}>
    <h2 className={styles.queryExamplesSectionTitle}>
      {title}
      {productStatus && (
        <>
          {' '}
          <Badge status={productStatus} />
        </>
      )}
    </h2>
    <ul className={classNames('list-unstyled', styles.queryExamplesItems)}>
      {queryExamples
        .filter(({ query }) => query.length > 0)
        .map(({ id, query, helperText, slug }) => (
          <QueryExampleChip
            id={id}
            key={query}
            query={query}
            slug={slug}
            helperText={helperText}
            onClick={onQueryExampleClick}
          />
        ))}
    </ul>
  </div>
)

export const QueryExamplesLayout: React.FunctionComponent<QueryExamplesLayout> = ({
  queryColumns,
  onQueryExampleClick,
}) => (
  <div className={styles.queryExamplesSectionsColumns}>
    {queryColumns.map((column, index) => (
      <div key={`column-${queryColumns[index][0].title}`}>
        {column.map(({ title, productStatus, queryExamples }) => (
          <ExamplesSection
            key={title}
            title={title}
            productStatus={productStatus}
            queryExamples={queryExamples}
            onQueryExampleClick={onQueryExampleClick}
          />
        ))}
      </div>
    ))}
  </div>
)