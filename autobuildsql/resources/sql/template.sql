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

 Date: 23/12/2021 15:24:27
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for template
-- ----------------------------
DROP TABLE IF EXISTS `template`;
CREATE TABLE `template` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `template_name` varchar(100) NOT NULL DEFAULT '' COMMENT '任务ID',
  `created_id` int(11) DEFAULT NULL COMMENT '创建者ID',
  `template_type` varchar(100) NOT NULL DEFAULT '' COMMENT '模板类型',
  `items` longtext NOT NULL COMMENT '模板内容',
  `is_owasp` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否OWASP',
  `report_language` varchar(100) NOT NULL DEFAULT '' COMMENT '报告语言',
  PRIMARY KEY (`id`),
  KEY `idx_template_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=55 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of template
-- ----------------------------
BEGIN;
INSERT INTO `template` VALUES (1, '2021-12-23 15:22:09.869', '2021-12-23 15:22:09.869', NULL, 'Android-全量模板', 0, 'ad', '[\"h5_storage\",\"h5_websql\",\"h5_innerhtml\",\"aud_register_receiver\",\"cvs_dataleak\",\"aud_coms_activity\",\"aud_coms_service\",\"aud_coms_receiver\",\"aud_coms_provider\",\"cvs_local_port\",\"aud_pendingintent\",\"aud_component_hijack\",\"cvs_intent_risk\",\"cvs_fragment_risk\",\"aud_reflect\",\"aud_webview_fileurl\",\"aud_inject\",\"cvs_wv_inject\",\"aud_webview_hide_interface\",\"aud_unzip\",\"cvs_anydown\",\"cvs_refuse_service\",\"aud_sdcard_loaddex\",\"aud_sdcard_loadso\",\"aud_stack_protect\",\"aud_random_space\",\"aud_root_device\",\"aud_risk_webBrowser\",\"aud_savepwd\",\"aud_webview_file\",\"aud_cert\",\"aud_logapi\",\"cvs_sqlinject\",\"cvs_encrypt_risk\",\"cvs_rsa_risk\",\"aud_key_risk\",\"aud_attack\",\"aud_webview_remote_debug\",\"aud_backup\",\"aud_sensapi\",\"aud_dbglobalrw\",\"cvs_globalrw\",\"cvs_sharedprefs\",\"cvs_sharedprefs_shareuserid\",\"cvs_internal_storage_mode\",\"aud_get_dir\",\"aud_ffmpeg_risk\",\"aud_debug\",\"aud_clipboard\",\"cvs_residua\",\"cvs_random\",\"aud_url\",\"aud_sensitive_account_password\",\"aud_sensitive_phone\",\"aud_shield\",\"aud_decompile\",\"aud_so_protect\",\"aud_tamper\",\"aud_signaturev2\",\"aud_res_protect\",\"aud_signature\",\"aud_code_proguard\",\"aud_certificate\",\"aud_JniRegisterNatives_risk\",\"aud_nolauncher_service_risk\",\"aud_sign_cert\",\"sec_perms\",\"sec_behavior\",\"sec_virus\",\"sec_words\",\"sec_other_sdk\",\"aud_res_apk\",\"aud_uihijack\",\"aud_kb_input\",\"aud_screen_shots\",\"aud_trans\",\"cvs_x509trust\",\"aud_host_name\",\"aud_intermediator_risk\",\"cvs_wv_sslerror\",\"cvs_packet_capture\"]', 2, 'zh_cn');
INSERT INTO `template` VALUES (2, '2021-12-23 15:23:40.312', '2021-12-23 15:23:40.312', NULL, 'IOS-全量模板', 0, 'ios', '[\"h5_ios_storage\",\"h5_ios_innerhtml\",\"ios_cvs_xcode_ghost\",\"ios_cvs_high_risk_api\",\"ios_cvs_private_api\",\"ios_cvs_zip_down\",\"ios_cvs_iback_door\",\"ios_compile_arc\",\"ios_compile_arc_api\",\"ios_custom_method_long\",\"ios_aud_code_proguard\",\"ios_compile_pie\",\"ios_compile_ssp\",\"ios_third_library_inject\",\"ios_maco_format\",\"ios_weak_encryption\",\"ios_cvs_weak_hash\",\"ios_cvs_random_risk\",\"ios_aud_debug\",\"ios_cvs_keyboard_hijack\",\"ios_cvs_debug_info\",\"ios_cvs_webView_access_file\",\"ios_cvs_prison_break\",\"ios_cvs_sqlite_risk\",\"ios_cvs_profile_leakage\",\"ios_sec_other_frameworks\",\"ios_cvs_http_protocol\",\"ios_cvs_https_auth\",\"ios_cvs_url_schemes\",\"ios_str_leakage\",\"ios_explicit_syscall\",\"ios_syscall\",\"ios_create_exec_mem\",\"ios_resign\",\"ios_sql_exec_code\",\"ios_str_format\",\"ios_sec_perms\",\"ios_sec_behavior\",\"ios_cvs_words\",\"ios_sec_other_sdk\",\"ios_check_certificate_type\",\"ios_appstore_risk\"]', 2, 'zh_cn');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
