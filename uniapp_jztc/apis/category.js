/**
 * 分类相关API
 */
import { get, post } from '@/utils/request'

// 获取内容分类列表
export function getCategories(type = 1) {
  return get('/categories', { type })
}

// 获取信息分类列表
export function getInfoCategories(type = 1) {
  return get('/wx/client/content/categories', { type }, true)
}

// 按分类获取内容列表
export function getCategoryList(params) {
  return get('/category/list', params)
}

// 获取闲置物品分类
export function getIdleCategories() {
  return getInfoCategories(2) // type=2表示闲置物品分类
}

// 默认导出所有API
export default {
  getCategories,
  getInfoCategories,
  getCategoryList,
  getIdleCategories
} 