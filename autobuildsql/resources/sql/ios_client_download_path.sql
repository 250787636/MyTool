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

 Date: 12/01/2022 17:21:27
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ios_client_download_path
-- ----------------------------
DROP TABLE IF EXISTS `ios_client_download_path`;
CREATE TABLE `ios_client_download_path`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '\'源码加固客户端id\'',
  `ios_client_edition` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '\'源码加固客户端版本\'',
  `ios_windows_download_path` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '\'源码加固客户端(windows)地址\'',
  `ios_mac_download_path` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '\'源码加固客户端(Mac)地址\'',
  `ios_windows_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '\'源码加固客户端(windows)名\'',
  `ios_mac_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '\'源码加固客户端(Mac)名\'',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ios_client_download_path
-- ----------------------------
INSERT INTO `ios_client_download_path` VALUES (1, '', '', '', NULL, NULL);

SET FOREIGN_KEY_CHECKS = 1;
