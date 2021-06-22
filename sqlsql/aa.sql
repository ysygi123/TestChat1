-- --------------------------------------------------------
-- 主机:                           192.168.3.36
-- 服务器版本:                        5.7.28-log - MySQL Community Server (GPL)
-- 服务器操作系统:                      linux-glibc2.12
-- HeidiSQL 版本:                  10.3.0.5827
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

-- 导出  表 mywork1.message 结构
CREATE TABLE IF NOT EXISTS `message` (
	`id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	`message` VARCHAR(1000) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
	`send_uid` INT(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '发送者uid',
	`receive_uid` INT(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '接受者uid',
	`created_time` INT(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '创建时间戳',
	PRIMARY KEY (`id`) USING BTREE
);

-- 正在导出表  mywork1.message 的数据：~3 rows (大约)
DELETE FROM `message`;
/*!40000 ALTER TABLE `message` DISABLE KEYS */;
INSERT INTO `message` (`id`, `message`, `send_uid`, `receive_uid`, `created_time`) VALUES
	(1, '去你妈几把', 1, 2, 1624370141),
	(2, '去你妈几把', 1, 2, 1624370148),
	(3, '去你妈几把', 1, 2, 1624370149);
/*!40000 ALTER TABLE `message` ENABLE KEYS */;

-- 导出  表 mywork1.message_list 结构
CREATE TABLE IF NOT EXISTS `message_list` (
	`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	`uid` INT(10) UNSIGNED NOT NULL DEFAULT '0',
	`message` VARCHAR(50) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
	`message_type` TINYINT(4) NOT NULL DEFAULT '0' COMMENT '1用户消息 2群消息',
	`created_time` INT(10) UNSIGNED NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `uid` (`uid`) USING BTREE
);

-- 正在导出表  mywork1.message_list 的数据：~0 rows (大约)
DELETE FROM `message_list`;
/*!40000 ALTER TABLE `message_list` DISABLE KEYS */;
/*!40000 ALTER TABLE `message_list` ENABLE KEYS */;

-- 导出  表 mywork1.user 结构
CREATE TABLE IF NOT EXISTS `user` (
	`uid` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(255) NOT NULL COMMENT '用户名' COLLATE 'utf8_unicode_ci',
	`rname` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '邮件' COLLATE 'utf8_unicode_ci',
	`mobile` VARCHAR(11) NOT NULL DEFAULT '' COMMENT '电话' COLLATE 'utf8_unicode_ci',
	`passwd` CHAR(32) NOT NULL DEFAULT '' COMMENT '密码' COLLATE 'utf8_unicode_ci',
	PRIMARY KEY (`uid`) USING BTREE,
	UNIQUE INDEX `username` (`username`) USING BTREE
);

-- 正在导出表  mywork1.user 的数据：~2 rows (大约)
DELETE FROM `user`;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` (`uid`, `username`, `rname`, `mobile`, `passwd`) VALUES
	(1, '10', '10', '10', 'd3d9446802a44259755d38e6d163e820'),
	(2, '11', '11', '11', '6512bd43d9caa6e02c990b0a82652dca');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
