"use strict";
const common_vendor = require("../../../common/vendor.js");
const mixins_device = require("../../../mixins/device.js");
const apis_content = require("../../../apis/content.js");
const utils_constants = require("../../../utils/constants.js");
const utils_auth = require("../../../utils/auth.js");
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
        // 默认为空，等待接口返回后设置
        title: "",
        description: "",
        price: "",
        originalPrice: "",
        location: "同城交易",
        latitude: void 0,
        // 位置纬度
        longitude: void 0,
        // 位置经度
        contact: {
          name: "",
          phone: "",
          wechat: ""
        }
      },
      idleTypes: []
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
    this.ensureRegionDataLoaded();
    this.loadIdleCategories();
    if (options.draftId) {
      this.loadDraftData(options.draftId);
    }
  },
  methods: {
    handleBack() {
      common_vendor.index.navigateBack();
    },
    // 加载闲置类型分类数据
    async loadIdleCategories() {
      if (this.loading)
        return;
      this.loading = true;
      try {
        common_vendor.index.showLoading({
          title: "加载中...",
          mask: true
        });
        const res = await apis_content.getInfoCategories(2);
        common_vendor.index.__f__("log", "at pages/publish/idle/index.vue:241", "闲置分类数据:", res);
        if (res.code === 0 && res.data && res.data.list) {
          this.idleTypes = res.data.list.map((item) => {
            return {
              label: item.name,
              value: item.id.toString(),
              iconUrl: item.icon,
              // 使用API返回的图标URL
              id: item.id
            };
          });
          if (this.idleTypes.length > 0) {
            this.formData.type = this.idleTypes[0].value;
          }
        } else {
          throw new Error("获取分类失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/publish/idle/index.vue:262", "加载闲置分类失败:", error);
        this.useDefaultCategories();
        common_vendor.index.showToast({
          title: "加载分类失败，使用默认分类",
          icon: "none"
        });
      } finally {
        this.loading = false;
        common_vendor.index.hideLoading();
      }
    },
    // 使用默认闲置分类数据（网络错误时的备选方案）
    useDefaultCategories() {
      this.idleTypes = [];
      this.formData.type = "";
      common_vendor.index.showToast({
        title: "分类加载失败，请刷新重试",
        icon: "none",
        duration: 2e3
      });
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
                common_vendor.index.__f__("error", "at pages/publish/idle/index.vue:315", "上传图片失败:", error);
                return null;
              }
            });
            const uploadResults = await Promise.all(uploadPromises);
            const successUrls = uploadResults.filter((url) => url !== null);
            this.images = [...this.images, ...successUrls];
          } catch (error) {
            common_vendor.index.__f__("error", "at pages/publish/idle/index.vue:329", "处理上传图片时出错:", error);
            common_vendor.index.showToast({
              title: "上传图片失败",
              icon: "none"
            });
          } finally {
            common_vendor.index.hideLoading();
          }
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/publish/idle/index.vue:339", "选择图片失败:", error);
        common_vendor.index.showToast({
          title: "选择图片失败",
          icon: "none"
        });
        if (this.loading) {
          common_vendor.index.hideLoading();
        }
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
              common_vendor.index.__f__("log", "at pages/publish/idle/index.vue:375", "上传响应:", data);
              if (data.code === 0 && data.data) {
                resolve(data.data.url);
              } else {
                reject(new Error(data.message || "上传失败"));
              }
            } catch (error) {
              common_vendor.index.__f__("error", "at pages/publish/idle/index.vue:383", "解析上传响应失败:", error);
              reject(error);
            }
          },
          fail: (err) => {
            common_vendor.index.__f__("error", "at pages/publish/idle/index.vue:388", "上传请求失败:", err);
            reject(err);
          }
        });
      });
    },
    deleteImage(index) {
      this.images.splice(index, 1);
    },
    showLocationPicker() {
      common_vendor.index.choosePoi({
        success: (res) => {
          common_vendor.index.__f__("log", "at pages/publish/idle/index.vue:401", "选择位置成功:", res);
          if (res.type === 1) {
            let cityName = res.city || "";
            if (cityName.indexOf("省") > 0) {
              cityName = cityName.split("省")[1];
            } else if (cityName.indexOf("自治区") > 0) {
              cityName = cityName.split("自治区")[1];
            }
            this.formData.location = cityName;
            this.formData.latitude = void 0;
            this.formData.longitude = void 0;
          } else if (res.type === 2) {
            let address = res.address || "";
            let name = res.name || "";
            let formattedAddress = "";
            if (address.indexOf("省") > 0) {
              address = address.split("省")[1];
            } else if (address.indexOf("自治区") > 0) {
              address = address.split("自治区")[1];
            }
            if (address.indexOf("市") > 0) {
              formattedAddress = address + " " + name;
            } else {
              formattedAddress = address + " " + name;
            }
            this.formData.location = formattedAddress;
            this.formData.latitude = res.latitude;
            this.formData.longitude = res.longitude;
          }
        },
        fail: (err) => {
          common_vendor.index.__f__("error", "at pages/publish/idle/index.vue:450", "选择位置失败:", err);
          if (err.errMsg && err.errMsg.indexOf("cancel") !== -1) {
            return;
          }
          common_vendor.index.showToast({
            title: "位置选择失败，请重试",
            icon: "none"
          });
        }
      });
    },
    selectIdleType(type) {
      this.formData.type = type;
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
        const draftId = utils_storage.saveIdleDraft(draftData);
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
        common_vendor.index.__f__("error", "at pages/publish/idle/index.vue:501", "保存草稿失败", e);
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
        return;
      }
      if (!this.validateForm(true))
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
        const submitData = {
          categoryId: this.formData.type ? parseInt(this.formData.type) : 0,
          title: this.formData.title,
          price: parseFloat(this.formData.price) || 0,
          originalPrice: parseFloat(this.formData.originalPrice) || 0,
          content,
          tradeMethod: "线下交易",
          tradePlace: this.formData.location || "同城交易",
          images: this.images,
          // 添加地区ID
          regionId: this.getRegionId()
        };
        common_vendor.index.__f__("log", "at pages/publish/idle/index.vue:562", "提交数据:", submitData);
        const res = await apis_content.createIdle(submitData);
        if (loadingShown) {
          common_vendor.index.hideLoading();
          loadingShown = false;
        }
        if (res.code === 0) {
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
        common_vendor.index.__f__("error", "at pages/publish/idle/index.vue:593", "发布闲置失败:", error);
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
    validatePriceInput(field) {
      const value = this.formData[field];
      if (!value)
        return;
      const regex = /^\d+(\.\d{0,2})?$/;
      if (!regex.test(value)) {
        const hasForbiddenChars = /[^\d.]/.test(value);
        const cleaned = value.replace(/[^\d.]/g, "").replace(/\.{2,}/g, ".").replace(/^(\d+\.\d{0,2}).*$/, "$1");
        this.formData[field] = cleaned;
        if (hasForbiddenChars) {
          common_vendor.index.showToast({
            title: "请只输入数字和小数点",
            icon: "none",
            duration: 1500
          });
        }
      }
    },
    validateForm(isPublish) {
      if (!this.formData.type) {
        common_vendor.index.showToast({
          title: "请选择闲置类型",
          icon: "none"
        });
        return false;
      }
      if (this.images.length === 0) {
        common_vendor.index.showToast({
          title: "请至少上传一张图片",
          icon: "none"
        });
        return false;
      }
      if (!this.formData.title.trim()) {
        common_vendor.index.showToast({
          title: "请填写商品标题",
          icon: "none"
        });
        return false;
      }
      if (isPublish) {
        if (!this.formData.price) {
          common_vendor.index.showToast({
            title: "请填写商品价格",
            icon: "none"
          });
          return false;
        }
        const priceRegex = /^\d+(\.\d{1,2})?$/;
        if (!priceRegex.test(this.formData.price)) {
          common_vendor.index.showToast({
            title: "请输入正确的价格格式",
            icon: "none"
          });
          return false;
        }
        if (this.formData.originalPrice !== void 0 && this.formData.originalPrice !== "") {
          if (!priceRegex.test(this.formData.originalPrice)) {
            common_vendor.index.showToast({
              title: "请输入正确的原价格式",
              icon: "none"
            });
            return false;
          }
          const hasForbiddenChars = /[^\d.]/.test(this.formData.originalPrice);
          if (hasForbiddenChars) {
            common_vendor.index.showToast({
              title: "原价只能输入数字",
              icon: "none"
            });
            return false;
          }
        }
        if (!this.formData.location) {
          common_vendor.index.showToast({
            title: "请选择交易地点",
            icon: "none"
          });
          return false;
        }
      }
      return true;
    },
    // 确保区域数据已加载
    ensureRegionDataLoaded() {
      if (this.$store && this.$store.state.region && (!this.$store.state.region.regionList || this.$store.state.region.regionList.length === 0)) {
        this.$store.dispatch("region/getRegionList").catch((error) => {
          common_vendor.index.__f__("error", "at pages/publish/idle/index.vue:729", "加载地区数据失败:", error);
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
          common_vendor.index.__f__("log", "at pages/publish/idle/index.vue:743", "找到当前地区ID:", regionId, "地区名称:", currentLocation);
        } else {
          common_vendor.index.__f__("log", "at pages/publish/idle/index.vue:745", "未找到当前地区ID，使用默认ID 0, 当前地区名称:", currentLocation);
        }
      }
      return regionId;
    },
    // 加载草稿数据
    loadDraftData(draftId) {
      try {
        const draftData = utils_storage.getIdleDraft(Number(draftId));
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
        if (draftData.price) {
          this.formData.price = draftData.price;
        }
        if (draftData.originalPrice) {
          this.formData.originalPrice = draftData.originalPrice;
        }
        if (draftData.location) {
          this.formData.location = draftData.location;
        }
        if (draftData.latitude) {
          this.formData.latitude = draftData.latitude;
        }
        if (draftData.longitude) {
          this.formData.longitude = draftData.longitude;
        }
        if (draftData.contact) {
          this.formData.contact = draftData.contact;
        }
        if (draftData.images && draftData.images.length > 0) {
          this.images = draftData.images;
        }
      } catch (e) {
        common_vendor.index.__f__("error", "at pages/publish/idle/index.vue:806", "加载草稿失败", e);
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
    f: common_vendor.f($data.idleTypes, (type, index, i0) => {
      return common_vendor.e({
        a: type.iconUrl
      }, type.iconUrl ? {
        b: type.iconUrl
      } : {
        c: "c1ab8c16-1-" + i0,
        d: common_vendor.p({
          type: "info",
          size: "24",
          color: $data.formData.type === type.value ? "#ffffff" : "#666666"
        })
      }, {
        e: common_vendor.t(type.label),
        f: index,
        g: $data.formData.type === type.value ? 1 : "",
        h: common_vendor.o(($event) => $options.selectIdleType(type.value), index)
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
        b: "c1ab8c16-2-" + i0,
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
    s: common_vendor.o([($event) => $data.formData.price = $event.detail.value, ($event) => $options.validatePriceInput("price")]),
    t: $data.formData.price,
    v: common_vendor.o([($event) => $data.formData.originalPrice = $event.detail.value, ($event) => $options.validatePriceInput("originalPrice")]),
    w: $data.formData.originalPrice,
    x: common_vendor.t($data.formData.location || "请选择交易地点"),
    y: common_vendor.p({
      type: "right",
      size: "16",
      color: "#999999"
    }),
    z: common_vendor.o((...args) => $options.showLocationPicker && $options.showLocationPicker(...args)),
    A: common_vendor.o((...args) => $options.saveDraft && $options.saveDraft(...args)),
    B: common_vendor.o((...args) => $options.publishItem && $options.publishItem(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/publish/idle/index.js.map
