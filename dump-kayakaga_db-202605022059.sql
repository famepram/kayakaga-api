-- MySQL dump 10.13  Distrib 8.0.44, for Win64 (x86_64)
--
-- Host: localhost    Database: kayakaga_db
-- ------------------------------------------------------
-- Server version	8.0.44

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
-- Dumping data for table `accounts`
--

LOCK TABLES `accounts` WRITE;
/*!40000 ALTER TABLE `accounts` DISABLE KEYS */;
INSERT INTO `accounts` VALUES (1,1,1,'BCA',12100000,'#2563EB',1,'2026-04-27 13:57:31','2026-04-27 13:57:31'),(2,1,1,'Jenius',8300000,'#7F77DD',0,'2026-04-27 13:57:31','2026-04-27 13:57:31'),(3,1,2,'GoPay',3000000,'#1D9E75',0,'2026-04-27 13:57:31','2026-04-27 13:57:31'),(4,1,1,'BCA',12100000,'#2563EB',1,'2026-04-27 14:31:36','2026-04-27 14:31:36'),(5,1,1,'Jenius',8300000,'#7F77DD',0,'2026-04-27 14:31:36','2026-04-27 14:31:36'),(6,1,2,'GoPay',3000000,'#1D9E75',0,'2026-04-27 14:31:36','2026-04-27 14:31:36');
/*!40000 ALTER TABLE `accounts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `goal_milestones`
--

LOCK TABLES `goal_milestones` WRITE;
/*!40000 ALTER TABLE `goal_milestones` DISABLE KEYS */;
INSERT INTO `goal_milestones` VALUES (1,1,10000000,'2024-03-01 00:00:00','2026-04-27 13:57:56'),(2,1,20000000,NULL,'2026-04-27 13:57:56'),(3,1,50000000,NULL,'2026-04-27 13:57:56'),(4,1,100000000,NULL,'2026-04-27 13:57:56'),(5,2,5000000,'2024-06-01 00:00:00','2026-04-27 13:57:56'),(6,2,10000000,'2025-01-01 00:00:00','2026-04-27 13:57:56'),(7,2,20000000,NULL,'2026-04-27 13:57:56'),(8,2,30000000,NULL,'2026-04-27 13:57:56'),(9,1,10000000,'2024-03-01 00:00:00','2026-04-27 14:31:46'),(10,1,20000000,NULL,'2026-04-27 14:31:46'),(11,1,50000000,NULL,'2026-04-27 14:31:46'),(12,1,100000000,NULL,'2026-04-27 14:31:46'),(13,2,5000000,'2024-06-01 00:00:00','2026-04-27 14:31:46'),(14,2,10000000,'2025-01-01 00:00:00','2026-04-27 14:31:46'),(15,2,20000000,NULL,'2026-04-27 14:31:46'),(16,2,30000000,NULL,'2026-04-27 14:31:46');
/*!40000 ALTER TABLE `goal_milestones` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `goals`
--

LOCK TABLES `goals` WRITE;
/*!40000 ALTER TABLE `goals` DISABLE KEYS */;
INSERT INTO `goals` VALUES (1,1,1,1,'DP Rumah',100000000,18000000,3000000,'2028-12-01','2026-04-27 13:57:53','2026-04-27 13:57:53'),(2,1,2,2,'Dana Darurat',30000000,12600000,1200000,'2027-06-01','2026-04-27 13:57:53','2026-04-27 13:57:53'),(3,1,1,1,'DP Rumah',100000000,18000000,3000000,'2028-12-01','2026-04-27 14:31:41','2026-04-27 14:31:41'),(4,1,2,2,'Dana Darurat',30000000,12600000,1200000,'2027-06-01','2026-04-27 14:31:41','2026-04-27 14:31:41');
/*!40000 ALTER TABLE `goals` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `m_account_types`
--

LOCK TABLES `m_account_types` WRITE;
/*!40000 ALTER TABLE `m_account_types` DISABLE KEYS */;
INSERT INTO `m_account_types` VALUES (1,'savings','Tabungan','Rekening tabungan bank','bank',1,'2026-04-27 13:55:28'),(2,'ewallet','E-Wallet','Dompet digital','wallet',1,'2026-04-27 13:55:28'),(3,'investment','Investasi','Rekening investasi / sekuritas','chart',1,'2026-04-27 13:55:28');
/*!40000 ALTER TABLE `m_account_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `m_dependents`
--

LOCK TABLES `m_dependents` WRITE;
/*!40000 ALTER TABLE `m_dependents` DISABLE KEYS */;
INSERT INTO `m_dependents` VALUES (1,'single','Belum Menikah','Tidak ada tanggungan',1,'2026-04-27 13:55:53'),(2,'married','Sudah Menikah','Tanggungan pasangan',1,'2026-04-27 13:55:53'),(3,'children_1_2','Punya Anak 1-2','Tanggungan pasangan + 1-2 anak',1,'2026-04-27 13:55:53'),(4,'children_3_plus','Punya Anak 3+','Tanggungan pasangan + 3 anak atau lebih',1,'2026-04-27 13:55:53');
/*!40000 ALTER TABLE `m_dependents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `m_goal_types`
--

LOCK TABLES `m_goal_types` WRITE;
/*!40000 ALTER TABLE `m_goal_types` DISABLE KEYS */;
INSERT INTO `m_goal_types` VALUES (1,'house_dp','DP Rumah','Uang muka pembelian rumah','home',1,'2026-04-27 13:55:42'),(2,'emergency_fund','Dana Darurat','Simpanan untuk keadaan darurat','shield',1,'2026-04-27 13:55:42'),(3,'vacation','Liburan','Tabungan untuk liburan','plane',1,'2026-04-27 13:55:42'),(4,'routine_investment','Investasi Rutin','Investasi berkala jangka panjang','trending-up',1,'2026-04-27 13:55:42'),(5,'education','Dana Pendidikan','Biaya pendidikan anak/diri','book',1,'2026-04-27 13:55:42'),(6,'business_capital','Modal Usaha','Modal untuk memulai bisnis','briefcase',1,'2026-04-27 13:55:42');
/*!40000 ALTER TABLE `m_goal_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `m_risk_profiles`
--

LOCK TABLES `m_risk_profiles` WRITE;
/*!40000 ALTER TABLE `m_risk_profiles` DISABLE KEYS */;
INSERT INTO `m_risk_profiles` VALUES (1,'undecided','Belum Ditentukan','Profil risiko belum diisi',1,'2026-04-27 13:55:48'),(2,'conservative','Konservatif','Prioritas keamanan, return stabil',1,'2026-04-27 13:55:48'),(3,'moderate','Moderat','Balance risiko dan return',1,'2026-04-27 13:55:48'),(4,'aggressive','Agresif','Toleransi risiko tinggi, return maksimal',1,'2026-04-27 13:55:48');
/*!40000 ALTER TABLE `m_risk_profiles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `m_transaction_categories`
--

LOCK TABLES `m_transaction_categories` WRITE;
/*!40000 ALTER TABLE `m_transaction_categories` DISABLE KEYS */;
INSERT INTO `m_transaction_categories` VALUES (1,'food_beverage','Makanan & Minuman',NULL,'utensils','#D85A30',1,1,'2026-04-27 13:55:33'),(2,'transport','Transportasi',NULL,'car','#7F77DD',1,1,'2026-04-27 13:55:33'),(3,'entertainment','Hiburan',NULL,'gamepad','#BA7517',1,1,'2026-04-27 13:55:33'),(4,'bills','Tagihan',NULL,'zap','#2563EB',1,1,'2026-04-27 13:55:33'),(5,'shopping','Belanja',NULL,'shopping-bag','#E24B4A',1,1,'2026-04-27 13:55:33'),(6,'health','Kesehatan',NULL,'heart','#1D9E75',1,1,'2026-04-27 13:55:33'),(7,'investment','Investasi',NULL,'trending-up','#085041',1,1,'2026-04-27 13:55:33'),(8,'other','Lainnya',NULL,'more-horizontal','#888780',1,1,'2026-04-27 13:55:33'),(9,'income','Pemasukan',NULL,'arrow-down-circle','#1D9E75',0,1,'2026-04-27 13:55:33');
/*!40000 ALTER TABLE `m_transaction_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `m_transaction_sources`
--

LOCK TABLES `m_transaction_sources` WRITE;
/*!40000 ALTER TABLE `m_transaction_sources` DISABLE KEYS */;
INSERT INTO `m_transaction_sources` VALUES (1,'manual','Manual','Diinput manual oleh user',1,'2026-04-27 13:55:37'),(2,'csv','Import CSV','Diimport dari file CSV bank',1,'2026-04-27 13:55:37'),(3,'receipt','Foto Struk','Diparsing dari foto struk via AI',1,'2026-04-27 13:55:37'),(4,'bank_sync','Bank Sync','Sinkronisasi otomatis dari bank API',1,'2026-04-27 13:55:37');
/*!40000 ALTER TABLE `m_transaction_sources` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `refresh_tokens`
--

LOCK TABLES `refresh_tokens` WRITE;
/*!40000 ALTER TABLE `refresh_tokens` DISABLE KEYS */;
INSERT INTO `refresh_tokens` VALUES (1,3,'29ecd6b6-2af9-46b6-b220-572997dcbb65','2026-05-28 14:53:20',0,'2026-04-28 14:53:20'),(2,3,'79fe0b93-0685-4f2e-907a-4284283741fc','2026-05-28 15:13:37',0,'2026-04-28 15:13:37');
/*!40000 ALTER TABLE `refresh_tokens` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `transactions`
--

LOCK TABLES `transactions` WRITE;
/*!40000 ALTER TABLE `transactions` DISABLE KEYS */;
INSERT INTO `transactions` VALUES (1,1,1,9,1,'2026-04-25','08:00:00','Gaji April 2026',12000000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(2,1,1,4,1,'2026-04-08','07:00:00','PLN Token Listrik',-350000,NULL,0,1,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(3,1,1,4,1,'2026-04-15','07:00:00','Internet Provider',-250000,NULL,0,1,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(4,1,1,4,1,'2026-04-22','07:00:00','BPJS Kesehatan',-150000,NULL,0,1,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(5,1,1,4,1,'2026-04-12','07:00:00','PDAM',-85000,NULL,0,1,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(6,1,1,3,1,'2026-04-26','00:00:00','Netflix',-186000,NULL,0,1,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(7,1,1,3,1,'2026-04-06','00:00:00','Spotify',-55000,NULL,0,1,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(8,1,1,1,1,'2026-04-01','12:30:00','Warung Makan Padang',-32000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(9,1,1,1,1,'2026-04-03','08:00:00','Kopi Kenangan',-38000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(10,1,1,1,1,'2026-04-07','13:00:00','McDonalds',-75000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(11,1,1,1,1,'2026-04-17','19:00:00','Restoran Sunda',-120000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(12,1,3,1,1,'2026-04-02','12:00:00','GoFood - Ayam Geprek',-35000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(13,1,3,1,1,'2026-04-09','19:00:00','GoFood - Mie Ayam',-28000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(14,1,3,1,1,'2026-04-16','12:30:00','GoFood - Nasi Padang',-32000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(15,1,3,1,1,'2026-04-23','20:00:00','GoFood - Pizza',-85000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(16,1,3,2,1,'2026-04-01','09:00:00','Grab',-45000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(17,1,3,2,1,'2026-04-08','08:30:00','Grab',-38000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(18,1,3,2,1,'2026-04-15','09:15:00','Grab',-42000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(19,1,3,2,1,'2026-04-22','08:45:00','Grab',-40000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(20,1,1,5,1,'2026-04-05','14:00:00','Indomaret',-89000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(21,1,1,5,1,'2026-04-15','16:00:00','Indomaret',-153000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(22,1,1,5,1,'2026-04-26','14:10:00','Indomaret',-145000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(23,1,1,6,1,'2026-04-10','10:00:00','Apotek K24',-85000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(24,1,1,6,1,'2026-04-18','11:00:00','Klinik Pratama',-165000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(25,1,1,3,1,'2026-04-11','14:00:00','CGV Cinema',-75000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(26,1,1,3,1,'2026-04-19','15:00:00','Steam Games',-450000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(27,1,1,1,1,'2026-04-20','20:00:00','Seafood Ancol',-350000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(28,1,2,7,1,'2026-04-05','10:00:00','Transfer Tabungan DP Rumah',-3000000,NULL,0,1,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(29,1,2,7,1,'2026-04-05','10:01:00','Transfer Dana Darurat',-1200000,NULL,0,1,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(30,1,2,5,1,'2026-04-14','13:00:00','Tokopedia',-280000,NULL,0,0,'2026-04-27 13:58:11','2026-04-27 13:58:11'),(31,1,1,9,1,'2026-04-25','08:00:00','Gaji April 2026',12000000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(32,1,1,4,1,'2026-04-08','07:00:00','PLN Token Listrik',-350000,NULL,0,1,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(33,1,1,4,1,'2026-04-15','07:00:00','Internet Provider',-250000,NULL,0,1,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(34,1,1,4,1,'2026-04-22','07:00:00','BPJS Kesehatan',-150000,NULL,0,1,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(35,1,1,4,1,'2026-04-12','07:00:00','PDAM',-85000,NULL,0,1,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(36,1,1,3,1,'2026-04-26','00:00:00','Netflix',-186000,NULL,0,1,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(37,1,1,3,1,'2026-04-06','00:00:00','Spotify',-55000,NULL,0,1,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(38,1,1,1,1,'2026-04-01','12:30:00','Warung Makan Padang',-32000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(39,1,1,1,1,'2026-04-03','08:00:00','Kopi Kenangan',-38000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(40,1,1,1,1,'2026-04-07','13:00:00','McDonalds',-75000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(41,1,1,1,1,'2026-04-17','19:00:00','Restoran Sunda',-120000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(42,1,3,1,1,'2026-04-02','12:00:00','GoFood - Ayam Geprek',-35000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(43,1,3,1,1,'2026-04-09','19:00:00','GoFood - Mie Ayam',-28000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(44,1,3,1,1,'2026-04-16','12:30:00','GoFood - Nasi Padang',-32000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(45,1,3,1,1,'2026-04-23','20:00:00','GoFood - Pizza',-85000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(46,1,3,2,1,'2026-04-01','09:00:00','Grab',-45000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(47,1,3,2,1,'2026-04-08','08:30:00','Grab',-38000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(48,1,3,2,1,'2026-04-15','09:15:00','Grab',-42000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(49,1,3,2,1,'2026-04-22','08:45:00','Grab',-40000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(50,1,1,5,1,'2026-04-05','14:00:00','Indomaret',-89000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(51,1,1,5,1,'2026-04-15','16:00:00','Indomaret',-153000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(52,1,1,5,1,'2026-04-26','14:10:00','Indomaret',-145000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(53,1,1,6,1,'2026-04-10','10:00:00','Apotek K24',-85000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(54,1,1,6,1,'2026-04-18','11:00:00','Klinik Pratama',-165000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(55,1,1,3,1,'2026-04-11','14:00:00','CGV Cinema',-75000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(56,1,1,3,1,'2026-04-19','15:00:00','Steam Games',-450000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(57,1,1,1,1,'2026-04-20','20:00:00','Seafood Ancol',-350000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(58,1,2,7,1,'2026-04-05','10:00:00','Transfer Tabungan DP Rumah',-3000000,NULL,0,1,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(59,1,2,7,1,'2026-04-05','10:01:00','Transfer Dana Darurat',-1200000,NULL,0,1,'2026-04-27 14:31:57','2026-04-27 14:31:57'),(60,1,2,5,1,'2026-04-14','13:00:00','Tokopedia',-280000,NULL,0,0,'2026-04-27 14:31:57','2026-04-27 14:31:57');
/*!40000 ALTER TABLE `transactions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `user_credentials`
--

LOCK TABLES `user_credentials` WRITE;
/*!40000 ALTER TABLE `user_credentials` DISABLE KEYS */;
INSERT INTO `user_credentials` VALUES (1,1,'andi@finai.dev','$2a$10$GANTI_DENGAN_BCRYPT_HASH_ASLI','2026-04-27 13:57:24','2026-04-27 13:57:24'),(3,3,'femmy.pramana@gmail.com','$2a$10$T1dypSnZegsnnnf56CUWXegFlL42qNzTnQ9Bh8XvgFaA3mgoBu2JG','2026-04-28 14:53:20','2026-04-28 14:53:20');
/*!40000 ALTER TABLE `user_credentials` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'Andi Pratama','Jakarta','Karyawan swasta',1,12000000,7000000,30000000,1,'IDR','2026-04-27 13:57:09','2026-04-27 13:57:09'),(2,'Andi Pratama','Jakarta','Karyawan swasta',1,12000000,7000000,30000000,1,'IDR','2026-04-27 14:30:56','2026-04-27 14:30:56'),(3,'Femmy Pramana','','',1,0,0,0,1,'IDR','2026-04-28 14:53:20','2026-04-28 14:53:20');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-05-02 20:59:53
