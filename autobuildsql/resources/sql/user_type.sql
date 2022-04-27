/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80026
 Source Host           : localhost:3306
 Source Schema         : middlegroundpcitc

 Target Server Type    : MySQL
 Target Server Version : 80026
 File Encoding         : 65001

 Date: 01/03/2022 14:53:58
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user_type
-- ----------------------------
DROP TABLE IF EXISTS `user_type`;
CREATE TABLE `user_type`  (
  `id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '用户类型 0-普通  1-安控管理员 2-安全监测员 3-沙箱管理员 4-沙箱审计员',
  `u_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '用户类型 普通  安控管理员 安全监测员 沙箱管理员 沙箱审计员',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_type
-- ----------------------------
INSERT INTO `user_type` VALUES ('0', '普通用户');
INSERT INTO `user_type` VALUES ('1', '安控管理员');
INSERT INTO `user_type` VALUES ('2', '安全监测员');
INSERT INTO `user_type` VALUES ('3', '沙箱管理员');
INSERT INTO `user_type` VALUES ('4', '沙箱审计员');

SET FOREIGN_KEY_CHECKS = 1;
