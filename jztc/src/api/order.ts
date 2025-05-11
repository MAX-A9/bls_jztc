import request from '../utils/request';

// 订单列表查询参数
export interface OrderListParams {
  page?: number;
  pageSize?: number;
  orderNo?: string;
  clientName?: string;
  status?: string;
  startTime?: string;
  endTime?: string;
  product?: string;
}

// 订单列表项
export interface OrderListItem {
  id: number;
  orderNo: string;
  clientName: string;
  contentId: number;
  productName: string;
  amount: number;
  status: number;
  statusText: string;
  paymentMethod: string;
  createdAt: string;
  payTime: string;
}

// 订单列表响应
export interface OrderListResponse {
  code?: number;
  message?: string;
  data?: {
    list: OrderListItem[];
    total: number;
    page: number;
  };
  list?: OrderListItem[];
  total?: number;
  page?: number;
}

// 订单详情响应
export interface OrderDetailResponse {
  code?: number;
  message?: string;
  data?: {
    id: number;
    orderNo: string;
    userId: number;
    clientName: string;
    contentId?: number;
    productName: string;
    amount: number;
    status: number;
    statusText: string;
    paymentMethod: string;
    transactionId: string;
    remark: string;
    createdAt: string;
    payTime: string;
  };
  id?: number;
  orderNo?: string;
  userId?: number;
  clientName?: string;
  contentId?: number;
  productName?: string;
  amount?: number;
  status?: number;
  statusText?: string;
  paymentMethod?: string;
  transactionId?: string;
  remark?: string;
  createdAt?: string;
  payTime?: string;
}

// 创建订单参数
export interface OrderCreateParams {
  userId: number;
  productName: string;
  amount: number;
  remark?: string;
}

// 创建订单响应
export interface OrderCreateResponse {
  code?: number;
  message?: string;
  data?: {
    orderNo: string;
  };
  orderNo?: string;
}

// 更新订单状态参数
export interface OrderUpdateStatusParams {
  orderNo: string;
  status: number;
  paymentMethod?: string;
  transactionId?: string;
  remark?: string;
}

/**
 * 获取订单列表
 * @param params 查询参数
 * @returns 订单列表
 */
export const getOrderList = (params: OrderListParams): Promise<OrderListResponse> => {
  return request.get('/order/list', { params });
};

/**
 * 获取订单详情
 * @param orderNo 订单号
 * @returns 订单详情
 */
export const getOrderDetail = (orderNo: string): Promise<OrderDetailResponse> => {
  return request.get('/order/detail', { params: { orderNo } });
};

/**
 * 创建订单
 * @param params 订单信息
 * @returns 创建结果
 */
export const createOrder = (params: OrderCreateParams): Promise<OrderCreateResponse> => {
  return request.post('/order/create', params);
};

/**
 * 取消订单
 * @param orderNo 订单号
 * @returns 取消结果
 */
export const cancelOrder = (orderNo: string): Promise<void> => {
  return request.post('/order/cancel', { orderNo });
};

/**
 * 删除订单
 * @param orderNo 订单号
 * @returns 删除结果
 */
export const deleteOrder = (orderNo: string): Promise<void> => {
  return request.post('/order/delete', { orderNo });
};

/**
 * 更新订单状态
 * @param params 更新参数
 * @returns 更新结果
 */
export const updateOrderStatus = (params: OrderUpdateStatusParams): Promise<void> => {
  return request.post('/order/update-status', params);
}; 