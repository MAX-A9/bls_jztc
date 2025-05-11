import request from '@/utils/request';

// 通用响应结构
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 分类项接口
export interface CategoryItem {
  id: number;
  name: string;
  isActive: boolean;
  sortOrder: number;
  icon?: string;
  createdAt: string;
  updatedAt: string;
}

// 首页分类列表响应
export interface HomeCategoryListRes {
  code: number;
  message: string;
  data: {
    list: CategoryItem[];
  };
}

// 创建首页分类请求参数
export interface HomeCategoryCreateReq {
  name: string;
  isActive?: boolean;
  sortOrder?: number;
  icon?: string;
}

// 创建首页分类响应
export interface HomeCategoryCreateRes {
  code: number;
  message: string;
  data: {
    id: number;
  };
}

// 更新首页分类请求参数
export interface HomeCategoryUpdateReq {
  id: number;
  name: string;
  isActive?: boolean;
  sortOrder?: number;
  icon?: string;
}

// 更新首页分类响应
export interface HomeCategoryUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 删除首页分类响应
export interface HomeCategoryDeleteRes {
  code: number;
  message: string;
  data: null;
}

/**
 * 获取首页分类列表
 * @returns 首页分类列表
 */
export const getHomeCategoryList = (): Promise<HomeCategoryListRes> => {
  return request.get('/home-category/list');
};

/**
 * 创建首页分类
 * @param data 分类数据
 * @returns 创建结果
 */
export const createHomeCategory = (data: HomeCategoryCreateReq): Promise<HomeCategoryCreateRes> => {
  return request({
    url: '/home-category/create',
    method: 'post',
    data,
  });
};

/**
 * 更新首页分类
 * @param data 分类数据
 * @returns 更新结果
 */
export const updateHomeCategory = (data: HomeCategoryUpdateReq): Promise<HomeCategoryUpdateRes> => {
  return request({
    url: '/home-category/update',
    method: 'put',
    data,
  });
};

/**
 * 删除首页分类
 * @param id 分类ID
 * @returns 删除结果
 */
export const deleteHomeCategory = (id: number): Promise<HomeCategoryDeleteRes> => {
  return request.delete('/home-category/delete', { params: { id } });
};

// ===== 闲置分类相关接口 =====

// 闲置分类列表响应
export interface IdleCategoryListRes {
  code: number;
  message: string;
  data: {
    list: CategoryItem[];
  };
}

// 创建闲置分类请求参数
export interface IdleCategoryCreateReq {
  name: string;
  isActive?: boolean;
  sortOrder?: number;
  icon?: string;
}

// 创建闲置分类响应
export interface IdleCategoryCreateRes {
  code: number;
  message: string;
  data: {
    id: number;
  };
}

// 更新闲置分类请求参数
export interface IdleCategoryUpdateReq {
  id: number;
  name: string;
  isActive?: boolean;
  sortOrder?: number;
  icon?: string;
}

// 更新闲置分类响应
export interface IdleCategoryUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 删除闲置分类响应
export interface IdleCategoryDeleteRes {
  code: number;
  message: string;
  data: null;
}

/**
 * 获取闲置分类列表
 * @returns 闲置分类列表
 */
export const getIdleCategoryList = (): Promise<IdleCategoryListRes> => {
  return request.get('/idle-category/list');
};

/**
 * 创建闲置分类
 * @param data 分类数据
 * @returns 创建结果
 */
export const createIdleCategory = (data: IdleCategoryCreateReq): Promise<IdleCategoryCreateRes> => {
  return request({
    url: '/idle-category/create',
    method: 'post',
    data,
  });
};

/**
 * 更新闲置分类
 * @param data 分类数据
 * @returns 更新结果
 */
export const updateIdleCategory = (data: IdleCategoryUpdateReq): Promise<IdleCategoryUpdateRes> => {
  return request({
    url: '/idle-category/update',
    method: 'put',
    data,
  });
};

/**
 * 删除闲置分类
 * @param id 分类ID
 * @returns 删除结果
 */
export const deleteIdleCategory = (id: number): Promise<IdleCategoryDeleteRes> => {
  return request.delete('/idle-category/delete', { params: { id } });
};

// ===== 所有分类相关接口 =====

/**
 * 所有分类项（包含分类类型）
 */
export interface CategoryWithType extends CategoryItem {
  type: 'home' | 'idle'; // home: 首页分类, idle: 闲置分类
}

/**
 * 获取所有分类响应
 */
export interface GetCategoriesRes {
  code: number;
  message: string;
  data: {
    homeCategories: CategoryItem[]; // 首页分类
    idleCategories: CategoryItem[]; // 闲置分类
  };
}

/**
 * 获取所有分类列表（包括首页和闲置）
 * @returns 所有分类列表
 */
export const getAllCategories = (): Promise<GetCategoriesRes> => {
  return request.get('/categories');
};

/**
 * 获取所有分类并转换为Select组件可用的格式
 * @returns 转换后的分类列表，用于Select组件
 */
export const getCategoriesForSelect = async (): Promise<{label: string, value: string, type: string}[]> => {
  const res = await getAllCategories();
  
  console.log('API返回的分类数据:', res);
  
  if (res.code === 0) {
    const { homeCategories, idleCategories } = res.data;
    
    // 合并并转换为Select组件可用的格式
    const options: {label: string, value: string, type: string}[] = [];
    
    // 添加首页分类
    if (homeCategories && homeCategories.length > 0) {
      options.push({ label: '首页分类', value: 'home-group', type: 'group' });
      homeCategories.forEach(item => {
        if (item.isActive) { // 只添加启用的分类
          options.push({
            label: item.name,
            value: `home:${item.name}`,
            type: 'home'
          });
        }
      });
    }
    
    // 添加闲置分类
    if (idleCategories && idleCategories.length > 0) {
      options.push({ label: '闲置分类', value: 'idle-group', type: 'group' });
      idleCategories.forEach(item => {
        if (item.isActive) { // 只添加启用的分类
          options.push({
            label: item.name,
            value: `idle:${item.name}`,
            type: 'idle'
          });
        }
      });
    }
    
    console.log('转换后的分类选项:', options);
    return options;
  }
  
  return [];
}; 