"use strict";
const common_vendor = require("../../common/vendor.js");
const utils_request = require("../../utils/request.js");
const utils_pinyin = require("../../utils/pinyin.js");
const state = {
  regionList: [],
  // 地区列表
  loading: false,
  // 加载状态
  error: null
  // 错误信息
};
const getters = {
  // 按字母分组的地区列表
  groupedRegionList: (state2) => {
    const result = {};
    if (!state2.regionList.length) {
      return result;
    }
    state2.regionList.forEach((region2) => {
      const firstLetter = utils_pinyin.getStringFirstLetter(region2.name);
      if (!result[firstLetter]) {
        result[firstLetter] = [];
      }
      result[firstLetter].push(region2);
    });
    return result;
  },
  // 获取热门地区
  hotRegions: (state2) => {
    const hotRegions = state2.regionList.filter((region2) => region2.hot === 1);
    return hotRegions.length > 0 ? hotRegions : state2.regionList.slice(0, 9);
  }
};
const mutations = {
  // 设置地区列表
  SET_REGION_LIST(state2, regionList) {
    state2.regionList = regionList;
  },
  // 设置加载状态
  SET_LOADING(state2, status) {
    state2.loading = status;
  },
  // 设置错误信息
  SET_ERROR(state2, error) {
    state2.error = error;
  }
};
const actions = {
  // 获取地区列表
  async getRegionList({ commit }, params = { status: 0 }) {
    commit("SET_LOADING", true);
    commit("SET_ERROR", null);
    try {
      const response = await utils_request.get("/wx/client/region/list", params, true);
      if (response && response.code === 0) {
        const regionList = response.data.list.map((item) => ({
          id: item.id,
          name: item.name,
          location: item.location,
          level: item.level,
          hot: item.hot || 0,
          // 添加热门标记，默认为0
          status: item.status,
          createdAt: item.createdAt,
          updatedAt: item.updatedAt
        }));
        commit("SET_REGION_LIST", regionList);
        try {
          common_vendor.index.setStorageSync("regionList", JSON.stringify(regionList));
          common_vendor.index.__f__("log", "at store/modules/region.js:89", "区域列表数据已保存到本地存储，共", regionList.length, "条记录");
        } catch (storageError) {
          common_vendor.index.__f__("error", "at store/modules/region.js:91", "保存区域列表到本地存储失败:", storageError);
        }
        return regionList;
      } else {
        commit("SET_ERROR", (response == null ? void 0 : response.message) || "获取地区列表失败");
        return [];
      }
    } catch (error) {
      common_vendor.index.__f__("error", "at store/modules/region.js:100", "获取地区列表出错:", error);
      commit("SET_ERROR", error.message || "获取地区列表失败");
      return [];
    } finally {
      commit("SET_LOADING", false);
    }
  }
};
const region = {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
};
exports.region = region;
//# sourceMappingURL=../../../.sourcemap/mp-weixin/store/modules/region.js.map
