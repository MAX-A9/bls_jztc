import { get } from '@/utils/request.js';
import { getStringFirstLetter } from '@/utils/pinyin.js';

// 初始状态
const state = {
  regionList: [], // 地区列表
  loading: false, // 加载状态
  error: null, // 错误信息
};

// getters
const getters = {
  // 按字母分组的地区列表
  groupedRegionList: (state) => {
    const result = {};
    
    if (!state.regionList.length) {
      return result;
    }
    
    // 按地区名称首字母进行分组（使用拼音首字母）
    state.regionList.forEach(region => {
      const firstLetter = getStringFirstLetter(region.name);
      if (!result[firstLetter]) {
        result[firstLetter] = [];
      }
      result[firstLetter].push(region);
    });
    
    return result;
  },
  
  // 获取热门地区
  hotRegions: (state) => {
    // 如果API返回的数据中包含hot标记，则使用标记的方式获取热门城市
    const hotRegions = state.regionList.filter(region => region.hot === 1);
    
    // 如果没有标记为热门的城市，则返回前9个作为热门城市
    return hotRegions.length > 0 ? hotRegions : state.regionList.slice(0, 9);
  }
};

// mutations
const mutations = {
  // 设置地区列表
  SET_REGION_LIST(state, regionList) {
    state.regionList = regionList;
  },
  
  // 设置加载状态
  SET_LOADING(state, status) {
    state.loading = status;
  },
  
  // 设置错误信息
  SET_ERROR(state, error) {
    state.error = error;
  }
};

// actions
const actions = {
  // 获取地区列表
  async getRegionList({ commit }, params = { status: 0 }) {
    commit('SET_LOADING', true);
    commit('SET_ERROR', null);
    
    try {
      const response = await get('/wx/client/region/list', params, true);
      
      if (response && response.code === 0) {
        // 处理接口返回的数据，注意数据在 response.data.list 中
        const regionList = response.data.list.map(item => ({
          id: item.id,
          name: item.name,
          location: item.location,
          level: item.level,
          hot: item.hot || 0, // 添加热门标记，默认为0
          status: item.status,
          createdAt: item.createdAt,
          updatedAt: item.updatedAt
        }));
        
        commit('SET_REGION_LIST', regionList);
        
        // 将数据保存到本地存储，方便下次使用
        try {
          uni.setStorageSync('regionList', JSON.stringify(regionList));
          console.log('区域列表数据已保存到本地存储，共', regionList.length, '条记录');
        } catch (storageError) {
          console.error('保存区域列表到本地存储失败:', storageError);
        }
        
        return regionList;
      } else {
        commit('SET_ERROR', response?.message || '获取地区列表失败');
        return [];
      }
    } catch (error) {
      console.error('获取地区列表出错:', error);
      commit('SET_ERROR', error.message || '获取地区列表失败');
      return [];
    } finally {
      commit('SET_LOADING', false);
    }
  }
};

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}; 