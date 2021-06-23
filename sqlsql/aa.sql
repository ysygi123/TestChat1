-- --------------------------------------------------------
-- 主机:                           192.168.199.112
-- 服务器版本:                        8.0.17 - Source distribution
-- 服务器操作系统:                      Linux
-- HeidiSQL 版本:                  10.3.0.5827
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- 导出 testchat1 的数据库结构
CREATE DATABASE IF NOT EXISTS `testchat1` /*!40100 DEFAULT CHARACTER SET utf8 */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `testchat1`;

-- 导出  表 testchat1.group 结构
DROP TABLE IF EXISTS `group`;
CREATE TABLE IF NOT EXISTS `group` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`group_name` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '群名' COLLATE 'utf8mb4_general_ci',
	`people_num` SMALLINT(5) NOT NULL DEFAULT '0' COMMENT '群人数',
	`created_time` INT(10) NOT NULL DEFAULT '0',
	`update_time` INT(10) NOT NULL DEFAULT '0',
	`is_del` TINYINT(3) NOT NULL DEFAULT '0' COMMENT '1不删2删',
	PRIMARY KEY (`id`) USING BTREE
);

-- 正在导出表  testchat1.group 的数据：~0 rows (大约)
DELETE FROM `group`;
/*!40000 ALTER TABLE `group` DISABLE KEYS */;
/*!40000 ALTER TABLE `group` ENABLE KEYS */;

