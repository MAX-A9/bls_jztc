import request from '@/utils/request';

// 通用响应结构
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 评论项接口
export interface CommentItem {
  id: number;
  contentId: number;
  contentTitle: string;
  clientId: number;
  realName: string;
  comment: string;
  status: string;
  statusText: string;
  createdAt: string;
  updatedAt: string;
}

// 评论列表请求参数
export interface CommentListReq {
  page?: number;
  pageSize?: number;
  contentId?: number;
  status?: string;
  realName?: string;
  comment?: string;
}

// 评论列表响应
export interface CommentListRes {
  code: number;
  message: string;
  data: {
    list: CommentItem[];
    total: number;
    page: number;
  };
}

// 评论详情响应
export interface CommentDetailItem {
  id: number;
  contentId: number;
  contentTitle: string;
  clientId: number;
  realName: string;
  comment: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

// 评论详情响应
export interface CommentDetailRes {
  code: number;
  message: string;
  data: CommentDetailItem;
}

// 创建评论请求参数
export interface CommentCreateReq {
  contentId: number;
  clientId: number;
  realName: string;
  comment: string;
  status?: string;
}

// 创建评论响应
export interface CommentCreateRes {
  code: number;
  message: string;
  data: {
    id: number;
  };
}

// 更新评论请求参数
export interface CommentUpdateReq {
  id: number;
  comment: string;
  realName: string;
  status?: string;
}

// 更新评论响应
export interface CommentUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 删除评论响应
export interface CommentDeleteRes {
  code: number;
  message: string;
  data: null;
}

// 更新评论状态请求参数
export interface CommentStatusUpdateReq {
  id: number;
  status: string;
  realName: string;
}

// 更新评论状态响应
export interface CommentStatusUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 内容评论列表响应
export interface ContentCommentsRes {
  code: number;
  message: string;
  data: {
    list: CommentItem[];
    total: number;
    page: number;
  };
}

/**
 * 获取评论列表
 * @param params 查询参数
 * @returns 评论列表
 */
export const getCommentList = (params: CommentListReq): Promise<CommentListRes> => {
  return request.get('/comment/list', { params });
};

/**
 * 获取指定内容的评论列表
 * @param contentId 内容ID
 * @param params 分页参数
 * @returns 评论列表
 */
export const getContentComments = (contentId: number, params?: { page?: number; pageSize?: number }): Promise<ContentCommentsRes> => {
  return request.get('/comment/content-comments', { 
    params: { 
      contentId,
      ...params
    } 
  });
};

/**
 * 获取评论详情
 * @param id 评论ID
 * @returns 评论详情
 */
export const getCommentDetail = (id: number): Promise<CommentDetailRes> => {
  return request.get('/comment/detail', { params: { id } });
};

/**
 * 创建评论
 * @param data 评论数据
 * @returns 创建结果
 */
export const createComment = (data: CommentCreateReq): Promise<CommentCreateRes> => {
  return request({
    url: '/comment/create',
    method: 'post',
    data,
  });
};

/**
 * 更新评论
 * @param data 评论数据
 * @returns 更新结果
 */
export const updateComment = (data: CommentUpdateReq): Promise<CommentUpdateRes> => {
  return request({
    url: '/comment/update',
    method: 'put',
    data,
  });
};

/**
 * 删除评论
 * @param id 评论ID
 * @returns 删除结果
 */
export const deleteComment = (id: number): Promise<CommentDeleteRes> => {
  return request.delete('/comment/delete', { params: { id } });
};

/**
 * 更新评论状态
 * @param data 状态数据
 * @returns 更新结果
 */
export const updateCommentStatus = (data: CommentStatusUpdateReq): Promise<CommentStatusUpdateRes> => {
  return request({
    url: '/comment/status/update',
    method: 'put',
    data,
  });
}; 