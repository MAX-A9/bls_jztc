"use strict";
const common_vendor = require("../../../common/vendor.js");
const mixins_device = require("../../../mixins/device.js");
const apis_content = require("../../../apis/content.js");
const utils_constants = require("../../../utils/constants.js");
const utils_auth = require("../../../utils/auth.js");
const utils_pay = require("../../../utils/pay.js");
const utils_storage = require("../../../utils/storage.js");
const utils_deviceInfo = require("../../../utils/device-info.js");
const _sfc_main = {
  mixins: [mixins_device.deviceMixin],
  data() {
    return {
      statusBarHeight: 0,
      navBarHeight: 44,
      images: [],
      loading: false,
      formData: {
        type: "",
        title: "",
        description: "",
        publishPackageId: null,
        // 发布套餐ID
        topPackageId: null
        // 置顶套餐ID，null表示不使用置顶
      },
      infoTypes: [],
      publishPackages: [],
      // 发布套餐列表
      topPackages: []
      // 置顶套餐列表
    };
  },
  computed: {
    ...common_vendor.mapState("region", ["regionList"])
  },
  onLoad(options) {
    if (!utils_auth.isLoggedIn()) {
      common_vendor.index.showToast({
        title: "请先登录",
        icon: "none",
        duration: 2e3
      });
      setTimeout(() => {
        common_vendor.index.navigateTo({
          url: "/pages/login/index"
        });
      }, 1500);
      return;
    }
    this.statusBarHeight = utils_deviceInfo.deviceInfo.getStatusBarHeight();
    this.loadCategories();
    this.loadPackages();
    this.ensureRegionDataLoaded();
    if (options.draftId) {
      this.loadDraftData(options.draftId);
    }
  },
  methods: {
    handleBack() {
      common_vendor.index.navigateBack();
    },
    // 加载分类数据
    async loadCategories() {
      if (this.loading)
        return;
      this.loading = true;
      try {
        common_vendor.index.showLoading({
          title: "加载中...",
          mask: true
        });
        const res = await apis_content.getInfoCategories(1);
        common_vendor.index.__f__("log", "at pages/publish/info/index.vue:227", "分类数据:", res);
        if (res.code === 0 && res.data && res.data.list) {
          this.infoTypes = res.data.list;
          if (this.infoTypes.length > 0) {
            this.formData.type = this.infoTypes[0].id.toString();
          }
        } else {
          throw new Error("获取分类失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/publish/info/index.vue:240", "加载分类失败:", error);
        this.useDefaultCategories();
        common_vendor.index.showToast({
          title: "加载分类失败，请重试",
          icon: "none"
        });
      } finally {
        this.loading = false;
        common_vendor.index.hideLoading();
      }
    },
    // 使用默认分类数据
    useDefaultCategories() {
      this.infoTypes = [];
      this.formData.type = "";
      common_vendor.index.showToast({
        title: "分类加载失败，请刷新重试",
        icon: "none",
        duration: 2e3
      });
    },
    handlePreview() {
      common_vendor.index.showToast({
        title: "预览功能开发中",
        icon: "none"
      });
    },
    selectInfoType(type) {
      this.formData.type = type;
    },
    async chooseImage() {
      try {
        const res = await new Promise((resolve, reject) => {
          common_vendor.index.chooseImage({
            count: 9 - this.images.length,
            sizeType: ["compressed"],
            sourceType: ["album", "camera"],
            success: resolve,
            fail: reject
          });
        });
        if (res.tempFilePaths && res.tempFilePaths.length > 0) {
          common_vendor.index.showLoading({
            title: "上传中...",
            mask: true
          });
          try {
            const uploadPromises = res.tempFilePaths.map(async (filePath) => {
              try {
                return await this.uploadFile(filePath);
              } catch (error) {
                common_vendor.index.__f__("error", "at pages/publish/info/index.vue:302", "上传图片失败:", error);
                return null;
              }
            });
            const uploadResults = await Promise.all(uploadPromises);
            const successUrls = uploadResults.filter((url) => url !== null);
            this.images = [...this.images, ...successUrls];
          } catch (error) {
            common_vendor.index.__f__("error", "at pages/publish/info/index.vue:316", "处理上传图片时出错:", error);
            common_vendor.index.showToast({
              title: "上传图片失败",
              icon: "none"
            });
          } finally {
            common_vendor.index.hideLoading();
          }
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/publish/info/index.vue:326", "选择图片失败:", error);
        common_vendor.index.showToast({
          title: "选择图片失败",
          icon: "none"
        });
      }
    },
    // 上传文件到服务器
    async uploadFile(filePath) {
      return new Promise((resolve, reject) => {
        const token = common_vendor.index.getStorageSync("token");
        if (!token) {
          common_vendor.index.showToast({
            title: "请先登录",
            icon: "none"
          });
          return reject(new Error("请先登录"));
        }
        common_vendor.index.uploadFile({
          url: utils_constants.API_BASE_URL + "/wx/upload/image",
          filePath,
          name: "file",
          header: {
            "Authorization": "Bearer " + token
          },
          success: (uploadRes) => {
            try {
              const data = JSON.parse(uploadRes.data);
              common_vendor.index.__f__("log", "at pages/publish/info/index.vue:358", "上传响应:", data);
              if (data.code === 0 && data.data) {
                resolve(data.data.url);
              } else {
                reject(new Error(data.message || "上传失败"));
              }
            } catch (error) {
              common_vendor.index.__f__("error", "at pages/publish/info/index.vue:366", "解析上传响应失败:", error);
              reject(error);
            }
          },
          fail: (err) => {
            common_vendor.index.__f__("error", "at pages/publish/info/index.vue:371", "上传请求失败:", err);
            reject(err);
          }
        });
      });
    },
    deleteImage(index) {
      this.images.splice(index, 1);
    },
    selectPublishPackage(packageId) {
      this.formData.publishPackageId = packageId;
    },
    selectTopPackage(packageId) {
      this.formData.topPackageId = packageId;
    },
    saveDraft() {
      if (!this.formData.title && !this.formData.description) {
        common_vendor.index.showToast({
          title: "请填写标题或描述",
          icon: "none"
        });
        return;
      }
      try {
        const draftData = {
          ...this.formData,
          images: this.images
        };
        const draftId = utils_storage.saveInfoDraft(draftData);
        if (draftId) {
          common_vendor.index.showToast({
            title: "草稿保存成功",
            icon: "success"
          });
          setTimeout(() => {
            common_vendor.index.navigateBack();
          }, 1500);
        } else {
          throw new Error("保存失败");
        }
      } catch (e) {
        common_vendor.index.__f__("error", "at pages/publish/info/index.vue:419", "保存草稿失败", e);
        common_vendor.index.showToast({
          title: "保存草稿失败",
          icon: "none"
        });
      }
    },
    async publishItem() {
      if (!utils_auth.isLoggedIn()) {
        common_vendor.index.showToast({
          title: "请先登录",
          icon: "none"
        });
        setTimeout(() => {
          common_vendor.index.navigateTo({
            url: "/pages/login/index"
          });
        }, 1500);
        return;
      }
      if (!this.validateForm())
        return;
      let loadingShown = false;
      try {
        common_vendor.index.showLoading({
          title: "正在发布...",
          mask: true
        });
        loadingShown = true;
        let content = "";
        const textParagraphs = this.formData.description.split("\n").filter((p) => p.trim());
        textParagraphs.forEach((paragraph) => {
          content += `<p>${paragraph}</p>`;
        });
        if (this.images.length > 0) {
          this.images.forEach((imgUrl) => {
            content += `<p><img src="${imgUrl}"></p>`;
          });
        }
        let regionId = this.getRegionId();
        const submitData = {
          categoryId: this.formData.type ? parseInt(this.formData.type) : 0,
          title: this.formData.title,
          content,
          publishPackageId: this.formData.publishPackageId,
          topPackageId: this.formData.topPackageId,
          isTopRequest: this.formData.topPackageId !== null,
          topDays: 0,
          // 现在通过topPackageId来处理
          regionId,
          images: this.images
        };
        common_vendor.index.__f__("log", "at pages/publish/info/index.vue:488", "提交数据:", submitData);
        const res = await apis_content.createInfo(submitData);
        if (loadingShown) {
          common_vendor.index.hideLoading();
          loadingShown = false;
        }
        if (res.code === 0) {
          const orderNo = res.data.orderNo;
          if (orderNo) {
            const publishPackage = this.publishPackages.find((p) => p.id === this.formData.publishPackageId);
            const topPackage = this.formData.topPackageId ? this.topPackages.find((p) => p.id === this.formData.topPackageId) : null;
            const publishPrice = publishPackage ? publishPackage.price : 0;
            const topPrice = topPackage ? topPackage.price : 0;
            const totalPrice = publishPrice + topPrice;
            if (totalPrice > 0) {
              try {
                const body = topPackage ? `${publishPackage.title}+${topPackage.title}` : publishPackage.title;
                const payResult = await utils_pay.requestWxPay({
                  body,
                  orderNo,
                  totalFee: totalPrice
                });
                if (payResult.success) {
                  common_vendor.index.showToast({
                    title: "支付成功",
                    icon: "success"
                  });
                } else {
                  common_vendor.index.showToast({
                    title: payResult.message || "支付已取消",
                    icon: "none"
                  });
                }
              } catch (payError) {
                common_vendor.index.__f__("error", "at pages/publish/info/index.vue:541", "支付过程出错:", payError);
                common_vendor.index.showToast({
                  title: payError.message || "支付失败",
                  icon: "none"
                });
              }
            }
          }
          common_vendor.index.showToast({
            title: "发布成功",
            icon: "success"
          });
          setTimeout(() => {
            common_vendor.index.navigateBack();
          }, 1500);
        } else {
          throw new Error(res.message || "发布失败");
        }
      } catch (error) {
        if (loadingShown) {
          common_vendor.index.hideLoading();
          loadingShown = false;
        }
        common_vendor.index.__f__("error", "at pages/publish/info/index.vue:569", "发布信息失败:", error);
        common_vendor.index.showToast({
          title: error.message || "发布失败，请重试",
          icon: "none"
        });
      } finally {
        if (loadingShown) {
          common_vendor.index.hideLoading();
        }
      }
    },
    validateForm() {
      if (!this.formData.type) {
        common_vendor.index.showToast({
          title: "请选择信息类型",
          icon: "none"
        });
        return false;
      }
      if (!this.formData.title.trim()) {
        common_vendor.index.showToast({
          title: "请填写信息标题",
          icon: "none"
        });
        return false;
      }
      if (!this.formData.description.trim()) {
        common_vendor.index.showToast({
          title: "请填写信息描述",
          icon: "none"
        });
        return false;
      }
      if (this.publishPackages.length > 0 && !this.formData.publishPackageId) {
        common_vendor.index.showToast({
          title: "请选择发布套餐",
          icon: "none"
        });
        return false;
      }
      return true;
    },
    // 确保区域数据已加载
    ensureRegionDataLoaded() {
      if (this.$store && this.$store.state.region && (!this.$store.state.region.regionList || this.$store.state.region.regionList.length === 0)) {
        this.$store.dispatch("region/getRegionList").catch((error) => {
          common_vendor.index.__f__("error", "at pages/publish/info/index.vue:628", "加载地区数据失败:", error);
        });
      }
    },
    getRegionId() {
      let regionId = 0;
      const currentLocation = common_vendor.index.getStorageSync("currentLocation");
      if (this.regionList && this.regionList.length > 0) {
        const currentRegion = this.regionList.find((region) => region.name === currentLocation);
        if (currentRegion) {
          regionId = currentRegion.id;
          common_vendor.index.__f__("log", "at pages/publish/info/index.vue:642", "找到当前地区ID:", regionId, "地区名称:", currentLocation);
        } else {
          common_vendor.index.__f__("log", "at pages/publish/info/index.vue:644", "未找到当前地区ID，使用默认ID 0, 当前地区名称:", currentLocation);
        }
      }
      return regionId;
    },
    // 加载套餐数据
    async loadPackages() {
      try {
        common_vendor.index.showLoading({
          title: "加载中...",
          mask: true
        });
        const res = await apis_content.getPackageList();
        if (res.code === 0 && res.data) {
          let topPackages = res.data.topPackages || [];
          let publishPackages = res.data.publishPackages || [];
          const topEnabled = res.data.topEnabled !== void 0 ? res.data.topEnabled : true;
          const publishEnabled = res.data.publishEnabled !== void 0 ? res.data.publishEnabled : true;
          if (topPackages.length > 0) {
            topPackages = topPackages.sort((a, b) => {
              if (a.sortOrder === b.sortOrder || a.sortOrder === void 0 || b.sortOrder === void 0) {
                return a.id - b.id;
              }
              return a.sortOrder - b.sortOrder;
            });
          }
          if (publishPackages.length > 0) {
            publishPackages = publishPackages.sort((a, b) => {
              if (a.sortOrder === b.sortOrder || a.sortOrder === void 0 || b.sortOrder === void 0) {
                return a.id - b.id;
              }
              return a.sortOrder - b.sortOrder;
            });
          }
          this.topPackages = topEnabled ? topPackages : [];
          this.publishPackages = publishEnabled ? publishPackages : [];
          if (this.publishPackages.length > 0) {
            this.formData.publishPackageId = this.publishPackages[0].id;
          } else {
            this.formData.publishPackageId = null;
          }
          if (!topEnabled) {
            this.formData.topPackageId = null;
          }
        } else {
          throw new Error("获取套餐数据失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/publish/info/index.vue:711", "加载套餐数据失败:", error);
        common_vendor.index.showToast({
          title: "加载套餐失败，请重试",
          icon: "none"
        });
      } finally {
        common_vendor.index.hideLoading();
      }
    },
    // 加载草稿数据
    loadDraftData(draftId) {
      try {
        const draftData = utils_storage.getInfoDraft(Number(draftId));
        if (!draftData) {
          common_vendor.index.showToast({
            title: "草稿不存在",
            icon: "none"
          });
          return;
        }
        if (draftData.type) {
          this.formData.type = draftData.type;
        }
        if (draftData.title) {
          this.formData.title = draftData.title;
        }
        if (draftData.description) {
          this.formData.description = draftData.description;
        }
        if (draftData.publishPackageId) {
          this.formData.publishPackageId = draftData.publishPackageId;
        }
        if (draftData.topPackageId) {
          this.formData.topPackageId = draftData.topPackageId;
        }
        if (draftData.images && draftData.images.length > 0) {
          this.images = draftData.images;
        }
      } catch (e) {
        common_vendor.index.__f__("error", "at pages/publish/info/index.vue:759", "加载草稿失败", e);
        common_vendor.index.showToast({
          title: "加载草稿失败",
          icon: "none"
        });
      }
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  _easycom_uni_icons2();
}
const _easycom_uni_icons = () => "../../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
if (!Math) {
  _easycom_uni_icons();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: common_vendor.p({
      type: "left",
      size: "20",
      color: "#333333"
    }),
    b: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    c: $data.navBarHeight + "px",
    d: $data.statusBarHeight + "px",
    e: `calc(${$data.statusBarHeight}px + ${$data.navBarHeight}px)`,
    f: common_vendor.f($data.infoTypes, (type, index, i0) => {
      return common_vendor.e({
        a: type.icon
      }, type.icon ? {
        b: type.icon
      } : {
        c: "d9e4d362-1-" + i0,
        d: common_vendor.p({
          type: "info",
          size: "24",
          color: $data.formData.type === type.id.toString() ? "#ffffff" : "#666666"
        })
      }, {
        e: common_vendor.t(type.name),
        f: index,
        g: $data.formData.type === type.id.toString() ? 1 : "",
        h: common_vendor.o(($event) => $options.selectInfoType(type.id.toString()), index)
      });
    }),
    g: $data.formData.title,
    h: common_vendor.o(($event) => $data.formData.title = $event.detail.value),
    i: common_vendor.t($data.formData.title.length),
    j: $data.formData.description,
    k: common_vendor.o(($event) => $data.formData.description = $event.detail.value),
    l: common_vendor.t($data.formData.description.length),
    m: common_vendor.f($data.images, (item, index, i0) => {
      return {
        a: item,
        b: "d9e4d362-2-" + i0,
        c: common_vendor.o(($event) => $options.deleteImage(index), index),
        d: index
      };
    }),
    n: common_vendor.p({
      type: "closeempty",
      size: "14",
      color: "#ffffff"
    }),
    o: $data.images.length < 9
  }, $data.images.length < 9 ? {
    p: common_vendor.p({
      type: "camera",
      size: "24",
      color: "#999999"
    }),
    q: common_vendor.t($data.images.length),
    r: common_vendor.o((...args) => $options.chooseImage && $options.chooseImage(...args))
  } : {}, {
    s: $data.publishPackages.length > 0
  }, $data.publishPackages.length > 0 ? {
    t: common_vendor.f($data.publishPackages, (item, k0, i0) => {
      return common_vendor.e({
        a: common_vendor.t(item.title),
        b: common_vendor.t(item.price === 0 ? "免费" : item.price + "元"),
        c: $data.formData.publishPackageId === item.id
      }, $data.formData.publishPackageId === item.id ? {} : {}, {
        d: $data.formData.publishPackageId === item.id ? 1 : "",
        e: "publish-" + item.id,
        f: common_vendor.o(($event) => $options.selectPublishPackage(item.id), "publish-" + item.id)
      });
    })
  } : {}, {
    v: $data.topPackages.length > 0
  }, $data.topPackages.length > 0 ? {
    w: common_vendor.f($data.topPackages, (item, k0, i0) => {
      return common_vendor.e({
        a: common_vendor.t(item.title),
        b: common_vendor.t(item.price === 0 ? "免费" : item.price + "元"),
        c: $data.formData.topPackageId === item.id
      }, $data.formData.topPackageId === item.id ? {} : {}, {
        d: $data.formData.topPackageId === item.id ? 1 : "",
        e: "top-" + item.id,
        f: common_vendor.o(($event) => $options.selectTopPackage(item.id), "top-" + item.id)
      });
    })
  } : {}, {
    x: common_vendor.o((...args) => $options.saveDraft && $options.saveDraft(...args)),
    y: common_vendor.o((...args) => $options.publishItem && $options.publishItem(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/publish/info/index.js.map
