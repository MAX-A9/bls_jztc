"use strict";
const common_vendor = require("../../../common/vendor.js");
const utils_deviceInfo = require("../../../utils/device-info.js");
const common_assets = require("../../../common/assets.js");
const _sfc_main = {
  data() {
    return {
      statusBarHeight: 0,
      navBarHeight: 44,
      formData: {
        name: "",
        phone: "",
        idCard: "",
        region: ["", "", ""],
        regionText: "",
        skills: [],
        experienceIndex: 0,
        experience: "",
        introduction: "",
        idCardFront: "",
        idCardBack: "",
        qualification: "",
        agreement: false
      },
      skillList: [
        { label: "家电维修", value: "appliance" },
        { label: "水电维修", value: "plumbing" },
        { label: "空调安装", value: "ac" },
        { label: "家具安装", value: "furniture" },
        { label: "管道疏通", value: "pipe" },
        { label: "电器安装", value: "electronics" },
        { label: "防水补漏", value: "waterproof" },
        { label: "开锁换锁", value: "lock" }
      ],
      experienceOptions: ["1年以下", "1-3年", "3-5年", "5-10年", "10年以上"]
    };
  },
  onLoad() {
    this.statusBarHeight = utils_deviceInfo.deviceInfo.getStatusBarHeight();
  },
  methods: {
    handleBack() {
      common_vendor.index.navigateBack();
    },
    handleRegionChange(e) {
      this.formData.region = e.detail.value;
      this.formData.regionText = e.detail.value.join(" ");
    },
    toggleSkill(value) {
      const index = this.formData.skills.indexOf(value);
      if (index === -1) {
        this.formData.skills.push(value);
      } else {
        this.formData.skills.splice(index, 1);
      }
    },
    handleExperienceChange(e) {
      this.formData.experienceIndex = e.detail.value;
      this.formData.experience = this.experienceOptions[e.detail.value];
    },
    uploadImage(field) {
      common_vendor.index.chooseImage({
        count: 1,
        sizeType: ["compressed"],
        sourceType: ["album", "camera"],
        success: (res) => {
          this.formData[field] = res.tempFilePaths[0];
        }
      });
    },
    toggleAgreement() {
      this.formData.agreement = !this.formData.agreement;
    },
    showAgreement() {
      common_vendor.index.showModal({
        title: "师傅入驻协议",
        content: "这里是师傅入驻协议的详细内容...",
        showCancel: false
      });
    },
    handleSubmit() {
      if (!this.formData.name) {
        return common_vendor.index.showToast({ title: "请输入姓名", icon: "none" });
      }
      if (!this.formData.phone || !/^1\d{10}$/.test(this.formData.phone)) {
        return common_vendor.index.showToast({ title: "请输入正确的手机号", icon: "none" });
      }
      if (!this.formData.idCard) {
        return common_vendor.index.showToast({ title: "请输入身份证号", icon: "none" });
      }
      if (!this.formData.regionText) {
        return common_vendor.index.showToast({ title: "请选择所在城市", icon: "none" });
      }
      if (this.formData.skills.length === 0) {
        return common_vendor.index.showToast({ title: "请选择至少一项专业技能", icon: "none" });
      }
      if (!this.formData.experience) {
        return common_vendor.index.showToast({ title: "请选择工作年限", icon: "none" });
      }
      if (!this.formData.idCardFront) {
        return common_vendor.index.showToast({ title: "请上传身份证正面照片", icon: "none" });
      }
      if (!this.formData.idCardBack) {
        return common_vendor.index.showToast({ title: "请上传身份证反面照片", icon: "none" });
      }
      common_vendor.index.showLoading({ title: "提交中..." });
      setTimeout(() => {
        common_vendor.index.hideLoading();
        common_vendor.index.showModal({
          title: "提交成功",
          content: "您的入驻申请已提交，我们将在1-3个工作日内审核，请留意短信通知",
          showCancel: false,
          success: () => {
            common_vendor.index.navigateBack();
          }
        });
      }, 1500);
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
  return {
    a: common_vendor.p({
      type: "left",
      size: "20",
      color: "#333333"
    }),
    b: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    c: $data.navBarHeight + "px",
    d: $data.statusBarHeight + "px",
    e: common_assets._imports_0$6,
    f: $data.formData.name,
    g: common_vendor.o(($event) => $data.formData.name = $event.detail.value),
    h: $data.formData.phone,
    i: common_vendor.o(($event) => $data.formData.phone = $event.detail.value),
    j: $data.formData.idCard,
    k: common_vendor.o(($event) => $data.formData.idCard = $event.detail.value),
    l: common_vendor.t($data.formData.regionText || "请选择所在城市"),
    m: common_vendor.o((...args) => $options.handleRegionChange && $options.handleRegionChange(...args)),
    n: $data.formData.region,
    o: common_vendor.p({
      type: "right",
      size: "16",
      color: "#999999"
    }),
    p: common_vendor.f($data.skillList, (item, index, i0) => {
      return {
        a: common_vendor.t(item.label),
        b: index,
        c: $data.formData.skills.includes(item.value) ? 1 : "",
        d: common_vendor.o(($event) => $options.toggleSkill(item.value), index)
      };
    }),
    q: common_vendor.t($data.formData.experience || "请选择工作年限"),
    r: $data.experienceOptions,
    s: common_vendor.o((...args) => $options.handleExperienceChange && $options.handleExperienceChange(...args)),
    t: $data.formData.experienceIndex,
    v: common_vendor.p({
      type: "right",
      size: "16",
      color: "#999999"
    }),
    w: $data.formData.introduction,
    x: common_vendor.o(($event) => $data.formData.introduction = $event.detail.value),
    y: common_vendor.t($data.formData.introduction.length),
    z: $data.formData.idCardFront || "/static/upload-id-front.png",
    A: common_vendor.o(($event) => $options.uploadImage("idCardFront")),
    B: $data.formData.idCardBack || "/static/upload-id-back.png",
    C: common_vendor.o(($event) => $options.uploadImage("idCardBack")),
    D: $data.formData.qualification || "/static/upload-qualification.png",
    E: common_vendor.o(($event) => $options.uploadImage("qualification")),
    F: common_vendor.p({
      type: $data.formData.agreement ? "checkbox-filled" : "circle",
      size: "20",
      color: $data.formData.agreement ? "#fc3e2b" : "#999999"
    }),
    G: common_vendor.o((...args) => $options.toggleAgreement && $options.toggleAgreement(...args)),
    H: common_vendor.o((...args) => $options.showAgreement && $options.showAgreement(...args)),
    I: !$data.formData.agreement,
    J: common_vendor.o((...args) => $options.handleSubmit && $options.handleSubmit(...args)),
    K: `calc(${$data.statusBarHeight}px + ${$data.navBarHeight}px)`
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/master/join/index.js.map
