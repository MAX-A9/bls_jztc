import request from '../utils/request';

// 产品列表查询参数
export interface ProductListParams {
  page?: number;
  size?: number;
  sortField?: 'id' | 'name' | 'price' | 'stock' | 'duration' | 'sales' | 'sort_order' | 'created_at';
  name?: string;
  categoryId?: number;
  status?: -1 | 0 | 1 | 2; // 状态 -1:全部 0:未上架 1:已上架 2:已售罄
  duration?: number;
  stock?: number;
  tags?: string;
}

// 产品数据接口
export interface ProductItem {
  id: number;
  name: string;
  description: string;
  thumbnail: string;
  images: string;
  price: number;
  stock: number;
  categoryId: number;
  categoryName: string;
  status: number;
  duration: number;
  tags?: string;
  sales: number;
  sortOrder: number;
  createdAt: string;
  updatedAt: string | null;
}

// 商品分类数据接口
export interface ShopCategory {
  id: number;
  name: string;
  productCount: number;
  sortOrder: number;
  status: number; // 状态 0:禁用 1:启用
  image: string; // 分类图片URL
  createdAt?: string;
  updatedAt?: string;
}

// 商品分类列表查询参数
export interface ShopCategoryListParams {
  page?: number;
  size?: number;
  name?: string;
  status?: number; // 状态 0:禁用 1:启用
}

// 商品分类列表响应
export interface ShopCategoryListRes {
  code: number;
  message: string;
  data: {
    list: ShopCategory[];
    total: number;
    page: number;
    size: number;
    pages: number;
  };
}

// 创建商品分类请求参数
export interface CreateShopCategoryReq {
  name: string;
  sortOrder: number;
  status: number;
  image: string;
  productCount: number;
}

// 更新商品分类请求参数
export interface UpdateShopCategoryReq {
  id: number;
  name: string;
  sortOrder: number;
  status: number;
  image: string;
  productCount: number;
}

// 更新商品分类状态请求参数
export interface UpdateShopCategoryStatusReq {
  id: number;
  status: number;
}

// 通用响应接口
export interface CommonResponse {
  code: number;
  message: string;
  data: any;
}

/**
 * 获取产品列表
 * @param params 查询参数
 * @returns 产品列表
 */
export const getProductList = (params: ProductListParams) => {
  return request.get('/product/list', { params });
};

/**
 * 获取产品详情
 * @param id 产品ID
 * @returns 产品详情
 */
export const getProductDetail = (id: number) => {
  return request.get(`/product/detail/${id}`);
};

/**
 * 创建产品
 * @param params 产品信息
 * @returns 创建结果
 */
export const createProduct = (params: any) => {
  return request.post('/product/create', params);
};

/**
 * 更新产品
 * @param params 产品信息
 * @returns 更新结果
 */
export const updateProduct = (params: any) => {
  return request.post('/product/update', params);
};

/**
 * 删除产品
 * @param id 产品ID
 * @returns 删除结果
 */
export const deleteProduct = (id: number) => {
  return request.post('/product/delete', { id });
};

/**
 * 更新产品状态
 * @param id 商品ID
 * @param status 状态
 * @returns 更新结果
 */
export const updateProductStatus = (id: number, status: number) => {
  return request.post('/product/status', { id, status });
};

/**
 * 获取商品分类列表
 * @param params 查询参数
 * @returns 商品分类列表
 */
export const getShopCategoryList = (params?: ShopCategoryListParams): Promise<ShopCategoryListRes> => {
  return request.get('/shop-category/list', { params });
};

/**
 * 获取商品分类详情
 * @param id 分类ID
 * @returns 分类详情
 */
export const getShopCategory = (id: number): Promise<CommonResponse> => {
  return request.get('/shop-category/get', { params: { id } });
};

/**
 * 创建商品分类
 * @param data 分类数据
 * @returns 创建结果
 */
export const createShopCategory = (data: CreateShopCategoryReq): Promise<CommonResponse> => {
  return request.post('/shop-category/create', data);
};

/**
 * 更新商品分类
 * @param data 分类数据
 * @returns 更新结果
 */
export const updateShopCategory = (data: UpdateShopCategoryReq): Promise<CommonResponse> => {
  return request.put('/shop-category/update', data);
};

/**
 * 删除商品分类
 * @param id 分类ID
 * @returns 删除结果
 */
export const deleteShopCategory = (id: number): Promise<CommonResponse> => {
  return request.delete('/shop-category/delete', { data: { id } });
};

/**
 * 更新商品分类状态
 * @param data 更新数据
 * @returns 更新结果
 */
export const updateShopCategoryStatus = (data: UpdateShopCategoryStatusReq): Promise<CommonResponse> => {
  return request.put('/shop-category/status/update', data);
}; 