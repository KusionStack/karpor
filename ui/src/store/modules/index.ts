/**
 * 合并所有slice
 */

// 多个Slice的引入；
import globalSlice from './globalSlice'

// globalSlice: 表示globalSlice的模块名称  store.globalSlice.xxx 就可以取到globalSlice管理的数据
export default {
  globalSlice: globalSlice,
}
