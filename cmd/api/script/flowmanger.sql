/*
 Navicat Premium Data Transfer

 Source Server         : l1.ttut.cc
 Source Server Type    : MySQL
 Source Server Version : 80031
 Source Host           : l1.ttut.cc:38583
 Source Schema         : flowmanger

 Target Server Type    : MySQL
 Target Server Version : 80031
 File Encoding         : 65001

 Date: 20/09/2024 16:54:36
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ip_account
-- ----------------------------
DROP TABLE IF EXISTS `ip_account`;
CREATE TABLE `ip_account` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '账号ID',
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '账号名称',
  `stop_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '账号到期时间',
  `total_flow` bigint unsigned NOT NULL DEFAULT '0' COMMENT '总流量',
  `used_traffic` bigint unsigned NOT NULL DEFAULT '0' COMMENT '已使用流量',
  `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '账号状态',
  `user_level` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '用户会员等级',
  `create_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '注册时间',
  `update_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '会员信息上次更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of ip_account
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for ip_announcements
-- ----------------------------
DROP TABLE IF EXISTS `ip_announcements`;
CREATE TABLE `ip_announcements` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '公告ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '标题',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '内容',
  `enable` tinyint NOT NULL DEFAULT '0' COMMENT '是否启用',
  `created_time` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_time` bigint NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of ip_announcements
-- ----------------------------
BEGIN;
INSERT INTO `ip_announcements` (`id`, `title`, `content`, `enable`, `created_time`, `updated_time`) VALUES (1, '测试', '测试', 1, 1709095, 17090900);
INSERT INTO `ip_announcements` (`id`, `title`, `content`, `enable`, `created_time`, `updated_time`) VALUES (2, '测试2', '测试2', 1, 1709096, 17090900);
INSERT INTO `ip_announcements` (`id`, `title`, `content`, `enable`, `created_time`, `updated_time`) VALUES (3, '测试3', '测试3', 1, 1709097, 17090900);
INSERT INTO `ip_announcements` (`id`, `title`, `content`, `enable`, `created_time`, `updated_time`) VALUES (4, '测试4', '测试4', 1, 1709098, 17090900);
INSERT INTO `ip_announcements` (`id`, `title`, `content`, `enable`, `created_time`, `updated_time`) VALUES (5, '测试5', '测试5', 1, 1709099, 17090900);
COMMIT;

-- ----------------------------
-- Table structure for ip_assets
-- ----------------------------
DROP TABLE IF EXISTS `ip_assets`;
CREATE TABLE `ip_assets` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '资产ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `commodity_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '套餐(商品)名称',
  `expire_time` bigint NOT NULL DEFAULT '0' COMMENT '到期时间戳',
  `type` tinyint NOT NULL DEFAULT '1' COMMENT '类型',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态:未使用/已使用/使用完毕',
  `total_count` bigint NOT NULL DEFAULT '0' COMMENT '资产资产数量',
  `used_count` bigint NOT NULL DEFAULT '0' COMMENT '资产已分配数量/已使用数量',
  `unit` tinyint NOT NULL DEFAULT '1' COMMENT '单位',
  `created_time` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_time` bigint NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_ip_assets_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of ip_assets
-- ----------------------------
BEGIN;
INSERT INTO `ip_assets` (`id`, `user_id`, `commodity_name`, `expire_time`, `type`, `status`, `total_count`, `used_count`, `unit`, `created_time`, `updated_time`) VALUES (1, 1, '测试商品', 1728489600, 1, 1, 10, 5, 1, 1725939360, 1725939360);
INSERT INTO `ip_assets` (`id`, `user_id`, `commodity_name`, `expire_time`, `type`, `status`, `total_count`, `used_count`, `unit`, `created_time`, `updated_time`) VALUES (2, 1, '测试商品', 1728489600, 2, 1, 10, 5, 1, 1725939360, 1725939360);
COMMIT;

-- ----------------------------
-- Table structure for ip_commodities
-- ----------------------------
DROP TABLE IF EXISTS `ip_commodities`;
CREATE TABLE `ip_commodities` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '名称',
  `price` decimal(10,2) DEFAULT '0.00' COMMENT '价格',
  `original_price` decimal(10,2) DEFAULT '0.00' COMMENT '原价',
  `enable` tinyint NOT NULL DEFAULT '0' COMMENT '是否启用',
  `duration` bigint NOT NULL DEFAULT '0' COMMENT '有效时长',
  `currency` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '货币',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '描述',
  `weight` bigint NOT NULL DEFAULT '0' COMMENT '权重',
  `created_time` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_time` bigint NOT NULL DEFAULT '0' COMMENT '更新时间',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT '商品类型',
  `unit` tinyint NOT NULL DEFAULT '1' COMMENT '单位',
  `total_count` bigint NOT NULL DEFAULT '0' COMMENT '对应资产数量',
  `duration_type_id` bigint NOT NULL DEFAULT '0' COMMENT '时长类型ID',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of ip_commodities
-- ----------------------------
BEGIN;
INSERT INTO `ip_commodities` (`id`, `name`, `price`, `original_price`, `enable`, `duration`, `currency`, `description`, `weight`, `created_time`, `updated_time`, `type`, `unit`, `total_count`, `duration_type_id`) VALUES (1, '测试商品一', 10000.00, 1000000.00, 1, 100000, 'RMB', '尝鲜入门级\\n不限制请求数量\\n自定义IP时效\\n高带宽配速\\n可指定国家/城市/运营商\\n30天时效\\n支持HTTP(S)/Socks5协议', 1, 0, 0, 1, 1, 1, 0);
INSERT INTO `ip_commodities` (`id`, `name`, `price`, `original_price`, `enable`, `duration`, `currency`, `description`, `weight`, `created_time`, `updated_time`, `type`, `unit`, `total_count`, `duration_type_id`) VALUES (2, '测试商品二', 10000.00, 1000000.00, 1, 100000, 'RMB', '尝鲜入门级\\n不限制请求数量\\n自定义IP时效\\n高带宽配速\\n可指定国家/城市/运营商\\n30天时效\\n支持HTTP(S)/Socks5协议', 1, 0, 0, 1, 1, 10, 0);
INSERT INTO `ip_commodities` (`id`, `name`, `price`, `original_price`, `enable`, `duration`, `currency`, `description`, `weight`, `created_time`, `updated_time`, `type`, `unit`, `total_count`, `duration_type_id`) VALUES (3, '测试商品三', 10000.00, 1000000.00, 1, 100000, 'RMB', '尝鲜入门级\\n不限制请求数量\\n自定义IP时效\\n高带宽配速\\n可指定国家/城市/运营商\\n30天时效\\n支持HTTP(S)/Socks5协议', 1, 0, 0, 1, 1, 100, 0);
INSERT INTO `ip_commodities` (`id`, `name`, `price`, `original_price`, `enable`, `duration`, `currency`, `description`, `weight`, `created_time`, `updated_time`, `type`, `unit`, `total_count`, `duration_type_id`) VALUES (4, '静态住宅测试商品一', 10000.00, 1000000.00, 1, 100000, 'RMB', '尝鲜入门级\\n不限制请求数量\\n自定义IP时效\\n高带宽配速\\n可指定国家/城市/运营商\\n30天时效\\n支持HTTP(S)/Socks5协议', 1, 0, 0, 2, 1, 100, 1);
INSERT INTO `ip_commodities` (`id`, `name`, `price`, `original_price`, `enable`, `duration`, `currency`, `description`, `weight`, `created_time`, `updated_time`, `type`, `unit`, `total_count`, `duration_type_id`) VALUES (5, '静态住宅测试商品一', 10000.00, 1000000.00, 1, 100000, 'RMB', '尝鲜入门级\\n不限制请求数量\\n自定义IP时效\\n高带宽配速\\n可指定国家/城市/运营商\\n30天时效\\n支持HTTP(S)/Socks5协议', 1, 0, 0, 2, 1, 100, 2);
INSERT INTO `ip_commodities` (`id`, `name`, `price`, `original_price`, `enable`, `duration`, `currency`, `description`, `weight`, `created_time`, `updated_time`, `type`, `unit`, `total_count`, `duration_type_id`) VALUES (6, '静态住宅测试商品一', 10000.00, 1000000.00, 1, 100000, 'RMB', '尝鲜入门级\\n不限制请求数量\\n自定义IP时效\\n高带宽配速\\n可指定国家/城市/运营商\\n30天时效\\n支持HTTP(S)/Socks5协议', 1, 0, 0, 2, 1, 100, 3);
COMMIT;

-- ----------------------------
-- Table structure for ip_duration_types
-- ----------------------------
DROP TABLE IF EXISTS `ip_duration_types`;
CREATE TABLE `ip_duration_types` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '时长类型ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '时长类型名称',
  `type` tinyint NOT NULL COMMENT '时长类型',
  `count` bigint NOT NULL COMMENT '时长数量',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
  `weight` tinyint NOT NULL DEFAULT '0' COMMENT '权重',
  `create_time` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` bigint NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of ip_duration_types
-- ----------------------------
BEGIN;
INSERT INTO `ip_duration_types` (`id`, `name`, `type`, `count`, `status`, `weight`, `create_time`, `update_time`) VALUES (1, '单月套餐', 5, 1, 2, 0, 0, 0);
INSERT INTO `ip_duration_types` (`id`, `name`, `type`, `count`, `status`, `weight`, `create_time`, `update_time`) VALUES (2, '双月套餐', 5, 2, 2, 1, 0, 0);
INSERT INTO `ip_duration_types` (`id`, `name`, `type`, `count`, `status`, `weight`, `create_time`, `update_time`) VALUES (3, '季度套餐', 5, 3, 2, 2, 0, 0);
COMMIT;

-- ----------------------------
-- Table structure for ip_traffic_country
-- ----------------------------
DROP TABLE IF EXISTS `ip_traffic_country`;
CREATE TABLE `ip_traffic_country` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '账号ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '国家名称',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '国家代码',
  `is_hot` tinyint NOT NULL DEFAULT '0' COMMENT '是否热门国家',
  `is_recommend` tinyint NOT NULL DEFAULT '0' COMMENT '是否为推荐国家',
  `is_enabled` tinyint NOT NULL DEFAULT '0' COMMENT '是否启用',
  `weight` tinyint NOT NULL DEFAULT '0' COMMENT '权重',
  `stock_quantity` bigint NOT NULL DEFAULT '0' COMMENT 'IP库存数量',
  `create_time` bigint NOT NULL COMMENT '创建时间',
  `update_time` bigint NOT NULL COMMENT '更新时间',
  `region_id` bigint NOT NULL COMMENT '区域ID',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of ip_traffic_country
-- ----------------------------
BEGIN;
INSERT INTO `ip_traffic_country` (`id`, `name`, `code`, `is_hot`, `is_recommend`, `is_enabled`, `weight`, `stock_quantity`, `create_time`, `update_time`, `region_id`) VALUES (1, '中国', 'CN', 1, 1, 1, 1, 100, 0, 0, 1);
INSERT INTO `ip_traffic_country` (`id`, `name`, `code`, `is_hot`, `is_recommend`, `is_enabled`, `weight`, `stock_quantity`, `create_time`, `update_time`, `region_id`) VALUES (2, '美国', 'US', 0, 0, 1, 2, 300, 0, 0, 1);
INSERT INTO `ip_traffic_country` (`id`, `name`, `code`, `is_hot`, `is_recommend`, `is_enabled`, `weight`, `stock_quantity`, `create_time`, `update_time`, `region_id`) VALUES (3, '美国', 'US', 0, 0, 1, 2, 300, 0, 0, 2);
INSERT INTO `ip_traffic_country` (`id`, `name`, `code`, `is_hot`, `is_recommend`, `is_enabled`, `weight`, `stock_quantity`, `create_time`, `update_time`, `region_id`) VALUES (4, '德国', 'DE', 0, 0, 1, 2, 300, 0, 0, 2);
COMMIT;

-- ----------------------------
-- Table structure for ip_traffic_country_commodites
-- ----------------------------
DROP TABLE IF EXISTS `ip_traffic_country_commodites`;
CREATE TABLE `ip_traffic_country_commodites` (
  `commodites_id` bigint NOT NULL COMMENT '国家商品ID',
  `country_id` bigint NOT NULL COMMENT '国家ID',
  KEY `idx_ip_traffic_country_commodites_commodites_id` (`commodites_id`) USING BTREE,
  KEY `idx_ip_traffic_country_commodites_country_id` (`country_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of ip_traffic_country_commodites
-- ----------------------------
BEGIN;
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (4, 2);
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (5, 1);
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (4, 1);
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (6, 2);
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (6, 1);
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (5, 2);
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (4, 3);
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (5, 3);
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (6, 3);
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (4, 4);
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (5, 4);
INSERT INTO `ip_traffic_country_commodites` (`commodites_id`, `country_id`) VALUES (6, 4);
COMMIT;

-- ----------------------------
-- Table structure for ip_traffic_region
-- ----------------------------
DROP TABLE IF EXISTS `ip_traffic_region`;
CREATE TABLE `ip_traffic_region` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '区域ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '区域名称',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '区域代码',
  `is_hot` tinyint NOT NULL DEFAULT '0' COMMENT '是否热门区域',
  `is_recommend` tinyint NOT NULL DEFAULT '0' COMMENT '是否为推荐区域',
  `is_enabled` tinyint NOT NULL DEFAULT '0' COMMENT '是否启用',
  `weight` tinyint NOT NULL DEFAULT '0' COMMENT '权重',
  `created_time` bigint NOT NULL COMMENT '创建时间',
  `updated_time` bigint NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of ip_traffic_region
-- ----------------------------
BEGIN;
INSERT INTO `ip_traffic_region` (`id`, `name`, `code`, `is_hot`, `is_recommend`, `is_enabled`, `weight`, `created_time`, `updated_time`) VALUES (1, '亚太', 'COD1', 0, 0, 1, 0, 0, 0);
INSERT INTO `ip_traffic_region` (`id`, `name`, `code`, `is_hot`, `is_recommend`, `is_enabled`, `weight`, `created_time`, `updated_time`) VALUES (2, '欧美', 'COD2', 0, 0, 1, 0, 0, 0);
COMMIT;

-- ----------------------------
-- Table structure for ip_user
-- ----------------------------
DROP TABLE IF EXISTS `ip_user`;
CREATE TABLE `ip_user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '账号ID',
  `create_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '注册时间',
  `update_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '会员信息上次更新时间',
  `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '账号状态',
  `money` double(10,2) NOT NULL DEFAULT '0.00' COMMENT '账号余额',
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '账号名称',
  `email` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `phone` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `password` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `salt` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码盐',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_ip_user_email` (`email`) USING BTREE,
  KEY `idx_ip_user_phone` (`phone`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of ip_user
-- ----------------------------
BEGIN;
INSERT INTO `ip_user` (`id`, `create_time`, `update_time`, `status`, `money`, `name`, `email`, `phone`, `password`, `salt`) VALUES (1, 1725877598, 1725877598, 1, 0.00, '', 'test@test.com', '', '66d7fbeb4907f525df686d99b53cc475', '77243ac5adda');
COMMIT;

-- ----------------------------
-- Table structure for ip_verification_codes
-- ----------------------------
DROP TABLE IF EXISTS `ip_verification_codes`;
CREATE TABLE `ip_verification_codes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `email` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `ip` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '',
  `created_time` bigint DEFAULT NULL,
  `expires_time` bigint DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_ip_verification_codes_email` (`email`) USING BTREE,
  KEY `idx_ip_verification_codes_ip` (`ip`) USING BTREE,
  KEY `idx_ip_verification_codes_expires_time` (`expires_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of ip_verification_codes
-- ----------------------------
BEGIN;
INSERT INTO `ip_verification_codes` (`id`, `code`, `email`, `ip`, `created_time`, `expires_time`) VALUES (1, 'B5V61Q', 'test@test.com', '192.168.1.199', 1725877030, 0);
INSERT INTO `ip_verification_codes` (`id`, `code`, `email`, `ip`, `created_time`, `expires_time`) VALUES (2, 'BFLMAG', 'test@test.com', '192.168.1.199', 1725877039, 0);
INSERT INTO `ip_verification_codes` (`id`, `code`, `email`, `ip`, `created_time`, `expires_time`) VALUES (3, '47C3Y8', 'test@test.com', '192.168.1.199', 1725877045, 0);
INSERT INTO `ip_verification_codes` (`id`, `code`, `email`, `ip`, `created_time`, `expires_time`) VALUES (4, 'L9W6TG', 'test@test.com', '192.168.1.199', 1725877096, 0);
INSERT INTO `ip_verification_codes` (`id`, `code`, `email`, `ip`, `created_time`, `expires_time`) VALUES (5, '52HPGU', 'test@test.com', '192.168.1.199', 1725877493, 0);
INSERT INTO `ip_verification_codes` (`id`, `code`, `email`, `ip`, `created_time`, `expires_time`) VALUES (6, 'GPHTQS', 'test@test.com', '192.168.1.199', 1725877579, 0);
INSERT INTO `ip_verification_codes` (`id`, `code`, `email`, `ip`, `created_time`, `expires_time`) VALUES (7, 'SQDC4E', 'test@test.com', '192.168.1.199', 1725881482, 0);
INSERT INTO `ip_verification_codes` (`id`, `code`, `email`, `ip`, `created_time`, `expires_time`) VALUES (8, '9BI0G4', 'test@test.com', '192.168.1.199', 1726039180, 0);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
