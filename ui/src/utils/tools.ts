import yaml from 'js-yaml'
import moment from 'moment'
import _ from 'lodash'

export function truncationPageData({ list, page, pageSize }) {
  return list && list?.length > 0
    ? list?.slice((page - 1) * pageSize, page * pageSize)
    : []
}

export function utcDateToLocalDate(date, FMT = 'YYYY-MM-DD HH:mm:ss') {
  return moment.utc(date).local().format(FMT)
}

function pow1024(num) {
  return Math.pow(1024, num)
}

/**
 * @param size
 * @returns {string|*}
 */
export const filterSize = size => {
  if (!size) return ''
  if (size < pow1024(1)) return size + ' B'
  if (size < pow1024(2)) return (size / pow1024(1)).toFixed(2) + ' KB'
  if (size < pow1024(3)) return (size / pow1024(2)).toFixed(2) + ' MB'
  if (size < pow1024(4)) return (size / pow1024(3)).toFixed(2) + ' GB'
  return (size / pow1024(4)).toFixed(2) + ' TB'
}

export function format_with_regex(number) {
  return !(number + '').includes('.')
    ? (number + '').replace(/\d{1,3}(?=(\d{3})+$)/g, match => match + ',')
    : (number + '').replace(/\d{1,3}(?=(\d{3})+(\.))/g, match => match + ',')
}

export const yaml2json = (yamlStr: string) => {
  try {
    return {
      data: yaml.load(yamlStr),
      error: false,
    }
  } catch (err) {
    return {
      data: '',
      error: true,
    }
  }
}

export function generateTopologyData(data) {
  const nodes = []
  const edges = []
  const edgeSet = new Set()

  for (const key in data) {
    nodes.push({
      id: key,
      label: key,
      data: {
        count: data[key].count,
        resourceGroup: data[key].resourceGroup,
      },
    })
  }

  function addEdge(source, target) {
    const edgeId = `${source}->${target}`
    if (!edgeSet.has(edgeId)) {
      edges.push({
        source,
        target,
      })
      edgeSet.add(edgeId)
    }
  }

  for (const key in data) {
    const relationships = data[key].relationship
    for (const targetKey in relationships) {
      const relationType = relationships[targetKey]
      if (relationType === 'child') {
        addEdge(key, targetKey)
      } else if (relationType === 'parent') {
        addEdge(targetKey, key)
      }
    }
  }

  return { nodes, edges }
}

export function generateResourceTopologyData(data) {
  const nodes = []
  const edges = []

  const addNode = (id, label, resourceGroup) => {
    nodes.push({ id, label, resourceGroup })
  }

  const uniqueEdges = new Set()
  const addEdge = (source, target) => {
    const edgeKey = `${source}=>${target}`
    if (!uniqueEdges.has(edgeKey)) {
      edges.push({ source, target })
      uniqueEdges.add(edgeKey)
    }
  }

  Object.keys(data).forEach(key => {
    const entity = data[key]

    addNode(key, key.split(':')[1].split('.')[1], entity?.resourceGroup)

    entity.children.forEach(child => {
      addEdge(key, child)
    })

    if (entity.Parents) {
      entity.Parents.forEach(parent => {
        addEdge(parent, key)
      })
    }
  })
  return { nodes, edges }
}

export function getDataType(data) {
  const map = new Map()
  map.set('[object String]', 'String')
  map.set('[object Number]', 'Number')
  map.set('[object Boolean]', 'Boolean')
  map.set('[object Symbol]', 'Symbol')
  map.set('[object Undefined]', 'Undefined')
  map.set('[object Null]', 'Null')
  map.set('[object Function]', 'Function')
  map.set('[object Date]', 'Date')
  map.set('[object Array]', 'Array')
  map.set('[object RegExp]', 'RegExp')
  map.set('[object Error]', 'Error')
  map.set('[object HTMLDocument]', 'HTMLDocument')
  const type = Object.prototype.toString.call(data)
  return map.get(type)
}

export function capitalized(word) {
  return word.charAt(0).toUpperCase() + word?.slice(1)
}

export function isEmptyObject(obj) {
  return obj && Object.keys(obj)?.length === 0
}

export function hasDuplicatesOfObjectArray(array) {
  const seen = new Set()
  for (const item of array) {
    const serialized = JSON.stringify(item)
    if (seen.has(serialized)) {
      return true
    }
    seen.add(serialized)
  }
  return false
}

export function filterKeywordsOfArray(list, keywords, attribute) {
  const result = []
  if (keywords?.length === 1) {
    list?.forEach((item: any) => {
      if (_.get(item, attribute)?.toLowerCase()?.includes(keywords?.[0])) {
        result.push(item)
      }
    })
  } else {
    list?.forEach((item: any) => {
      if (
        keywords?.every((innerValue: string) =>
          _.get(item, attribute)?.toLowerCase()?.includes(innerValue),
        )
      ) {
        result.push(item)
      }
    })
  }
  return result
}

export function getTextSizeByCanvas(
  text: string,
  fontSize: number,
  fontFamily?: string,
) {
  const canvas = document.createElement('canvas')
  const context = canvas.getContext('2d')
  context.font = `${fontSize}px ${fontFamily}`
  const metrics = context.measureText(text)
  const width = Math.ceil(metrics.width)
  canvas.remove()
  return width + 2
}

export function getHistoryList(key) {
  return localStorage?.getItem(key)
    ? JSON.parse(localStorage?.getItem(key))
    : []
}

export function deleteHistoryByItem(key, val: string) {
  const lastHistory: any = localStorage.getItem(key)
  const tmp = lastHistory ? JSON.parse(lastHistory) : []
  if (tmp?.length > 0 && tmp?.includes(val)) {
    const newList = tmp?.filter(item => item !== val)
    localStorage.setItem(key, JSON.stringify(newList))
  }
}

export function cacheHistory(key, val: string) {
  const lastHistory: any = localStorage.getItem(key)
  const tmp = lastHistory ? JSON.parse(lastHistory) : []
  const newList = [val, ...tmp?.filter(item => item !== val)]
  localStorage.setItem(key, JSON.stringify(newList))
  return getHistoryList(key)
}
