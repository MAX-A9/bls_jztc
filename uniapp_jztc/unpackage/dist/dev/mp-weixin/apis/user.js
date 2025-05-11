"use strict";
const common_vendor = require("../common/vendor.js");
const utils_request = require("../utils/request.js");
const utils_constants = require("../utils/constants.js");
function wxappLogin(data) {
  return utils_request.post(utils_constants.API_PATHS.WXAPP_LOGIN, data).catch((err) => {
    throw err;
  });
}
function getClientInfo() {
  return utils_request.get(utils_constants.API_PATHS.CLIENT_INFO).catch((err) => {
    throw err;
  });
}
function updateClientProfile(data) {
  return utils_request.put(utils_constants.API_PATHS.CLIENT_UPDATE_PROFILE, data).catch((err) => {
    throw err;
  });
}
function getWxLoginCode() {
  return new Promise((resolve, reject) => {
    common_vendor.index.login({
      provider: "weixin",
      success: (res) => {
        if (res.code) {
          resolve(res.code);
        } else {
          const error = { message: "获取微信登录凭证失败", res };
          reject(error);
        }
      },
      fail: (err) => {
        reject(err || { message: "微信登录失败" });
      }
    });
  });
}
function getWxUserInfo() {
  return new Promise((resolve, reject) => {
    common_vendor.index.getUserProfile({
      desc: "用于完善用户资料",
      success: (res) => {
        if (res.userInfo) {
          resolve(res.userInfo);
        } else {
          reject({ message: "获取用户信息失败" });
        }
      },
      fail: (err) => {
        reject(err || { message: "用户拒绝授权" });
      }
    });
  });
}
function getImageBase64(imageUrl) {
  return new Promise((resolve, reject) => {
    common_vendor.index.getImageInfo({
      src: imageUrl,
      success: (imgInfo) => {
        const canvasSize = 200;
        const ctx = common_vendor.index.createCanvasContext("avatarCanvas");
        ctx.clearRect(0, 0, canvasSize, canvasSize);
        ctx.drawImage(imgInfo.path, 0, 0, canvasSize, canvasSize);
        let drawTimeout = setTimeout(() => {
          reject(new Error("Canvas绘制超时"));
        }, 2e3);
        ctx.draw(false, () => {
          clearTimeout(drawTimeout);
          common_vendor.index.canvasToTempFilePath({
            canvasId: "avatarCanvas",
            fileType: "jpg",
            quality: 0.8,
            success: (res) => {
              common_vendor.index.getFileSystemManager().readFile({
                filePath: res.tempFilePath,
                encoding: "base64",
                success: (base64Res) => {
                  const base64Data = "data:image/jpeg;base64," + base64Res.data;
                  resolve(base64Data);
                },
                fail: (err) => {
                  reject(new Error("读取图片文件失败: " + JSON.stringify(err)));
                }
              });
            },
            fail: (err) => {
              reject(new Error("导出图片失败: " + JSON.stringify(err)));
            }
          });
        });
      },
      fail: (err) => {
        reject(new Error("获取图片信息失败: " + JSON.stringify(err)));
      }
    });
  });
}
function uploadAvatar(avatarUrl) {
  return new Promise(async (resolve, reject) => {
    try {
      const base64Image = await getImageBase64(avatarUrl);
      const result = await utils_request.put(utils_constants.API_PATHS.CLIENT_UPDATE_PROFILE, {
        avatarUrl: base64Image
      });
      resolve(result);
    } catch (err) {
      reject(err);
    }
  });
}
function getOrderList(params) {
  return utils_request.get("/wx/client/order/list", params);
}
function getOrderDetail(params) {
  return utils_request.get("/wx/client/order/detail", params);
}
function cancelOrder(data) {
  return utils_request.post("/wx/client/order/cancel", data);
}
function getButlerImage() {
  return utils_request.get("/wx/client/butler/image");
}
const userApi = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  cancelOrder,
  getButlerImage,
  getClientInfo,
  getImageBase64,
  getOrderDetail,
  getOrderList,
  getWxLoginCode,
  getWxUserInfo,
  updateClientProfile,
  uploadAvatar,
  wxappLogin
}, Symbol.toStringTag, { value: "Module" }));
exports.getButlerImage = getButlerImage;
exports.getWxLoginCode = getWxLoginCode;
exports.userApi = userApi;
exports.wxappLogin = wxappLogin;
//# sourceMappingURL=../../.sourcemap/mp-weixin/apis/user.js.map
