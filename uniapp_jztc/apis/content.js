import { get, post, put, del } from '@/utils/request'

// 通用处理函数：确保页码不小于1
function ensureValidPageNumber(params) {
  if (params && params.page !== undefined) {
    params.page = Math.max(1, parseInt(params.page) || 1);
  }
  return params;
}

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

// 获取内容详情
export function getContentDetail(id) {
  return get('/detail', { id })
}

// 获取公开内容详情
export function getPublicContentDetail(id) {
  return get('/wx/client/content/public/detail', { id }, true)
}

// 发布信息
export function createInfo(data) {
  return post('/wx/client/content/info/create', data)
}

// 更新内容
export function updateContent(id, data) {
  return put(`/update/${id}`, data)
}

// 删除内容
export function deleteContent(id) {
  return del(`/delete/${id}`)
}

// 收藏内容
export function collectContent(id) {
  return post(`/collect/${id}`)
}

// 取消收藏
export function uncollectContent(id) {
  return del(`/collect/${id}`)
}

// 关注发布人
export const followUser = (data) => {
	return post('/wx/publisher/follow', {
		publisher_id: data.publisher_id
	})
}

// 取消关注发布人
export const unfollowUser = (data) => {
	return post('/wx/publisher/unfollow', {
		publisher_id: data.publisher_id
	})
}

// 评论内容
export function commentContent(id, data) {
  return post(`/comment/${id}`, data)
}

// 举报内容
export function reportContent(id, data) {
  return post(`/report/${id}`, data)
}

// 发布闲置物品
export function createIdle(data) {
  return post('/wx/client/content/idle/create', data)
}

// 获取区域内容列表
export function getRegionContentList(params) {
  return get('/wx/client/content/region/list', ensureValidPageNumber(params), true)
}

// 创建评论
export function createComment(data) {
  return post('/wx/client/comment/create', data)
}

// 获取评论列表
export function getCommentList(contentId, params) {
  return get('/wx/client/comment/list', { contentId, ...params }, true)
}

// 添加收藏
export function addFavorite(contentId) {
  return post('/wx/favorite/add', { contentId })
}

// 取消收藏
export function cancelFavorite(contentId) {
  return post('/wx/favorite/cancel', { contentId })
}

// 获取收藏状态
export function getFavoriteStatus(contentId) {
  return get('/wx/favorite/status', { contentId })
}

// 获取收藏列表
export function getFavoriteList(params) {
  return get('/wx/favorite/list', params)
}

// 获取收藏数量
export function getFavoriteCount() {
  return get('/wx/favorite/count')
}

// 添加浏览记录
export function addBrowseRecord(contentId, contentType = 'article') {
  return post('/wx/client/browse-history/add', { contentId, contentType })
}

// 获取浏览历史列表
export function getBrowseHistoryList(params) {
  return get('/wx/client/browse-history/list', params)
}

// 清空浏览历史
export function clearBrowseHistory(timeType = 'all') {
  return post('/wx/client/browse-history/clear', { timeType })
}

// 获取浏览记录数量
export function getBrowseHistoryCount() {
  return get('/wx/client/browse-history/count')
}

// 获取发布者信息
export function getPublisherInfo(publisherId) {
	return get('/wx/publisher/info', {
		publisher_id: publisherId
	}, true)
}

// 获取发布者关注状态
export function getPublisherFollowStatus(publisherId) {
	return get('/wx/publisher/follow/status', {
		publisher_id: publisherId
	})
}

// 获取我关注的发布者列表
export function getFollowingList(params) {
  return get('/wx/publisher/following/list', {
    page: params.page || 1,
    size: params.size || 10
  })
}

// 获取我关注的发布者数量
export function getFollowingCount() {
  return get('/wx/publisher/following/count')
}

// 获取区域闲置列表
export function getRegionIdleList(params) {
  return get('/wx/client/content/region/idle/list', ensureValidPageNumber(params), true)
}

// 获取套餐列表
export function getPackageList() {
  return get('/wx/client/package/list', {}, true)
}

// 获取小程序列表
export function getMiniProgramList() {
  return get('/wx/mini-program/list', {}, true)
} 