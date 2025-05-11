/**
 * 发布相关API
 */
import { get, post, put, del } from '@/utils/request'

// 发布信息
export function createInfo(data) {
  return post('/wx/client/content/info/create', data)
}

// 发布闲置物品
export function createIdle(data) {
  return post('/wx/client/content/idle/create', data)
}

// 获取套餐列表
export function getPackageList() {
  return get('/wx/client/package/list', {}, true)
}

// 更新内容
export function updateContent(id, data) {
  return put(`/update/${id}`, data)
}

// 删除内容
export function deleteContent(id) {
  return del(`/delete/${id}`)
}

// 获取用户发布的内容列表
export function getUserPublishList(params) {
  return get('/wx/client/content/user/list', params)
}

// 获取我的发布列表
export function getMyPublishList(params) {
  return get('/wx/client/publish/list', params)
}

// 获取发布数量
export function getPublishCount() {
  return get('/wx/client/publish/count')
}

// 默认导出所有API
export default {
  createInfo,
  createIdle,
  getPackageList,
  updateContent,
  deleteContent,
  getUserPublishList,
  getMyPublishList,
  getPublishCount
} 