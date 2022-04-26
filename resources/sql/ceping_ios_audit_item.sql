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
-- Table structure for table `ceping_ios_audit_item`
--

DROP TABLE IF EXISTS `ceping_ios_audit_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ceping_ios_audit_item` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) DEFAULT NULL,
  `item_key` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `level` varchar(10) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `score` int(11) DEFAULT NULL,
  `sort` int(11) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `solution` longtext COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ceping_ios_audit_item`
--

LOCK TABLES `ceping_ios_audit_item` WRITE;
/*!40000 ALTER TABLE `ceping_ios_audit_item` DISABLE KEYS */;
INSERT INTO `ceping_ios_audit_item` VALUES (1,1,'ios_sec_perms','权限信息','L',0,1,1,''),(2,1,'ios_sec_behavior','行为信息','L',0,2,1,''),(3,1,'ios_cvs_words','敏感词','L',0,3,1,''),(4,1,'ios_sec_other_sdk','第三方SDK检测','L',0,4,1,''),(5,1,'ios_check_certificate_type','证书类型检测','L',2,5,1,''),(6,1,'ios_appstore_risk','无法上架Appstore风险','L',2,6,1,''),(7,2,'ios_aud_code_proguard','代码未混淆风险','M',4,1,1,''),(8,2,'ios_compile_pie','未使用地址空间随机化技术风险','L',3,2,1,''),(9,2,'ios_compile_ssp','未使用编译器堆栈保护技术风险','L',3,3,1,''),(10,2,'ios_third_library_inject','注入攻击风险','H',6,4,1,''),(11,2,'ios_maco_format','可执行文件被篡改风险','M',4,5,1,''),(12,3,'ios_aud_debug','动态调试攻击风险','H',6,1,1,''),(13,3,'ios_cvs_keyboard_hijack','输入监听风险','M',4,2,1,''),(14,3,'ios_cvs_debug_info','调试日志函数调用风险','L',2,3,1,''),(15,3,'ios_cvs_webView_access_file','Webview组件跨域访问风险','H',6,4,1,''),(16,3,'ios_cvs_prison_break','越狱设备运行风险','L',3,5,1,''),(17,3,'ios_cvs_sqlite_risk','数据库明文存储风险','M',4,6,1,''),(18,3,'ios_cvs_profile_leakage','配置文件信息明文存储风险','H',5,7,1,''),(19,3,'ios_sec_other_frameworks','动态库信息泄露风险','L',3,8,1,''),(20,4,'ios_cvs_http_protocol','HTTP传输数据风险','L',3,1,1,''),(21,4,'ios_cvs_https_auth','HTTPS未校验服务器证书漏洞','L',3,2,1,''),(22,4,'ios_cvs_url_schemes','URL Schemes劫持漏洞','M',4,3,1,''),(23,5,'ios_weak_encryption','AES/DES加密算法不安全使用漏洞','L',3,1,1,''),(24,5,'ios_cvs_weak_hash','弱哈希算法使用漏洞','L',3,2,1,''),(25,5,'ios_cvs_random_risk','随机数不安全使用漏洞','L',3,3,1,''),(26,6,'ios_cvs_xcode_ghost','XcodeGhost感染漏洞','H',6,1,1,''),(27,6,'ios_cvs_high_risk_api','不安全的API函数引用风险','H',5,2,1,''),(28,6,'ios_cvs_private_api','Private Methods使用检测','H',5,3,1,''),(29,6,'ios_cvs_zip_down','ZipperDown解压漏洞','M',4,4,1,''),(30,6,'ios_cvs_iback_door','iBackDoor控制漏洞','H',6,5,1,''),(31,6,'ios_compile_arc','未使用自动管理内存技术风险','L',3,6,1,''),(32,6,'ios_compile_arc_api','内存分配函数不安全风险','L',3,7,1,''),(33,6,'ios_custom_method_long','自定义函数逻辑过于复杂风险','M',3,8,1,''),(34,7,'ios_str_leakage','明文字符串泄露风险','L',2,1,1,''),(35,7,'ios_explicit_syscall','外部函数显式调用风险','L',2,2,1,''),(36,7,'ios_syscall','系统调用暴露风险','L',2,3,1,''),(37,7,'ios_create_exec_mem','创建可执行权限内存风险','M',4,4,1,''),(38,7,'ios_resign','篡改和二次打包风险','H',6,5,1,''),(39,7,'ios_sql_exec_code','SQLite内存破坏漏洞','M',4,6,1,''),(40,7,'ios_str_format','格式化字符串漏洞','M',4,7,1,''),(41,8,'h5_ios_storage','Web Storage数据泄露风险','L',0,1,1,''),(42,8,'h5_ios_innerhtml','InnerHTML的XSS攻击漏洞','H',5,2,1,'');
/*!40000 ALTER TABLE `ceping_ios_audit_item` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-12-16 14:44:51
