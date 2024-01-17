import yaml from 'js-yaml';
import moment from 'moment';


export function truncationPageData({ list, page, pageSize }) {
  return list && list?.length > 0 ? list?.slice((page - 1) * pageSize, page * pageSize) : [];
}

export function utcDateToLocalDate(date, FMT = 'YYYY-MM-DD HH:mm:ss') {
  return moment.utc(date).local().format(FMT)
}

// 求次幂
function pow1024(num) {
  return Math.pow(1024, num)
}

/**
 * 文件大小 字节转换单位
 * @param size
 * @returns {string|*}
 */
export const filterSize = (size) => {

  if (!size) return '';
  if (size < pow1024(1)) return size + ' B';
  if (size < pow1024(2)) return (size / pow1024(1)).toFixed(2) + ' KB';
  if (size < pow1024(3)) return (size / pow1024(2)).toFixed(2) + ' MB';
  if (size < pow1024(4)) return (size / pow1024(3)).toFixed(2) + ' GB';
  return (size / pow1024(4)).toFixed(2) + ' TB'
}




export function format_with_regex(number) {
  return !(number + '').includes('.')
    ? // 就是说1-3位后面一定要匹配3位
    (number + '').replace(/\d{1,3}(?=(\d{3})+$)/g, (match) => {
      return match + ',';
    })
    : (number + '').replace(/\d{1,3}(?=(\d{3})+(\.))/g, (match) => {
      return match + ',';
    });
}

export const yaml2json = (yamlStr: string) => {
  try {
    return {
      data: yaml.load(yamlStr),
      error: false,
    };
  } catch (err) {
    return {
      data: '',
      error: true,
    };
  }
};


export function generateTopologyData(data) {
  const nodes = [];
  const edges = [];
  const edgeSet = new Set(); // 用来记录已经创建的边，避免重复

  // 创建节点
  for (const key in data) {
    nodes.push({
      id: key,
      label: key,
      data: {
        count: data[key].count,
        locator: data[key].locator,
      },
    });
  }

  // 创建边，同时防止重复
  function addEdge(source, target) {
    const edgeId = `${source}->${target}`;
    if (!edgeSet.has(edgeId)) {
      edges.push({
        source,
        target,
      });
      edgeSet.add(edgeId); // 添加到集合中标记为已创建
    }
  }

  // 遍历数据生成边
  for (const key in data) {
    const relationships = data[key].relationship;
    for (const targetKey in relationships) {
      const relationType = relationships[targetKey];
      if (relationType === "child") {
        addEdge(key, targetKey); // 添加父到子的边
      } else if (relationType === "parent") {
        addEdge(targetKey, key); // 添加子到父的边
      }
    }
  }


  return { nodes, edges }
}

export function generateResourceTopologyData(data) {
  // Initialize nodes and edges arrays
  const nodes = [];
  const edges = [];

  // A helper function to add nodes
  const addNode = (id, label, locator) => {
    nodes.push({ id, label, locator });
  };

  // A helper function to add edges with duplication check
  const uniqueEdges = new Set();
  const addEdge = (source, target) => {
    const edgeKey = `${source}=>${target}`;
    if (!uniqueEdges.has(edgeKey)) {
      edges.push({ source, target });
      uniqueEdges.add(edgeKey);
    }
  };

  // Iterate over the data and populate nodes and edges
  Object.keys(data).forEach(key => {
    const entity = data[key];

    // Add the entity as a node
    addNode(key, key.split(':')[1].split('.')[1], entity?.locator);

    // Add edges for all children with duplication check
    entity.children.forEach(child => {
      addEdge(key, child);
    });

    // Add edges for all parents with duplication check
    if (entity.Parents) {
      entity.Parents.forEach(parent => {
        addEdge(parent, key);
      });
    }
  });
  return { nodes, edges }
}

export function getDataType(data) {
  const map = new Map();
  map.set('[object String]', 'String');
  map.set('[object Number]', 'Number');
  map.set('[object Boolean]', 'Boolean');
  map.set('[object Symbol]', 'Symbol');
  map.set('[object Undefined]', 'Undefined');
  map.set('[object Null]', 'Null');
  map.set('[object Function]', 'Function');
  map.set('[object Date]', 'Date');
  map.set('[object Array]', 'Array');
  map.set('[object RegExp]', 'RegExp');
  map.set('[object Error]', 'Error');
  map.set('[object HTMLDocument]', 'HTMLDocument');
  const type = Object.prototype.toString.call(data);
  return map.get(type);
}

