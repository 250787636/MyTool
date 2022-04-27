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
-- Table structure for table `ceping_audit_category`
--

DROP TABLE IF EXISTS `ceping_audit_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ceping_audit_category` (
  `type` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
  `id` int(11) NOT NULL,
  `category_key` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `category_name` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `category_sort` int(11) DEFAULT NULL,
  PRIMARY KEY (`type`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ceping_audit_category`
--

LOCK TABLES `ceping_audit_category` WRITE;
/*!40000 ALTER TABLE `ceping_audit_category` DISABLE KEYS */;
INSERT INTO `ceping_audit_category` VALUES ('ad',1,'sec','自身安全',1),('ad',2,'src_check','程序源文件安全',2),('ad',3,'local_data','本地数据存储安全',3),('ad',4,'trans_safe','通信数据传输安全',4),('ad',5,'identity_check','身份认证安全',5),('ad',6,'internal_data_exchange','内部数据交互安全',6),('ad',7,'prevention_attack','恶意攻击防范能力',7),('ad',8,'h5_hybrid','HTML5安全',8),('ad',9,'other_sdk_detect','第三方SDK检测',9),('ad',10,'content_security','内容安全',10),('ad',11,'optimization_suggestion','优化建议',11),('ios',1,'ios_sec','自身安全',1),('ios',2,'ios_src_check','二进制代码保护',2),('ios',3,'ios_local_data','客户端数据存储安全',3),('ios',4,'ios_trans_safe','数据传输安全',4),('ios',5,'ios_encryption','加密算法及密码安全',5),('ios',6,'ios_app_sec','iOS应用安全规范',6),('ios',7,'ios_source_safe','程序源文件安全',7),('ios',8,'h5_hybrid','HTML5安全',8),('ios',9,'other_sdk_detect','第三方SDK检测',9),('ios',10,'content_security','内容安全',10),('mp',1,'mp_sec','自身安全',1),('mp',2,'mp_comm_xfe','通信传输安全检测',2),('mp',3,'mp_data_reveal','数据泄漏检测',3),('mp',4,'mp_com_risk','组件漏洞检测',4),('mp',5,'mp_unsafe_setting','HTTP不安全配置检测',5),('sdk',1,'sdk_sec','自身安全',1),('sdk',2,'sdk_risk','风险漏洞',2);
/*!40000 ALTER TABLE `ceping_audit_category` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-12-02 21:49:54