-- 导出  表 testchat1.group_message 结构
DROP TABLE IF EXISTS `group_message`;
CREATE TABLE IF NOT EXISTS `group_message` (
	`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	`group_id` INT(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '群id',
	`message_content` VARCHAR(1000) NOT NULL DEFAULT '0' COLLATE 'utf8mb4_general_ci',
	`created_time` INT(10) UNSIGNED NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE
);

-- 正在导出表  testchat1.group_message 的数据：~0 rows (大约)
DELETE FROM `group_message`;
/*!40000 ALTER TABLE `group_message` DISABLE KEYS */;
/*!40000 ALTER TABLE `group_message` ENABLE KEYS */;

-- 导出  表 testchat1.group_users 结构
DROP TABLE IF EXISTS `group_users`;
CREATE TABLE IF NOT EXISTS `group_users` (
	`id` INT(10) NOT NULL AUTO_INCREMENT,
	`created_time` INT(10) NOT NULL DEFAULT '0',
	`update_time` INT(10) NOT NULL DEFAULT '0',
	`uid` INT(10) NOT NULL DEFAULT '0',
	`is_del` TINYINT(3) NOT NULL DEFAULT '0' COMMENT '1不删2删',
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `uid` (`uid`) USING BTREE
);

-- 正在导出表  testchat1.group_users 的数据：~0 rows (大约)
DELETE FROM `group_users`;
/*!40000 ALTER TABLE `group_users` DISABLE KEYS */;
/*!40000 ALTER TABLE `group_users` ENABLE KEYS */;

-- 导出  表 testchat1.message 结构
DROP TABLE IF EXISTS `message`;
CREATE TABLE IF NOT EXISTS `message` (
	`id` INT(10) NOT NULL AUTO_INCREMENT,
	`message_content` VARCHAR(1000) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
	`send_uid` INT(10) NOT NULL DEFAULT '0' COMMENT '发送者uid',
	`receive_uid` INT(10) NOT NULL DEFAULT '0' COMMENT '接受者uid',
	`created_time` INT(10) NOT NULL DEFAULT '0' COMMENT '创建时间戳',
	PRIMARY KEY (`id`) USING BTREE
);

-- 正在导出表  testchat1.message 的数据：~3 rows (大约)
DELETE FROM `message`;
/*!40000 ALTER TABLE `message` DISABLE KEYS */;
/*!40000 ALTER TABLE `message` ENABLE KEYS */;

-- 导出  表 testchat1.message_list 结构
DROP TABLE IF EXISTS `message_list`;
CREATE TABLE IF NOT EXISTS `message_list` (
	`id` INT(10) NOT NULL AUTO_INCREMENT,
	`uid` INT(10) NOT NULL DEFAULT '0',
	`from_id` INT(10) NOT NULL DEFAULT '0' COMMENT '来源id',
	`message_content` VARCHAR(50) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
	`message_type` TINYINT(4) NOT NULL DEFAULT '0' COMMENT '1用户消息message 2群消息 3加好友请求 4群邀请',
	`created_time` INT(10) UNSIGNED NOT NULL DEFAULT '0',
	`update_time` INT(10) UNSIGNED NOT NULL DEFAULT '0',
	`message_num` TINYINT(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '消息数量',
	`is_del` TINYINT(3) NOT NULL DEFAULT '1' COMMENT '1不删 2删',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `from_id_message_type_is_del` (`from_id`, `message_type`, `is_del`) USING BTREE,
	INDEX `uid` (`uid`) USING BTREE
);

-- 正在导出表  testchat1.message_list 的数据：~0 rows (大约)
DELETE FROM `message_list`;
/*!40000 ALTER TABLE `message_list` DISABLE KEYS */;
/*!40000 ALTER TABLE `message_list` ENABLE KEYS */;

-- 导出  表 testchat1.user 结构
DROP TABLE IF EXISTS `user`;
CREATE TABLE IF NOT EXISTS `user` (
	`uid` INT(10) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(255) NOT NULL COMMENT '用户名' COLLATE 'utf8_unicode_ci',
	`rname` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '邮件' COLLATE 'utf8_unicode_ci',
	`mobile` VARCHAR(11) NOT NULL DEFAULT '' COMMENT '电话' COLLATE 'utf8_unicode_ci',
	`passwd` CHAR(32) NOT NULL DEFAULT '' COMMENT '密码' COLLATE 'utf8_unicode_ci',
	`created_time` INT(10) UNSIGNED NOT NULL DEFAULT '0',
	`update_time` INT(10) UNSIGNED NOT NULL DEFAULT '0',
	PRIMARY KEY (`uid`) USING BTREE,
	UNIQUE INDEX `username` (`username`) USING BTREE
);

-- 正在导出表  testchat1.user 的数据：~2 rows (大约)
DELETE FROM `user`;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` (`uid`, `username`, `rname`, `mobile`, `passwd`, `created_time`, `update_time`) VALUES
	(1, '10', '10', '10', 'd3d9446802a44259755d38e6d163e820', 0, 0),
	(2, '11', '11', '11', '6512bd43d9caa6e02c990b0a82652dca', 0, 0);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;

-- 导出  表 testchat1.user_add_friend_request 结构
DROP TABLE IF EXISTS `user_add_friend_request`;
CREATE TABLE IF NOT EXISTS `user_add_friend_request` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`request_uid` INT(11) NOT NULL DEFAULT '0',
	`receive_uid` INT(11) NOT NULL DEFAULT '0',
	`created_time` INT(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `request_uid` (`request_uid`) USING BTREE,
	INDEX `receive_uid` (`receive_uid`) USING BTREE
);

-- 正在导出表  testchat1.user_add_friend_request 的数据：~0 rows (大约)
DELETE FROM `user_add_friend_request`;
/*!40000 ALTER TABLE `user_add_friend_request` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_add_friend_request` ENABLE KEYS */;

-- 导出  表 testchat1.user_friends 结构
DROP TABLE IF EXISTS `user_friends`;
CREATE TABLE IF NOT EXISTS `user_friends` (
	`id` INT(10) NOT NULL AUTO_INCREMENT,
	`uid` INT(10) NOT NULL DEFAULT '0',
	`friend_uid` INT(10) NOT NULL DEFAULT '0',
	`created_time` INT(10) UNSIGNED NOT NULL DEFAULT '0',
	`update_time` INT(10) UNSIGNED NOT NULL DEFAULT '0',
	`is_del` TINYINT(3) NOT NULL DEFAULT '1' COMMENT '1不删2删',
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `uid` (`uid`) USING BTREE
);

-- 正在导出表  testchat1.user_friends 的数据：~0 rows (大约)
DELETE FROM `user_friends`;
/*!40000 ALTER TABLE `user_friends` DISABLE KEYS */;
INSERT INTO `user_friends` (`id`, `uid`, `friend_uid`, `created_time`, `update_time`, `is_del`) VALUES
	(1, 1, 2, 0, 0, 1),
	(2, 2, 1, 0, 0, 1);
/*!40000 ALTER TABLE `user_friends` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
