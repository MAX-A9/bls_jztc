"use strict";
const common_vendor = require("../../common/vendor.js");
const apis_content = require("../../apis/content.js");
const common_assets = require("../../common/assets.js");
const _sfc_main = {
  name: "CommunityNews",
  props: {
    categoryId: {
      type: [Number, String],
      default: 0
    }
  },
  data() {
    return {
      newsList: [],
      loading: false,
      hasMore: true,
      page: 1,
      pageSize: 10,
      regionId: 0,
      // 当前区域ID
      currentCategoryId: 0,
      // 当前选中的分类ID
      defaultImage: "/static/images/default-idle.png"
      // 默认图片路径，使用统一的默认图片
    };
  },
  created() {
    this.currentCategoryId = this.categoryId;
    const currentLocationId = common_vendor.index.getStorageSync("currentLocationId") || 0;
    if (currentLocationId) {
      this.regionId = currentLocationId;
      common_vendor.index.__f__("log", "at components/community-news/index.vue:84", "闲置列表组件创建时，从本地存储初始化区域ID:", this.regionId);
    }
    common_vendor.index.$on("idleCategoryChanged", this.handleCategoryChanged);
    common_vendor.index.$on("locationChanged", this.handleLocationChanged);
    this.loadData();
  },
  beforeDestroy() {
    common_vendor.index.$off("idleCategoryChanged", this.handleCategoryChanged);
    common_vendor.index.$off("locationChanged", this.handleLocationChanged);
  },
  methods: {
    // 处理分类变更
    handleCategoryChanged(data) {
      common_vendor.index.__f__("log", "at components/community-news/index.vue:104", "接收到分类变更事件:", data);
      if (data && data.categoryId !== void 0) {
        this.currentCategoryId = data.categoryId;
        this.resetAndLoad();
      }
    },
    // 处理位置变更
    handleLocationChanged(data) {
      common_vendor.index.__f__("log", "at components/community-news/index.vue:114", "位置变更事件触发，接收到的regionId:", data.regionId);
      const newRegionId = data.regionId || 0;
      if (!newRegionId) {
        common_vendor.index.__f__("error", "at components/community-news/index.vue:119", "接收到无效的regionId:", newRegionId);
        return;
      }
      if (this.regionId !== newRegionId) {
        common_vendor.index.__f__("log", "at components/community-news/index.vue:125", "地区ID已更新:", this.regionId, "->", newRegionId);
        this.regionId = newRegionId;
        this.resetAndLoad();
      } else {
        common_vendor.index.__f__("log", "at components/community-news/index.vue:130", "地区ID未变化，保持当前值:", this.regionId);
      }
    },
    // 重置列表并重新加载数据
    resetAndLoad() {
      this.page = 1;
      this.newsList = [];
      this.hasMore = true;
      this.loadData();
    },
    // 加载闲置物品数据
    async loadData() {
      if (this.loading)
        return;
      try {
        this.loading = true;
        const currentLocationId = common_vendor.index.getStorageSync("currentLocationId") || 0;
        if (!this.regionId && currentLocationId) {
          this.regionId = currentLocationId;
        }
        const params = {
          regionId: this.regionId,
          page: this.page,
          // 后端API要求页码最小值为1
          pageSize: this.pageSize,
          category: this.currentCategoryId
          // 使用当前分类ID
        };
        common_vendor.index.__f__("log", "at components/community-news/index.vue:165", "加载闲置物品列表，参数:", params);
        const result = await apis_content.getRegionIdleList(params);
        if (result.code === 0 && result.data) {
          const { list = [], total = 0 } = result.data;
          const formattedList = list.map((item) => {
            return {
              id: item.id,
              title: item.title || "",
              desc: item.summary || "暂无描述",
              location: item.tradePlace || "未知地点",
              price: item.price || 0,
              image: item.image || this.defaultImage,
              wantCount: item.likes || 0
            };
          });
          if (this.page === 1) {
            this.newsList = formattedList;
          } else {
            this.newsList = [...this.newsList, ...formattedList];
          }
          this.hasMore = this.newsList.length < total;
          common_vendor.index.__f__("log", "at components/community-news/index.vue:197", `加载了${list.length}条闲置物品数据，总共${total}条`);
          this.$parent && this.$parent.isRefreshing === true && this.notifyRefreshComplete();
        } else {
          throw new Error(result.message || "获取闲置物品列表失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at components/community-news/index.vue:205", "加载闲置物品列表失败:", error);
        if (this.page === 1 && this.newsList.length === 0) {
          common_vendor.index.__f__("log", "at components/community-news/index.vue:209", "使用默认数据");
          this.useDefaultData();
        }
        common_vendor.index.showToast({
          title: "加载失败，请重试",
          icon: "none",
          duration: 2e3
        });
      } finally {
        this.loading = false;
      }
    },
    // 通知父组件刷新完成
    notifyRefreshComplete() {
      if (this.$parent && this.$parent.isRefreshing === true) {
        this.$parent.isRefreshing = false;
        common_vendor.index.__f__("log", "at components/community-news/index.vue:228", "内容加载完成，关闭刷新状态");
      }
    },
    // 加载更多数据
    loadMore() {
      if (this.loading || !this.hasMore)
        return;
      this.page++;
      this.loadData();
    },
    // 图片加载错误处理
    handleImageError(index) {
      this.newsList[index].image = this.defaultImage;
    },
    // 使用默认数据（接口失败时的备选方案，实际应用中应删除）
    useDefaultData() {
      this.newsList = [];
      this.hasMore = false;
    },
    handleNewsClick(news) {
      this.recordBrowseHistory(news.id);
      common_vendor.index.navigateTo({
        url: `/pages/community/detail?id=${news.id}`,
        success: () => {
          common_vendor.index.__f__("log", "at components/community-news/index.vue:259", "跳转到闲置详情页:", news);
        },
        fail: (err) => {
          common_vendor.index.__f__("error", "at components/community-news/index.vue:262", "跳转失败:", err);
        }
      });
    },
    // 记录浏览历史
    async recordBrowseHistory(contentId) {
      if (!contentId) {
        common_vendor.index.__f__("error", "at components/community-news/index.vue:270", "记录浏览历史失败: 内容ID为空");
        return;
      }
      try {
        const result = await apis_content.addBrowseRecord(contentId, "idle");
        if (result.code !== 0) {
          if (true) {
            common_vendor.index.__f__("error", "at components/community-news/index.vue:280", "记录闲置浏览历史失败:", result.message);
          }
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at components/community-news/index.vue:284", "记录闲置浏览历史异常:", error);
      }
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: $data.loading && $data.newsList.length === 0
  }, $data.loading && $data.newsList.length === 0 ? {} : {}, {
    b: !$data.loading || $data.newsList.length > 0
  }, !$data.loading || $data.newsList.length > 0 ? {
    c: common_vendor.f($data.newsList, (item, index, i0) => {
      return {
        a: item.image || $data.defaultImage,
        b: common_vendor.o(($event) => $options.handleImageError(index), index),
        c: common_vendor.t(item.title),
        d: common_vendor.t(item.desc),
        e: common_vendor.t(item.location),
        f: common_vendor.t(item.price),
        g: common_vendor.t(item.wantCount),
        h: index,
        i: common_vendor.o(($event) => $options.handleNewsClick(item), index)
      };
    })
  } : {}, {
    d: $data.newsList.length > 0
  }, $data.newsList.length > 0 ? common_vendor.e({
    e: $data.hasMore && !$data.loading
  }, $data.hasMore && !$data.loading ? {} : {}, {
    f: $data.loading && $data.newsList.length > 0
  }, $data.loading && $data.newsList.length > 0 ? {} : {}, {
    g: !$data.hasMore && $data.newsList.length > 0
  }, !$data.hasMore && $data.newsList.length > 0 ? {} : {}) : {}, {
    h: !$data.loading && $data.newsList.length === 0
  }, !$data.loading && $data.newsList.length === 0 ? {
    i: common_assets._imports_0$1
  } : {});
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createComponent(Component);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/components/community-news/index.js.map
