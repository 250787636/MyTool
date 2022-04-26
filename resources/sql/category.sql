/*
 Navicat Premium Data Transfer

 Source Server         : 172.16.38.215
 Source Server Type    : MySQL
 Source Server Version : 50733
 Source Host           : 172.16.38.215:33306
 Source Schema         : middlegroundabc

 Target Server Type    : MySQL
 Target Server Version : 50733
 File Encoding         : 65001

 Date: 16/12/2021 11:27:39
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `category_name` longtext,
  `ce_ping_type` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of category
-- ----------------------------
BEGIN;
INSERT INTO `category` VALUES (1, '自身安全', 'ios');
INSERT INTO `category` VALUES (2, '二进制代码保护', 'ios');
INSERT INTO `category` VALUES (3, '客户端数据存储安全', 'ios');
INSERT INTO `category` VALUES (4, '数据传输安全', 'ios');
INSERT INTO `category` VALUES (5, '加密算法及密码安全', 'ios');
INSERT INTO `category` VALUES (6, 'iOS应用安全规范', 'ios');
INSERT INTO `category` VALUES (7, '程序源文件安全', 'ios');
INSERT INTO `category` VALUES (8, 'HTML5安全', 'ios');
INSERT INTO `category` VALUES (9, '自身安全', 'ad');
INSERT INTO `category` VALUES (10, '程序源文件安全', 'ad');
INSERT INTO `category` VALUES (11, '本地数据存储安全', 'ad');
INSERT INTO `category` VALUES (12, '通信数据传输安全', 'ad');
INSERT INTO `category` VALUES (13, '身份认证安全', 'ad');
INSERT INTO `category` VALUES (14, '内部数据交互安全', 'ad');
INSERT INTO `category` VALUES (15, '恶意攻击防范能力', 'ad');
INSERT INTO `category` VALUES (16, 'HTML5安全', 'ad');
INSERT INTO `category` VALUES (17, '自身安全', 'mp');
INSERT INTO `category` VALUES (18, '通信传输安全检测', 'mp');
INSERT INTO `category` VALUES (19, '数据泄漏检测', 'mp');
INSERT INTO `category` VALUES (20, '组件漏洞检测', 'mp');
INSERT INTO `category` VALUES (21, 'HTTP不安全配置检测', 'mp');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
