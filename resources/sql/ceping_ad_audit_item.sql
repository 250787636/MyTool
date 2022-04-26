-- MySQL dump 10.13  Distrib 8.0.22, for Win64 (x86_64)
--
-- Host: 172.16.102.58    Database: aimrsk
-- ------------------------------------------------------
-- Server version	8.0.18

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `ceping_ad_audit_item`
--

DROP TABLE IF EXISTS `ceping_ad_audit_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ceping_ad_audit_item` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) DEFAULT NULL,
  `item_key` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `level` varchar(10) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `score` int(11) DEFAULT NULL,
  `is_dynamic` tinyint(4) DEFAULT NULL,
  `sort` int(11) DEFAULT NULL,
  `status` bigint(20) DEFAULT NULL,
  `solution` longtext COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=83 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ceping_ad_audit_item`
--

LOCK TABLES `ceping_ad_audit_item` WRITE;
/*!40000 ALTER TABLE `ceping_ad_audit_item` DISABLE KEYS */;
INSERT INTO `ceping_ad_audit_item` VALUES (1,1,'sec_perms','权限信息','L',0,0,1,1,''),(2,1,'sec_behavior','行为信息','L',0,0,2,1,''),(3,1,'sec_virus','病毒扫描','H',6,0,3,1,''),(4,1,'sec_words','敏感词信息','L',0,0,4,1,''),(5,1,'sec_other_sdk','第三方SDK检测','L',0,0,5,1,''),(6,1,'aud_res_apk','资源文件中的Apk文件','L',0,0,6,1,''),(7,2,'aud_shield','加固壳识别','H',5,0,1,1,''),(8,2,'aud_decompile','Java代码反编译风险','H',6,0,2,1,''),(9,2,'aud_so_protect','So文件破解风险','L',3,0,3,1,''),(10,2,'aud_tamper','篡改和二次打包风险','H',6,1,4,1,''),(11,2,'aud_signaturev2','Janus签名机制漏洞','H',6,0,5,1,''),(12,2,'aud_res_protect','资源文件泄露风险','L',2,0,6,1,''),(13,2,'aud_signature','应用签名未校验风险','H',6,1,7,1,''),(14,2,'aud_code_proguard','代码未混淆风险','M',4,0,8,1,''),(15,2,'aud_certificate','使用调试证书发布应用风险','L',3,0,9,1,''),(16,2,'aud_JniRegisterNatives_risk','仅使用Java代码风险','L',3,0,10,1,''),(17,2,'aud_nolauncher_service_risk','启动隐藏服务风险','L',3,0,11,1,''),(18,2,'aud_sign_cert','应用签名算法不安全风险','L',3,0,12,1,''),(19,3,'aud_savepwd','Webview明文存储密码风险','H',5,0,1,1,''),(20,3,'aud_webview_file','Webview File同源策略绕过漏洞','H',5,0,2,1,''),(21,3,'aud_cert','明文数字证书风险','M',4,0,3,1,''),(22,3,'aud_logapi','调试日志函数调用风险','L',2,0,4,1,''),(23,3,'cvs_sqlinject','数据库注入漏洞','H',5,1,5,1,''),(24,3,'cvs_encrypt_risk','AES/DES加密方法不安全使用漏洞','L',3,0,6,1,''),(25,3,'cvs_rsa_risk','RSA加密算法不安全使用漏洞','M',4,0,7,1,''),(26,3,'aud_key_risk','密钥硬编码漏洞','H',5,0,8,1,''),(27,3,'aud_attack','动态调试攻击风险','H',6,1,9,1,''),(28,3,'aud_webview_remote_debug','Webview远程调试风险','M',4,0,10,1,''),(29,3,'aud_backup','应用数据任意备份风险','M',4,0,11,1,''),(30,3,'aud_sensapi','敏感函数调用风险','L',2,0,12,1,''),(31,3,'aud_dbglobalrw','数据库文件任意读写漏洞','M',4,0,13,1,''),(32,3,'cvs_globalrw','全局可读写的内部文件漏洞','M',4,0,14,1,''),(33,3,'cvs_sharedprefs','SharedPreferences数据全局可读写漏洞','M',4,0,15,1,''),(34,3,'cvs_sharedprefs_shareuserid','SharedUserId属性设置漏洞','M',4,0,16,1,''),(35,3,'cvs_internal_storage_mode','Internal Storage数据全局可读写漏洞','M',4,0,17,1,''),(36,3,'aud_get_dir','getDir数据全局可读写漏洞','M',4,0,18,1,''),(37,3,'aud_ffmpeg_risk','FFmpeg文件读取漏洞','H',5,0,19,1,''),(38,3,'aud_debug','Java层动态调试风险','M',4,0,20,1,''),(39,3,'aud_clipboard','剪切板敏感信息泄露漏洞','M',4,0,21,1,''),(40,3,'cvs_residua','内网测试信息残留漏洞','L',2,0,22,1,''),(41,3,'cvs_random','随机数不安全使用漏洞','L',3,0,23,1,''),(42,3,'aud_url','代码残留URL信息检测','L',2,0,24,1,''),(43,3,'aud_sensitive_account_password','残留账户密码信息检测','L',2,0,25,1,''),(44,3,'aud_sensitive_phone','残留手机号信息检测','L',2,0,26,1,''),(45,4,'aud_trans','HTTP传输数据风险','L',2,0,1,1,''),(46,4,'cvs_x509trust','HTTPS未校验服务器证书漏洞','M',4,0,2,1,''),(47,4,'aud_host_name','HTTPS未校验主机名漏洞','M',4,0,3,1,''),(48,4,'aud_intermediator_risk','HTTPS允许任意主机名漏洞','M',4,0,4,1,''),(49,4,'cvs_wv_sslerror','Webview绕过证书校验漏洞','M',4,0,5,1,''),(50,4,'cvs_packet_capture','HTTP报文信息泄漏风险','L',2,1,6,1,''),(51,5,'aud_uihijack','界面劫持风险','M',4,1,1,1,''),(52,5,'aud_kb_input','输入监听风险','M',4,0,2,1,''),(53,5,'aud_screen_shots','截屏攻击风险','M',4,0,3,1,''),(54,6,'aud_register_receiver','动态注册Receiver风险','M',4,0,1,1,''),(55,6,'cvs_dataleak','Content Provider数据泄露漏洞','H',5,1,2,1,''),(56,6,'aud_coms_activity','Activity组件导出风险','M',4,0,3,1,''),(57,6,'aud_coms_service','Service组件导出风险','M',4,0,4,1,''),(58,6,'aud_coms_receiver','Broadcast Receiver组件导出风险','M',4,0,5,1,''),(59,6,'aud_coms_provider','Content Provider组件导出风险','M',4,0,6,1,''),(60,6,'cvs_local_port','本地端口开放越权漏洞','M',4,1,7,1,''),(61,6,'aud_pendingintent','PendingIntent错误使用Intent风险','L',3,0,8,1,''),(62,6,'aud_component_hijack','Intent组件隐式调用风险','L',2,0,9,1,''),(63,6,'cvs_intent_risk','Intent Scheme URL攻击漏洞','M',4,0,10,1,''),(64,6,'cvs_fragment_risk','Fragment注入攻击漏洞','M',4,0,11,1,''),(65,6,'aud_reflect','反射调用风险','L',3,0,12,1,''),(66,7,'aud_webview_fileurl','“应用克隆”漏洞攻击风险','H',6,0,1,1,''),(67,7,'aud_inject','动态注入攻击风险','H',6,1,2,1,''),(68,7,'cvs_wv_inject','Webview远程代码执行漏洞','H',5,0,3,1,''),(69,7,'aud_webview_hide_interface','未移除有风险的Webview系统隐藏接口漏洞','H',5,0,4,1,''),(70,7,'aud_unzip','zip文件解压目录遍历漏洞','M',4,0,5,1,''),(71,7,'cvs_anydown','下载任意apk漏洞','H',5,0,6,1,''),(72,7,'cvs_refuse_service','拒绝服务攻击漏洞','M',4,1,7,1,''),(73,7,'aud_sdcard_loaddex','从sdcard加载dex风险','M',4,0,8,1,''),(74,7,'aud_sdcard_loadso','从sdcard加载so风险','M',4,0,9,1,''),(75,7,'aud_stack_protect','未使用编译器堆栈保护技术风险','L',3,0,10,1,''),(76,7,'aud_random_space','未使用地址空间随机化技术风险','L',3,0,11,1,''),(77,7,'cvs_emulator_run','模拟器运行风险','L',3,1,12,1,''),(78,7,'aud_root_device','Root设备运行风险','L',3,1,13,1,''),(79,7,'aud_risk_webBrowser','不安全的浏览器调用漏洞','M',4,0,14,1,''),(80,8,'h5_storage','Web Storage数据泄露风险','L',0,0,1,1,''),(81,8,'h5_websql','WebSQL注入漏洞','H',5,0,2,1,''),(82,8,'h5_innerhtml','InnerHTML的XSS攻击漏洞','H',5,0,3,1,'');
/*!40000 ALTER TABLE `ceping_ad_audit_item` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-12-16 14:44:37
