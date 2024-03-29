/*
 Navicat Premium Data Transfer

 Source Server         : veweiyi.cn-mysql8.0
 Source Server Type    : MySQL
 Source Server Version : 80034
 Source Host           : veweiyi.cn:3306
 Source Schema         : blog-veweiyi

 Target Server Type    : MySQL
 Target Server Version : 80034
 File Encoding         : 65001

 Date: 23/02/2024 15:21:47
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for api
-- ----------------------------
DROP TABLE IF EXISTS `api`;
CREATE TABLE `api` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'api名称',
  `path` varchar(128) NOT NULL DEFAULT '' COMMENT 'api路径',
  `method` varchar(16) NOT NULL DEFAULT '' COMMENT 'api请求方法',
  `parent_id` int NOT NULL DEFAULT '0' COMMENT '分组id',
  `traceable` tinyint NOT NULL DEFAULT '0' COMMENT '是否追溯操作记录 0需要，1是',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1开，2关',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_path_method` (`path`,`method`,`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=166 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='api路由';

-- ----------------------------
-- Records of api
-- ----------------------------
BEGIN;
INSERT INTO `api` VALUES (1, 'AI', '', '', 0, 0, 1, '2024-01-18 19:51:35', '2024-01-18 19:51:35');
INSERT INTO `api` VALUES (2, 'Api', '', '', 0, 0, 1, '2024-01-18 19:51:35', '2024-01-18 19:51:35');
INSERT INTO `api` VALUES (3, 'Article', '', '', 0, 0, 1, '2024-01-18 19:51:36', '2024-01-18 19:51:36');
INSERT INTO `api` VALUES (4, 'Auth', '', '', 0, 0, 1, '2024-01-18 19:51:36', '2024-01-18 19:51:36');
INSERT INTO `api` VALUES (5, 'Captcha', '', '', 0, 0, 1, '2024-01-18 19:51:36', '2024-01-18 19:51:36');
INSERT INTO `api` VALUES (6, 'Category', '', '', 0, 0, 1, '2024-01-18 19:51:36', '2024-01-18 19:51:36');
INSERT INTO `api` VALUES (7, 'Comment', '', '', 0, 0, 1, '2024-01-18 19:51:37', '2024-01-18 19:51:37');
INSERT INTO `api` VALUES (8, 'FriendLink', '', '', 0, 0, 1, '2024-01-18 19:51:37', '2024-01-18 19:51:37');
INSERT INTO `api` VALUES (9, 'Menu', '', '', 0, 0, 1, '2024-01-18 19:51:37', '2024-01-18 19:51:37');
INSERT INTO `api` VALUES (10, 'OperationLog', '', '', 0, 0, 1, '2024-01-18 19:51:38', '2024-01-18 19:51:38');
INSERT INTO `api` VALUES (11, 'Page', '', '', 0, 0, 1, '2024-01-18 19:51:38', '2024-01-18 19:51:38');
INSERT INTO `api` VALUES (12, 'Photo', '', '', 0, 0, 1, '2024-01-18 19:51:38', '2024-01-18 19:51:38');
INSERT INTO `api` VALUES (13, 'PhotoAlbum', '', '', 0, 0, 1, '2024-01-18 19:51:38', '2024-01-18 19:51:38');
INSERT INTO `api` VALUES (14, 'Remark', '', '', 0, 0, 1, '2024-01-18 19:51:39', '2024-01-18 19:51:39');
INSERT INTO `api` VALUES (15, 'Role', '', '', 0, 0, 1, '2024-01-18 19:51:39', '2024-01-18 19:51:39');
INSERT INTO `api` VALUES (16, 'Tag', '', '', 0, 0, 1, '2024-01-18 19:51:39', '2024-01-18 19:51:39');
INSERT INTO `api` VALUES (17, 'Talk', '', '', 0, 0, 1, '2024-01-18 19:51:39', '2024-01-18 19:51:39');
INSERT INTO `api` VALUES (18, 'Upload', '', '', 0, 0, 1, '2024-01-18 19:51:40', '2024-01-18 19:51:40');
INSERT INTO `api` VALUES (19, 'User', '', '', 0, 0, 1, '2024-01-18 19:51:40', '2024-01-18 19:51:40');
INSERT INTO `api` VALUES (20, 'Website', '', '', 0, 0, 1, '2024-01-18 19:51:40', '2024-01-18 19:51:40');
INSERT INTO `api` VALUES (21, 'Websocket', '', '', 0, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (22, '和Chatgpt聊天', '/api/v1/ai/chat', 'POST', 1, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (23, 'Chatgpt扮演角色', '/api/v1/ai/cos', 'POST', 1, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (24, '创建接口', '/api/v1/api', 'POST', 2, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (25, '更新接口', '/api/v1/api', 'PUT', 2, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (26, '批量删除接口', '/api/v1/api/batch_delete', 'DELETE', 2, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (27, '获取api列表', '/api/v1/api/details_list', 'POST', 2, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (28, '分页获取接口列表', '/api/v1/api/list', 'POST', 2, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (29, '同步api列表', '/api/v1/api/sync', 'POST', 2, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (30, '删除接口', '/api/v1/api/{id}', 'DELETE', 2, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (31, '查询接口', '/api/v1/api/{id}', 'GET', 2, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (32, '保存文章', '/api/v1/admin/article', 'POST', 3, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (33, '删除文章-逻辑删除', '/api/v1/admin/article/delete', 'PUT', 3, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (34, '分页获取文章列表', '/api/v1/admin/article/list', 'POST', 3, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (35, '更新文章', '/api/v1/admin/article/top', 'PUT', 3, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (36, '删除文章', '/api/v1/admin/article/{id}', 'DELETE', 3, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (37, '查询文章', '/api/v1/admin/article/{id}', 'GET', 3, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (38, '文章归档(时间轴)', '/api/v1/article/archives', 'POST', 3, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (39, '分页获取文章列表', '/api/v1/article/list', 'POST', 3, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (40, '通过标签或者id获取文章列表', '/api/v1/article/series', 'POST', 3, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (41, '文章相关推荐', '/api/v1/article/{id}/details', 'GET', 3, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (42, '点赞文章', '/api/v1/article/{id}/like', 'PUT', 3, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (43, '发送忘记密码邮件', '/api/v1/forget/password', 'POST', 4, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (44, '重置密码', '/api/v1/forget/reset_password', 'POST', 4, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (45, '登录', '/api/v1/login', 'POST', 4, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (46, '注销', '/api/v1/logoff', 'POST', 4, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (47, '登出', '/api/v1/logout', 'GET', 4, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (48, '获取授权地址', '/api/v1/oauth/login', 'POST', 4, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (49, '获取授权地址', '/api/v1/oauth/url', 'POST', 4, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (50, '注册', '/api/v1/register', 'POST', 4, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (51, '发送注册邮件', '/api/v1/register/email', 'POST', 4, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (52, '发送验证码', '/api/v1/captcha/email', 'POST', 5, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (53, '生成验证码', '/api/v1/captcha/image', 'POST', 5, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (54, '检验验证码', '/api/v1/captcha/verify', 'POST', 5, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (55, '创建文章分类', '/api/v1/category', 'POST', 6, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (56, '更新文章分类', '/api/v1/category', 'PUT', 6, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (57, '批量删除文章分类', '/api/v1/category/batch_delete', 'DELETE', 6, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (58, '分页获取文章分类详情列表', '/api/v1/category/details_list', 'POST', 6, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (59, '分页获取文章分类列表', '/api/v1/category/list', 'POST', 6, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (60, '删除文章分类', '/api/v1/category/{id}', 'DELETE', 6, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (61, '查询文章分类', '/api/v1/category/{id}', 'GET', 6, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (62, '创建评论', '/api/v1/comment', 'POST', 7, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (63, '更新评论', '/api/v1/comment', 'PUT', 7, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (64, '批量删除评论', '/api/v1/comment/batch_delete', 'DELETE', 7, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (65, '分页获取评论列表', '/api/v1/comment/details_list', 'POST', 7, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (66, '分页获取评论列表', '/api/v1/comment/list', 'POST', 7, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (67, '获取用户评论列表', '/api/v1/comment/list/back', 'POST', 7, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (68, '删除评论', '/api/v1/comment/{id}', 'DELETE', 7, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (69, '查询评论', '/api/v1/comment/{id}', 'GET', 7, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (70, '点赞评论', '/api/v1/comment/{id}/like', 'POST', 7, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (71, '查询评论回复列表', '/api/v1/comment/{id}/reply_list', 'POST', 7, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (72, '创建友链', '/api/v1/friend_link', 'POST', 8, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (73, '更新友链', '/api/v1/friend_link', 'PUT', 8, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (74, '批量删除友链', '/api/v1/friend_link/batch_delete', 'DELETE', 8, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (75, '分页获取友链列表', '/api/v1/friend_link/list', 'POST', 8, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (76, '删除友链', '/api/v1/friend_link/{id}', 'DELETE', 8, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (77, '查询友链', '/api/v1/friend_link/{id}', 'GET', 8, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (78, '创建菜单', '/api/v1/menu', 'POST', 9, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (79, '更新菜单', '/api/v1/menu', 'PUT', 9, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (80, '批量删除菜单', '/api/v1/menu/batch_delete', 'DELETE', 9, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (81, '获取菜单列表', '/api/v1/menu/details_list', 'POST', 9, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (82, '分页获取菜单列表', '/api/v1/menu/list', 'POST', 9, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (83, '删除菜单', '/api/v1/menu/{id}', 'DELETE', 9, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (84, '查询菜单', '/api/v1/menu/{id}', 'GET', 9, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (85, '创建操作记录', '/api/v1/operation_log', 'POST', 10, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (86, '更新操作记录', '/api/v1/operation_log', 'PUT', 10, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (87, '批量删除操作记录', '/api/v1/operation_log/batch_delete', 'DELETE', 10, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (88, '分页获取操作记录列表', '/api/v1/operation_log/list', 'POST', 10, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (89, '删除操作记录', '/api/v1/operation_log/{id}', 'DELETE', 10, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (90, '查询操作记录', '/api/v1/operation_log/{id}', 'GET', 10, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (91, '创建页面', '/api/v1/page', 'POST', 11, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (92, '更新页面', '/api/v1/page', 'PUT', 11, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (93, '批量删除页面', '/api/v1/page/batch_delete', 'DELETE', 11, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (94, '分页获取页面列表', '/api/v1/page/list', 'POST', 11, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (95, '删除页面', '/api/v1/page/{id}', 'DELETE', 11, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (96, '查询页面', '/api/v1/page/{id}', 'GET', 11, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (97, '创建相片', '/api/v1/photo', 'POST', 12, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (98, '更新相片', '/api/v1/photo', 'PUT', 12, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (99, '批量删除相片', '/api/v1/photo/batch_delete', 'DELETE', 12, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (100, '分页获取相片列表', '/api/v1/photo/list', 'POST', 12, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (101, '删除相片', '/api/v1/photo/{id}', 'DELETE', 12, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (102, '查询相片', '/api/v1/photo/{id}', 'GET', 12, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (103, '创建相册', '/api/v1/photo_album', 'POST', 13, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (104, '更新相册', '/api/v1/photo_album', 'PUT', 13, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (105, '批量删除相册', '/api/v1/photo_album/batch_delete', 'DELETE', 13, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (106, '获取相册详情列表', '/api/v1/photo_album/details_list', 'POST', 13, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (107, '分页获取相册列表', '/api/v1/photo_album/list', 'POST', 13, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (108, '删除相册', '/api/v1/photo_album/{id}', 'DELETE', 13, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (109, '查询相册', '/api/v1/photo_album/{id}', 'GET', 13, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (110, '获取相册详情', '/api/v1/photo_album/{id}/details', 'GET', 13, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (111, '创建留言', '/api/v1/remark', 'POST', 14, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (112, '更新留言', '/api/v1/remark', 'PUT', 14, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (113, '批量删除留言', '/api/v1/remark/batch_delete', 'DELETE', 14, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (114, '分页获取留言列表', '/api/v1/remark/list', 'POST', 14, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (115, '删除留言', '/api/v1/remark/{id}', 'DELETE', 14, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (116, '查询留言', '/api/v1/remark/{id}', 'GET', 14, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (117, '创建角色', '/api/v1/role', 'POST', 15, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (118, '更新角色', '/api/v1/role', 'PUT', 15, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (119, '批量删除角色', '/api/v1/role/batch_delete', 'DELETE', 15, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (120, '获取角色列表', '/api/v1/role/details_list', 'POST', 15, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (121, '分页获取角色列表', '/api/v1/role/list', 'POST', 15, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (122, '更新角色菜单', '/api/v1/role/update_menus', 'POST', 15, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (123, '更新角色资源', '/api/v1/role/update_resources', 'POST', 15, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (124, '删除角色', '/api/v1/role/{id}', 'DELETE', 15, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (125, '查询角色', '/api/v1/role/{id}', 'GET', 15, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (126, '创建文章标签', '/api/v1/tag', 'POST', 16, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (127, '更新文章标签', '/api/v1/tag', 'PUT', 16, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (128, '批量删除文章标签', '/api/v1/tag/batch_delete', 'DELETE', 16, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (129, '分页获取文章分类详情列表', '/api/v1/tag/details_list', 'POST', 16, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (130, '分页获取文章标签列表', '/api/v1/tag/list', 'POST', 16, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (131, '删除文章标签', '/api/v1/tag/{id}', 'DELETE', 16, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (132, '查询文章标签', '/api/v1/tag/{id}', 'GET', 16, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (133, '创建说说', '/api/v1/talk', 'POST', 17, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (134, '更新说说', '/api/v1/talk', 'PUT', 17, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (135, '批量删除说说', '/api/v1/talk/batch_delete', 'DELETE', 17, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (136, '分页获取说说详情列表', '/api/v1/talk/details_list', 'POST', 17, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (137, '分页获取说说列表', '/api/v1/talk/list', 'POST', 17, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (138, '删除说说', '/api/v1/talk/{id}', 'DELETE', 17, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (139, '查询说说', '/api/v1/talk/{id}', 'GET', 17, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (140, '分页获取说说详情列表', '/api/v1/talk/{id}/details', 'GET', 17, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (141, '点赞说说', '/api/v1/talk/{id}/like', 'PUT', 17, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (142, '上传文件', '/api/v1/upload/{label}', 'POST', 18, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (143, '上传语言', '/api/v1/voice', 'POST', 18, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (144, '获取用户接口权限', '/api/v1/user/apis', 'GET', 19, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (145, '获取用户地区列表', '/api/v1/user/area_list', 'POST', 19, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (146, '更换用户头像', '/api/v1/user/avatar', 'POST', 19, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (147, '获取用户信息', '/api/v1/user/info', 'GET', 19, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (148, '修改用户信息', '/api/v1/user/info', 'POST', 19, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (149, '获取用户列表', '/api/v1/user/list', 'POST', 19, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (150, '获取用户登录历史', '/api/v1/user/login_history', 'POST', 19, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (151, '批量删除登录历史', '/api/v1/user/login_history/batch_delete', 'DELETE', 19, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (152, '获取用户菜单权限', '/api/v1/user/menus', 'GET', 19, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (153, '获取在线用户列表', '/api/v1/user/online_list', 'POST', 19, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (154, '修改用户角色', '/api/v1/user/update_roles', 'POST', 19, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (155, '修改用户状态', '/api/v1/user/update_status', 'POST', 19, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (156, '获取博客前台首页信息', '/api/v1/', 'GET', 20, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (157, '关于我', '/api/v1/about/me', 'GET', 20, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (158, '获取后台首页信息', '/api/v1/admin', 'GET', 20, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (159, '更新我的信息', '/api/v1/admin/about/me', 'POST', 20, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (160, '获取配置', '/api/v1/admin/config', 'POST', 20, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (161, '更新配置', '/api/v1/admin/config', 'PUT', 20, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (162, '获取服务器信息', '/api/v1/admin/system/state', 'GET', 20, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (163, '查询聊天记录', '/api/v1/chat/records', 'POST', 20, 1, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (164, '获取网站配置', '/api/v1/website/config', 'GET', 20, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
INSERT INTO `api` VALUES (165, '查询聊天记录', '/api/v1/ws', 'GET', 21, 0, 1, '2024-01-18 19:51:41', '2024-01-18 19:51:41');
COMMIT;

-- ----------------------------
-- Table structure for article
-- ----------------------------
DROP TABLE IF EXISTS `article`;
CREATE TABLE `article` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '作者',
  `category_id` int NOT NULL DEFAULT '0' COMMENT '文章分类',
  `article_cover` varchar(1024) NOT NULL DEFAULT '' COMMENT '文章缩略图',
  `article_title` varchar(64) NOT NULL DEFAULT '' COMMENT '标题',
  `article_content` longtext NOT NULL COMMENT '内容',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT '文章类型 1原创 2转载 3翻译',
  `original_url` varchar(255) NOT NULL DEFAULT '' COMMENT '原文链接',
  `is_top` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否置顶 0否 1是',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除  0否 1是',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态值 1公开 2私密 3评论可见',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=84 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='文章';

-- ----------------------------
-- Records of article
-- ----------------------------
BEGIN;
INSERT INTO `article` VALUES (54, 0, 187, 'https://veport.oss-cn-beijing.aliyuncs.com/articles/3a4b4e40fb8aa5fcc016f0228938d321.jpg', '测试文章', '恭喜你成功运行博客！\n', 1, '', 0, 0, 0, '2022-01-18 00:29:02', '2023-11-08 20:28:31');
INSERT INTO `article` VALUES (56, 0, 188, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/c9787ff560eb5d94e5ce0e1dd2cecd12.jpg', '2022-01-20', '第一次接触这个网站，在手机上看见主页面很好看，会滚动不一样的诗句，配合背景很有意境。', 1, '', 0, 0, 0, '2022-01-20 21:40:29', '2023-11-02 20:13:33');
INSERT INTO `article` VALUES (59, 0, 188, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/8b03884995623eab1a76772f23b58875.jpg', '网站搭建过程中的一些感想', '\n&emsp;&emsp;前言：感觉自己的大学过得挺混的，这与学校的教学方式有很大的关系。我们学校比较注重课程基础，教了很多其实没必要学的东西(没错，模电、数电说的就是你)，而真正对程序员核心的Java语言到大三上学期才开始学，而且还是阉割版的。我以为这样的教学方式，是学院为了培养研究生而制定的。真正想要毕业后就找工作的，都是自己学的。\n\n&emsp;&emsp;大一大二挺忙的，基本每天都是在为微积分、概率论、模电、数电(还要提一下)烦恼。大三了之后课开始变少了，开始真正的接触到互联网行业。\n\n&emsp;&emsp;首先是学会的是写Android应用软件，并学会了kotlin语言。起因是三四月份的时候，我当时想准备考研，所以打算暑假在实验室混。实验室的学长问我会不会写Android应用，我说会(其实我不会)。于是当天下午就去书店买了一本《Android 第一行代码》，利用课余学了三个月以后，总算可以自己独立写自己的应用了。并且熟悉了kotlin语言。之后一段时间里面试都是在和面试官聊Android。 \n\n&emsp;&emsp;写Android应用的时候，一直都是用别人的api，当时特别希望可以自定义自己的api。但奈于对后端不熟悉，所以一直也没实现。\n\n&emsp;&emsp;去年10月底，学院开始选毕设题目。我眼疾手快抢了一个与Android有关的题目，题目的要求中需要把信息加密保存起来。我想“如果用手机本地保存的方式实现，那么这应用也太low了吧。我不能用这么low的方式完成，我要把数据存在云端。”于是，我下定决心开始学习后端技术。\n\n&emsp;&emsp;刚学习springboot，买了一本书。但是发现看不懂，因为里面全是使用@注解的形式定义类，而且作者常常跳过一些步骤(可能他认为不重要，但对没有基础的来说看着看着就一头雾水)。后来在b站上看了一些spring基础的视频，才开始入门。第一个项目是乐字节的云e办，其中security+jwt权限认证框架比较难理解，但也不得不说这是一个很好的入门项目，视频讲解也很详细。\n\n&emsp;&emsp;云e办的后端学了三周，每天有8个小时以上都是在敲代码。学成之后，自己也写了毕设的接口，基本上已经完成功能，还有待优化和测试。这个时候，我其实已经通过移动端+后端完成毕设了。\n\n&emsp;&emsp;本来不打算学前端的。直到在github上发现了风丶神大佬写的vue+springboot博客页面，当时就觉得“我靠！这玩意儿好牛，我一定要拥有！”。\n\n&emsp;&emsp;但是这个博客只有源码，没有说明，下载了发现根本看不懂vue项目。那怎么办，还是得学呗，先学vue。于是想到了有教程视频也有文档说明的云e办，只要学会了云e办的vue前端，那么起码也能看得懂源码了。vue需要前端基础，于是又去b站跟着乐字节的念安老师学了html+css+JavaScript。学会了之后跟着视频学云e办前端。最后终于在云端部署了自己的个人网站。\n\n&emsp;&emsp;关于网站搭建，之前物联网课设的时候也有用阿里云的ECS服务器。只会便遭到了阿里云腾讯云百度云各种云的电话骚扰，向我推销云服务器。终于禁不住买了腾讯云的一个域名(10元/年)和轻量应用服务器(74元/年)。\n\n&emsp;&emsp;云服务器买了之后，还要备案域名，申请SSL证书，DNS域名解析到服务器，开通的阿里云的oss用于存储上传的图片，部署Java、tomcat、mysql、redis、rabbitmq、nginx……有时间再总结一下遇到的坑。\n\n&emsp;&emsp;这段时间很累，收获也很多。其实一个人有了目标之后，学习的潜能和动力会有很大的提升。在在求知的途中，我也接触到了很多自己未知的领域，也明白自己所知甚浅。今后的光阴，也希望自己少打游戏，多学习，不要虚度光阴。\n&emsp;&emsp;路漫漫其修远兮，吾将上下而求索。\n::: hljs-center\n\n![1897CFE31F692AA278E02620E8021357.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/1897cfe31f692aa278e02620e8021357.png)\n\n:::\n', 1, '', 0, 0, 0, '2022-01-21 12:04:02', '2023-11-02 20:13:09');
INSERT INTO `article` VALUES (60, 1, 191, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/cfeb11ab6be04ca78f24a0d8974a296d.png', '博客技术总结', '## 1.技术介绍\n\n前端：vue + vuex + vue-router + axios + vuetify + element + echarts\n\n后端：SpringBoot + nginx + docker + SpringSecurity + Swagger2 + MyBatisPlus + Mysql + Redis + elasticsearch + rabbitMQ + MaxWell\n\n其他：接入QQ，微博第三方登录\n\n## 2.运行环境\n\n开发工具：IDEA\n\n服务器：腾讯云2核4G CentOS 8.2 64bit\n\n对象存储：阿里云OSS\n\n## 3.功能展示\n### 移动端前台\n\n<img src=\" https://veport.oss-cn-beijing.aliyuncs.com/articles/3d6d442430f9de01c2e7cbcb914c1a2a.jpg\" alt=\"1754BCE1444B06E890FFAC61A1BF9BD5.jpg\" style=\"zoom:20%;\" />\n&emsp;\n<img src=\" https://veport.oss-cn-beijing.aliyuncs.com/articles/12431b4262b7e9c2fe1c10d2a9797ad1.jpg\" alt=\"12431B4262B7E9C2FE1C10D2A9797AD1.jpg\" style=\"zoom:20%;\" />\n\n\n### PC端前台\n![截屏20220219 22.54.49.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/f5e11e5a67fd322be0487d1b520ec9a4.png)\n\n![截屏20220219 23.06.58.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/f36f002e31660912dcd76d34dcca538b.png)\n\n![截屏20220219 23.10.43.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/eac9f4a203d6377efa37edb912166e60.png)\n\n### PC端后台\n![admin1.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/6ec21b7f32199a0c1418a9968d0b44f6.png)\n![admin2.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/0e7fb892d411d34da39c3d4ab1142ed6.png)\n\n## 4.项目总结\n\n感谢风丶神提供的项目源码，项目源码地址[GitHub](https://github.com/X1192176811/blog)。\n感兴趣的可以去点个star，搭建一个自己的博客吧。!\n\n本博客前台地址：[https://ve77.cn/blog](https://ve77.cn/blog)\n本博客后台地址：[https://ve77.cn/admin](https://ve77.cn/admin)\n项目源码地址：[https://github.com/7914-ve/ve-blog](https://github.com/7914-ve/ve-blog)\nlanguage', 1, '', 1, 0, 1, '2022-01-21 12:21:31', '2022-05-01 22:47:31');
INSERT INTO `article` VALUES (62, 2, 188, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/048dee462bb06ee4de2bf5907fb44c31.jpg', '算法-刷题日记', '春招临近，该准备刷算法题了。 \n \n2022-1-22   今日写了8道算法题，都是些字符串处理基础题。\n\n习惯了用本地编译器，使用在线编译器时常常出现各种语法错误。\n\n今天写了一下代码，让大家见笑了。\n```Java\n1.\npublic class Main(){ \n}\n2.\nint[] arr=new int[](n);\n3.\nchart c=str.chartAt(i);\n4.\n if(c=>\'a\'&&c<=\'f\')\n     count+=(9+c-\'a\');\n5.\nSystem.out.print(i+‘ ’);\n```\n\n常用的双指针遍历数组去重的算法：\n```java\n     /**\n     * 遍历数组并去掉其中重复的元素，相同元素保存一个\n     * @param arr 排序前的数组\n     * @return 去重且排序后的数组\n     */\n    public static ArrayList simpleArray(int[] arr){\n        Arrays.sort(arr);\n        ArrayList res=new ArrayList();\n        int left=0;\n        for(int i=0;i<arr.length;i++){\n            // System.out.println(\"--->\"+arr[left]);\n\n            //左指针指向等待加入的元素下标，当遇到与比他大的元素时，加入数组\n            if(arr[left]!=arr[i]){\n                res.add(arr[left]);\n                System.out.println(arr[left]);\n                left=i;\n            }\n        }\n        // 最大的元素没有比较元素，需要加入数组\n        System.out.println(arr[left]);\n        res.add(arr[left]);\n        return res;\n    }\n```\n2022-1-23   今天刷了7题，依然在牛客网上刷华为的机试题。\n\n一直以为String类中有reverse()函数，结果编译器报错了。reverse()是StringBuffer的函数，下次一定注意。\n\n做了一题字符串字典排序，自己用TreeMap做的。做完看了看题解，原来字符串也可以用Arrays.sort()；\n\n2022-1-24 今天偷懒。', 1, '', 0, 0, 2, '2022-01-22 12:59:49', '2022-02-11 12:15:46');
INSERT INTO `article` VALUES (64, 0, 189, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/07afc2963f27e63239e50bc65bed6a6f.jpg', 'Tomcat 10 配置默认访问端口为443(使用https访问)', '\n# 前言\ntomcat配置好了以后默认是使用8080端口访问的，也就是需要在使用\"域名.com:8080\"才能访问。这篇总结一下如何修改tomcat配置，使可以用\"http://域名.com\"或\"https://域名.com\" 访问。\n## 前期准备\n环境配置：\n 1. 腾讯云轻量应用服务器: CentOS 8.2 64bit\n 2. 远程访问推荐使用图形化界面(Mac 建议Royal TSX,Windows建议Mobaxterm)\n 3. Tomcat 10.0.4 ；\n 4. Java 1.8 ；\n \n 前提条件:\n\n 配置访问80端口即\"http://域名.com\"不需要证书\n 配置访问443端口即\"https://域名.com\" 需要SSL证书，证书可以从你购买服务器的运营商那里获取\n\n## 具体操作步骤\n话不多说，直接进入正题\n编辑在 /usr/tomcat/*/conf 目录(这个目录是你安装tomcat的目录)下的 server.xml 文件。添加如下内容：\n```xml\n// An highlighted block\n<Connector port=\"443\" protocol=\"HTTP/1.1\" SSLEnabled=\"true\"\n  maxThreads=\"150\" scheme=\"https\" secure=\"true\"\n#证书保存的路径\n  keystoreFile=\"/usr/*/conf/域名.com.jks\" \n#密钥库密码\n  keystorePass=\"******\"\n  clientAuth=\"false\"/>\n```\n详细 server.xml 文件和一些参数解释如下(可以直接复制过去)：\n```xml\n<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!--\nServer 根元素，创建⼀个Server实例，⼦标签有 Listener、GlobalNamingResources、Service\nport：关闭服务器的监听端⼝\nshutdown：关闭服务器的指令字符串\n-->\n<Server port=\"8005\" shutdown=\"SHUTDOWN\">\n\n    <!-- 创建 5 个监听器  start -->\n    <!-- 以⽇志形式输出服务器 、操作系统、JVM的版本信息 -->\n    <Listener className=\"org.apache.catalina.startup.VersionLoggerListener\"/>\n    <!-- 加载（服务器启动） 和 销毁 （服务器停⽌） APR。 如果找不到APR库， 则会输出⽇志， 并不影响 Tomcat启动 -->\n    <Listener className=\"org.apache.catalina.core.AprLifecycleListener\" SSLEngine=\"on\"/>\n    <!-- 避免JRE内存泄漏问题 -->\n    <Listener className=\"org.apache.catalina.core.JreMemoryLeakPreventionListener\"/>\n    <!-- 加载（服务器启动） 和 销毁（服务器停⽌） 全局命名服务 -->\n    <Listener className=\"org.apache.catalina.mbeans.GlobalResourcesLifecycleListener\"/>\n    <!-- 在Context停⽌时重建 Executor 池中的线程， 以避免ThreadLocal 相关的内存泄漏 -->\n    <Listener className=\"org.apache.catalina.core.ThreadLocalLeakPreventionListener\"/>\n    <!-- 创建 5 个监听器  end -->\n\n\n    <!--\n         定义服务器全局的JNDI 资源 命名服务\n    -->\n    <GlobalNamingResources>\n        <Resource name=\"UserDatabase\" auth=\"Container\"\n                  type=\"org.apache.catalina.UserDatabase\"\n                  description=\"User database that can be updated and saved\"\n                  factory=\"org.apache.catalina.users.MemoryUserDatabaseFactory\"\n                  pathname=\"conf/tomcat-users.xml\"/>\n    </GlobalNamingResources>\n\n    <!--\n            该标签⽤于创建 Service 实例，默认使⽤ org.apache.catalina.core.StandardService。\n       默认情况下，Tomcat 仅指定了Service 的名称， 值为 \"Catalina\"。\n       Service ⼦标签为 ： Listener、Executor、Connector、Engine，\n       其中：\n       Listener ⽤于为Service添加⽣命周期监听器，\n       Executor ⽤于配置Service 共享线程池，(可以给多个 Connector连接器使用)\n       Connector ⽤于配置Service 包含的链接器，\n       Engine ⽤于配置Service中链接器对应的Servlet 容器引擎\n     -->\n    <Service name=\"Catalina\">\n\n        <!-- 默认情况下，Service 并未添加共享线程池配置。 如果我们想添加⼀个线程池， 可以在<Executor> 下添加如下配置：\n              name：线程池名称，⽤于 Connector中指定\n              namePrefix：所创建的每个线程的名称前缀，⼀个单独的线程名称为：namePrefix+线程编号\n              maxThreads：池中最⼤线程数\n              minSpareThreads：活跃线程数，也就是核⼼池线程数，这些线程不会被销毁，会⼀直存在\n              maxIdleTime：线程空闲时间，超过该时间后，空闲线程会被销毁，默认值为6000（1分钟），单位毫秒\n              maxQueueSize：在被执⾏前最⼤线程排队数⽬，默认为Int的最⼤值，也就是⼴义的⽆限。除⾮特殊情况，这个值 不需要更改，否则会有请求不会被处理的情况发⽣\n              prestartminSpareThreads：启动线程池时是否启动 minSpareThreads部分线程。默认值为false，即不启动\n              threadPriority：线程池中线程优先级，默认值为5，值从1到10\n              className：线程池实现类，未指定情况下，默认实现类为\n              org.apache.catalina.core.StandardThreadExecutor。\n              如果想使⽤⾃定义线程池⾸先需要实现org.apache.catalina.Executor接⼝-->\n        <Executor name=\"tomcatThreadPool\"\n                  namePrefix=\"catalina-exec-\"\n                  maxThreads=\"200\"\n                  minSpareThreads=\"100\"\n                  maxIdleTime=\"60000\"\n                  maxQueueSize=\"Integer.MAX_VALUE\"\n                  prestartminSpareThreads=\"true\"\n                  threadPriority=\"5\"\n                  className=\"org.apache.catalina.core.StandardThreadExecutor\"/>\n\n        <!--\n           Connector 标签⽤于创建链接器实例，默认情况下，server.xml 配置了两个链接器，⼀个⽀持HTTP协议，⼀个⽀持AJP协议\n           ⼤多数情况下，我们并不需要新增链接器配置，只是根据需要对已有链接器进⾏优化\n                port：\n                     端⼝号，Connector ⽤于创建服务端Socket 并进⾏监听， 以等待客户端请求链接。如果该属性设置为0， Tomcat将会随机选择⼀个可⽤的端⼝号给当前Connector 使⽤\n                protocol：\n                     当前Connector ⽀持的访问协议。 默认为 HTTP/1.1 ， 并采⽤⾃动切换机制选择⼀个基于 JAVA NIO 的链接器或者基于本地APR的链接器（根据本地是否含有Tomcat的本地库判定）\n                connectionTimeOut:\n                     Connector 接收链接后的等待超时时间， 单位为 毫秒。 -1 表示不超时。\n                redirectPort：\n                     如果当前接收的是一个 https 的请求，那么tomcat 会将请求转发到 redirectPort指定的端口。\n                     比如现在设定的：8443 端口当前Connector 不⽀持SSL请求， 接收到了⼀个请求， 并且也符合security-constraint 约束，需要SSL传输，Catalina⾃动将请求重定向到指定的端⼝。\n                executor：\n                     指定共享线程池的名称， 也可以通过maxThreads、minSpareThreads 等属性配置内部线程池。\n                URIEncoding:\n                     ⽤于指定编码URI的字符编码， Tomcat8.x版本默认的编码为 UTF-8 , Tomcat7.x版本默认为ISO8859-1\n -->\n        <!--org.apache.coyote.http11.Http11NioProtocol， ⾮阻塞式 Java NIO 链接器，tomcat8配置nio会报错，可能是已经集成了nio的原因-->\n        <Connector port=\"80\"\n                   protocol=\"HTTP/1.1\"\n                   connectionTimeout=\"20000\"\n                   redirectPort=\"443\"\n                   executor=\"tomcatThreadPool\"\n                   URIEncoding=\"utf-8\"/>\n\n\n        <!-- certificateKeystoreFile 用于指定证书所在的目录 ；\n                        certificateKeystorePassword 用于指定证书的密码;type是使用的加密算法-->\n        <Connector port=\"443\" protocol=\"org.apache.coyote.http11.Http11NioProtocol\"\n                   maxThreads=\"150\" schema=\"https\" secure=\"true\" SSLEnabled=\"true\">\n            <SSLHostConfig>\n                <Certificate\n                        certificateKeystoreFile=\"conf/你的域名.cn.jks\"\n                        certificateKeystorePassword=\"你申请证书时提交密码\"\n                        type=\"RSA\" />\n            </SSLHostConfig>\n        </Connector>\n\n\n        <!-- Define an AJP 1.3 Connector on port 8009 -->\n\n        <Connector protocol=\"AJP/1.3\"\n                   address=\"::1\"\n                   port=\"8009\"\n                   redirectPort=\"443\" />\n\n\n        <!--name： ⽤于指定Engine 的名称， 默认为Catalina\n         defaultHost：默认使⽤的虚拟主机名称， 当客户端请求指向的主机⽆效时， 将交由默认的虚拟主机处\n              理， 默认为localhost-->\n        <Engine name=\"Catalina\" defaultHost=\"localhost\">\n            <Realm className=\"org.apache.catalina.realm.LockOutRealm\">\n                <Realm className=\"org.apache.catalina.realm.UserDatabaseRealm\"\n                       resourceName=\"UserDatabase\"/>\n            </Realm>\n\n            <!--Host 标签⽤于配置⼀个虚拟主机\n                      name：该host的名称\n                      appBase ：指定 war包放置的路径，可以是绝对路径，也可以是相对路径（相对路径，相对的就是tomcat的安装目录\n                      unpackWARs ：是否自动解压 war包\n                      autoDeploy：是否自动部署 （有点热部署的效果）-->\n            <Host name=\"localhost\" appBase=\"webapps\"\n                  unpackWARs=\"true\" autoDeploy=\"true\">\n\n                <!-- 记录当前 host 处理请求的日志 -->\n                <Valve className=\"org.apache.catalina.valves.AccessLogValve\" directory=\"logs\"\n                       prefix=\"localhost_access_log\" suffix=\".txt\"\n                       pattern=\"%h %l %u %t &quot;%r&quot; %s %b\"/>\n            </Host>\n        </Engine>\n    </Service>\n</Server>\n\n```\n其中有一个需要注意的地方就是，证书的位置certificateKeystoreFile可以填绝对路径，也可以填相对路径。如果填写的是相对路径，那地址应该是conf的上一层目录(如果你把jks文件放在server.xml的同级目录下，此处应该填\"conf/域名.jks\")，我因为这个踩过一些坑。\n## HTTP 自动跳转 HTTPS 的安全配置（可选）\n如果您需要将 HTTP 请求自动重定向到 HTTPS。您可以通过以下操作设置：\n\n编辑 /usr/*/conf 目录下的 web.xml 文件，找到 </welcome-file-list> 标签。\n请在结束标签 </welcome-file-list> 后面换行，并添加以下内容：\n```xml\n	<login-config>\n    <!-- Authorization setting for SSL -->\n    <auth-method>CLIENT-CERT</auth-method>\n    <realm-name>Client Cert Users-only Area</realm-name>\n    </login-config>\n    \n    <security-constraint>\n    <!-- Authorization setting for SSL -->\n    <web-resource-collection>\n    <web-resource-name>SSL</web-resource-name>\n    <url-pattern>/*</url-pattern>\n    </web-resource-collection>\n    <user-data-constraint>\n    <transport-guarantee>CONFIDENTIAL</transport-guarantee>\n    </user-data-constraint>\n    </security-constraint>\n```\n## 如何检验配置是否成功\n\n修改server.xml文件后，停止tomcat服务,在/usr/tomcat/*/bin目录下输入：\n```linux\n./shutdown.sh\n```\n然后以下命令检查配置文件是否有误，如果有报错信息就在网上搜索一下或者自己解决就好了：\n```linux\n./configtest.sh\n```\n\n以上步骤没有问题以后，输入以下命令就可以使用\"https://域名.com\"访问tomcat了。\n```linux\n./startup.sh\n```\n\n网上的教程大部分都只说./shutdown.sh和./startup.sh两个命令重启tomcat，但是有时候重启时有一些报错信息并不显示，就是访问80端口没有问题，但是443端口配置有错无法访问。\n\n还要注意的是，./configtest.sh命令一定要在./shutdown.sh停止tomcat之后执行，要不然会出现端口已被占用的错误。\n\n## 结语\n最近在忙最近的毕业设计，最近做后端接口和最近写网页，然后在腾讯云上面买了一台服务器还有一个域名(总共花了80块钱，租了一年)。花了半个月终于备案完成，于是迫不及待的把自己写的网页传上去。其中配置docker、tomcat、mysql、rabbitmq、redis都遇到一些坑。希望可以帮到大家吧。\n\n [1]: 腾讯云Tomcat 服务器 SSL 证书安装部署（JKS 格式）https://cloud.tencent.com/document/product/400/35224\n ', 1, '', 0, 0, 0, '2022-01-22 13:29:27', '2023-11-02 20:10:53');
INSERT INTO `article` VALUES (66, 2, 191, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/f9e9490619d664167657258c21db086a.jpeg', 'vue+nginx 打包发布刷新页面空白问题', '最近遇到一个问题，网站页面可以通过路由点击进入，但是直接在浏览器url处刷新页面之后会是空白页。\n\n原以为是nginx反向代理出现了问题，然后改了又改。后来又仔细想了想，会不会是vue的history模式出现问题。但是测试了之后发现，vue在本地运行时页面的跳转是正常的，也可以通过url刷新访问到页面。\n\n最后才醒悟，应该是vue打包之后环境问题。在一篇文章发现了这个bug的解决方案。\n\n来看看怎么说：\n\n```JavaScript\n因为我们的应用是单页客户端应用，当使用 history 模式时，URL 就像正常的 url，可以直接访问http://www.xxx.com/user/id，但是因为vue-router设置的路径不是真实存在的路径，所以刷新就会返回404错误。\n\n想要history模式正常访问，还需要后台配置支持。要在服务端增加一个覆盖所有情况的候选资源：如果 URL 匹配不到任何静态资源，则应该返回同一个 index.html 页面，这个页面就是你 app 依赖的页面。\n\n也就是在服务端修改404错误页面的配置路径，让其指向到index.html。\n\n**拓展**\n部署后，当访问一些页面的时候，报错 Uncaught SyntaxError: Unexpected token ‘＜’。\n\n解决方案：\n\n刚开始publicPath是’./’，需要改成’/’，即在vue.config.js中修改配置\n\nmodule.exports = {\n  ...\n  publicPath: \'/\',\n}\n\n```\n重要的是这句话”刚开始publicPath是’./’，需要改成’/’，即在vue.config.js中修改配置“,于是我打包时把路径改为‘/’。\n\n但是还是不行，浏览器页面报错找到“https://ve77.cn/static/***.css\"文件。\n\n看到这个问题我就想明白了,因为我是放在/blog路径下的,使用此处改成‘/blog’才对。果然，发布之后可以通过url访问了！！！\n\n\n分享一下自己的nginx配置\n```xml\n#user  nobody;\nworker_processes  2;\n#日志位置和日志级别\nerror_log /usr/local/webserver/nginx/logs/nginx_error.txt;\npid /usr/local/webserver/nginx/nginx.pid;\n\n#pid        logs/nginx.pid;\n\nevents {\n    worker_connections  1024;\n}\n\n\nhttp {\n    include       mime.types;\n    default_type  application/octet-stream;\n    sendfile        on;\n    keepalive_timeout  65;\n\n    client_max_body_size     50m;\n    client_body_buffer_size  10m;\n    client_header_timeout    1m;\n    client_body_timeout      1m;\n\n    gzip on;\n    gzip_min_length  1k;\n    gzip_buffers     4 16k;\n    gzip_comp_level  4;\n    gzip_types text/plain application/javascript application/x-javascript text/css application/xml text/javascript application/x-httpd-php image/jpeg image/gif image/png;\n    gzip_vary on;\n\n    access_log  logs/access.log ;\n\n    # http默认端口，转发到https\n    # 接口只能http访问，此时被重定向了。\n    server {\n        listen       80;\n        server_name  ve77.com  www.ve77.com  static.ve77.com  www.static.ve77.com;\n        rewrite ^(.*) https://$host$request_uri;\n        access_log  logs/host80.access.txt ;\n    }\n\n    # HTTPS server\n    server {\n        #SSL 访问端口号为 443\n        listen       443 ssl;\n        #填写绑定证书的域名\n        server_name  ve77.cn;  #填写绑定证书的域名\n        #证书文件名称\n        ssl_certificate      /jks/ve77.cn_nginx/ve77.cn_bundle.crt;\n        #私钥文件名称\n        ssl_certificate_key  /jks/ve77.cn_nginx/ve77.cn.key;\n        #请按照以下协议配置\n        ssl_protocols TLSv1.2 TLSv1.3;\n        #请按照以下套件配置，配置加密套件，写法遵循 openssl 标准。\n        ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:HIGH:!aNULL:!MD5:!RC4:!DHE;\n        ssl_prefer_server_ciphers  on;\n\n        ssl_session_cache    shared:SSL:1m;\n        ssl_session_timeout  5m;\n        access_log  logs/host443.access.txt ;\n\n        #末位别加/ 要不然会路径错误\n        location  /blog {\n            root   /usr/local/vue/ ; #站点目录+/blog\n            index  index.html index.htm;\n            try_files  $uri $uri/ /blog/index.html;\n        }\n\n        location /admin {\n            root   /usr/local/vue/;\n            index  index.html index.htm;\n            try_files $uri $uri/ /admin/index.html;\n        }\n        #重定向\n        location ^~ /api {\n            proxy_pass https://ve77.cn:8088;\n            proxy_set_header Host $host:8088; #这里是重点,这样配置才不会丢失端口\n            proxy_set_header   X-Real-IP        $remote_addr;\n            proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;\n        }\n\n        location ^~ /websocket {\n            proxy_pass https://ve77.cn:8088/api;\n            proxy_http_version 1.1;\n            proxy_set_header Upgrade $http_upgrade;\n            proxy_set_header Connection \"Upgrade\";\n            proxy_set_header Host $host:8088;\n            proxy_set_header X-Real-IP $remote_addr;\n            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;\n            proxy_set_header X-Forwarded-Proto $scheme;\n        }\n\n    }\n\n            #重点 start\n            # 这一点写location里面是一样的 但是提在外面就不用重复写了\n            # 如果内网nginx监听端口与外网访问的端口不一致 需要配置成这样\n            #	proxy_set_header Host $host:$server_port;\n            # 	proxy_set_header X-Real-IP $remote_addr;\n            #	proxy_set_header REMOTE-HOST $remote_addr;\n            # 	proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;\n            #重点 end\n            #在 Nginx 根目录下，通过执行以下命令验证配置文件问题。\n            #./sbin/nginx -t\n}\n\n```\n\n常用的nginx配置命令\n```xml\n重新加载/重启 nginx服务\n\n/usr/local/webserver/nginx/sbin/nginx -s reload\n\n/usr/local/webserver/nginx/sbin/nginx -s reopen\n\n \n\n验证nginx配置文件\n\n/usr/local/webserver/nginx/sbin/nginx -t\n```\n参考文献： \n\n[vue history模型下的问题](https://www.jb51.net/article/119075.htm)\n\n[Vue路由为history模式的nginx配置](https://blog.csdn.net/kiscon/article/details/115416832)\n', 1, '', 0, 0, 1, '2022-01-22 23:55:55', '2022-02-11 23:26:12');
INSERT INTO `article` VALUES (68, 2, 189, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/f9fa18da262910eb13f802b003147915.jpg', 'WebSocket 结合 Nginx 实现域名及 WSS 协议访问', '# [WebSocket 结合 Nginx 实现域名及 WSS 协议访问](https://www.cnblogs.com/mafly/p/websocket.html)\n\n## 简单了解一下 WebSocket\n\n现在，很多网站为了实现推送技术，所用的技术都是轮询。轮询是在特定的的时间间隔（如每1秒），由浏览器对服务器发出HTTP请求，然后由服务器返回最新的数据给客户端的浏览器。这种传统的模式带来很明显的缺点，即浏览器需要不断的向服务器发出请求，然而HTTP请求可能包含较长的头部，其中真正有效的数据可能只是很小的一部分，显然这样会浪费很多的带宽等资源。\n在这种情况下，HTML5定义了WebSocket协议，能更好的节省服务器资源和带宽，并且能够更实时地进行通讯。\nWebSocket一种在单个 TCP 连接上进行全双工通讯的协议。使得客户端和服务器之间的数据交换变得更加简单，允许服务端主动向客户端推送数据。在 WebSocket API 中，浏览器和服务器只需要完成一次握手，两者之间就直接可以创建持久性的连接，并进行双向数据传输。\n\n> 以上信息摘自维基百科（https://zh.wikipedia.org/wiki/WebSocket）\n\n简单点说，WebSocket 就是减小客户端与服务器端建立连接的次数，减小系统资源开销，只需要一次 HTTP 握手，整个通讯过程是建立在一次连接/状态中，也就避免了HTTP的非状态性，服务端会一直与客户端保持连接，直到你关闭请求，同时由原本的客户端主动询问，转换为服务器有信息的时候推送。当然，它还能做实时通信、更好的二进制支持、支持扩展、更好的压缩效果等这些优点。\n\n推荐一个知乎上叫 Ovear 的网友关于 WebSocket 原理的回答，嘻哈风格科普文，简直不要更赞了！地址：https://www.zhihu.com/question/20215561/answer/40316953\n\n## ws 和 wss 又是什么鬼？\n\nWebsocket使用 `ws` 或 `wss` 的统一资源标志符，类似于 `HTTP` 或 `HTTPS`，其中 `wss` 表示在 TLS 之上的 Websocket ，相当于 HTTPS 了。如：\n\n```bash\nws://example.com/chat\nwss://example.com/chat\n```\n\n默认情况下，Websocket 的 ws 协议使用 80 端口；运行在TLS之上时，wss 协议默认使用 443 端口。其实说白了，wss 就是 ws 基于 SSL 的安全传输，与 HTTPS 一样样的道理。\n\n如果你的网站是 HTTPS 协议的，那你就不能使用 `ws://` 了，浏览器会 block 掉连接，和 HTTPS 下不允许 HTTP 请求一样，如下图：\n![image.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/374bca19e5fc1b45850781e4131d3d9c.png)\n```rust\nMixed Content: The page at \'https://domain.com/\' was loaded over HTTPS, but attempted to connect to the insecure WebSocket endpoint \'ws://x.x.x.x:xxxx/\'. This request has been blocked; this endpoint must be available over WSS.\n```\n\n这种情况，毫无疑问我们就需要使用 `wss:\\\\` 安全协议了，我们是不是简单的把 `ws:\\\\` 改为 `wss:\\\\` 就行了？那试试呗。\n\n改好了，报错啦！！！\n![image.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/449049fb73844e5210bd5f20cfec9582.png)\n\n```vbnet\nVM512:35 WebSocket connection to \'wss://IP地址:端口号/websocket\' failed: Error in connection establishment: net::ERR_SSL_PROTOCOL_ERROR\n```\n\n很明显 SSL 协议错误，说明就是证书问题了。记着，这时候我们一直拿的是 `IP地址 + 端口号` 这种方式连接 WebSocket 的，这根本就没有证书存在好么，况且生成环境你也要用 `IP地址 + 端口号` 这种方式连接 WebSocket 吗？肯定不行阿，要用域名方式连接 WebSocket 阿。\n\n## Nginx 配置域名支持 WSS\n\n不用废话，直接在配置 HTTPS 域名位置加入如下配置：\n\n```bash\nlocation /websocket {\n    proxy_pass http://backend;\n    proxy_http_version 1.1;\n    proxy_set_header Upgrade $http_upgrade;\n    proxy_set_header Connection \"upgrade\";\n}\n```\n\n接着拿域名再次连接试一下，不出意外会看 101 状态码：\n![upgrade_101](https://images2015.cnblogs.com/blog/539095/201706/539095-20170622132017570-2009453161.png)\n\n这样就完成了，在 HTTPPS 下以域名方式连接 WebSocket ，可以开心的玩耍了。\n\n**稍微解释一下 Nginx 配置**\nNginx 自从 1.3 版本就开始支持 WebSocket 了，并且可以为 WebSocket 应用程序做反向代理和负载均衡。\nWebSocket 和 HTTP 协议不同，但是 WebSocket 中的握手和 HTTP 中的握手兼容，它使用 HTTP 中的 Upgrade 协议头将连接从 HTTP 升级到 WebSocket，当客户端发过来一个 `Connection: Upgrade`请求头时，Nginx 是不知道的，所以，当 Nginx 代理服务器拦截到一个客户端发来的 `Upgrade` 请求时，需要显式来设置`Connection` 、`Upgrade` 头信息，并使用 101（交换协议）返回响应，在客户端和代理服务器、后端服务器之间建立隧道来支持 WebSocket。\n\n当然，还需要注意一下，WebSockets 仍然受到 Nginx 缺省为60秒的 proxy_read_timeout 的影响。这意味着，如果你有一个程序使用了 WebSockets，但又可能超过60秒不发送任何数据的话，那你要么需要增加超时时间，要么实现一个 ping 的消息以保持联系。使用 ping 的解决方法有额外的好处，可以发现连接是否被意外关闭。\n\n更具体文档详见 Nginx 官方文档：http://nginx.org/en/docs/http/websocket.html\n\n## [总结一下](http://blog.mayongfa.cn/291.html)\n\n这一篇文章主要了解一下 WebSocket 基本原理和一些使用用途，并解决在实际开发使用过程中遇到的坑，HTTPS 下使用 wss 协议的问题，以及配合 Nginx 使用域名方式建立连接，不使用 `IP地址 + 端口号` 连接 WebSocket，因为这种方式不够优雅。\n\n原文链接:https://www.cnblogs.com/mafly/p/websocket.html', 1, '', 0, 0, 1, '2022-02-09 23:52:33', '2022-02-11 23:32:46');
INSERT INTO `article` VALUES (69, 2, 191, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/f9fa18da262910eb13f802b003147915.jpg', '网站搭建之https配置websocket连接(wss)', '## 前言\nws和wss的关系，就像http和https的关系。\n当网站发布上线后，由于websocket问题，聊天室功能不能正常使用。具体问题如下：使用http非安全协议访问网站时，可以访问ws://localhost:8088/api/websocket，聊天室功能正常。但是使用https安全协议访问网站时，可以访问ws://localhost:8088/api/websocket失败\n\n![BAC09972C5BC2F039B371679953FCDC7.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/bac09972c5bc2f039b371679953fcdc7.png)\n\n\n查看报错信息可以知道，因为https不支持ws协议，这个请求被拦截了。所以我们应该想到https只支持wss协议。那么我们只简单的将websocket地址前缀改为wss？修改后依然报错：\n\n![04B43675B833412CB6580A279614331C.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/028181b8fe4cd46c4b13b3267d3d703d.png)\n\n出现这个问题的原因是，后端springboot并不支持https访问，因此我们需要给后端添加SSL证书。\n\n```xml\n在终端输入以下命令，生成keystore.jks证书：\nkeytool -genkey -alias tomcat -keyalg RSA -keystore ./keystore.jks -validity 3650 -ext san=dns:spam,ip:110.42.180.40\n```\n\n![3386BE795E832853F5A0BDC23B19C265.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/ee48a2d53494dffe911025b3c1bd1bd7.png)\n\n其中 -keystore是证书路径，-validity 是有效期(可选)，单位为天，-ext 是额外信息(可选)，这里填了发布单位。\n\n然后在项目的配置文件yml里设置\n\n![image20220210104501701.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/4648a2b0bcf1f0417fd56ab5cef7eefa.png)\n\n\n当可以使用https访问接口时，说明已经配置成功\n\n![B683E1968A284CD698005E6C84ED786A.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/2a93889977a7a0a8053ca8fc15afb0ae.png)\n\n这个时候就可以使用wss连接websocket啦，[websocket在线测试地址](http://www.jsons.cn/websocket/)。  \n(注意！！不要使用Chrome浏览器测试，因为我们自己颁发的证书是不被Chrome认可的！！如果想使用)  \n\n**推荐使用官方颁发的SLL证书部署websocket和springboot。具体实现方式可查阅参考链接4和参考链接5。\n**\n\n参考链接:\n\n1.[spring boot 集成 websocket 的四种方式](https://www.cnblogs.com/kiwifly/p/11729304.html)\n2.[Springboot2.1集成WebSocket配置wss访问](https://blog.csdn.net/u012977315/article/details/84944708)\n3.[websocket在线测试地址](http://www.jsons.cn/websocket/)\n4.[Nginx 服务器 SSL 证书安装部署[腾讯云]](https://cloud.tencent.com/document/product/400/35244)\n5.[在Spring Boot中配置ssl证书实现https](https://www.jianshu.com/p/eb52e0f5ee85)', 1, '', 0, 0, 1, '2022-02-10 10:57:51', '2022-02-11 23:31:32');
INSERT INTO `article` VALUES (70, 2, 192, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/f9fa18da262910eb13f802b003147915.jpg', '算法日记之动态规划-背包问题', '## 动态规划:\n### **算法模板：**\n```java\nimport java.util.*;\npublic class Main{\n    public static void main(String[] args)  {\n        Scanner input=new Scanner(System.in);\n        int m=input.nextInt();//横行 → 容量大小\n        int n=input.nextInt();//竖行 ↓ 物品个数\n        int dp[][]=new int[n+1][m+1];\n\n        for(int i=0;i<=m;i++){\n            dp[i][1] = 1;//初始化，当0个物品放入容量为1的背包时的收益\n        }\n        for(int j=0;j<=n;j++){\n            dp[1][j] = 1;//初始化，当1个物品放入容量为0的背包时的收益\n        }\n        for(int i=1;i<=n;i++){\n            for(int j=1;j<=m;j++){\n                //递推公式，不同问题公式不同。注意此处判断条件是>=\n                if(j>=i) {\n                    dp[i][j]=Math.max(dp[i-1][j],dp[i-1][j-i]);\n                }\n            }\n        }\n        System.out.println(dp[n][m]);\n    }\n}\n```\n\n\n\n\n### **0-1背包问题**\n\n问题描述：有一个背包可以装物品的总重量为W，现有N个物品，每个物品中w[i]，价值v[i]，用背包装物品，能装的最大价值是多少？\n\n物品只有一件，装下后剩余为0，未装下剩余是1.因此称为01背包问题。\n\n**定义状态转移数组dp[i][j]，表示前i个物品，背包重量为j的情况下能装的最大价值。**\n\n例如，dp[3][4]=6，表示用前3个物品装入重量为4的背包所能获得的最大价值为6，此时并不是3个物品全部装入，而是3个物品满足装入背包的条件下的最大价值。\n\n状态转移方程：\ndp[i][j] = max(dp[i-1][j], dp[i-1][j-w[i]]+v[i])\n\ndp[i-1][j]表示当前物品不放入背包，背包空间剩余仍为j，获得收益不变；dp[i-1][j-w[i]]+v[i]表示当前物品放入背包，背包空间剩余为j-w[i]，获得收益为dp[i-1][j-w[i]]+v[i]；  \n**即当前第i个物品要么放入背包，要么不放入背包**。\n\n----\n**个人理解**\n拿a[1][1]来说，它的值就是背包容量为1，只考虑编号0，1的物品时，背包所能装入的最大价值；\n然后既然是动态规划，那就一定有初值，也就是a[0][j] = 0;  a[i][0] = 0;即第一行和第一列都为0；\n然后根据初值来推后面的值；\n\n首先要判断本行所对应的物品是否能装入背包，\n拿a[1][1]来说，首先要判断，若只考虑编号为1的物品，它是否可以装入背包，此时的背包容量为1，而编号为1的物品的体积为2，故它无法装入背包，那么a[1][1]的值和背包容量为1，只考虑编号为0的物品时，背包所能装入的最大价值(即a[0][1])是相等的；\n\n若能装入背包；那么有两种选择:\n(1)装入本行物品，即先装入本行的物品，然后剩下背包容量装其他价值之和最大的物品\n(2)不装本行物品，即背包容量都用来装除了本行物品之外的其他物品(即本行前面几行的物品)\n然后比较(1)(2)选择较大者；\n\n拿a[2][4]来说，此时的背包容量为4,编号为2的物品的体积为3，故2号物品能装入背包，然后两种选择：\n(1)装入2号物品，此时背包剩余容量为1，此时只剩下两个物品，那就是编号为0和1的物品，查表得a[1][1]=0\n故此时的最大价值为a[1][1]加上2号物品的价值，也就是4\n(2)不装2号物品，即背包容量都用来装除了本行物品之外的其他物品(即本行前面几行的物品)\n\n由于不装入2号物品，此时的最大价值和只考虑编号为0，1物品，背包容量为4的情况的最大价值(即a[1][4])是相等的，\n也就是3；\n故选择(1)(2)中较大者，a[2][4]=4;\n\n----\n\n算法：\n```Java\n    /***\n     * 01背包问题\n     * dp[i][j] = max(dp[i-1][j], dp[i-1][j-w[i]] +v[i])\n     * dp[i-1][j]表示当前物品不放入背包，背包空间剩余仍为j，获得收益不变；\n     * dp[i-1][j-w[i]]+v[i]表示当前物品放入背包，背包空间剩余为j-w[i]，获得收益为dp[i-1][j-w[i]]+v[i]；\n     * @param n 背包容量\n     * @param m 物品个数\n     * @param v 物品体积数组\n     * @param w 物品价值数组\n     * @return 最大收益\n     */\n    public static int Knapsack(int n, int m, int[] v, int[] w){\n\n        int dp[][] = new int[m + 1][n + 1];//m件物品放入n容量背包\n\n        for (int i = 1; i <= m; i++) {\n            for (int j = 1; j <= n; j++) {\n                //可以装下\n                if (j > v[i]) {\n                    dp[i][j] = Math.max(dp[i - 1][j], dp[i - 1][j - v[i]] + v[i] * w[i]);\n                    //System.out.println(\"dp[\" + i + \"][\" + j + \"]=\" + dp[i][j]);\n                }\n            }\n        }\n        //体积volume\n        int volume = n;\n        //输出装入背包的物品,回溯\n        System.out.println(\"编号---重量---价值---收益\");\n        for (int i = m; i >= 1; i--) {\n            if (dp[i][volume] == dp[i - 1][volume]) {\n                //没有装入该物品，所以容量不变\n            } else {\n                //装入了该物品，容量减少v[i]\n                volume = volume - v[i];\n                System.out.println(i+\"---\"+v[i] + \"---\" + w[i]+\"---\" + w[i]*v[i]);\n            }\n        }\n        return dp[m][n];\n    }\n\n```\n测试：\n![C91B2F4DEAEF451CB6616056DB86FC88.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/d0c1e2267c325ae496e92411418a1ad3.png)\n\n----\n### 放苹果\n## 描述\n把m个同样的苹果放在n个同样的盘子里，允许有的盘子空着不放，问共有多少种不同的分法？（用K表示）5，1，1和1，5，1 是同一种分法。\n数据范围：0 \\le m \\le 10 \\0≤*m*≤10 ，1 \\le n \\le 10 \\1≤*n*≤10 。\n本题含有多组样例输入。\n\n### 输入描述：\n输入两个int整数\n\n### 输出描述：\n输出结果，int型\n\n## 示例1\n输入：\n```bash\n7 3\n```\n输出：\n```bash\n8\n```\n\n算法：\n```Java\npublic class Main {\n    public static void main(String[] args)  {\n        Scanner scanner = new Scanner(System.in);\n        while (scanner.hasNext()){\n            int apples = scanner.nextInt();//苹果0 - 10最多（下同，最多）\n            int panels = scanner.nextInt();//1-10\n            int[][] dp = new int[apples + 1][panels + 1];\n\n            for(int i = 0; i <= apples; i++){\n                dp[i][1] = 1;//多少个苹果放在一个盘子上，都只有一种方法\n            }\n            for(int j = 1; j <= panels; j++){\n                dp[1][j] = 1;//只有一个苹果，不管有多少盘子都只有一种方法\n                dp[0][j] = 1;//0个苹果，不管有多少盘子都只有一种方法\n            }\n\n            for(int i = 2; i <= apples; i++){//2个苹果到10个苹果(因为0-1苹果的情况已经在上面处理)\n                for(int j = 1; j <= panels; j++){//1个盘子到10个盘子\n                    //i个苹果放到j个盘子里的方法数，等于所有盘子都有苹果的方法 + 部分盘子没有苹果的方法\n                    //所有盘子都有苹果，等于每个盘子去掉一个苹果，\n                    //部分盘子没有苹果，先假设一个盘子没有苹果，剩下的盘子到底怎么分配苹果，慢慢算\n                    if(i>=j) {\n                        dp[i][j] = dp[i][j - 1] +  dp[i - j][j];\n                    }else{\n                        dp[i][j] = dp[i][j - 1];\n                    }\n                }\n            }\n            System.out.println(dp[apples][panels]);\n        }\n    }\n}\n```\n', 1, '', 0, 0, 1, '2022-02-12 16:28:13', '2022-02-19 21:54:48');
INSERT INTO `article` VALUES (71, 2, 189, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/f9fa18da262910eb13f802b003147915.jpg', 'cookie、session与token的真正区别', '## **发展史**\n\n1、很久很久以前，Web 基本上就是文档的浏览而已， 既然是浏览，作为服务器， 不需要记录谁在某一段时间里都浏览了什么文档，每次请求都是一个新的HTTP协议， 就是请求加响应， 尤其是我不用记住是谁刚刚发了HTTP请求， 每个请求对我来说都是全新的。这段时间很嗨皮。\n\n2、但是随着交互式Web应用的兴起，像在线购物网站，需要登录的网站等等，马上就面临一个问题，那就是要管理会话，必须记住哪些人登录系统， 哪些人往自己的购物车中放商品， 也就是说我必须把每个人区分开，这就是一个不小的挑战，因为HTTP请求是无状态的，所以想出的办法就是给大家发一个会话标识([session](https://so.csdn.net/so/search?q=session&spm=1001.2101.3001.7020) id), 说白了就是一个随机的字串，每个人收到的都不一样， 每次大家向我发起HTTP请求的时候，把这个字符串给一并捎过来， 这样我就能区分开谁是谁了\n\n3、这样大家很嗨皮了，可是服务器就不嗨皮了，每个人只需要保存自己的session id，而服务器要保存所有人的session id ！如果访问服务器多了， 就得由成千上万，甚至几十万个。\n\n这对服务器说是一个巨大的开销 ， 严重的限制了服务器扩展能力， 比如说我用两个机器组成了一个集群， 小F通过机器A登录了系统， 那session id会保存在机器A上， 假设小F的下一次请求被转发到机器B怎么办？机器B可没有小F的 session id啊。\n\n有时候会采用一点小伎俩： **session sticky** ， 就是让小F的请求一直粘连在机器A上， 但是这也不管用， 要是机器A挂掉了， 还得转到机器B去。\n\n那只好做session 的复制了， 把session id 在两个机器之间搬来搬去， 快累死了。\n\n![img](https://img-blog.csdnimg.cn/20190509111418335.jpg?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3dobDE5MDQxMg==,size_16,color_FFFFFF,t_70)\n\n后来有个叫Memcached的支了招：把session id 集中存储到一个地方， 所有的机器都来访问这个地方的数据， 这样一来，就不用复制了， 但是增加了单点失败的可能性， 要是那个负责session 的机器挂了， 所有人都得重新登录一遍， 估计得被人骂死。\n\n![img](https://img-blog.csdnimg.cn/2019050911155517.jpg?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3dobDE5MDQxMg==,size_16,color_FFFFFF,t_70)\n\n也尝试把这个单点的机器也搞出集群，增加可靠性， 但不管如何， 这小小的session 对我来说是一个沉重的负担\n\n4、于是有人就一直在思考， 我为什么要保存这可恶的session呢， 只让每个客户端去保存该多好？\n\n可是如果不保存这些session id , 怎么验证客户端发给我的session id 的确是我生成的呢？ 如果不去验证，我们都不知道他们是不是合法登录的用户， 那些不怀好意的家伙们就可以伪造session id , 为所欲为了。\n\n嗯，对了，关键点就是验证 ！\n\n比如说， 小F已经登录了系统， 我给他发一个令牌(token)， 里边包含了小F的 user id， 下一次小F 再次通过Http 请求访问我的时候， 把这个token 通过Http header 带过来不就可以了。\n\n不过这和session id没有本质区别啊， 任何人都可以可以伪造， 所以我得想点儿办法， 让别人伪造不了。\n\n那就对数据做一个签名吧， 比如说我用HMAC-SHA256 算法，加上一个只有我才知道的密钥， 对数据做一个签名， 把这个签名和数据一起作为token ， 由于密钥别人不知道， 就无法伪造token了。\n\n![img](https://img-blog.csdnimg.cn/20190509111646773.jpg?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3dobDE5MDQxMg==,size_16,color_FFFFFF,t_70)\n\n这个token 我不保存， 当小F把这个token 给我发过来的时候，我再用同样的HMAC-SHA256 算法和同样的密钥，对数据再计算一次签名， 和token 中的签名做个比较， 如果相同， 我就知道小F已经登录过了，并且可以直接取到小F的user id , 如果不相同， 数据部分肯定被人篡改过， 我就告诉发送者：对不起，没有认证。\n\n![img](https://img-blog.csdnimg.cn/20190509111736141.jpg?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3dobDE5MDQxMg==,size_16,color_FFFFFF,t_70)\n\nToken 中的数据是明文保存的（虽然我会用Base64做下编码， 但那不是加密）， 还是可以被别人看到的， 所以我不能在其中保存像密码这样的敏感信息。\n\n当然， 如果一个人的token 被别人偷走了， 那我也没办法， 我也会认为小偷就是合法用户， 这其实和一个人的session id 被别人偷走是一样的。\n\n这样一来， 我就不保存session id 了， 我只是生成token , 然后验证token ， 我用我的CPU计算时间获取了我的session 存储空间 ！\n\n解除了session id这个负担， 可以说是无事一身轻， 我的机器集群现在可以轻松地做水平扩展， 用户访问量增大， 直接加机器就行。这种无状态的感觉实在是太好了！\n\n## **Cookie**\n\ncookie 是一个非常具体的东西，指的就是浏览器里面能永久存储的一种数据，仅仅是浏览器实现的一种数据存储功能。\n\ncookie由服务器生成，发送给浏览器，浏览器把cookie以kv形式保存到某个目录下的文本文件内，下一次请求同一网站时会把该cookie发送给服务器。由于cookie是存在客户端上的，所以浏览器加入了一些限制确保cookie不会被恶意使用，同时不会占据太多磁盘空间，所以每个域的cookie数量是有限的。\n\n## **Session**\n\nsession 从字面上讲，就是会话。这个就类似于你和一个人交谈，你怎么知道当前和你交谈的是张三而不是李四呢？对方肯定有某种特征（长相等）表明他就是张三。\n\nsession 也是类似的道理，服务器要知道当前发请求给自己的是谁。为了做这种区分，服务器就要给每个客户端分配不同的“身份标识”，然后客户端每次向服务器发请求的时候，都带上这个“身份标识”，服务器就知道这个请求来自于谁了。至于客户端怎么保存这个“身份标识”，可以有很多种方式，对于浏览器客户端，大家都默认采用 cookie 的方式。\n\n服务器使用session把用户的信息临时保存在了服务器上，用户离开网站后session会被销毁。这种用户信息存储方式相对cookie来说更安全，可是session有一个缺陷：如果web服务器做了负载均衡，那么下一个操作请求到了另一台服务器的时候session会丢失。\n\n## **Token**\n\n在Web领域基于Token的身份验证随处可见。在大多数使用Web API的互联网公司中，tokens 是多用户下处理认证的最佳方式。\n\n以下几点特性会让你在程序中使用基于Token的身份验证\n\n1. 无状态、可扩展\n2. 支持移动设备\n3. 跨程序调用\n4. 安全\n\n那些使用基于Token的身份验证的大佬们\n\n大部分你见到过的API和Web应用都使用tokens。例如Facebook, Twitter, Google+, GitHub等。\n\n### **Token的起源**\n\n在介绍基于Token的身份验证的原理与优势之前，不妨先看看之前的认证都是怎么做的。\n\n**基于服务器的验证**\n\n我们都是知道HTTP协议是无状态的，这种无状态意味着程序需要验证每一次请求，从而辨别客户端的身份。\n\n在这之前，程序都是通过在服务端存储的登录信息来辨别请求的。这种方式一般都是通过存储Session来完成。\n\n随着Web，应用程序，已经移动端的兴起，这种验证的方式逐渐暴露出了问题。尤其是在可扩展性方面。\n\n**基于服务器验证方式暴露的一些问题**\n\n1. **Seesion：**每次认证用户发起请求时，服务器需要去创建一个记录来存储信息。当越来越多的用户发请求时，内存的开销也会不断增加。\n2. **可扩展性：**在服务端的内存中使用Seesion存储登录信息，伴随而来的是可扩展性问题。\n3. **CORS(跨域资源共享)：**当我们需要让数据跨多台移动设备上使用时，跨域资源的共享会是一个让人头疼的问题。在使用Ajax抓取另一个域的资源，就可以会出现禁止请求的情况。\n4. **CSRF(跨站请求伪造)：**用户在访问银行网站时，他们很容易受到跨站请求伪造的攻击，并且能够被利用其访问其他的网站。\n\n在这些问题中，可扩展行是最突出的。因此我们有必要去寻求一种更有行之有效的方法。\n\n**基于Token的验证原理**\n\n基于Token的身份验证是无状态的，我们不将用户信息存在服务器或Session中。\n\n这种概念解决了在服务端存储信息时的许多问题\n\n> NoSession意味着你的程序可以根据需要去增减机器，而不用去担心用户是否登录。\n\n基于Token的身份验证的过程如下:\n\n1. 用户通过用户名和密码发送请求。\n2. 程序验证。\n3. 程序返回一个签名的token 给客户端。\n4. 客户端储存token,并且每次用于每次发送请求。\n5. 服务端验证token并返回数据。\n\n每一次请求都需要token。token应该在HTTP的头部发送从而保证了Http请求无状态。我们同样通过设置服务器属性Access-Control-Allow-Origin:* ，让服务器能接受到来自所有域的请求。\n\n需要主要的是，在ACAO头部标明(designating)*时，不得带有像HTTP认证，客户端SSL证书和cookies的证书。\n\n实现思路：\n\n![img](https://img-blog.csdnimg.cn/20190509111955368.jpg?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3dobDE5MDQxMg==,size_16,color_FFFFFF,t_70)\n\n1. 用户登录校验，校验成功后就返回Token给客户端。\n2. 客户端收到数据后保存在客户端\n3. 客户端每次访问API是携带Token到服务器端。\n4. 服务器端采用filter过滤器校验。校验成功则返回请求数据，校验失败则返回错误码\n\n当我们在程序中认证了信息并取得token之后，我们便能通过这个Token做许多的事情。\n\n我们甚至能基于创建一个基于权限的token传给第三方应用程序，这些第三方程序能够获取到我们的数据（当然只有在我们允许的特定的token）\n\n### **Token的优势**\n\n**无状态、可扩展**\n\n在客户端存储的Tokens是无状态的，并且能够被扩展。基于这种无状态和不存储Session信息，负载负载均衡器能够将用户信息从一个服务传到其他服务器上。\n\n如果我们将已验证的用户的信息保存在Session中，则每次请求都需要用户向已验证的服务器发送验证信息(称为Session亲和性)。用户量大时，可能会造成一些拥堵。\n\n但是不要着急。使用tokens之后这些问题都迎刃而解，因为tokens自己hold住了用户的验证信息。\n\n**安全性**\n\n请求中发送token而不再是发送cookie能够防止CSRF(跨站请求伪造)。即使在客户端使用cookie存储token，cookie也仅仅是一个存储机制而不是用于认证。不将信息存储在Session中，让我们少了对session操作。\n\ntoken是有时效的，一段时间之后用户需要重新验证。我们也不一定需要等到token自动失效，token有撤回的操作，通过token revocataion可以使一个特定的token或是一组有相同认证的token无效。\n\n**可扩展性**\n\nTokens能够创建与其它程序共享权限的程序。例如，能将一个随便的社交帐号和自己的大号(Fackbook或是Twitter)联系起来。当通过服务登录Twitter(我们将这个过程Buffer)时，我们可以将这些Buffer附到Twitter的数据流上(we are allowing Buffer to post to our Twitter stream)。\n\n使用tokens时，可以提供可选的权限给第三方应用程序。当用户想让另一个应用程序访问它们的数据，我们可以通过建立自己的API，得出特殊权限的tokens。\n\n**多平台跨域**\n\n我们提前先来谈论一下CORS(跨域资源共享)，对应用程序和服务进行扩展的时候，需要介入各种各种的设备和应用程序。\n\n> Having our API just serve data, we can also make the design choice to serve assets from a CDN. This eliminates the issues that CORS brings up after we set a quick header configuration for our application.\n\n只要用户有一个通过了验证的token，数据和资源就能够在任何域上被请求到。\n\nAccess-Control-Allow-Origin: *   \n\n基于标准创建token的时候，你可以设定一些选项。我们在后续的文章中会进行更加详尽的描述，但是标准的用法会在JSON Web Tokens体现。\n\n最近的程序和文档是供给JSON Web Tokens的。它支持众多的语言。这意味在未来的使用中你可以真正的转换你的认证机制。\n\n-----\n\n# [token和session的区别](https://www.cnblogs.com/belongs-to-qinghua/p/11353228.html)\n\n　　session和token都是用来保持会话，功能相同\n\n# 一、**session机制，原理**\n\n![image.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/b09d9fde1415d82d3bf4d0c79d6e1737.png)\n\n- session是服务端存储的一个对象，主要用来存储所有访问过该服务端的客户端的用户信息（也可以存储其他信息），从而实现保持用户会话状态。但是服务器重启时，内存会被销毁，存储的用户信息也就消失了。\n\n　　　　不同的用户访问服务端的时候会在session对象中存储键值对，“键”用来存储开启这个用户信息的“钥匙”，在登录成功后，“钥匙”通过cookie返回给客户端，客户端存储为sessionId记录在cookie中。当客户端再次访问时，会**默认携带**cookie中的sessionId来实现会话机制。\n\n- session是基于cookie的。\n\n1. cookie的数据4k左右\n2. cookie存储数据的格式：字符串key=value\n3. cookie存储有效期：可以自行通过expires进行具体的日期设置，如果没设置，默认是关闭浏览器时失效。\n4. cookie有效范围：当前域名下有效。所以**session这种会话存储方式方式只适用于客户端代码和服务端代码运行在同一台服务器上**（前后端项目协议、域名、端口号都一致，即在一个项目下）\n\n- session持久化\n\n　　　　用于解决重启服务器后session就消失的问题。在数据库中存储session，而不是存储在内存中。通过包：express-mysql-session\n\n- 其它\n\n　　　　当客户端存储的cookie失效后，服务端的session不会立即销毁，会有一个延时，服务端会定期清理无效session，不会造成无效数据占用存储空间的问题。\n\n# 二、**token机制，原理**\n\n![image.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/d33e38ca01de0b5579421eae1d6d6128.png)\n- **适用于项目级的前后端分离（前后端代码运行在不同的服务器下）**\n\n　　　　请求登录时，token和sessionId原理相同，是对key和key对应的用户信息进行加密后的加密字符，登录成功后，会在响应主体中将{token：\'字符串\'}返回给客户端。客户端通过cookie、sessionStorage、localStorage都可以进行存储。再次请求时**不会默认携带**，需要在请求拦截器位置给请求头中添加认证字段Authorization携带token信息，服务器端就可以通过token信息查找用户登录状态。', 2, 'https://blog.csdn.net/whl190412/article/details/90024671', 0, 0, 1, '2022-02-15 17:25:15', '2022-02-19 10:29:08');
INSERT INTO `article` VALUES (72, 2, 192, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/f9fa18da262910eb13f802b003147915.jpg', '排序算法', '![0C4C13288D00FC4FDAC24D8BA0936E4F.jpg]( https://veport.oss-cn-beijing.aliyuncs.com/articles/0c4c13288d00fc4fdac24d8ba0936e4f.jpg)\n\n# 算法1：最快最简单的排序——桶排序\n\n在我们生活的这个世界中到处都是被排序过的。站队的时候会按照身高排序，考试的名次需要按照分数排序，网上购物的时候会按照价格排序，电子邮箱中的邮件按照时间排序……总之很多东西都需要排序，可以说排序是无处不在。现在我们举个具体的例子来介绍一下排序算法。\n\n![picture1.1](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/1.1.png)\n\n首先出场的我们的主人公小哼，上面这个可爱的娃就是啦。期末考试完了老师要将同学们的分数按照从高到低排序。小哼的班上只有 5 个同学，这 5 个同学分别考了 5 分、3 分、5 分、2 分和 8 分，哎考的真是惨不忍睹（满分是 10 分）。接下来将分数进行从大到小排序，排序后是 8 5 5 3 2。你有没有什么好方法编写一段程序，让计算机随机读入 5 个数然后将这 5 个数从大到小输出？请先想一想，至少想 15 分钟再往下看吧(*^__^*) 。\n\n![picture1.2](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/1.2.png)\n\n我们这里只需借助一个一维数组就可以解决这个问题。请确定你真的仔细想过再往下看哦。\n\n首先我们需要申请一个大小为 11 的数组 int a[11]。OK 现在你已经有了 11 个变量，编号从 a[0]~a[10]。刚开始的时候，我们将 a[0]~a[10]都初始化为 0，表示这些分数还都没有人得过。例如 a[0]等于 0 就表示目前还没有人得过 0 分，同理 a[1]等于 0 就表示目前还没有人得过 1 分……a[10]等于 0 就表示目前还没有人得过 10 分。\n\n![picture1.3](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/1.3.png)\n\n下面开始处理每一个人的分数，第一个人的分数是 5 分，我们就将相对应 a[5]的值在原来的基础增加 1，即将 a[5]的值从 0 改为 1，表示 5 分出现过了一次。\n\n![picture1.4](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/1.4.png)\n\n第二个人的分数是 3 分，我们就把相对应 a[3]的值在原来的基础上增加 1，即将 a[3]的值从 0 改为 1，表示 3 分出现过了一次。\n\n![picture1.5](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/1.5.png)\n\n注意啦！第三个人的分数也是“5 分”，所以a[5]的值需要在此基础上再增加 1，即将 a[5]的值从 1 改为 2。表示 5 分出现过了两次。\n\n![picture1.6](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/1.6.png)\n\n按照刚才的方法处理第四个和第五个人的分数。最终结果就是下面这个图啦。\n\n![picture1.7](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/1.7.png)\n\n\n\n你发现没有，a[0]~a[10]中的数值其实就是 0 分到 10 分每个分数出现的次数。接下来，我们只需要将出现过的分数打印出来就可以了，出现几次就打印几次，具体如下。 　　a[0]为 0，表示“0”没有出现过，不打印。\n　　a[1]为 0，表示“1”没有出现过，不打印。\n　　a[2]为 1，表示“2”出现过 1 次，打印 2。\n　　a[3]为 1，表示“3”出现过 1 次，打印 3。\n　　a[4]为 0，表示“4”没有出现过，不打印。\n　　a[5]为 2，表示“5”出现过 2 次，打印5 5。\n　　a[6]为 0，表示“6”没有出现过，不打印。\n　　a[7]为 0，表示“7”没有出现过，不打印。\n　　a[8]为 1，表示“8”出现过 1 次，打印 8。\n　　a[9]为 0，表示“9”没有出现过，不打印。\n　　a[10]为 0，表示“10”没有出现过，不打印。\n　　最终屏幕输出“2 3 5 5 8”，另外此处的每一个桶的作用其实就是“标记”每个数出现的次数，因此我喜欢将之前的数组 a 换个更贴切的名字 book（book 这个单词有记录、标记的意思），代码实现如下。。\n\n```c\n    #include <stdio.h>\n    int main()\n    {\n        int a[11],i,j,t;\n        for(i=0;i<=10;i++)\n            a[i]=0;  //初始化为0\n\n        for(i=1;i<=5;i++)  //循环读入5个数\n        {\n            scanf(\"%d\",&t);  //把每一个数读到变量t中\n            a[t]++;  //进行计数\n        }\n        for(i=0;i<=10;i++)  //依次判断a[0]~a[10]\n            for(j=1;j<=a[i];j++)  //出现了几次就打印几次\n                printf(\"%d \",i);\n        getchar();getchar();\n        //这里的getchar();用来暂停程序，以便查看程序输出的内容\n        //也可以用system(\"pause\");等来代替\n        return 0;\n    }\n```\n\n输入数据为\n\n```sh\n5 3 5 2 8 \n```\n\n这种排序方法我们暂且叫他“桶排序”。因为其实真正的桶排序要比这个复杂一些，以后再详细讨论，目前此算法已经能够满足我们的需求了。\n\n这个算法就好比有 11 个桶，编号从 0~10。每出现一个数，就将对应编号的桶中的放一个小旗子，最后只要数数每个桶中有几个小旗子就 OK 了。例如 2 号桶中有 1 个小旗子，表示 2 出现了一次；3 号桶中有 1 个小旗子，表示 3 出现了一次；5 号桶中有 2 个小旗子，表示 5 出现了两次；8 号桶中有 1 个小旗子，表示 8 出现了一次。\n\n![picture1.8](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/1.8.png)\n\n&emsp;&emsp;这是一个非常快的排序算法。桶排序从 1956 年就开始被使用，该算法的基本思想是由 E.J.Issac R.C.Singleton 提出来。之前说过，其实这并不是真正的桶排序算法，真正的桶排序算法要比这个更加复杂。\n\n&emsp;&emsp;真正的桶排序思想是**划分多个范围相同的区间，每个子区间自排序，最后合并**。\n\n三、核心代码\n\n```java\n    public static void bucketSort(int[] arr){\n        \n        // 计算最大值与最小值\n        int max = Integer.MIN_VALUE;\n        int min = Integer.MAX_VALUE;\n        for(int i = 0; i < arr.length; i++){\n            max = Math.max(max, arr[i]);\n            min = Math.min(min, arr[i]);\n        }\n        \n        // 计算桶的数量\n        int bucketNum = (max - min) / arr.length + 1;\n        ArrayList<ArrayList<Integer>> bucketArr = new ArrayList<>(bucketNum);\n        for(int i = 0; i < bucketNum; i++){\n            bucketArr.add(new ArrayList<Integer>());\n        }\n        \n        // 将每个元素放入桶\n        for(int i = 0; i < arr.length; i++){\n            int num = (arr[i] - min) / (arr.length);\n            bucketArr.get(num).add(arr[i]);\n        }\n        \n        // 对每个桶进行排序\n        for(int i = 0; i < bucketArr.size(); i++){\n            Collections.sort(bucketArr.get(i));\n        }\n        \n        // 将桶中的元素赋值到原序列\n    	int index = 0;\n    	for(int i = 0; i < bucketArr.size(); i++){\n    		for(int j = 0; j < bucketArr.get(i).size(); j++){\n    			arr[index++] = bucketArr.get(i).get(j);\n    		}\n    	}  \n    }\n```\n\n四、复杂度分析\n\n1. 时间复杂度：O(N + C)\n对于待排序序列大小为 N，共分为 M 个桶，主要步骤有：\n   N 次循环，将每个元素装入对应的桶中\n   M 次循环，对每个桶中的数据进行排序（平均每个桶有 N/M 个元素）\n   一般使用较为快速的排序算法，时间复杂度为 O ( N l o g N )            O(NlogN)O(NlogN)，实际的桶排序过程是以链表形式插入的。\n\n整个桶排序的时间复杂度为：\n&emsp;O ( N ) + O ( M ∗ ( N / M ∗ l o g ( N / M ) ) ) = O ( N ∗ ( l o g ( N / M ) + 1 ) ) O(N)+O(M*(N/M*log(N/M)))=O(N*(log(N/M)+1))O(N)+O(M∗(N/M∗log(N/M))) =O(N∗(log(N/M)+1))\n\n当 N = M 时，复杂度为 O ( N ) O(N)O(N)\n\n&emsp;2. 额外空间复杂度：O(N + M)\n五、稳定性分析\n&emsp;桶排序的稳定性取决于桶内排序使用的算法。\n# 算法 2：邻居好说话：冒泡排序\n**1.冒泡的基本思想**\n\n冒泡排序是一种交换排序，核心是冒泡，把数组中最大的那个往上冒，冒的过程就是和他相邻的元素交换。\n\n重复走访要排序的数列，通过两两比较相邻记录的排序码。排序过程中每次从后往前冒一个最小值，且每次能确定一个数在序列中的最终位置。若发生逆序，则交换；有俩种方式进行冒泡，一种是先把小的冒泡到前边去，另一种是把大的元素冒泡到后边。\n\n**2.实现逻辑：**\n\n比较相邻的元素。如果第一个比第二个大，就交换他们两个。\n\n对每一对相邻元素作同样的工作，从开始第一对到结尾的最后一对。在这一点，最后的元素应该会是最大的数。\n\n针对所有的元素重复以上的步骤，除了最后一个。\n\n持续每次对越来越少的元素重复上面的步骤，直到没有任何一对数字需要比较。\n\n通过两层循环控制：\n\n\\- 第一个循环（外循环），负责把需要冒泡的那个数字排除在外；\n\n\\- 第二个循环（内循环），负责两两比较交换。\n\n**动图演示bubble_sort**\n\n![img](http://tukuimg.bdstatic.com/scrop/f3d6d6b1b26cc9d7a4ff1dc840537b77.gif)\n\n算法：\n\n\n```c\n    #include <stdio.h>\n    int main()\n    {\n      int a[100],i,j,t,n;\n        scanf(\"%d\",&n);  //输入一个数n，表示接下来有n个数\n        for(i=1;i<=n;i++)  //循环读入n个数到数组a中\n            scanf(\"%d\",&a[i]);\n        //冒泡排序的核心部分\n        for(i=1;i<=n-1;i++) //n个数排序，只用进行n-1趟\n        {\n            for(j=1;j<=n-i;j++) //从第1位开始比较直到最后一个尚未归位的数，想一想为什么到n-i就可以了。\n            {\n                if(a[j]<a[j+1]) //比较大小并交换\n                {  t=a[j]; a[j]=a[j+1]; a[j+1]=t;  }\n            }\n        }\n        for(i=1;i<=n;i++)  //输出结果\n            printf(\"%d \",a[i]);\n        getchar();getchar();\n        return 0;\n    }\n```\n\n**4.总结**\n\n冒泡排序毕竟是一种效率低下的排序方法，在数据规模很小时，可以采用。数据规模比较大时，建议采用其它排序方法。其他相关排序算法会在后续系列中逐一来分析说明，敬请关注哦！\n\n# 算法 3：最常用的排序——快速排序\n\n上一节的冒泡排序可以说是我们学习第一个真正的排序算法，并且解决了桶排序浪费空间的问题，但在算法的执行效率上却牺牲了很多，它的时间复杂度达到了 **O(N2)**。假如我们的计算机每秒钟可以运行 **10** 亿次，那么对 **1** 亿个数进行排序，桶排序则只需要 **0.1** 秒，而冒泡排序则需要 **1** 千万秒，达到 **115** 天之久，是不是很吓人。那有没有既不浪费空间又可以快一点的排序算法呢？那就是“快速排序”啦！光听这个名字是不是就觉得很高端呢。\n\n假设我们现在对“**6 1 2 7 9 3 4 5 10 8**”这个 10 个数进行排序。首先在这个序列中随便找一个数作为基准数（不要被这个名词吓到了，就是一个用来参照的数，待会你就知道它用来做啥的了）。为了方便，就让第一个数 **6** 作为基准数吧。接下来，需要将这个序列中所有比基准数大的数放在 **6** 的右边，比基准数小的数放在 **6** 的左边，类似下面这种排列。\n\n3 1 2 5 4 **6** 9 7 10 8\n\n在初始状态下，数字 **6** 在序列的第 **1** 位。我们的目标是将 **6** 挪到序列中间的某个位置，假设这个位置是 **k**。现在就需要寻找这个 **k**，并且以第 **k** 位为分界点，左边的数都小于等于 **6**，右边的数都大于等于 **6**。想一想，你有办法可以做到这点吗？\n\n给你一个提示吧。请回忆一下冒泡排序，是如何通过“交换”，一步步让每个数归位的。此时你也可以通过“交换”的方法来达到目的。具体是如何一步步交换呢？怎样交换才既方便又节省时间呢？先别急着往下看，拿出笔来，在纸上画画看。我高中时第一次学习冒泡排序算法的时候，就觉得冒泡排序很浪费时间，每次都只能对相邻的两个数进行比较，这显然太不合理了。于是我就想了一个办法，后来才知道原来这就是“快速排序”，请允许我小小的自恋一下(^o^)。\n\n方法其实很简单：分别从初始序列“**6 1 2 7 9 3 4 5 10 8**”两端开始“探测”。先从**右**往**左**找一个小于 **6** 的数，再从**左**往**右**找一个大于 **6** 的数，然后交换他们。这里可以用两个变量 **i** 和 **j**，分别指向序列最左边和最右边。我们为这两个变量起个好听的名字“哨兵 i”和“哨兵 j”。刚开始的时候让哨兵 i 指向序列的最左边（即 **i=1**），指向数字 **6**。让哨兵 **j** 指向序列的最右边（即 **j=10**），指向数字 **8**。\n\n![picture3.1](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/3.1.png)\n\n首先哨兵 **j** 开始出动。因为此处设置的基准数是最左边的数，所以需要让哨兵 **j** 先出动，这一点非常重要（请自己想一想为什么）。哨兵 **j** 一步一步地向左挪动（即 **j--**），直到找到一个小于 **6** 的数停下来。接下来哨兵 **i** 再一步一步向右挪动（即 **i++**），直到找到一个数大于 **6** 的数停下来。最后哨兵 **j** 停在了数字 **5** 面前，哨兵 **i** 停在了数字 **7** 面前。\n\n![picture3.2](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/3.2.png)\n\n![picture3.3](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/3.3.png)\n\n现在交换哨兵 **i** 和哨兵 **j** 所指向的元素的值。交换之后的序列如下。\n\n6 1 2 **5** 9 3 4 **7** 10 8\n\n![picture3.4](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/3.4.png)\n\n![picture3.5](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/3.5.png)\n\n到此，第一次交换结束。接下来开始哨兵 **j** 继续向左挪动（再友情提醒，每次必须是哨兵 **j** 先出发）。他发现了 **4**（比基准数 **6** 要小，满足要求）之后停了下来。哨兵 **i** 也继续向右挪动的，他发现了 **9**（比基准数 **6** 要大，满足要求）之后停了下来。此时再次进行交换，交换之后的序列如下。\n\n6 1 2 5 **4** 3 **9** 7 10 8\n\n第二次交换结束，“探测”继续。哨兵 **j** 继续向左挪动，他发现了 **3**（比基准数 **6** 要小，满足要求）之后又停了下来。哨兵 **i** 继续向右移动，糟啦！此时哨兵 **i** 和哨兵 **j** 相遇了，哨兵 **i** 和哨兵 **j** 都走到 **3** 面前。说明此时“探测”结束。我们将基准数 **6** 和 **3** 进行交换。交换之后的序列如下。\n\n**3** 1 2 5 4 **6** 9 7 10 8\n\n![picture3.6](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/3.6.png)\n\n![picture3.7](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/3.7.png)\n\n![picture3.8](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/3.8.png)\n\n到此第一轮“探测”真正结束。此时以基准数 **6** 为分界点，**6** 左边的数都小于等于 **6**，**6** 右边的数都大于等于 **6**。回顾一下刚才的过程，其实哨兵 **j** 的使命就是要找小于基准数的数，而哨兵 **i** 的使命就是要找大于基准数的数，直到 **i** 和 **j** 碰头为止。\n\nOK，解释完毕。现在基准数 **6** 已经归位，它正好处在序列的第 **6** 位。此时我们已经将原来的序列，以 **6** 为分界点拆分成了两个序列，左边的序列是“**3 1 2 5 4**”，右边的序列是“ **9 7 10 8** ”。接下来还需要分别处理这两个序列。因为 **6** 左边和右边的序列目前都还是很混乱的。不过不要紧，我们已经掌握了方法，接下来只要模拟刚才的方法分别处理 **6** 左边和右边的序列即可。现在先来处理 **6** 左边的序列现吧。\n\n左边的序列是“**3 1 2 5 4**”。请将这个序列以 **3** 为基准数进行调整，使得 **3** 左边的数都小于等于 **3**，**3** 右边的数都大于等于 **3**。好了开始动笔吧。\n\n如果你模拟的没有错，调整完毕之后的序列的顺序应该是。\n\n2 1 **3** 5 4\n\nOK，现在 **3** 已经归位。接下来需要处理 **3** 左边的序列“ **2 1** ”和右边的序列“**5 4**”。对序列“ **2 1** ”以 **2** 为基准数进行调整，处理完毕之后的序列为“**1 2**”，到此 **2** 已经归位。序列“**1**”只有一个数，也不需要进行任何处理。至此我们对序列“ **2 1** ”已全部处理完毕，得到序列是“**1 2**”。序列“**5 4**”的处理也仿照此方法，最后得到的序列如下。\n\n1 2 3 4 5 6 9 7 10 8\n\n对于序列“**9 7 10 8**”也模拟刚才的过程，直到不可拆分出新的子序列为止。最终将会得到这样的序列，如下。\n\n1 2 3 4 5 6 7 8 9 10\n\n到此，排序完全结束。细心的同学可能已经发现，快速排序的每一轮处理其实就是将这一轮的基准数归位，直到所有的数都归位为止，排序就结束了。下面上个霸气的图来描述下整个算法的处理过程。\n\n![picture3.9](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/images/3.9.png)\n\n快速排序之所比较快，因为相比冒泡排序，每次交换是跳跃式的。每次排序的时候设置一个基准点，将小于等于基准点的数全部放到基准点的左边，将大于等于基准点的数全部放到基准点的右边。这样在每次交换的时候就不会像冒泡排序一样每次只能在相邻的数之间进行交换，交换的距离就大的多了。因此总的比较和交换次数就少了，速度自然就提高了。当然在最坏的情况下，仍可能是相邻的两个数进行了交换。因此快速排序的最差时间复杂度和冒泡排序是一样的都是 **O(N2)**，它的平均时间复杂度为 **O(NlogN)**。其实快速排序是基于一种叫做“二分”的思想。我们后面还会遇到“二分”思想，到时候再聊。先上代码，如下。\n\n```java\nimport java.util.*;\npublic class Sort {\n    public static void main(String[] args) {\n        int[] array = new int[]{6 ,1 ,2 ,7 ,9 ,3 ,4 ,5 ,10 ,8};\n        System.out.println(\"-->\"+Arrays.toString(array));\n        quickSort(array, 0, array.length - 1);\n        System.out.println(\"-->\"+Arrays.toString(array));\n        System.out.println(Arrays.toString(array));\n    }\n    /**\n     * 快速排序：换位法\n     */\n    public static void quickSort(int[] array, int i, int j) {\n        int start=i,end=j;\n        int value = array[start];//哨兵\n        if (i >= j)\n            return;//递归出口\n\n        while (i < j) {\n            //从右边找比value小的元素\n            while (i < j && array[j] >= value) {\n                j--;\n            }\n            //从左边找比value大的元素\n            while (i < j && array[i] <= value) {\n                i++;\n            }\n            //交换位置\n            if(i<j){\n                System.out.println(String.format(\"array[i]:%d, array[j]:%d,\", array[i], array[j]));\n                int t=array[i];\n                array[i]=array[j];\n                array[j]=t;\n                System.out.println(\"-->\"+Arrays.toString(array));\n            }\n        }\n        //交换哨兵元素与相遇位置：一遍交换位置之后i一定等于j,且i的位置是右边第一个比value小的元素，交换位置即可\n        array[start]=array[i];\n        array[i] = value;// 放置哨兵元素\n        System.out.println(\"-->\"+Arrays.toString(array));\n        System.out.println(String.format(\",i:%d, j:%d, start:%d, end:%d\", i, j,start,end));\n        quickSort(array, start, i-1);//对基准元素左边的元素进行递归排序\n        quickSort(array, i+1, end);//对基准元素右边的元素进行递归排序\n    }\n}\n```\n\n可以输入以下数据进行验证\n\n```sh\n1061279345108\n```\n\n运行结果是\n\n```sh\n12345678910\n```\n\n下面是程序执行过程中数组 **a** 的变化过程，带下划线的数表示的已归位的基准数。\n\n```sh\n    1 2 7 9 3 4 5 10 8\n    1 2 5 4 6 9 7 10 8\n    1 3 5 4 6 9 7 10 8\n    2 3 5 4 6 9 7 10 8\n    2 3 5 4 6 9 7 10 8\n    2 3 4 5 6 9 7 10 8\n    2 3 4 5 6 9 7 10 8\n    2 3 4 5 6 8 7 9 10\n    2 3 4 5 6 7 8 9 10\n    2 3 4 5 6 7 8 9 10\n    2 3 4 5 6 7 8 9 10\n```\n\n快速排序除了交换两个数的位置之外，还有一种拆补法。同样可以达比基准元素大的元素放在它的右边，比其小的放在它的左边的效果。\n\n运行过程如下：\n![42EF8E48D02EE22C29D77D97C1723A15.jpg]( https://veport.oss-cn-beijing.aliyuncs.com/articles/42ef8e48d02ee22c29d77d97c1723a15.jpg)\n```Java\n    /**\n     * 快速排序：拆补法\n     * 把快速排序联想成东拆西补或西拆东补，一边拆一边补，直到所有元素达到有序状态。\n     */\n    public static void quickSort1(int[] array, int i, int j) {\n        int start=i,end=j;\n        int value = array[i];//哨兵\n        if (i >= j)\n            return;//递归出口\n\n        while (i < j) {\n            //从右往左扫描，找到第一个比基准元素小的元素\n            while (i < j && array[j] >= value) {\n                j--;\n            }\n            //找到这种元素arr[right]后与arr[left]交换\n            array[i] = array[j];\n            System.out.println(Arrays.toString(array));\n\n            //从左往右扫描，找到第一个比基准元素大的元素\n            while (i < j && array[i] <= value) {\n                i++;\n            }\n            //找到这种元素arr[left]后，与arr[right]交换\n            array[j] = array[i];\n            System.out.println(Arrays.toString(array));\n        }\n        //哨兵元素归位\n        array[i] = value;\n        System.out.println(\"-->\"+Arrays.toString(array));\n        System.out.println(String.format(\",i:%d, j:%d, start:%d, end:%d\", i, j,start,end));\n        quickSort1(array, start, i-1);//对基准元素左边的元素进行递归排序\n        quickSort1(array, i+1, end);//对基准元素右边的元素进行递归排序\n    }\n```\n\n\n快速排序由 C. A. R. Hoare（东尼霍尔，Charles Antony Richard Hoare）在 1960 年提出，之后又有许多人做了进一步的优化。如果你对快速排序感兴趣可以去看看东尼霍尔 1962 年在 Computer Journal 发表的论文“Quicksort”以及《算法导论》的第七章。快速排序算法仅仅是东尼霍尔在计算机领域才能的第一次显露，他在 1980 年获得了图灵奖。\n\n**时间复杂度**\n&emsp;&emsp;快速排序涉及到递归调用，所以该算法的时间复杂度还需要从递归算法的复杂度开始说起；\n&emsp;&emsp;递归算法的时间复杂度公式：T[n] = aT[n/b] + f(n)  ；对于递归算法的时间复杂度这里就不展开来说了；\n\n**最优情况下时间复杂度**\n&emsp;&emsp;快速排序最优的情况就是每一次取到的元素都刚好平分整个数组(很显然我上面的不是)；\n&emsp;&emsp;此时的时间复杂度公式则为：T[n] = 2T[n/2] + f(n)；T[n/2]为平分后的子数组的时间复杂度，f[n] 为平分这个数组时所花的时间；\n\n综上所述：快速排序最优的情况下时间复杂度为：O( nlogn )\n\n\n\n**最差情况下时间复杂度**\n&emsp;&emsp;最差的情况就是每一次取到的元素就是数组中最小/最大的，这种情况其实就是冒泡排序了(每一次都排好一个元素的顺序)。\n&emsp;&emsp;这种情况时间复杂度就好计算了，就是冒泡排序的时间复杂度：T[n] = n * (n-1) = n^2 + n;\n\n 综上所述：快速排序最差的情况下时间复杂度为：O( n^2 )\n\n**平均时间复杂度**\n       快速排序的平均时间复杂度也是：O(nlogn)\n\n**空间复杂度**\n&emsp;&emsp;其实这个空间复杂度不太好计算，因为有的人使用的是非就地排序，那样就不好计算了（因为有的人用到了辅助数组，所以这就要计算到你的元素个数了）；我就分析下就地快速排序的空间复杂度吧；\n&emsp;&emsp;首先就地快速排序使用的空间是O(1)的，也就是个常数级；而真正消耗空间的就是递归调用了，因为每次递归就要保持一些数据；\n     最优的情况下空间复杂度为：O(logn)  ；每一次都平分数组的情况\n     最差的情况下空间复杂度为：O( n )      ；退化为冒泡排序的情况\n\n\n参考链接：\n\n[最快最简单的排序——桶排序](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/bucket-sort.html)  \n\n[【排序】图解桶排序](https://blog.csdn.net/qq_27124771/article/details/87651495)  \n\n[排序算法之冒泡排序，轻松追求“有序”之美](https://baijiahao.baidu.com/s?id=1663297052531965367&wfr=spider&for=pc)  \n\n[最常用的排序——快速排序](https://wiki.jikexueyuan.com/project/easy-learn-algorithm/fast-sort.html)\n\n', 1, '', 0, 0, 1, '2022-02-19 10:58:00', '2022-03-02 04:12:36');
INSERT INTO `article` VALUES (79, 2, 193, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/f9fa18da262910eb13f802b003147915.jpg', '面试杂谈', '## [浅谈C++中指针和引用的区别](https://www.cnblogs.com/dolphin0520/archive/2011/04/03/2004869.html)\n\n指针和引用在C++中很常用，但是对于它们之间的区别很多初学者都不是太熟悉，下面来谈谈他们2者之间的区别和用法。\n\n1.指针和引用的定义和性质区别：\n\n(1)指针：指针是一个变量，只不过这个变量存储的是一个地址，指向内存的一个存储单元；而引用跟原来的变量实质上是同一个东西，只不过是原变量的一个别名而已。如：\n\nint a=1;int *p=&a;\n\nint a=1;int &b=a;\n\n上面定义了一个整形变量和一个指针变量p，该指针变量指向a的存储单元，即p的值是a存储单元的地址。\n\n而下面2句定义了一个整形变量a和这个整形a的引用b，事实上a和b是同一个东西，在内存占有同一个存储单元。\n\n(2)可以有const指针，但是没有const引用；\n\n(3)指针可以有多级，但是引用只能是一级（int **p；合法 而 int &&a是不合法的）\n\n(4)指针的值可以为空，但是引用的值不能为NULL，并且引用在定义的时候必须初始化；\n\n(5)指针的值在初始化后可以改变，即指向其它的存储单元，而引用在进行初始化后就不会再改变了。\n\n(6)\"sizeof引用\"得到的是所指向的变量(对象)的大小，而\"sizeof指针\"得到的是指针本身的大小；\n\n(7)指针和引用的自增(++)运算意义不一样；\n## C/C++之内存四区的详细介绍\nC/C++编译器会把代码直接分为四个小区：\n\n(1).栈区(stack)：由编译器自动分配和释放，存放函数的参数值，局部变量的值等。\n\n(2).堆区(heap)：一般是由程序员分配释放(动态内存申请与释放)，若程序员不释放，程序结束时可能由操作系统回收\n\n(3).全局区(静态区，常量区)：全局变量和静态变量的存储是放在一起的，初始化的全局变量和静态变量在一块区域，未初始化的全局变量和未初始化的静态变量在相邻的另一块区域，该区域在程序结束后由操作系统释放。常量区是字符串常量和其他常量的存储位置，程序结束后由系统自动释放。\n\n(4).程序代码区：存放函数体的二进制代码。\n\n## qt/C++从源文件到可执行文件的编译过程\n\n从源文件(.c，.cpp，.h)到可执行文件(.exe，.dll——不是只有exe才叫可执行文件)依次经历下面几个过程\n\n![img](https://img-blog.csdnimg.cn/20191219211854678.png)\n\n### **预处理阶段**\n\n   尽管现在编译器都包含了预处理器，但是通常预处理是独立编译阶段的。也有的称这一阶段为预编译阶段。\n\n 预处理主要对条件编译指令及对宏定义的展开(替换)和对#include的处理等(注意：保留所有的#pragma编译指令，因为编译器需要使用它们)，同时也会删除程序中的注释和多余的空白符，经过预处理后生成一个没有宏定义、条件编译指令、没有特殊符号的文件，这个文件与与源文件并无本质区别\n\n### **编译阶段**\n\n  编译阶段是将.i文件生成机器语言.s。这一阶段主要是与语法上的检查和代码优化。语法上的检查有静态和动态。对于静态检查是指不用经过运算即可完成的检查，如类型转换；动态检查是指在运行阶段才能完成的，如除0操作。\n\n**优化阶段**\n\n  现在的优化主要有两类，一类是与硬件相关，另一类是与硬件无关。与硬件相关的优化是：借助硬件的性能，减少内存访问次数以及硬件执行指令的特点（如流水线、RISC、CISC、VLIW等）而对指令进行一些调整使目标代码比较短，从而提高效率；与硬件无关主要表现在优化循环(削弱强度，代码外提)、删除无用变量等\n\n### **汇编阶段**\n\n  汇编阶段就是将机器语言转换成机器指令，如mov指令等，最终生成.o文件\n\n### **链接阶段**\n\n  链接阶段就是将各个.o文件链接成可执行文件，.o是各个独立的模块。其中链接分为静态和动态链接。静态和动态区别主要表现在几个方面：\n\n 1、链接时机。静态是在形成可执行文件前就被安排的妥妥的了，动态是在执行时才知道谁找谁\n\n 2、内存持有方式。静态是对于每个需要依赖的地方都copy一份副本，而动态则是内存中只保留一份副本，大家共享这份副本。举例子说就是：静态就是老婆，人人手中有一辆；动态就是公交车，就那么一辆，谁用谁来取\n\n 3、升级影响。静态既然是在需要的地方都会保留一份副本，当所依赖的.o改变时，就会牵一发动全身；而动态则不会，默默的做好自己就行了。\n\n \n\n那么,qt的编译过程是怎样的呢?\n\n**qt的编译经过两个阶段,qmake和make**.\n\n**qmake阶段**.使用qmake.exe生成MakeFile、Makefile.Debug和Makefile.Release。也就是一个debug和release版本的Makefile文件。打开任意版本的Makefile可以看到，里面定义了各种语句，类似moc_xxx.cpp等。其实也就是说通过qmake.exe创建了一个依赖规则的文件，让编译器根据这个规则进行链接和生成目标文件\n\n**make阶段**.这一阶段其实就是普通的make文件的过程，具体可以看上一部分\n\n## ICMP协议是什么协议？\n\n### **ICMP是啥？**\n\nICMP，全称为Internet Control Message Protocol，即为因特网控制报文协议。\n\n### **ICMP解决了啥？**\n\nIP协议本身即没有为终端系统提供直接的方式来发现那些发往目的地址失败的IP数据包，也没有提供直接的方式来获取诊断信息。而ICMP就是为了解决IP协议的不足引入的网络协议。\n\n### **ICMP属于TCP/IP四层模型的哪一层？**\n\nICMP既不是网络层协议，也不是传输层协议，但是通常ICMP被认为是IP层的一部分。\n\n### **ICMP长啥样？**\n\n以ICMP分为ICMPv4和ICMPv6版本，分别对应IPv4协议和IPv6协议。以简单的ICMPv4为例，一条ICMPv4消息分为IPv4头部（20～60字节）、ICMP头部（又分为1字节的类型位、1字节的代码位、2字节的校验和）以及ICMP数据，其中IPv4头部中的协议位数值为1。\n![img](https://imgconvert.csdnimg.cn/aHR0cHM6Ly9tbWJpei5xcGljLmNuL21tYml6X3BuZy9KRHRiSG5ySFJuUGJpY0JVU21pYUs0anhJWGdTRDZGcmozdnAxSzZIVGJiM0NYbG10VFd2czluZ3JxZUFreTNtZEtYMmwxekNCZWg3cHhiTEpqT2V0cElRLzY0MA?x-oss-process=image/format,png)\n\n## 路由器和交换机的区别\n\n一、指代不同\n\n1、路由器：是连接两个或多个网络的硬件设备，在网络间起网关的作用，是读取每一个数据包中的地址然后决定如何传送的专用智能性的网络设备。\n\n2、交换机：是一种用于电（光）信号转发的网络设备。\n\n二、功能不同\n\n1、路由器：最主要的功能可以理解为实现信息的转送。把这个过程称之为寻址过程。因为在路由器处在不同网络之间，但并不一定是信息的最终接收地址。所以在路由器中, 通常存在着一张路由表。\n\n2、交换机：交换机有带宽很高的内部交换矩阵和背部总线，并且这个背部总线上挂接了所有的端口，通过内部交换矩阵，就能够把数据包直接而迅速地传送到目的节点而非所有节点， 这样就不会浪费网络资源，从而产生非常高的效率。\n\n三、特点不同\n\n1、路由器：核心是背板，高效率的背板有助于提高路由器的性能。由于传统的共享总线式背板无法满足路由器的需要，所以采用结构可以用不同技术实现的交换式背板。\n\n2、交换机：交换机在同一时刻可进行多个端口对之间的数据传输。每一端口都可视为独立的物理网段（注：非IP网段），连接在其上的网络设备独自享有全部的带宽，无须同其他设备竞争使用。\n\n参考资料来源：[百度百科-交换机](https://baike.baidu.com/item/交换机/103532?fr=aladdin)\n\n参考资料来源：[百度百科-路由器](https://baike.baidu.com/item/路由器/108294?fr=aladdin)\n\n## 区块链相关疑问解析\nhttps://www.runoob.com/w3cnote/question-and-analysis.html', 1, '', 0, 0, 2, '2022-02-19 11:44:02', '2022-03-03 23:03:20');
INSERT INTO `article` VALUES (80, 2, 191, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/4e319068b295ca52080979d5653c334d.jpg', '网站功能展示', '网站前端vue和后端springboot由IDEA编译打包，并发布到服务器。\n\n![code.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/5a0632f5fac4b8ce373a069253f9d1e2.png)\n\n## 前台页面展示\n### 移动端适配\n<img src=\" https://veport.oss-cn-beijing.aliyuncs.com/articles/3d6d442430f9de01c2e7cbcb914c1a2a.jpg\" alt=\"1754BCE1444B06E890FFAC61A1BF9BD5.jpg\" style=\"zoom: 20%;\" />\n&emsp;\n<img src=\" https://veport.oss-cn-beijing.aliyuncs.com/articles/12431b4262b7e9c2fe1c10d2a9797ad1.jpg\" alt=\"12431B4262B7E9C2FE1C10D2A9797AD1.jpg\" style=\"zoom:20%;\" />\n\n<img src=\" https://veport.oss-cn-beijing.aliyuncs.com/articles/d4febe8505c3e3d9f71493972b0cd6ed.jpg\" alt=\"D4FEBE8505C3E3D9F71493972B0CD6ED.jpg\" style=\"zoom:20%;\" />\n&emsp;\n<img src=\" https://veport.oss-cn-beijing.aliyuncs.com/articles/c41b4df3ead2d3dd6684a363e8ea4106.png\" alt=\"F949FB50529649899AF2093D2D61FEA0.png\" style=\"zoom:20%;\" />\n\n### PC端适配\n![blog1.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/b0321e4e4833f12250a2d0bbaf40889b.png)\n![blog2.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/96fb1393414a402d7e7d5e46f1fb6204.png)\n![blog4.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/9e7675104afe9bf523aa7862803bde4f.png)\n![blog5.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/28d8cefb6fc706ad673ef8656a631ecb.png)\n\n## 后台管理系统展示\n![admin1.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/37119259612910a845c10517c4c82900.png)![admin2.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/106c49bef28297ffe184257e6bae8c5e.png)![admin3.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/ba11c04eda104aeb922ac145fc40a983.png)![admin4.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/95ac11550856e34974fdeaf245bb7f81.png)![admin5.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/07e6cafa2ebae034501d58e4a4bb6d3c.png)![admin6.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/b27720c9fb435206db88491ca0c1a117.png)', 1, '', 0, 1, 1, '2022-02-19 22:49:28', '2022-02-19 22:51:34');
INSERT INTO `article` VALUES (81, 2, 188, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/f9fa18da262910eb13f802b003147915.jpg', 'Java 线程池', '### 1 ThreadPoolExecutor 的构造方法中各个参数的含义是什么？\n\n```java\npublic ThreadPoolExecutor(int corePoolSize,						// 第 1 个参数\n                          int maximumPoolSize,					// 第 2 个参数\n                          long keepAliveTime,					// 第 3 个参数\n                          TimeUnit unit,						// 第 4 个参数\n                          BlockingQueue<Runnable> workQueue,	// 第 5 个参数\n                          ThreadFactory threadFactory,			// 第 6 个参数\n                          RejectedExecutionHandler handler) {   // 第 7 个参数\n\n```\n\n| 序号 |      参数名       |          参数类型          | 参数含义                           |        取值范围         | 解释说明                       |\n| ---- | :---------------: | :------------------------: | ---------------------------------- | :---------------------: | :----------------------------------------------------------- |\n| 1    |  `corePoolSize`   |           `int`            | 核心线程数                         |           >=0           | 如果等于 0，则任务执行完成后，没有任何请求进入时就销毁线程池的线程；如果大于 0，即使本地任务执行完毕，核心线程也不会销毁。这个值的设置非常关键，设置过大会浪费资源，设置过小会导致线程频繁地创建或销毁。 |\n| 2    | `maximumPoolSize` |           `int`            | 线程池能够容纳同时执行的最大线程数 | >0并且>= `corePoolSize` | 如果任务数超过第 5 个参数 `workQueue` 的任务缓存上限且待执行的线程数小于 `maximumPoolSize` 时，需要借助第 5 个参数的帮助，缓存在队列中。如果 `maximumPoolSize` 与 `corePoolSize` 相等，则是固定大小线程池。**最大线程数 = 核心线程数 + 非核心线程数。** |\n| 3    |  `keepAliveTime`  |           `long`           | 线程池中的线程空闲时间             |           >=0           | 当空闲时间达到 `keepAliveTime` 值时，非核心线程会被销毁，直到只剩下 `corePoolSize` 个线程为止，避免浪费内存和句柄资源。在默认情况下，当线程池的线程数大于 `corePoolSize` 时，`keepAliveTime` 才会起作用。但是当 `ThreadPoolExecutor` 的 `allowCoreThreadTimeOut` 变量设置为 `true` 时，核心线程超时后也会被回收。 |\n| 4    |      `unit`       |         `TimeUnit`         | 时间单位                           |                         | `keepAliveTime` 的时间单位通常是 `TimeUnit.SECONDS`。        |\n| 5    |    `workQueue`    | `BlockingQueue<Runnable>`  | 任务缓存队列                       |     不可以为 `null`     | 当请求的线程数大于等于 `corePoolSize` 时，任务会进入 `BlockingQueue` 阻塞队列等待执行。 |\n| 6    |  `threadFactory`  |      `ThreadFactory`       | 线程工厂                           |     不可以为 `null`     | 用来生成一组相同任务的线程。线程池的命名是通过给这个 factory 增加组名前缀来实现的。在虚拟机栈分析时，就可以知道线程任务是由哪个线程工厂产生的。 |\n| 7    |     `handler`     | `RejectedExecutionHandler` | 执行拒绝策略的对象                 |     不可以为 `null`     | 当待执行的线程数大于等于 `maximumPoolSize` 时，就可以通过该策略处理请求，这是一种简单的限流保护。 |\n\n### 2 ThreadPoolExecutor 执行任务的规则是什么？\n\n1. 如果线程池中的线程数量未达到核心线程的数量，那么会直接启动一个核心线程来执行任务；\n2. 如果线程池中的线程数量已经达到或者超过核心线程的数量，那么任务会被插入到任务队列中排队等待执行；\n3. 如果在步骤 2 中无法将任务插入到任务队列中（原因往往是任务队列已满），这个时候如果线程数量未达到线程池规定的最大值，那么会启动一个非核心线程来执行任务；\n4. 如果步骤 3 中线程数量已经达到线程池规定的最大值，那么就拒绝执行此任务。\n\n绘制流程图如下：\n\n![在这里插入图片描述](https://img-blog.csdnimg.cn/710c04bdff7d4cc8b711d3c1b5964d01.png?x-oss-process=image/watermark,type_d3F5LXplbmhlaQ,shadow_50,text_Q1NETiBAd2lsbHdheXdhbmc2,size_16,color_FFFFFF,t_70,g_se,x_16)\n\n### 3 Executors 里的线程池有哪些？这些线程池有什么特点？\n\n| 比较项          | newCachedThreadPool            | newFixedThreadPool(int nThreads)                             | newSingleThreadExecutor        | newScheduledThreadPool(int corePoolSize)                     |\n| --------------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |\n| corePoolSize    | `0`                            | `nThreads`                     | 1    | corePoolSize                   |\n| maximumPoolSize | `Integer.MAX_VALUE`            | `nThreads`                     | 1    | Integer.MAX_VALUE              |\n| keepAliveTime   | `60L`                          | `0L`                           | `0L`                           | 0    |\n| unit            | `TimeUnit.SECONDS`             | `TimeUnit.MILLISECONDS`        | `TimeUnit.MILLISECONDS`        | TimeUnit.NANOSECONDS           |\n| workQueue       | `new SynchronousQueue<Runnable>()`                           | `new LinkedBlockingQueue<Runnable>()`                        | `new LinkedBlockingQueue<Runnable>()`                        | new DelayedWorkQueue()         |\n| 线程池名称      | 无界线程池，可以进行自动线程回收                             | 固定大小线程池                 | 单线程线程池                   | 执行定时任务，重复任务线程池       |\n| 特点            | 没有核心线程，只有非核心线程（最大为`Integer.MAX_VALUE`），超过 60s 的空闲线程会被回收，SynchronousQueue 会将任务直接提交给线程而不保持它们，所以任务会立即执行。 | 只有固定个数的核心线程，没有非核心线程，核心线程不会被回收，任务队列大小没有限制。当线程处于空闲状态时，不会被回收；当所有的线程都处于活动状态时，新任务到达就处于等待状态，直到有线程空闲出来。 | 只有一个核心线程，没有非核心线程，核心线程不会被回收，任务队列大小没有限制。可以保证所有的任务都在同一个线程中按顺序执行。 | 核心线程数是固定的，非核心线程数是没有限制的，非核心线程空闲时会被立即回收。 |\n\n\n\n\n\n### 4 最后\n\n对线程池的简单理解：\n\n　　假如有一个工厂，工厂里面有10个工人，每个工人同时只能做一件任务。\n\n　　因此只要当10个工人中有工人是空闲的，来了任务就分配给空闲的工人做；\n\n　　当10个工人都有任务在做时，如果还来了任务，就把任务进行排队等待；\n\n　　如果说新任务数目增长的速度远远大于工人做任务的速度，那么此时工厂主管可能会想补救措施，比如重新招4个临时工人进来；\n\n　　然后就将任务也分配给这4个临时工人做；\n\n　　如果说着14个工人做任务的速度还是不够，此时工厂主管可能就要考虑不再接收新的任务或者抛弃前面的一些任务了。\n\n　　当这14个工人当中有人空闲时，而新任务增长的速度又比较缓慢，工厂主管可能就考虑辞掉4个临时工了，只保持原来的10个工人，毕竟请额外的工人是要花钱的。\n\n \n\n　　这个例子中的corePoolSize就是10，而maximumPoolSize就是14（10+4）。\n\n　　也就是说corePoolSize就是线程池大小，maximumPoolSize在我看来是线程池的一种补救措施，即任务量突然过大时的一种补救措施。\n\n### 参考\n\n- [《Android开发艺术探索》第11章-Android 的线程和线程池](https://blog.csdn.net/willway_wang/article/details/122632665?spm=1001.2014.3001.5502)\n- [ Java并发编程：线程池的使用](https://www.cnblogs.com/dolphin0520/p/3932921.html)\n\n\n\n### 测试：\n```Java\npackage threadpool;\n\n\nimport java.util.concurrent.*;\n\n/**\n * @Description create for javaCourse .\n * 这里要重点解释一下corePoolSize、maximumPoolSize、largestPoolSize三个变量。\n *\n * 　　corePoolSize在很多地方被翻译成核心池大小，其实我的理解这个就是线程池的大小。举个简单的例子：\n * 　　假如有一个工厂，工厂里面有10个工人，每个工人同时只能做一件任务。\n *\n * 　　因此只要当10个工人中有工人是空闲的，来了任务就分配给空闲的工人做；\n *\n * 　　当10个工人都有任务在做时，如果还来了任务，就把任务进行排队等待；\n *\n * 　　如果说新任务数目增长的速度远远大于工人做任务的速度，那么此时工厂主管可能会想补救措施，比如重新招4个临时工人进来；\n *\n * 　　然后就将任务也分配给这4个临时工人做；\n *\n * 　　如果说着14个工人做任务的速度还是不够，此时工厂主管可能就要考虑不再接收新的任务或者抛弃前面的一些任务了。\n *\n * 　　当这14个工人当中有人空闲时，而新任务增长的速度又比较缓慢，工厂主管可能就考虑辞掉4个临时工了，只保持原来的10个工人，毕竟请额外的工人是要花钱的。\n *\n *\n *\n * 　　这个例子中的corePoolSize就是10，而maximumPoolSize就是14（10+4）。\n *\n * 　　也就是说corePoolSize就是线程池大小，maximumPoolSize在我看来是线程池的一种补救措施，即任务量突然过大时的一种补救措施。\n * @Author weiyi\n * @Date 2022/3/11\n */\npublic class ThreadPoolTest {\n    public static void main(String[] args) {\n\n        //60s不执行就删，不够再加，无上限\n        //cacheThreadPool();\n        //一直添加线程，直到保持在固定数量\n       // fixTheadPoolTest();\n        //核心线程只有一个，保证先进先出\n        //singleTheadPoolTest();\n        //固定大小线程池，定时任务：循环执行、延迟执行\n        //scheduleThreadPool();\n\n        newThreadPool();\n    }\n\n    public static void newThreadPool() {\n        /**\n         * ArrayBlockingQueue; 规定大小的BlockingQueue，其构造必须指定大小。其所含的对象是FIFO顺序排序的。\n         * LinkedBlockingQueue; 大小不固定的BlockingQueue，若其构造时指定大小，生成的BlockingQueue有大小限制，不指定大小，其大小有Integer.MAX_VALUE来决定。其所含的对象是FIFO顺序排序的。\n         * PriorityBlockingQueue：类似于LinkedBlockingQueue，但是其所含对象的排序不是FIFO，而是依据对象的自然顺序或者构造函数的Comparator决定。\n         * SynchronousQueue; 特殊的BlockingQueue，对其的操作必须是放和取交替完成。\n         */\n        BlockingQueue workQueue=new SynchronousQueue<Runnable>();\n\n        /**\n         * ThreadPoolExecutor.AbortPolicy:丢弃任务并抛出RejectedExecutionException异常。\n         * ThreadPoolExecutor.DiscardPolicy：也是丢弃任务，但是不抛出异常。\n         * ThreadPoolExecutor.DiscardOldestPolicy：丢弃队列最前面的任务，然后重新尝试执行任务（重复此过程）\n         * ThreadPoolExecutor.CallerRunsPolicy：由调用线程处理该任务\n         */\n        ThreadFactory threadFactory=new ThreadFactory() {\n            @Override\n            public Thread newThread(Runnable r) {\n               System.out.println(\"自定义线程工厂创建\"+r.toString());\n                Thread t = new Thread(r);\n                t.setDaemon(true);\n                return t;\n            }\n        };\n\n        RejectedExecutionHandler rejectedExecutionHandler=new ThreadPoolExecutor.CallerRunsPolicy();\n        ThreadPoolExecutor threadPoolExecutor = new ThreadPoolExecutor(\n                3,\n                5,\n                60L,\n                TimeUnit.SECONDS,\n                workQueue,\n                threadFactory,\n                rejectedExecutionHandler);\n\n        for (int i = 0; i < 5; i++) {\n            Runnable r1 = new Runnable() {\n                @Override\n                public void run() {\n                   System.out.println(\"Runnable 线程名称：\" + Thread.currentThread().getName() + \"，执行:3秒后执行\");\n                }\n            };\n\n            Callable c1=new Callable() {\n                @Override\n                public Object call() throws Exception {\n                    System.out.println(\"Callable 线程名称：\" + Thread.currentThread().getName() + \"，执行:3秒后执行\");\n                    return 1234567;\n                }\n            };\n            //execute无返回值 ——实现Runnable接口\n            threadPoolExecutor.execute(r1);\n            //submit有返回值 ——实现Callable接口\n           System.out.println(\"------->\"+threadPoolExecutor.submit(c1).toString());\n        }\n    }\n    /**\n     * newCachedThreadPool：\n     * <p>\n     * 底层：返回ThreadPoolExecutor实例，corePoolSize为0；\n     * maximumPoolSize为Integer.MAX_VALUE；keepAliveTime为60L；时间单位TimeUnit.SECONDS；\n     * workQueue为SynchronousQueue(同步队列)\n     *\n     * 通俗：当有新任务到来，则插入到SynchronousQueue中，由于SynchronousQueue是同步队列，因此会在池中寻找可用线程来执行，若有可以线程则执行，若没有可用线程则创建一个线程来执行该任务；若池中线程空闲时间超过指定时间，则该线程会被销毁。\n     * 适用：执行很多短期的异步任务\n     * 1.创建一个可缓存的线程池。如果线程池的大小超过了处理任务所需要的线程，那么就会回收部分空闲（60秒不执行任务）的线程<br>\n     * 2.当任务数增加时，此线程池又可以智能的添加新线程来处理任务<br>\n     * 3.此线程池不会对线程池大小做限制，线程池大小完全依赖于操作系统（或者说JVM）能够创建的最大线程大小<br>\n     */\n    public static void cacheThreadPool() {\n        ExecutorService cachedThreadPool = Executors.newCachedThreadPool();\n        for (int i = 1; i <= 10; i++) {\n            final int ii = i*10;\n            try {\n                Thread.sleep(ii);\n            } catch (InterruptedException e) {\n                e.printStackTrace();\n            }\n            cachedThreadPool.execute(new Runnable() {\n                @Override\n                public void run() {\n                   System.out.println(\"线程名称：\" + Thread.currentThread().getName() + \"，执行\" + ii);\n                }\n            });\n        }\n    }\n\n    /**\n     * newFixedThreadPool：\n     *\n     * 底层：返回ThreadPoolExecutor实例，接收参数为所设定线程数量n，\n     * corePoolSize和maximumPoolSize均为n；keepAliveTime为0L；\n     * 时间单位TimeUnit.MILLISECONDS；WorkQueue为：new LinkedBlockingQueue<Runnable>() 无界阻塞队列\n     *\n     * 通俗：创建可容纳固定数量线程的池子，每个线程的存活时间是无限的，当池子满了就不再添加线程了；如果池中的所有线程均在繁忙状态，对于新任务会进入阻塞队列中(无界的阻塞队列)\n     * 适用：执行长期任务\n     * 1.创建固定大小的线程池。每次提交一个任务就创建一个线程，直到线程达到线程池的最大大小<br>\n     * 2.线程池的大小一旦达到最大值就会保持不变，如果某个线程因为执行异常而结束，那么线程池会补充一个新线程<br>\n     * 3.因为线程池大小为3，每个任务输出index后sleep 2秒，所以每两秒打印3个数字，和线程名称<br>\n     */\n    public static void fixTheadPoolTest() {\n        ExecutorService fixedThreadPool = Executors.newFixedThreadPool(3);\n        for (int i = 0; i < 10; i++) {\n            final int ii = i;\n            fixedThreadPool.execute(() -> {\n               System.out.println(\"线程名称：\" + Thread.currentThread().getName() + \"，执行\" + ii);\n                try {\n                    Thread.sleep(2000);\n                } catch (InterruptedException e) {\n                    e.printStackTrace();\n                }\n            });\n        }\n    }\n\n    /**\n     * newSingleThreadExecutor:\n     *\n     * 底层：FinalizableDelegatedExecutorService包装的ThreadPoolExecutor实例，corePoolSize为1；maximumPoolSize为1；\n     * keepAliveTime为0L；时间单位TimeUnit.MILLISECONDS；workQueue为：\n     * new LinkedBlockingQueue<Runnable>() 无解阻塞队列\n     * 通俗：创建只有一个线程的线程池，当该线程正繁忙时，对于新任务会进入阻塞队列中(无界的阻塞队列)\n     * 适用：按顺序执行任务的场景\n     * 创建一个单线程化的线程池，它只会用唯一的工作线程来执行任务，保证所有任务按照指定顺序(FIFO, LIFO, 优先级)执行\n     */\n    public static void singleTheadPoolTest() {\n        ExecutorService pool = Executors.newSingleThreadExecutor();\n        for (int i = 0; i < 10; i++) {\n            final int ii = i;\n            pool.execute(() ->System.out.println(Thread.currentThread().getName() + \"=>\" + ii));\n        }\n    }\n\n    /**\n     * NewScheduledThreadPool:\n     *\n     * 底层：创建ScheduledThreadPoolExecutor实例，该对象继承了ThreadPoolExecutor，\n     * corePoolSize为传递来的参数，maximumPoolSize为Integer.MAX_VALUE；keepAliveTime为0；\n     * 时间单位TimeUnit.NANOSECONDS；workQueue为：new DelayedWorkQueue() 一个按超时时间升序排序的队列\n     *\n     * 通俗：创建一个固定大小的线程池，线程池内线程存活时间无限制，线程池可以支持定时及周期性任务执行，如果所有线程均处于繁忙状态，对于新任务会进入DelayedWorkQueue队列中，这是一种按照超时时间排序的队列结构\n     * 适用：执行周期性任务\n     * 创建一个定长线程池，支持定时及周期性任务执行。延迟执行\n     */\n    public static void scheduleThreadPool() {\n        ScheduledExecutorService scheduledThreadPool = Executors.newScheduledThreadPool(5);\n\n        Runnable r1 = new Runnable() {\n            @Override\n            public void run() {\n               System.out.println(\"线程名称：\" + Thread.currentThread().getName() + \"，执行:3秒后执行\");\n            }\n        };\n        scheduledThreadPool.schedule(r1, 3, TimeUnit.SECONDS);\n\n        Runnable r2 = new Runnable() {\n            @Override\n            public void run() {\n               System.out.println(\"线程名称：\" + Thread.currentThread().getName() + \"，执行:延迟2秒后每3秒执行一次\");\n            }\n        };\n        scheduledThreadPool.scheduleAtFixedRate(r2, 2, 3, TimeUnit.SECONDS);\n\n        Runnable r3 = new Runnable() {\n            @Override\n            public void run() {\n               System.out.println(\"线程名称：\" + Thread.currentThread().getName() + \"，执行:普通任务\");\n            }\n        };\n        for (int i = 0; i < 5; i++) {\n            scheduledThreadPool.execute(r3);\n        }\n    }\n}\n\n```\n\n', 1, '', 0, 0, 2, '2022-03-11 02:59:54', '2024-01-15 19:43:08');
INSERT INTO `article` VALUES (82, 1, 188, ' https://veport.oss-cn-beijing.aliyuncs.com/articles/f9fa18da262910eb13f802b003147915.jpg', 'HashMap实现原理及源码分析', '# [HashMap实现原理及源码分析](https://www.cnblogs.com/chengxiao/p/6059914.html)\n\n　　哈希表（hash table）也叫散列表，是一种非常重要的数据结构，应用场景及其丰富，许多缓存技术（比如memcached）的核心其实就是在内存中维护一张大的哈希表，而HashMap的实现原理也常常出现在各类的面试题中，重要性可见一斑。本文会对java集合框架中的对应实现HashMap的实现原理进行讲解，然后会对JDK7的HashMap源码进行分析。\n\n**目录**\n\n　　**一、[什么是哈希表](https://www.cnblogs.com/chengxiao/p/6059914.html#t1)**\n\n　　**二、[HashMap实现原理](https://www.cnblogs.com/chengxiao/p/6059914.html#t2)**\n\n　　**三、[为何HashMap的数组长度一定是2的次幂？](https://www.cnblogs.com/chengxiao/p/6059914.html#t3)**\n\n　　**四、[重写equals方法需同时重写hashCode方法](https://www.cnblogs.com/chengxiao/p/6059914.html#t4)**\n\n　　**五、[总结](https://www.cnblogs.com/chengxiao/p/6059914.html#t5)**\n\n# 一、什么是哈希表\n\n　　在讨论哈希表之前，我们先大概了解下其他数据结构在新增，查找等基础操作执行性能\n\n　　**数组**：采用一段连续的存储单元来存储数据。对于指定下标的查找，时间复杂度为O(1)；通过给定值进行查找，需要遍历数组，逐一比对给定关键字和数组元素，时间复杂度为O(n)，当然，对于有序数组，则可采用二分查找，插值查找，斐波那契查找等方式，可将查找复杂度提高为O(logn)；对于一般的插入删除操作，涉及到数组元素的移动，其平均复杂度也为O(n)\n\n　　**线性链表**：对于链表的新增，删除等操作（在找到指定操作位置后），仅需处理结点间的引用即可，时间复杂度为O(1)，而查找操作需要遍历链表逐一进行比对，复杂度为O(n)\n\n　　**二叉树**：对一棵相对平衡的有序二叉树，对其进行插入，查找，删除等操作，平均复杂度均为O(logn)。\n\n　　**哈希表**：相比上述几种数据结构，在哈希表中进行添加，删除，查找等操作，性能十分之高，不考虑哈希冲突的情况下，仅需一次定位即可完成，时间复杂度为O(1)，接下来我们就来看看哈希表是如何实现达到惊艳的常数阶O(1)的。\n\n　　我们知道，数据结构的物理存储结构只有两种：**顺序存储结构**和**链式存储结构**（像栈，队列，树，图等是从逻辑结构去抽象的，映射到内存中，也这两种物理组织形式），而在上面我们提到过，在数组中根据下标查找某个元素，一次定位就可以达到，哈希表利用了这种特性，**哈希表的主干就是数组**。\n\n　　比如我们要新增或查找某个元素，我们通过把当前元素的关键字 通过某个函数映射到数组中的某个位置，通过数组下标一次定位就可完成操作。**\n**\n\n　　　　　　　　**存储位置 = f(关键字)**\n\n　　其中，这个函数f一般称为**哈希函数**，这个函数的设计好坏会直接影响到哈希表的优劣。举个例子，比如我们要在哈希表中执行插入操作：\n\n![image.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/3c29c753f9584477b0f85483007b9d26.png)\n\n　　查找操作同理，先通过哈希函数计算出实际存储地址，然后从数组中对应地址取出即可。\n\n　　**哈希冲突**\n\n　　然而万事无完美，如果两个不同的元素，通过哈希函数得出的实际存储地址相同怎么办？也就是说，当我们对某个元素进行哈希运算，得到一个存储地址，然后要进行插入的时候，发现已经被其他元素占用了，其实这就是所谓的**哈希冲突**，也叫哈希碰撞。前面我们提到过，哈希函数的设计至关重要，好的哈希函数会尽可能地保证 **计算简单**和**散列地址分布均匀,**但是，我们需要清楚的是，数组是一块连续的固定长度的内存空间，再好的哈希函数也不能保证得到的存储地址绝对不发生冲突。那么哈希冲突如何解决呢？哈希冲突的解决方案有多种:开放定址法（发生冲突，继续寻找下一块未被占用的存储地址），再散列函数法，链地址法，而HashMap即是采用了链地址法，也就是**数组+链表**的方式，\n\n# 二、HashMap实现原理\n\n　HashMap的主干是一个Entry数组。Entry是HashMap的基本组成单元，每一个Entry包含一个key-value键值对。\n\n```\n//HashMap的主干数组，可以看到就是一个Entry数组，初始值为空数组{}，主干数组的长度一定是2的次幂，至于为什么这么做，后面会有详细分析。\ntransient Entry<K,V>[] table = (Entry<K,V>[]) EMPTY_TABLE;\n```\n\n Entry是HashMap中的一个静态内部类。代码如下\n\n```\n    static class Entry<K,V> implements Map.Entry<K,V> {\n        final K key;\n        V value;\n        Entry<K,V> next;//存储指向下一个Entry的引用，单链表结构\n        int hash;//对key的hashcode值进行hash运算后得到的值，存储在Entry，避免重复计算\n\n        /**\n         * Creates new entry.\n         */\n        Entry(int h, K k, V v, Entry<K,V> n) {\n            value = v;\n            next = n;\n            key = k;\n            hash = h;\n        } \n```\n\n\n\n 所以，HashMap的整体结构如下\n\n![image.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/991c29f4486243ca7369fd7645a6d321.png)\n\n　　**简单来说，HashMap由数组+链表组成的，数组是HashMap的主体，链表则是主要为了解决哈希冲突而存在的，如果定位到的数组位置不含链表（当前entry的next指向null）,那么对于查找，添加等操作很快，仅需一次寻址即可；如果定位到的数组包含链表，对于添加操作，其时间复杂度为O(n)，首先遍历链表，存在即覆盖，否则新增；对于查找操作来讲，仍需遍历链表，然后通过key对象的equals方法逐一比对查找。所以，性能考虑，HashMap中的链表出现越少，性能才会越好。**\n\n其他几个重要字段\n\n\n\n```\n//实际存储的key-value键值对的个数\ntransient int size;\n//阈值，当table == {}时，该值为初始容量（初始容量默认为16）；当table被填充了，也就是为table分配内存空间后，threshold一般为 capacity*loadFactory。HashMap在进行扩容时需要参考threshold，后面会详细谈到\nint threshold;\n//负载因子，代表了table的填充度有多少，默认是0.75\nfinal float loadFactor;\n//用于快速失败，由于HashMap非线程安全，在对HashMap进行迭代时，如果期间其他线程的参与导致HashMap的结构发生变化了（比如put，remove等操作），需要抛出异常ConcurrentModificationException\ntransient int modCount;\n```\n\n\n\nHashMap有4个构造器，其他构造器如果用户没有传入initialCapacity 和loadFactor这两个参数，会使用默认值\n\ninitialCapacity默认为16，loadFactory默认为0.75\n\n我们看下其中一个\n\n\n\n```\npublic HashMap(int initialCapacity, float loadFactor) {\n　　　　　//此处对传入的初始容量进行校验，最大不能超过MAXIMUM_CAPACITY = 1<<30(230)\n        if (initialCapacity < 0)\n            throw new IllegalArgumentException(\"Illegal initial capacity: \" +\n                 initialCapacity);\n        if (initialCapacity > MAXIMUM_CAPACITY)\n            initialCapacity = MAXIMUM_CAPACITY;\n        if (loadFactor <= 0 || Float.isNaN(loadFactor))\n            throw new IllegalArgumentException(\"Illegal load factor: \" +\n                 loadFactor);\n\n        this.loadFactor = loadFactor;\n        threshold = initialCapacity;\n　　　　　\n        init();//init方法在HashMap中没有实际实现，不过在其子类如 linkedHashMap中就会有对应实现\n    }\n```\n\n\n\n　　从上面这段代码我们可以看出，**在常规构造器中，没有为数组table分配内存空间（有一个入参为指定Map的构造器例外），而是在执行put操作的时候才真正构建table数组**\n\n　　OK,接下来我们来看看put操作的实现吧\n\n\n\n```\n    public V put(K key, V value) {\n        //如果table数组为空数组{}，进行数组填充（为table分配实际内存空间），入参为threshold，此时threshold为initialCapacity 默认是1<<4(24=16)\n        if (table == EMPTY_TABLE) {\n            inflateTable(threshold);\n        }\n       //如果key为null，存储位置为table[0]或table[0]的冲突链上\n        if (key == null)\n            return putForNullKey(value);\n        int hash = hash(key);//对key的hashcode进一步计算，确保散列均匀\n        int i = indexFor(hash, table.length);//获取在table中的实际位置\n        for (Entry<K,V> e = table[i]; e != null; e = e.next) {\n        //如果该对应数据已存在，执行覆盖操作。用新value替换旧value，并返回旧value\n            Object k;\n            if (e.hash == hash && ((k = e.key) == key || key.equals(k))) {\n                V oldValue = e.value;\n                e.value = value;\n                e.recordAccess(this);\n                return oldValue;\n            }\n        }\n        modCount++;//保证并发访问时，若HashMap内部结构发生变化，快速响应失败\n        addEntry(hash, key, value, i);//新增一个entry\n        return null;\n    }    \n```\n\n\n\n 先来看看inflateTable这个方法\n\n\n\n```\nprivate void inflateTable(int toSize) {\n        int capacity = roundUpToPowerOf2(toSize);//capacity一定是2的次幂\n        threshold = (int) Math.min(capacity * loadFactor, MAXIMUM_CAPACITY + 1);//此处为threshold赋值，取capacity*loadFactor和MAXIMUM_CAPACITY+1的最小值，capaticy一定不会超过MAXIMUM_CAPACITY，除非loadFactor大于1\n        table = new Entry[capacity];\n        initHashSeedAsNeeded(capacity);\n    }\n```\n\n\n\n　　inflateTable这个方法用于为主干数组table在内存中分配存储空间，通过roundUpToPowerOf2(toSize)可以确保capacity为大于或等于toSize的最接近toSize的二次幂，比如toSize=13,则capacity=16;to_size=16,capacity=16;to_size=17,capacity=32.\n\n\n\n```\n private static int roundUpToPowerOf2(int number) {\n        // assert number >= 0 : \"number must be non-negative\";\n        return number >= MAXIMUM_CAPACITY\n                ? MAXIMUM_CAPACITY\n                : (number > 1) ? Integer.highestOneBit((number - 1) << 1) : 1;\n    }\n```\n\n\n\nroundUpToPowerOf2中的这段处理使得数组长度一定为2的次幂，Integer.highestOneBit是用来获取最左边的bit（其他bit位为0）所代表的数值.\n\nhash函数\n\n\n\n```\n//这是一个神奇的函数，用了很多的异或，移位等运算，对key的hashcode进一步进行计算以及二进制位的调整等来保证最终获取的存储位置尽量分布均匀\nfinal int hash(Object k) {\n        int h = hashSeed;\n        if (0 != h && k instanceof String) {\n            return sun.misc.Hashing.stringHash32((String) k);\n        }\n\n        h ^= k.hashCode();\n\n        h ^= (h >>> 20) ^ (h >>> 12);\n        return h ^ (h >>> 7) ^ (h >>> 4);\n    }\n```\n\n\n\n以上hash函数计算出的值，通过indexFor进一步处理来获取实际的存储位置\n\n\n\n```\n　　/**\n     * 返回数组下标\n     */\n    static int indexFor(int h, int length) {\n        return h & (length-1);\n    }\n```\n\n\n\nh&（length-1）保证获取的index一定在数组范围内，举个例子，默认容量16，length-1=15，h=18,转换成二进制计算为\n\n```\n        1  0  0  1  0\n    &   0  1  1  1  1\n    __________________\n        0  0  0  1  0    = 2\n```\n\n　　最终计算出的index=2。有些版本的对于此处的计算会使用 取模运算，也能保证index一定在数组范围内，不过位运算对计算机来说，性能更高一些（HashMap中有大量位运算）\n\n所以最终存储位置的确定流程是这样的：\n\n![image.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/d669d9304ba570e375512d4948087e45.png)\n\n再来看看addEntry的实现：\n\n\n\n```\nvoid addEntry(int hash, K key, V value, int bucketIndex) {\n        if ((size >= threshold) && (null != table[bucketIndex])) {\n            resize(2 * table.length);//当size超过临界阈值threshold，并且即将发生哈希冲突时进行扩容\n            hash = (null != key) ? hash(key) : 0;\n            bucketIndex = indexFor(hash, table.length);\n        }\n\n        createEntry(hash, key, value, bucketIndex);\n    }\n```\n\n\n\n　　通过以上代码能够得知，当发生哈希冲突并且size大于阈值的时候，需要进行数组扩容，扩容时，需要新建一个长度为之前数组2倍的新的数组，然后将当前的Entry数组中的元素全部传输过去，扩容后的新数组长度为之前的2倍，所以扩容相对来说是个耗资源的操作。\n\n# 三、为何HashMap的数组长度一定是2的次幂？\n\n我们来继续看上面提到的resize方法\n\n\n\n```\n void resize(int newCapacity) {\n        Entry[] oldTable = table;\n        int oldCapacity = oldTable.length;\n        if (oldCapacity == MAXIMUM_CAPACITY) {\n            threshold = Integer.MAX_VALUE;\n            return;\n        }\n\n        Entry[] newTable = new Entry[newCapacity];\n        transfer(newTable, initHashSeedAsNeeded(newCapacity));\n        table = newTable;\n        threshold = (int)Math.min(newCapacity * loadFactor, MAXIMUM_CAPACITY + 1);\n    }\n```\n\n\n\n如果数组进行扩容，数组长度发生变化，而存储位置 index = h&(length-1),index也可能会发生变化，需要重新计算index，我们先来看看transfer这个方法\n\n\n\n```\nvoid transfer(Entry[] newTable, boolean rehash) {\n        int newCapacity = newTable.length;\n　　　　　//for循环中的代码，逐个遍历链表，重新计算索引位置，将老数组数据复制到新数组中去（数组不存储实际数据，所以仅仅是拷贝引用而已）\n        for (Entry<K,V> e : table) {\n            while(null != e) {\n                Entry<K,V> next = e.next;\n                if (rehash) {\n                    e.hash = null == e.key ? 0 : hash(e.key);\n                }\n                int i = indexFor(e.hash, newCapacity);\n　　　　　　　　　 //将当前entry的next链指向新的索引位置,newTable[i]有可能为空，有可能也是个entry链，如果是entry链，直接在链表头部插入。\n                e.next = newTable[i];\n                newTable[i] = e;\n                e = next;\n            }\n        }\n    }\n```\n\n\n\n　　这个方法将老数组中的数据逐个链表地遍历，扔到新的扩容后的数组中，我们的数组索引位置的计算是通过 对key值的hashcode进行hash扰乱运算后，再通过和 length-1进行位运算得到最终数组索引位置。\n\n　　hashMap的数组长度一定保持2的次幂，比如16的二进制表示为 10000，那么length-1就是15，二进制为01111，同理扩容后的数组长度为32，二进制表示为100000，length-1为31，二进制表示为011111。从下图可以我们也能看到这样会保证低位全为1，而扩容后只有一位差异，也就是多出了最左位的1，这样在通过 h&(length-1)的时候，只要h对应的最左边的那一个差异位为0，就能保证得到的新的数组索引和老数组索引一致(大大减少了之前已经散列良好的老数组的数据位置重新调换)，个人理解。\n\n![image.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/d38c9ea1ecdf22629d8f79ebb8679f05.png)\n\n 还有，数组长度保持2的次幂，length-1的低位都为1，会使得获得的数组索引index更加均匀，比如：\n\n![image.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/6f19110cb2a9b956502437c74a52e792.png)\n\n　　我们看到，上面的&运算，高位是不会对结果产生影响的（hash函数采用各种位运算可能也是为了使得低位更加散列），我们只关注低位bit，如果低位全部为1，那么对于h低位部分来说，任何一位的变化都会对结果产生影响，也就是说，要得到index=21这个存储位置，h的低位只有这一种组合。这也是数组长度设计为必须为2的次幂的原因。\n\n![image.png]( https://veport.oss-cn-beijing.aliyuncs.com/articles/ebdf005737a5d7daf25b2dbdb763e8d8.png)\n\n　　如果不是2的次幂，也就是低位不是全为1此时，要使得index=21，h的低位部分不再具有唯一性了，哈希冲突的几率会变的更大，同时，index对应的这个bit位无论如何不会等于1了，而对应的那些数组位置也就被白白浪费了。\n\nget方法\n\n\n\n```\n public V get(Object key) {\n　　　　 //如果key为null,则直接去table[0]处去检索即可。\n        if (key == null)\n            return getForNullKey();\n        Entry<K,V> entry = getEntry(key);\n        return null == entry ? null : entry.getValue();\n }\n```\n\n\n\nget方法通过key值返回对应value，如果key为null，直接去table[0]处检索。我们再看一下getEntry这个方法\n\n\n\n```\nfinal Entry<K,V> getEntry(Object key) {\n            \n        if (size == 0) {\n            return null;\n        }\n        //通过key的hashcode值计算hash值\n        int hash = (key == null) ? 0 : hash(key);\n        //indexFor (hash&length-1) 获取最终数组索引，然后遍历链表，通过equals方法比对找出对应记录\n        for (Entry<K,V> e = table[indexFor(hash, table.length)];\n             e != null;\n             e = e.next) {\n            Object k;\n            if (e.hash == hash && \n                ((k = e.key) == key || (key != null && key.equals(k))))\n                return e;\n        }\n        return null;\n    }    \n```\n\n\n\n　　可以看出，get方法的实现相对简单，key(hashcode)-->hash-->indexFor-->最终索引位置，找到对应位置table[i]，再查看是否有链表，遍历链表，通过key的equals方法比对查找对应的记录。要注意的是，有人觉得上面在定位到数组位置之后然后遍历链表的时候，e.hash == hash这个判断没必要，仅通过equals判断就可以。其实不然，试想一下，如果传入的key对象重写了equals方法却没有重写hashCode，而恰巧此对象定位到这个数组位置，如果仅仅用equals判断可能是相等的，但其hashCode和当前对象不一致，这种情况，根据Object的hashCode的约定，不能返回当前对象，而应该返回null，后面的例子会做出进一步解释。\n\n# 四、重写equals方法需同时重写hashCode方法\n\n　　关于HashMap的源码分析就介绍到这儿了，最后我们再聊聊老生常谈的一个问题，各种资料上都会提到，“重写equals时也要同时覆盖hashcode”，我们举个小例子来看看，如果重写了equals而不重写hashcode会发生什么样的问题\n\n\n\n```\n/**\n * Created by chengxiao on 2016/11/15.\n */\npublic class MyTest {\n    private static class Person{\n        int idCard;\n        String name;\n\n        public Person(int idCard, String name) {\n            this.idCard = idCard;\n            this.name = name;\n        }\n        @Override\n        public boolean equals(Object o) {\n            if (this == o) {\n                return true;\n            }\n            if (o == null || getClass() != o.getClass()){\n                return false;\n            }\n            Person person = (Person) o;\n            //两个对象是否等值，通过idCard来确定\n            return this.idCard == person.idCard;\n        }\n\n    }\n    public static void main(String []args){\n        HashMap<Person,String> map = new HashMap<Person, String>();\n        Person person = new Person(1234,\"乔峰\");\n        //put到hashmap中去\n        map.put(person,\"天龙八部\");\n        //get取出，从逻辑上讲应该能输出“天龙八部”\n        System.out.println(\"结果:\"+map.get(new Person(1234,\"萧峰\")));\n    }\n}\n```\n\n\n\n实际输出结果：\n\n```\n结果：null\n```\n\n　　如果我们已经对HashMap的原理有了一定了解，这个结果就不难理解了。尽管我们在进行get和put操作的时候，使用的key从逻辑上讲是等值的（通过equals比较是相等的），但由于没有重写hashCode方法，所以put操作时，key(hashcode1)-->hash-->indexFor-->最终索引位置 ，而通过key取出value的时候 key(hashcode1)-->hash-->indexFor-->最终索引位置，由于hashcode1不等于hashcode2，导致没有定位到一个数组位置而返回逻辑上错误的值null（也有可能碰巧定位到一个数组位置，但是也会判断其entry的hash值是否相等，上面get方法中有提到。）\n\n　　所以，在重写equals的方法的时候，必须注意重写hashCode方法，同时还要保证通过equals判断相等的两个对象，调用hashCode方法要返回同样的整数值。而如果equals判断不相等的两个对象，其hashCode可以相同（只不过会发生哈希冲突，应尽量避免）。\n\n# 五、总结\n\n　　本文描述了HashMap的实现原理，并结合源码做了进一步的分析，也涉及到一些源码细节设计缘由，最后简单介绍了为什么重写equals的时候需要重写hashCode方法。希望本篇文章能帮助到大家，同时也欢迎讨论指正，谢谢支持！\n\n作者： [dreamcatcher-cx](http://www.cnblogs.com/chengxiao/)\n\n出处： <http://www.cnblogs.com/chengxiao/>\n\n本文版权归作者和博客园共有，欢迎转载，但未经作者同意必须保留此段声明，且在页面明显位置给出原文链接。', 2, '', 0, 0, 2, '2022-03-11 03:24:24', '2023-11-10 14:52:11');
INSERT INTO `article` VALUES (83, 3, 190, 'https://static.veweiyi.cn/blog/3/upload/article/cover/wusheng_20231205203537.jpg', '我刚刚购买了新的服务器', '我刚刚购买了新的服务器1', 1, '', 1, 0, 1, '2023-11-09 18:46:15', '2023-12-20 12:21:14');
COMMIT;

-- ----------------------------
-- Table structure for article_tag
-- ----------------------------
DROP TABLE IF EXISTS `article_tag`;
CREATE TABLE `article_tag` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `article_id` int NOT NULL DEFAULT '0' COMMENT '文章id',
  `tag_id` int NOT NULL DEFAULT '0' COMMENT '标签id',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_article_tag_1` (`article_id`) USING BTREE,
  KEY `fk_article_tag_2` (`tag_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1064 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='文章-标签关联';

-- ----------------------------
-- Records of article_tag
-- ----------------------------
BEGIN;
INSERT INTO `article_tag` VALUES (854, 54, 29);
INSERT INTO `article_tag` VALUES (869, 59, 30);
INSERT INTO `article_tag` VALUES (870, 59, 33);
INSERT INTO `article_tag` VALUES (875, 56, 29);
INSERT INTO `article_tag` VALUES (928, 62, 30);
INSERT INTO `article_tag` VALUES (929, 62, 35);
INSERT INTO `article_tag` VALUES (932, 66, 32);
INSERT INTO `article_tag` VALUES (933, 66, 37);
INSERT INTO `article_tag` VALUES (934, 64, 33);
INSERT INTO `article_tag` VALUES (935, 64, 36);
INSERT INTO `article_tag` VALUES (936, 69, 38);
INSERT INTO `article_tag` VALUES (937, 69, 39);
INSERT INTO `article_tag` VALUES (938, 69, 40);
INSERT INTO `article_tag` VALUES (939, 68, 30);
INSERT INTO `article_tag` VALUES (940, 68, 37);
INSERT INTO `article_tag` VALUES (980, 71, 44);
INSERT INTO `article_tag` VALUES (981, 71, 45);
INSERT INTO `article_tag` VALUES (982, 71, 46);
INSERT INTO `article_tag` VALUES (1002, 70, 41);
INSERT INTO `article_tag` VALUES (1003, 70, 42);
INSERT INTO `article_tag` VALUES (1004, 70, 43);
INSERT INTO `article_tag` VALUES (1011, 80, 33);
INSERT INTO `article_tag` VALUES (1030, 72, 35);
INSERT INTO `article_tag` VALUES (1031, 72, 47);
INSERT INTO `article_tag` VALUES (1035, 79, 48);
INSERT INTO `article_tag` VALUES (1036, 81, 49);
INSERT INTO `article_tag` VALUES (1038, 82, 50);
INSERT INTO `article_tag` VALUES (1057, 60, 32);
INSERT INTO `article_tag` VALUES (1058, 60, 33);
INSERT INTO `article_tag` VALUES (1059, 60, 39);
INSERT INTO `article_tag` VALUES (1062, 83, 30);
INSERT INTO `article_tag` VALUES (1063, 83, 31);
COMMIT;

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) NOT NULL DEFAULT '',
  `v0` varchar(100) NOT NULL DEFAULT '',
  `v1` varchar(100) NOT NULL DEFAULT '',
  `v2` varchar(100) NOT NULL DEFAULT '',
  `v3` varchar(100) NOT NULL DEFAULT '',
  `v4` varchar(100) NOT NULL DEFAULT '',
  `v5` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `category_name` varchar(32) NOT NULL DEFAULT '' COMMENT '分类名',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=194 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='文章分类';

-- ----------------------------
-- Records of category
-- ----------------------------
BEGIN;
INSERT INTO `category` VALUES (187, '测试分类', '2022-01-18 00:29:02', '2024-01-15 19:43:08');
INSERT INTO `category` VALUES (188, '学习', '2022-01-19 22:35:50', '2024-01-15 19:43:08');
INSERT INTO `category` VALUES (189, '技术', '2022-01-21 12:21:31', '2024-01-15 19:43:08');
INSERT INTO `category` VALUES (190, 'bug修复', '2022-01-22 23:55:55', '2024-01-15 19:43:08');
INSERT INTO `category` VALUES (191, '网站搭建', '2022-02-11 23:26:12', '2024-01-15 19:43:08');
INSERT INTO `category` VALUES (192, '算法日记', '2022-02-12 16:28:13', '2024-01-15 19:43:08');
INSERT INTO `category` VALUES (193, '面试笔记', '2022-02-19 11:31:18', '2024-01-15 19:43:08');
COMMIT;

-- ----------------------------
-- Table structure for chat_message
-- ----------------------------
DROP TABLE IF EXISTS `chat_message`;
CREATE TABLE `chat_message` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户id',
  `chat_id` int NOT NULL DEFAULT '0' COMMENT '群聊id',
  `reply_msg_id` int NOT NULL DEFAULT '0' COMMENT '回复消息id',
  `content` varchar(1024) NOT NULL DEFAULT '' COMMENT '聊天内容',
  `ip_address` varchar(64) NOT NULL DEFAULT '' COMMENT 'ip地址',
  `ip_source` varchar(255) NOT NULL DEFAULT '' COMMENT 'ip来源',
  `type` int NOT NULL DEFAULT '0' COMMENT '类型',
  `status` int NOT NULL DEFAULT '0' COMMENT '0正常 1撤回 2已编辑',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='聊天消息';

-- ----------------------------
-- Records of chat_message
-- ----------------------------
BEGIN;
INSERT INTO `chat_message` VALUES (1, 3, 3, 0, '你好', '127.0.0.1', '本机地址', 0, 0, '2024-02-23 14:57:10.0', '2024-02-23 14:57:10.0');
INSERT INTO `chat_message` VALUES (2, 0, 3, 0, '你好！有什么可以帮助你的吗？', '', '', 1, 0, '2024-02-23 14:57:10.0', '2024-02-23 14:57:10.0');
INSERT INTO `chat_message` VALUES (3, 3, 3, 0, '你叫什么名字', '127.0.0.1', '本机地址', 0, 0, '2024-02-23 14:57:46.0', '2024-02-23 14:57:46.0');
INSERT INTO `chat_message` VALUES (4, 0, 3, 0, '我是Assistant，很高兴为您服务。您有什么问题或需要帮助吗？', '', '', 1, 0, '2024-02-23 14:57:46.0', '2024-02-23 14:57:46.0');
INSERT INTO `chat_message` VALUES (5, 3, 3, 0, '你可以帮我做什么', '127.0.0.1', '本机地址', 0, 0, '2024-02-23 14:58:15.0', '2024-02-23 14:58:15.0');
INSERT INTO `chat_message` VALUES (6, 0, 3, 0, '我是一个语言模型AI助手，可以回答你的问题，提供信息，帮助解决问题，甚至进行一些娱乐性的对话。有什么我可以帮助你的吗？', '', '', 1, 0, '2024-02-23 14:58:15.0', '2024-02-23 14:58:15.0');
INSERT INTO `chat_message` VALUES (7, 3, 3, 0, '你能告诉我我们第一次见面说的话吗', '127.0.0.1', '本机地址', 0, 0, '2024-02-23 14:58:41.0', '2024-02-23 14:58:41.0');
INSERT INTO `chat_message` VALUES (8, 0, 3, 0, '当然，我可以告诉你我们第一次见面时的对话。不过请注意，我无法实际记忆过去的对话，这只是一种模拟。你愿意听吗？', '', '', 1, 0, '2024-02-23 14:58:41.0', '2024-02-23 14:58:41.0');
INSERT INTO `chat_message` VALUES (9, 3, 3, 0, '我的第一个问题是什么', '127.0.0.1', '本机地址', 0, 0, '2024-02-23 14:58:52.0', '2024-02-23 14:58:52.0');
INSERT INTO `chat_message` VALUES (10, 0, 3, 0, '你的第一个问题是\"你叫什么名字？\"', '', '', 1, 0, '2024-02-23 14:58:52.0', '2024-02-23 14:58:52.0');
COMMIT;

-- ----------------------------
-- Table structure for chat_record
-- ----------------------------
DROP TABLE IF EXISTS `chat_record`;
CREATE TABLE `chat_record` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户id',
  `nickname` varchar(64) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `content` varchar(1000) NOT NULL DEFAULT '' COMMENT '聊天内容',
  `ip_address` varchar(64) NOT NULL DEFAULT '' COMMENT 'ip地址',
  `ip_source` varchar(255) NOT NULL DEFAULT '' COMMENT 'ip来源',
  `type` int NOT NULL DEFAULT '0' COMMENT '类型',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2938 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='聊天记录';

-- ----------------------------
-- Records of chat_record
-- ----------------------------
BEGIN;
INSERT INTO `chat_record` VALUES (2894, 0, '未知ip', 'https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', '<img src= \'https://static.talkxj.com/emoji/cy.jpg\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '未知ip', '', 3, '2022-02-09 23:08:55', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2895, 0, '未知ip', 'https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', '<img src= \'https://static.talkxj.com/emoji/goutou.jpg\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '未知ip', '', 3, '2022-02-10 19:16:49', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2896, 0, '未知ip', 'https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', '<img src= \'https://static.talkxj.com/emoji/daxiao.jpg\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '未知ip', '', 3, '2022-02-10 21:27:47', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2897, 0, '未知ip', 'https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', '<img src= \'https://static.talkxj.com/emoji/goutou.jpg\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '未知ip', '', 3, '2022-02-11 11:15:10', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2898, 0, '未知ip', 'https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', '<img src= \'https://static.talkxj.com/emoji/smile.jpg\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '未知ip', '', 3, '2022-02-11 15:58:31', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2899, 0, '未知ip', 'https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', '<img src= \'https://static.talkxj.com/emoji/goutou.jpg\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '未知ip', '', 3, '2022-02-11 16:07:35', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2900, 981, '静闻弦语', 'http://thirdqq.qlogo.cn/g?b=oidb&k=wwQSEfJO8IdQbaPPATTUWg&s=40&t=1644466423', '<img src= \'https://static.talkxj.com/emoji/daxiao.jpg\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '未知ip', '', 3, '2022-02-11 17:14:19', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2901, 0, '未知ip', 'https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', '？？', '未知ip', '', 3, '2022-02-18 23:07:10', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2902, 0, '未知ip', 'https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', '<img src= \'https://static.talkxj.com/emoji/tiaopi.jpg\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '未知ip', '', 3, '2022-02-19 22:26:40', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2903, 0, '未知ip', 'https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', '<img src= \'https://static.talkxj.com/emoji/xxy.jpg\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '未知ip', '', 3, '2022-02-19 22:29:53', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2904, 0, '', 'https://static.ve77.cn/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', 'https://static.ve77.cn/blog/voice/8d727c08f80af4b8ce15c4add0300ee7.wav', '', '', 5, '2023-06-25 09:31:07', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2905, 0, '', 'https://static.ve77.cn/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', 'https://static.ve77.cn/blog/voice/446ce581209b26ffa508347cde1766b6.wav', '', '', 5, '2023-06-25 09:31:09', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2906, 0, '', 'https://static.ve77.cn/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', 'https://static.ve77.cn/blog/voice/d112f7691d4be8c1c4c2d0ae1c403fb8.wav', '', '', 5, '2023-06-25 09:31:14', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2907, 0, '', 'https://static.ve77.cn/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', 'https://static.ve77.cn/blog/voice/f6326fd0c5cc6bf32421c9ce6ef27211.wav', '', '', 5, '2023-06-25 09:31:22', '2024-01-15 19:43:08');
INSERT INTO `chat_record` VALUES (2908, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', 's', '', '', 3, '2023-07-05 20:55:08', '2023-07-05 20:55:08');
INSERT INTO `chat_record` VALUES (2916, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', 'hi', '', '', 3, '2023-07-05 21:03:05', '2023-07-05 21:03:05');
INSERT INTO `chat_record` VALUES (2917, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', 'hell', '', '', 3, '2023-07-05 21:03:13', '2023-07-05 21:03:13');
INSERT INTO `chat_record` VALUES (2919, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', 'hi', '', '', 3, '2023-07-05 21:03:51', '2023-07-05 21:03:51');
INSERT INTO `chat_record` VALUES (2920, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', 'hello', '', '', 3, '2023-07-05 21:04:01', '2023-07-05 21:04:01');
INSERT INTO `chat_record` VALUES (2924, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', 'hi', '', '', 3, '2023-07-10 05:57:29', '2023-07-10 05:57:29');
INSERT INTO `chat_record` VALUES (2925, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', 'ss', '', '', 3, '2023-07-10 05:57:32', '2023-07-10 05:57:32');
INSERT INTO `chat_record` VALUES (2926, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', 's', '', '', 3, '2023-07-10 05:57:43', '2023-07-10 05:57:43');
INSERT INTO `chat_record` VALUES (2927, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', '你好', '', '', 3, '2023-07-10 05:57:55', '2023-07-10 05:57:55');
INSERT INTO `chat_record` VALUES (2928, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', 'ss', '', '', 3, '2023-07-10 05:58:01', '2023-07-10 05:58:01');
INSERT INTO `chat_record` VALUES (2929, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', '<img src= \'http://localhost:8888/blog/src/assets/emojis/qq/newemoji_007.gif\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '', '', 3, '2023-11-08 14:48:07', '2023-11-08 14:48:07');
INSERT INTO `chat_record` VALUES (2930, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', '1111', '', '', 3, '2023-11-08 14:57:40', '2023-11-08 14:57:40');
INSERT INTO `chat_record` VALUES (2931, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', '111', '', '', 3, '2023-11-08 14:57:45', '2023-11-08 14:57:45');
INSERT INTO `chat_record` VALUES (2932, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', '111', '', '', 3, '2023-11-08 14:57:46', '2023-11-08 14:57:46');
INSERT INTO `chat_record` VALUES (2933, 0, '', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', '<img src= \'http://localhost:8888/blog/src/assets/emojis/qq/28@2x.gif\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '', '', 3, '2023-11-08 15:05:16', '2023-11-08 15:05:16');
INSERT INTO `chat_record` VALUES (2934, 0, 'admin@qq.com', 'http://rxb1y0x1n.hn-bkt.clouddn.com/blog/avatar/43b90920409618f188bfc6923f16b9fa_20230713143425.jpg', '<img src= \'http://localhost:8888/blog/undefined\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '', '', 3, '2023-12-04 21:21:13', '2023-12-04 21:21:13');
INSERT INTO `chat_record` VALUES (2935, 0, 'admin@qq.com', 'http://rxb1y0x1n.hn-bkt.clouddn.com/blog/avatar/43b90920409618f188bfc6923f16b9fa_20230713143425.jpg', '<img src= \'https://static.veweiyi.cn/emoji/qq/27@2x.gif\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '', '', 3, '2023-12-04 21:23:06', '2023-12-04 21:23:06');
INSERT INTO `chat_record` VALUES (2936, 0, 'admin@qq.com', 'http://rxb1y0x1n.hn-bkt.clouddn.com/blog/avatar/43b90920409618f188bfc6923f16b9fa_20230713143425.jpg', '<img src= \'https://static.veweiyi.cn/emoji/qq/104@2x.gif\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '', '', 3, '2023-12-04 21:24:08', '2023-12-04 21:24:08');
INSERT INTO `chat_record` VALUES (2937, 0, 'admin@qq.com', 'http://rxb1y0x1n.hn-bkt.clouddn.com/blog/avatar/43b90920409618f188bfc6923f16b9fa_20230713143425.jpg', '<img src= \'https://static.veweiyi.cn/emoji/qq/26@2x.gif\' width=\'24\'height=\'24\' style=\'margin: 0 1px;vertical-align: text-bottom\'/>', '', '', 3, '2023-12-05 09:45:23', '2023-12-05 09:45:23');
COMMIT;

-- ----------------------------
-- Table structure for chat_session
-- ----------------------------
DROP TABLE IF EXISTS `chat_session`;
CREATE TABLE `chat_session` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `chat_title` varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
  `type` varchar(128) NOT NULL DEFAULT '' COMMENT '类型',
  `status` int NOT NULL DEFAULT '0' COMMENT '0正常 1删除',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='聊天会话';

-- ----------------------------
-- Records of chat_session
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '评论用户Id',
  `topic_id` int NOT NULL DEFAULT '0' COMMENT '评论主题id',
  `comment_content` text NOT NULL COMMENT '评论内容',
  `reply_user_id` int NOT NULL DEFAULT '0' COMMENT '回复用户id',
  `parent_id` int NOT NULL DEFAULT '0' COMMENT '父评论id',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT '评论类型 1.文章 2.友链 3.说说',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除  0否 1是',
  `is_review` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否审核',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评论时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_comment_user` (`user_id`) USING BTREE,
  KEY `fk_comment_parent` (`parent_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=754 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='评论';

-- ----------------------------
-- Records of comment
-- ----------------------------
BEGIN;
INSERT INTO `comment` VALUES (722, 1, 54, '测试评论', 0, 0, 1, 0, 1, '2022-01-24 23:34:25', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (723, 1, 54, '测试回复', 1, 722, 1, 0, 1, '2022-01-24 23:34:30', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (724, 1, 49, '测试评论', 0, 0, 3, 0, 1, '2022-01-24 23:35:25', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (725, 981, 51, '或许你想一起打个游戏？', 0, 0, 3, 0, 1, '2022-02-11 17:15:04', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (726, 981, 51, '或许你想一起打个游戏？', 981, 725, 3, 0, 1, '2022-02-11 17:19:20', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (727, 1, 51, '测试回复', 981, 725, 3, 0, 1, '2022-04-09 18:16:40', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (728, 985, 51, 'hello', 1, 725, 3, 0, 1, '2022-04-09 18:18:56', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (729, 1, 51, 'hello', 981, 725, 3, 0, 1, '2022-04-09 18:19:37', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (730, 985, 68, '很牛', 0, 0, 1, 0, 1, '2022-04-09 18:36:21', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (733, 986, 71, '牛啊牛啊', 0, 0, 1, 0, 1, '2022-05-06 22:45:14', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (734, 984, 71, '是真的牛', 986, 733, 1, 0, 1, '2022-05-06 22:49:14', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (735, 986, 71, '我也觉得', 984, 733, 1, 0, 1, '2022-05-06 22:49:41', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (736, 985, 71, '非常牛', 984, 733, 1, 0, 1, '2022-05-06 22:51:10', '2024-01-15 19:43:08');
INSERT INTO `comment` VALUES (737, 985, 71, 'nice', 986, 733, 1, 0, 1, '2022-05-06 22:51:24', '2024-01-15 19:43:08');
COMMIT;

-- ----------------------------
-- Table structure for friend_link
-- ----------------------------
DROP TABLE IF EXISTS `friend_link`;
CREATE TABLE `friend_link` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `link_name` varchar(32) NOT NULL DEFAULT '' COMMENT '链接名',
  `link_avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '链接头像',
  `link_address` varchar(64) NOT NULL DEFAULT '' COMMENT '链接地址',
  `link_intro` varchar(100) NOT NULL DEFAULT '' COMMENT '链接介绍',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_friend_link_user` (`link_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='友链';

-- ----------------------------
-- Records of friend_link
-- ----------------------------
BEGIN;
INSERT INTO `friend_link` VALUES (27, 'Tcefrep的个人博客', 'https://unsplash.it/100/100?image=295', 'https://www.tcefrep.site', '这是一篇简单的个人博客，也是一个我记录笔记的地方，欢迎各位到访', '2022-01-23 12:56:57', '2023-08-29 15:20:14');
INSERT INTO `friend_link` VALUES (29, '码霸霸', 'https://ld246.com/images/favicon.ico', 'https://blog.lupf.cn/', 'Java后端面试题刷题手册', '2022-02-15 16:45:54', '2023-09-01 15:32:37');
COMMIT;

-- ----------------------------
-- Table structure for menu
-- ----------------------------
DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '菜单名',
  `path` varchar(64) NOT NULL DEFAULT '' COMMENT '菜单路径',
  `title` varchar(32) NOT NULL DEFAULT '' COMMENT '菜单标题',
  `component` varchar(64) NOT NULL DEFAULT '' COMMENT '组件',
  `icon` varchar(64) NOT NULL DEFAULT '' COMMENT '菜单icon',
  `rank` int NOT NULL DEFAULT '0' COMMENT '排序',
  `parent_id` int NOT NULL DEFAULT '0' COMMENT '父id',
  `is_hidden` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否隐藏  0否1是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_path` (`path`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=48 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='菜单';

-- ----------------------------
-- Records of menu
-- ----------------------------
BEGIN;
INSERT INTO `menu` VALUES (1, 'Home', '/', 'menus.hshome', '', 'homeFilled', 0, 0, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (2, 'Home', '/home', 'menus.hshome', '', '', 0, 1, 1, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (3, 'Welcome', '/welcome', 'menus.hshome', '', '', 1, 1, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (4, 'article-management', '/article', '文章管理', '', 'document', 10, 0, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (5, 'article/publish', '/article/publish/:articleId', '文章管理', '', '', 0, 4, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (6, 'article/list', '/article/list', '文章管理', '', '', 1, 4, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (7, 'article/category', '/article/category', '文章管理', '', '', 2, 4, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (8, 'article/tag', '/article/tag', '文章管理', '', '', 3, 4, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (9, 'message-management', '/message', '消息管理', '', 'message', 20, 0, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (10, 'comment', '/message/comment', '消息管理', '', '', 0, 9, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (11, 'remark', '/message/remark', '消息管理', '', '', 1, 9, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (12, 'links', '/message/links', '消息管理', '', '', 2, 9, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (13, 'talk-management', '/talk', '说说管理', '', 'ChatDotSquare', 30, 0, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (14, 'talk_publish', '/talk/publish', '说说管理', '', '', 0, 13, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (15, 'talk_edit', '/talk/edit/:talkId', '说说管理', '', '', 1, 13, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (16, 'talk_list', '/talk/list', '说说管理', '', '', 2, 13, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (17, 'album-management', '/albums', '相册管理', '', 'picture', 40, 0, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (18, 'Albums', '/albums/list', '相册管理', '', '', 0, 17, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (19, 'Photo', '/albums/:id', '相册管理', '', '', 1, 17, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (20, 'Delete', '/albums/photo/delete', '相册管理', '', '', 2, 17, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (21, 'user-management', '/user', '用户管理', '', 'user', 50, 0, 0, '2024-01-18 19:41:11', '2024-01-18 19:41:11');
INSERT INTO `menu` VALUES (22, 'list', '/user/list', '用户管理', '', 'user-filled', 0, 21, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (23, 'online', '/user/online', '用户管理', '', 'user-filled', 1, 21, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (24, 'authority-management', '/authority', '权限管理', '', 'lock', 60, 0, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (25, 'Role', '/authority/role', '权限管理', '', '', 0, 24, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (26, 'Menu', '/authority/menu', '权限管理', '', '', 1, 24, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (27, 'Resource', '/authority/resource', '权限管理', '', '', 2, 24, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (28, 'log-management', '/log', '日志管理', '', 'document', 70, 0, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (29, 'operation', '/log/operation', '日志管理', '', '', 0, 28, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (30, 'website-management', '/website', '网站管理', '', 'setting', 80, 0, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (31, 'Config', '/website/config', '网站管理', '', '', 0, 30, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (32, 'Page', '/website/page', '网站管理', '', '', 1, 30, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (33, 'About', '/website/about', '网站管理', '', '', 2, 30, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (34, 'mine-management', '/mine', '个人中心', '', 'user', 90, 0, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (35, '/mine/userinfo', '/mine/userinfo', '个人中心', '', '', 0, 34, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (36, '/mine/history', '/mine/history', '个人中心', '', '', 1, 34, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (37, 'system-management', '/system', '系统管理', '', 'system', 90, 0, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (38, '/website/state', '/website/state', '系统管理', '', '', 0, 37, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (39, 'error', '/error', 'menus.hsabnormal', '', 'informationLine', 99, 0, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (40, '403', '/error/403', 'menus.hsabnormal', '', '', 0, 39, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (41, '404', '/error/404', 'menus.hsabnormal', '', '', 1, 39, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (42, '500', '/error/500', 'menus.hsabnormal', '', '', 2, 39, 0, '2024-01-18 19:41:12', '2024-01-18 19:41:12');
INSERT INTO `menu` VALUES (43, '', '/about', 'menus.hsAbout', '', 'info', 999, 0, 0, '2024-01-18 19:41:13', '2024-01-18 19:41:13');
INSERT INTO `menu` VALUES (44, 'About', '/about/index', 'menus.hsAbout', '', '', 0, 43, 0, '2024-01-18 19:41:13', '2024-01-18 19:41:13');
INSERT INTO `menu` VALUES (45, 'Login', '/login', 'menus.hslogin', '', '', 101, 0, 0, '2024-01-18 19:41:13', '2024-01-18 19:41:13');
INSERT INTO `menu` VALUES (46, '', '/redirect', 'status.hsLoad', '', '', 102, 0, 0, '2024-01-18 19:41:13', '2024-01-18 19:41:13');
INSERT INTO `menu` VALUES (47, 'Redirect', '/redirect/:path(.*)', 'status.hsLoad', '', '', 0, 46, 0, '2024-01-18 19:41:13', '2024-01-18 19:41:13');
COMMIT;

-- ----------------------------
-- Table structure for operation_log
-- ----------------------------
DROP TABLE IF EXISTS `operation_log`;
CREATE TABLE `operation_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户id',
  `nickname` varchar(64) DEFAULT '' COMMENT '用户昵称',
  `ip_address` varchar(255) DEFAULT '' COMMENT '操作ip',
  `ip_source` varchar(255) DEFAULT '' COMMENT '操作地址',
  `opt_module` varchar(32) DEFAULT '' COMMENT '操作模块',
  `opt_desc` varchar(255) DEFAULT '' COMMENT '操作描述',
  `request_url` varchar(255) DEFAULT '' COMMENT '请求地址',
  `request_method` varchar(32) DEFAULT '' COMMENT '请求方式',
  `request_header` varchar(1024) DEFAULT '' COMMENT '请求头参数',
  `request_data` varchar(4096) DEFAULT '' COMMENT '请求参数',
  `response_data` varchar(4096) DEFAULT '' COMMENT '返回数据',
  `response_status` int NOT NULL DEFAULT '0' COMMENT '响应状态码',
  `cost` varchar(32) NOT NULL DEFAULT '' COMMENT '耗时（ms）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='操作记录';

-- ----------------------------
-- Records of operation_log
-- ----------------------------
BEGIN;
INSERT INTO `operation_log` VALUES (1, 3, 'admin@qq.com', '127.0.0.1', '', 'Api', '获取api列表', '/api/v1/api/details_list', 'POST', '', '{\"page\":0,\"page_size\":10,\"sorts\":[{\"field\":\"id\",\"order\":\"desc\"}],\"conditions\":[]}', '{\n \"code\": 0,\n \"message\": \"\",\n \"data\": null,\n \"trace_id\": \"\"\n}', 200, '10.747791ms', '2024-01-18 19:41:39', '2024-01-18 19:41:38');
INSERT INTO `operation_log` VALUES (2, 3, 'admin@qq.com', '127.0.0.1', '', 'Menu', '获取菜单列表', '/api/v1/menu/details_list', 'POST', '', '{\"sorts\":[{\"field\":\"id\",\"order\":\"asc\"}],\"conditions\":[],\"page\":1,\"page_size\":10}', '{\n \"code\": 0,\n \"message\": \"\",\n \"data\": null,\n \"trace_id\": \"\"\n}', 200, '13.15975ms', '2024-01-18 19:42:18', '2024-01-18 19:42:18');
INSERT INTO `operation_log` VALUES (3, 3, 'admin@qq.com', '127.0.0.1', '', 'Api', '获取api列表', '/api/v1/api/details_list', 'POST', '', '{\"sorts\":[{\"field\":\"id\",\"order\":\"desc\"}],\"conditions\":[],\"page\":1,\"page_size\":10}', '{\n \"code\": 0,\n \"message\": \"\",\n \"data\": null,\n \"trace_id\": \"\"\n}', 200, '11.610792ms', '2024-01-18 19:42:19', '2024-01-18 19:42:19');
INSERT INTO `operation_log` VALUES (4, 3, 'admin@qq.com', '127.0.0.1', '', 'Role', '更新角色', '/api/v1/role', 'PUT', '', '{\"id\":4,\"is_default\":1,\"created_at\":\"2023-05-30T20:53:04+08:00\",\"updated_at\":\"2024-01-16T14:46:54+08:00\",\"role_pid\":0,\"role_domain\":\"blog\",\"role_name\":\"super-admin\",\"role_comment\":\"超级管理员\",\"is_disable\":1,\"menu_id_list\":null,\"resource_id_list\":null}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"id\\\":4,\\\"role_domain\\\":\\\"blog\\\",\\\"is_default\\\":1,\\\"created_at\\\":\\\"2023-05-30T20:53:04+08:00\\\",\\\"role_pid\\\":0,\\\"role_name\\\":\\\"super-admin\\\",\\\"role_comment\\\":\\\"超级管理员\\\",\\\"is_disable\\\":1,\\\"updated_at\\\":\\\"2024-01-18T19:52:18.507+08:00\\\"},\\\"trace_id\\\":\\\"5d7f8de0-9209-411d-99a1-702cb5191631\\\"}\"', 200, '23.742708ms', '2024-01-18 19:52:19', '2024-01-18 19:52:18');
INSERT INTO `operation_log` VALUES (5, 3, 'admin@qq.com', '127.0.0.1', '', 'Role', '更新角色', '/api/v1/role', 'PUT', '', '{\"is_default\":0,\"menu_id_list\":null,\"role_pid\":0,\"role_domain\":\"blog\",\"role_comment\":\"超级管理员\",\"created_at\":\"2023-05-30T20:53:04+08:00\",\"updated_at\":\"2024-01-16T14:46:54+08:00\",\"resource_id_list\":null,\"id\":4,\"role_name\":\"super-admin\",\"is_disable\":1}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"is_disable\\\":1,\\\"created_at\\\":\\\"2023-05-30T20:53:04+08:00\\\",\\\"id\\\":4,\\\"role_name\\\":\\\"super-admin\\\",\\\"role_comment\\\":\\\"超级管理员\\\",\\\"is_default\\\":0,\\\"updated_at\\\":\\\"2024-01-18T19:52:18.738+08:00\\\",\\\"role_pid\\\":0,\\\"role_domain\\\":\\\"blog\\\"},\\\"trace_id\\\":\\\"48b461de-3f61-4864-a75c-3d3bf055fd2f\\\"}\"', 200, '30.794083ms', '2024-01-18 19:52:19', '2024-01-18 19:52:18');
INSERT INTO `operation_log` VALUES (6, 3, 'admin@qq.com', '127.0.0.1', '', 'User', '获取用户登录历史', '/api/v1/user/login_history', 'POST', '', '{\"page_size\":10,\"sorts\":[{\"field\":\"id\",\"order\":\"desc\"}],\"conditions\":[],\"page\":1}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"list\\\":null,\\\"page\\\":1,\\\"page_size\\\":10,\\\"total\\\":0},\\\"trace_id\\\":\\\"ad9a2869-24f5-42ef-b5cb-85efe1d16f68\\\"}\"', 200, '56.164083ms', '2024-01-18 20:34:56', '2024-01-18 20:34:56');
INSERT INTO `operation_log` VALUES (7, 3, 'admin@qq.com', '127.0.0.1', '', 'Website', '获取配置', '/api/v1/admin/config', 'POST', '', '{\"key\":\"website_config\"}', '\"{\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":\\\"{\\\\\\\"alipay_qr_code\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/17f234dc487c1bb5bbb732869be0eb53.jpg\\\\\\\",\\\\\\\"is_chat_room\\\\\\\":1,\\\\\\\"is_reward\\\\\\\":1,\\\\\\\"tourist_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif\\\\\\\",\\\\\\\"website_create_time\\\\\\\":\\\\\\\"2022-01-19\\\\\\\",\\\\\\\"gitee\\\\\\\":\\\\\\\"https://gitee.com/wy791422171\\\\\\\",\\\\\\\"is_music_player\\\\\\\":0,\\\\\\\"website_author\\\\\\\":\\\\\\\"与梦\\\\\\\",\\\\\\\"website_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/84aa08357bf6e74fc1d4f33552475f91.gif\\\\\\\",\\\\\\\"website_intro\\\\\\\":\\\\\\\"分享美好生活。\\\\\\\",\\\\\\\"wei_xin_qr_code\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/6bed8a1130b170546341ece729e8819f.jpg\\\\\\\",\\\\\\\"github\\\\\\\":\\\\\\\"https://github.com/7914-ve\\\\\\\",\\\\\\\"is_email_notice\\\\\\\":1,\\\\\\\"is_message_review\\\\\\\":0,\\\\\\\"social_url_list\\\\\\\":[\\\\\\\"qq\\\\\\\",\\\\\\\"github\\\\\\\",\\\\\\\"gitee\\\\\\\"],\\\\\\\"user_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg\\\\\\\",\\\\\\\"website_notice\\\\\\\":\\\\\\\"用户需要查看、发表文章、修改其他信息请登录后台管理系统。网站后台管理系统->https://veweiyi.cn/admin。     \\\\\\\\n网站搭建问题请联系站长QQ791422171。\\\\\\\",\\\\\\\"is_comment_review\\\\\\\":0,\\\\\\\"qq\\\\\\\":\\\\\\\"791422171\\\\\\\",\\\\\\\"social_login_list\\\\\\\":[\\\\\\\"weibo\\\\\\\",\\\\\\\"qq\\\\\\\",\\\\\\\"feishu\\\\\\\"],\\\\\\\"website_name\\\\\\\":\\\\\\\"静闻弦语\\\\\\\",\\\\\\\"website_record_no\\\\\\\":\\\\\\\"桂ICP备2023013735号-1\\\\\\\",\\\\\\\"websocket_url\\\\\\\":\\\\\\\"wss://ve77.cn:8088/api/websocket\\\\\\\"}\\\",\\\"trace_id\\\":\\\"6eb133d9-801a-4ab8-ae30-be0ca1d8b0ec\\\",\\\"code\\\":200}\"', 200, '17.39525ms', '2024-01-21 20:42:02', '2024-01-21 20:42:01');
INSERT INTO `operation_log` VALUES (8, 3, 'admin@qq.com', '127.0.0.1', '', 'Website', '获取配置', '/api/v1/admin/config', 'POST', '', '{\"key\":\"website_config\"}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":\\\"{\\\\\\\"alipay_qr_code\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/17f234dc487c1bb5bbb732869be0eb53.jpg\\\\\\\",\\\\\\\"is_chat_room\\\\\\\":1,\\\\\\\"is_reward\\\\\\\":1,\\\\\\\"tourist_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif\\\\\\\",\\\\\\\"website_create_time\\\\\\\":\\\\\\\"2022-01-19\\\\\\\",\\\\\\\"gitee\\\\\\\":\\\\\\\"https://gitee.com/wy791422171\\\\\\\",\\\\\\\"is_music_player\\\\\\\":0,\\\\\\\"website_author\\\\\\\":\\\\\\\"与梦\\\\\\\",\\\\\\\"website_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/84aa08357bf6e74fc1d4f33552475f91.gif\\\\\\\",\\\\\\\"website_intro\\\\\\\":\\\\\\\"分享美好生活。\\\\\\\",\\\\\\\"wei_xin_qr_code\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/6bed8a1130b170546341ece729e8819f.jpg\\\\\\\",\\\\\\\"github\\\\\\\":\\\\\\\"https://github.com/7914-ve\\\\\\\",\\\\\\\"is_email_notice\\\\\\\":1,\\\\\\\"is_message_review\\\\\\\":0,\\\\\\\"social_url_list\\\\\\\":[\\\\\\\"qq\\\\\\\",\\\\\\\"github\\\\\\\",\\\\\\\"gitee\\\\\\\"],\\\\\\\"user_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg\\\\\\\",\\\\\\\"website_notice\\\\\\\":\\\\\\\"用户需要查看、发表文章、修改其他信息请登录后台管理系统。网站后台管理系统->https://veweiyi.cn/admin。     \\\\\\\\n网站搭建问题请联系站长QQ791422171。\\\\\\\",\\\\\\\"is_comment_review\\\\\\\":0,\\\\\\\"qq\\\\\\\":\\\\\\\"791422171\\\\\\\",\\\\\\\"social_login_list\\\\\\\":[\\\\\\\"weibo\\\\\\\",\\\\\\\"qq\\\\\\\",\\\\\\\"feishu\\\\\\\"],\\\\\\\"website_name\\\\\\\":\\\\\\\"静闻弦语\\\\\\\",\\\\\\\"website_record_no\\\\\\\":\\\\\\\"桂ICP备2023013735号-1\\\\\\\",\\\\\\\"websocket_url\\\\\\\":\\\\\\\"wss://ve77.cn:8088/api/websocket\\\\\\\"}\\\",\\\"trace_id\\\":\\\"352858f0-d2e5-48d9-8e4b-c89f97c5e551\\\"}\"', 200, '21.286167ms', '2024-01-21 20:42:36', '2024-01-21 20:42:35');
INSERT INTO `operation_log` VALUES (9, 3, 'admin@qq.com', '127.0.0.1', '', 'User', '获取用户登录历史', '/api/v1/user/login_history', 'POST', '', '{\"conditions\":[],\"page\":1,\"page_size\":10,\"sorts\":[{\"field\":\"id\",\"order\":\"desc\"}]}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"list\\\":[{\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:47:51 +0800 CST\\\",\\\"id\\\":13,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\"},{\\\"id\\\":12,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:35:00 +0800 CST\\\"},{\\\"login_time\\\":\\\"2024-01-19 17:29:40 +0800 CST\\\",\\\"id\\\":11,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\"},{\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:29:13 +0800 CST\\\",\\\"id\\\":10,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\"},{\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:28:25 +0800 CST\\\",\\\"id\\\":9},{\\\"id\\\":8,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 15:38:21 +0800 CST\\\"},{\\\"login_time\\\":\\\"2024-01-19 11:33:58 +0800 CST\\\",\\\"id\\\":7,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\"},{\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:27:05 +0800 CST\\\",\\\"id\\\":6,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\"},{\\\"login_time\\\":\\\"2024-01-19 10:17:16 +0800 CST\\\",\\\"id\\\":5,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\"},{\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:16:10 +0800 CST\\\",\\\"id\\\":4,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\"}],\\\"page\\\":1,\\\"page_size\\\":10,\\\"total\\\":13},\\\"trace_id\\\":\\\"94a96ca5-d292-45c3-a924-4a4aefeb0e9d\\\"}\"', 200, '46.038125ms', '2024-01-21 20:44:18', '2024-01-21 20:44:17');
INSERT INTO `operation_log` VALUES (10, 3, 'admin@qq.com', '127.0.0.1', '', 'User', '获取用户登录历史', '/api/v1/user/login_history', 'POST', '', '{\"page\":2,\"page_size\":10,\"sorts\":[{\"field\":\"id\",\"order\":\"desc\"}],\"conditions\":[]}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"list\\\":[{\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-18 20:51:52 +0800 CST\\\",\\\"id\\\":3,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\"},{\\\"id\\\":2,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-18 20:44:30 +0800 CST\\\"},{\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-18 20:43:46 +0800 CST\\\",\\\"id\\\":1}],\\\"page\\\":2,\\\"page_size\\\":10,\\\"total\\\":13},\\\"trace_id\\\":\\\"31029065-abc4-49e6-bb0b-e249bbb55a45\\\"}\"', 200, '60.042791ms', '2024-01-21 20:44:59', '2024-01-21 20:44:58');
INSERT INTO `operation_log` VALUES (11, 3, 'admin@qq.com', '127.0.0.1', '', 'User', '获取用户登录历史', '/api/v1/user/login_history', 'POST', '', '{\"page\":1,\"page_size\":10,\"sorts\":[{\"field\":\"id\",\"order\":\"desc\"}],\"conditions\":[]}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"list\\\":[{\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:47:51 +0800 CST\\\",\\\"id\\\":13},{\\\"id\\\":12,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:35:00 +0800 CST\\\"},{\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:29:40 +0800 CST\\\",\\\"id\\\":11},{\\\"login_time\\\":\\\"2024-01-19 17:29:13 +0800 CST\\\",\\\"id\\\":10,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\"},{\\\"login_time\\\":\\\"2024-01-19 17:28:25 +0800 CST\\\",\\\"id\\\":9,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\"},{\\\"id\\\":8,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 15:38:21 +0800 CST\\\"},{\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 11:33:58 +0800 CST\\\",\\\"id\\\":7,\\\"login_type\\\":\\\"email\\\"},{\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:27:05 +0800 CST\\\",\\\"id\\\":6,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\"},{\\\"id\\\":5,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:17:16 +0800 CST\\\"},{\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:16:10 +0800 CST\\\",\\\"id\\\":4}],\\\"page\\\":1,\\\"page_size\\\":10,\\\"total\\\":13},\\\"trace_id\\\":\\\"1417a90e-c4c0-443e-9285-933fdd461283\\\"}\"', 200, '46.952125ms', '2024-01-21 20:45:03', '2024-01-21 20:45:02');
INSERT INTO `operation_log` VALUES (12, 3, 'admin@qq.com', '127.0.0.1', '', 'User', '获取用户登录历史', '/api/v1/user/login_history', 'POST', '', '{\"page_size\":10,\"sorts\":[{\"field\":\"id\",\"order\":\"desc\"}],\"conditions\":[],\"page\":1}', '\"{\\\"trace_id\\\":\\\"9a34f4d7-2f7a-4648-8eac-fcba1a82a0b8\\\",\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"list\\\":[{\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:47:51 +0800 CST\\\",\\\"id\\\":13,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\"},{\\\"login_time\\\":\\\"2024-01-19 17:35:00 +0800 CST\\\",\\\"id\\\":12,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\"},{\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:29:40 +0800 CST\\\",\\\"id\\\":11,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\"},{\\\"id\\\":10,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:29:13 +0800 CST\\\"},{\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:28:25 +0800 CST\\\",\\\"id\\\":9,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\"},{\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 15:38:21 +0800 CST\\\",\\\"id\\\":8,\\\"login_type\\\":\\\"email\\\"},{\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 11:33:58 +0800 CST\\\",\\\"id\\\":7,\\\"login_type\\\":\\\"email\\\"},{\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:27:05 +0800 CST\\\",\\\"id\\\":6,\\\"login_type\\\":\\\"email\\\"},{\\\"id\\\":5,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:17:16 +0800 CST\\\"},{\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:16:10 +0800 CST\\\",\\\"id\\\":4}],\\\"page\\\":1,\\\"page_size\\\":10,\\\"total\\\":13}}\"', 200, '48.113709ms', '2024-01-21 20:45:18', '2024-01-21 20:45:18');
INSERT INTO `operation_log` VALUES (13, 3, 'admin@qq.com', '127.0.0.1', '', 'User', '获取用户登录历史', '/api/v1/user/login_history', 'POST', '', '{\"page_size\":10,\"sorts\":[{\"field\":\"id\",\"order\":\"desc\"}],\"conditions\":[],\"page\":1}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"list\\\":[{\\\"id\\\":13,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:47:51 +0800 CST\\\"},{\\\"id\\\":12,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:35:00 +0800 CST\\\"},{\\\"login_time\\\":\\\"2024-01-19 17:29:40 +0800 CST\\\",\\\"id\\\":11,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\"},{\\\"login_time\\\":\\\"2024-01-19 17:29:13 +0800 CST\\\",\\\"id\\\":10,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\"},{\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:28:25 +0800 CST\\\",\\\"id\\\":9,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\"},{\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 15:38:21 +0800 CST\\\",\\\"id\\\":8,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\"},{\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 11:33:58 +0800 CST\\\",\\\"id\\\":7,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\"},{\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:27:05 +0800 CST\\\",\\\"id\\\":6},{\\\"id\\\":5,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:17:16 +0800 CST\\\"},{\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:16:10 +0800 CST\\\",\\\"id\\\":4,\\\"login_type\\\":\\\"email\\\"}],\\\"page\\\":1,\\\"page_size\\\":10,\\\"total\\\":13},\\\"trace_id\\\":\\\"51160d81-93a0-4e35-a619-ed6f023a94cd\\\"}\"', 200, '47.505167ms', '2024-01-21 20:45:20', '2024-01-21 20:45:19');
INSERT INTO `operation_log` VALUES (14, 3, 'admin@qq.com', '127.0.0.1', '', 'User', '获取用户登录历史', '/api/v1/user/login_history', 'POST', '', '{\"page_size\":10,\"sorts\":[{\"field\":\"id\",\"order\":\"desc\"}],\"conditions\":[],\"page\":1}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"list\\\":[{\\\"login_time\\\":\\\"2024-01-19 17:47:51 +0800 CST\\\",\\\"id\\\":13,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\"},{\\\"id\\\":12,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:35:00 +0800 CST\\\"},{\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:29:40 +0800 CST\\\",\\\"id\\\":11},{\\\"id\\\":10,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:29:13 +0800 CST\\\"},{\\\"id\\\":9,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 17:28:25 +0800 CST\\\"},{\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 15:38:21 +0800 CST\\\",\\\"id\\\":8,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\"},{\\\"id\\\":7,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 11:33:58 +0800 CST\\\"},{\\\"id\\\":6,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:27:05 +0800 CST\\\"},{\\\"id\\\":5,\\\"login_type\\\":\\\"email\\\",\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:17:16 +0800 CST\\\"},{\\\"agent\\\":\\\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"login_time\\\":\\\"2024-01-19 10:16:10 +0800 CST\\\",\\\"id\\\":4,\\\"login_type\\\":\\\"email\\\"}],\\\"page\\\":1,\\\"page_size\\\":10,\\\"total\\\":13},\\\"trace_id\\\":\\\"3846e8a1-416b-4017-ad0e-1767eb7a3610\\\"}\"', 200, '45.356042ms', '2024-01-21 20:45:22', '2024-01-21 20:45:22');
INSERT INTO `operation_log` VALUES (15, 3, 'admin@qq.com', '127.0.0.1', '', 'Website', '获取配置', '/api/v1/admin/config', 'POST', '', '{\"key\":\"website_config\"}', '\"{\\\"trace_id\\\":\\\"237697a8-f5ae-4635-bfab-6f916491facf\\\",\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":\\\"{\\\\\\\"alipay_qr_code\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/17f234dc487c1bb5bbb732869be0eb53.jpg\\\\\\\",\\\\\\\"is_chat_room\\\\\\\":1,\\\\\\\"is_reward\\\\\\\":1,\\\\\\\"tourist_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif\\\\\\\",\\\\\\\"website_create_time\\\\\\\":\\\\\\\"2022-01-19\\\\\\\",\\\\\\\"gitee\\\\\\\":\\\\\\\"https://gitee.com/wy791422171\\\\\\\",\\\\\\\"is_music_player\\\\\\\":0,\\\\\\\"website_author\\\\\\\":\\\\\\\"与梦\\\\\\\",\\\\\\\"website_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/84aa08357bf6e74fc1d4f33552475f91.gif\\\\\\\",\\\\\\\"website_intro\\\\\\\":\\\\\\\"分享美好生活。\\\\\\\",\\\\\\\"wei_xin_qr_code\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/6bed8a1130b170546341ece729e8819f.jpg\\\\\\\",\\\\\\\"github\\\\\\\":\\\\\\\"https://github.com/7914-ve\\\\\\\",\\\\\\\"is_email_notice\\\\\\\":1,\\\\\\\"is_message_review\\\\\\\":0,\\\\\\\"social_url_list\\\\\\\":[\\\\\\\"qq\\\\\\\",\\\\\\\"github\\\\\\\",\\\\\\\"gitee\\\\\\\"],\\\\\\\"user_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg\\\\\\\",\\\\\\\"website_notice\\\\\\\":\\\\\\\"用户需要查看、发表文章、修改其他信息请登录后台管理系统。网站后台管理系统->https://veweiyi.cn/admin。     \\\\\\\\n网站搭建问题请联系站长QQ791422171。\\\\\\\",\\\\\\\"is_comment_review\\\\\\\":0,\\\\\\\"qq\\\\\\\":\\\\\\\"791422171\\\\\\\",\\\\\\\"social_login_list\\\\\\\":[\\\\\\\"weibo\\\\\\\",\\\\\\\"qq\\\\\\\",\\\\\\\"feishu\\\\\\\"],\\\\\\\"website_name\\\\\\\":\\\\\\\"静闻弦语\\\\\\\",\\\\\\\"website_record_no\\\\\\\":\\\\\\\"桂ICP备2023013735号-1\\\\\\\",\\\\\\\"websocket_url\\\\\\\":\\\\\\\"wss://ve77.cn:8088/api/websocket\\\\\\\"}\\\"}\"', 200, '16.445ms', '2024-01-21 20:45:30', '2024-01-21 20:45:30');
INSERT INTO `operation_log` VALUES (16, 3, 'admin@qq.com', '127.0.0.1', '', 'User', '修改用户状态', '/api/v1/user/update_status', 'POST', '', '{\"website\":\"\",\"created_at\":\"2023-07-04T11:50:22+08:00\",\"id\":11,\"username\":\"ou_ef9abf85b8905510f005587d9717a596\",\"nickname\":\"游客d21e88ea\",\"avatar\":\"https://mms1.baidu.com/it/u=2815887849,1501151317&fm=253&app=138&f=JPEG\",\"intro\":\"这个人很神秘，什么都没有写！\",\"status\":1,\"register_type\":\"feishu\",\"updated_at\":\"2024-01-14T01:14:13+08:00\",\"email\":\"\",\"ip_address\":\"127.0.0.1\",\"ip_source\":\"广西壮族自治区梧州市 移动\",\"roles\":[{\"role_name\":\"user\",\"role_comment\":\"用户\"}]}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"id\\\":11,\\\"username\\\":\\\"ou_ef9abf85b8905510f005587d9717a596\\\",\\\"status\\\":1,\\\"register_type\\\":\\\"feishu\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"created_at\\\":\\\"2023-07-04T11:50:22+08:00\\\",\\\"updated_at\\\":\\\"2024-01-21T20:45:59.977+08:00\\\",\\\"password\\\":\\\"6\\\",\\\"ip_source\\\":\\\"广西壮族自治区梧州市 移动\\\"},\\\"trace_id\\\":\\\"8801942f-07b9-431a-bc71-17ed67e4e3bd\\\"}\"', 200, '55.202042ms', '2024-01-21 20:46:00', '2024-01-21 20:46:00');
INSERT INTO `operation_log` VALUES (17, 3, 'admin@qq.com', '127.0.0.1', '', 'User', '修改用户状态', '/api/v1/user/update_status', 'POST', '', '{\"updated_at\":\"2024-01-14T01:14:13+08:00\",\"nickname\":\"游客d21e88ea\",\"intro\":\"这个人很神秘，什么都没有写！\",\"website\":\"\",\"register_type\":\"feishu\",\"created_at\":\"2023-07-04T11:50:22+08:00\",\"status\":0,\"roles\":[{\"role_name\":\"user\",\"role_comment\":\"用户\"}],\"username\":\"ou_ef9abf85b8905510f005587d9717a596\",\"email\":\"\",\"ip_address\":\"127.0.0.1\",\"id\":11,\"avatar\":\"https://mms1.baidu.com/it/u=2815887849,1501151317&fm=253&app=138&f=JPEG\",\"ip_source\":\"广西壮族自治区梧州市 移动\"}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"id\\\":11,\\\"username\\\":\\\"ou_ef9abf85b8905510f005587d9717a596\\\",\\\"register_type\\\":\\\"feishu\\\",\\\"updated_at\\\":\\\"2024-01-21T20:46:01.594+08:00\\\",\\\"password\\\":\\\"6\\\",\\\"status\\\":0,\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"广西壮族自治区梧州市 移动\\\",\\\"created_at\\\":\\\"2023-07-04T11:50:22+08:00\\\"},\\\"trace_id\\\":\\\"5eb7cf42-67e4-4b0c-900d-1f7f7492658e\\\"}\"', 200, '44.675583ms', '2024-01-21 20:46:02', '2024-01-21 20:46:01');
INSERT INTO `operation_log` VALUES (18, 3, 'admin@qq.com', '127.0.0.1', '', 'Website', '获取配置', '/api/v1/admin/config', 'POST', '', '{\"key\":\"website_config\"}', '\"{\\\"data\\\":\\\"{\\\\\\\"alipay_qr_code\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/17f234dc487c1bb5bbb732869be0eb53.jpg\\\\\\\",\\\\\\\"is_chat_room\\\\\\\":1,\\\\\\\"is_reward\\\\\\\":1,\\\\\\\"tourist_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif\\\\\\\",\\\\\\\"website_create_time\\\\\\\":\\\\\\\"2022-01-19\\\\\\\",\\\\\\\"gitee\\\\\\\":\\\\\\\"https://gitee.com/wy791422171\\\\\\\",\\\\\\\"is_music_player\\\\\\\":0,\\\\\\\"website_author\\\\\\\":\\\\\\\"与梦\\\\\\\",\\\\\\\"website_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/84aa08357bf6e74fc1d4f33552475f91.gif\\\\\\\",\\\\\\\"website_intro\\\\\\\":\\\\\\\"分享美好生活。\\\\\\\",\\\\\\\"wei_xin_qr_code\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/6bed8a1130b170546341ece729e8819f.jpg\\\\\\\",\\\\\\\"github\\\\\\\":\\\\\\\"https://github.com/7914-ve\\\\\\\",\\\\\\\"is_email_notice\\\\\\\":1,\\\\\\\"is_message_review\\\\\\\":0,\\\\\\\"social_url_list\\\\\\\":[\\\\\\\"qq\\\\\\\",\\\\\\\"github\\\\\\\",\\\\\\\"gitee\\\\\\\"],\\\\\\\"user_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg\\\\\\\",\\\\\\\"website_notice\\\\\\\":\\\\\\\"用户需要查看、发表文章、修改其他信息请登录后台管理系统。网站后台管理系统->https://veweiyi.cn/admin。     \\\\\\\\n网站搭建问题请联系站长QQ791422171。\\\\\\\",\\\\\\\"is_comment_review\\\\\\\":0,\\\\\\\"qq\\\\\\\":\\\\\\\"791422171\\\\\\\",\\\\\\\"social_login_list\\\\\\\":[\\\\\\\"weibo\\\\\\\",\\\\\\\"qq\\\\\\\",\\\\\\\"feishu\\\\\\\"],\\\\\\\"website_name\\\\\\\":\\\\\\\"静闻弦语\\\\\\\",\\\\\\\"website_record_no\\\\\\\":\\\\\\\"桂ICP备2023013735号-1\\\\\\\",\\\\\\\"websocket_url\\\\\\\":\\\\\\\"wss://ve77.cn:8088/api/websocket\\\\\\\"}\\\",\\\"trace_id\\\":\\\"a0a951c4-e6b6-4e82-bdb4-bd4d962a3dc1\\\",\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\"}\"', 200, '21.936958ms', '2024-01-21 20:47:24', '2024-01-21 20:47:24');
INSERT INTO `operation_log` VALUES (19, 3, 'admin@qq.com', '127.0.0.1', '', 'Website', '获取配置', '/api/v1/admin/config', 'POST', '', '{\"key\":\"website_config\"}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":\\\"{\\\\\\\"alipay_qr_code\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/17f234dc487c1bb5bbb732869be0eb53.jpg\\\\\\\",\\\\\\\"is_chat_room\\\\\\\":1,\\\\\\\"is_reward\\\\\\\":1,\\\\\\\"tourist_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif\\\\\\\",\\\\\\\"website_create_time\\\\\\\":\\\\\\\"2022-01-19\\\\\\\",\\\\\\\"gitee\\\\\\\":\\\\\\\"https://gitee.com/wy791422171\\\\\\\",\\\\\\\"is_music_player\\\\\\\":0,\\\\\\\"website_author\\\\\\\":\\\\\\\"与梦\\\\\\\",\\\\\\\"website_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/84aa08357bf6e74fc1d4f33552475f91.gif\\\\\\\",\\\\\\\"website_intro\\\\\\\":\\\\\\\"分享美好生活。\\\\\\\",\\\\\\\"wei_xin_qr_code\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/6bed8a1130b170546341ece729e8819f.jpg\\\\\\\",\\\\\\\"github\\\\\\\":\\\\\\\"https://github.com/7914-ve\\\\\\\",\\\\\\\"is_email_notice\\\\\\\":1,\\\\\\\"is_message_review\\\\\\\":0,\\\\\\\"social_url_list\\\\\\\":[\\\\\\\"qq\\\\\\\",\\\\\\\"github\\\\\\\",\\\\\\\"gitee\\\\\\\"],\\\\\\\"user_avatar\\\\\\\":\\\\\\\"https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg\\\\\\\",\\\\\\\"website_notice\\\\\\\":\\\\\\\"用户需要查看、发表文章、修改其他信息请登录后台管理系统。网站后台管理系统->https://veweiyi.cn/admin。     \\\\\\\\n网站搭建问题请联系站长QQ791422171。\\\\\\\",\\\\\\\"is_comment_review\\\\\\\":0,\\\\\\\"qq\\\\\\\":\\\\\\\"791422171\\\\\\\",\\\\\\\"social_login_list\\\\\\\":[\\\\\\\"weibo\\\\\\\",\\\\\\\"qq\\\\\\\",\\\\\\\"feishu\\\\\\\"],\\\\\\\"website_name\\\\\\\":\\\\\\\"静闻弦语\\\\\\\",\\\\\\\"website_record_no\\\\\\\":\\\\\\\"桂ICP备2023013735号-1\\\\\\\",\\\\\\\"websocket_url\\\\\\\":\\\\\\\"wss://ve77.cn:8088/api/websocket\\\\\\\"}\\\",\\\"trace_id\\\":\\\"3a6e7dc7-7a0f-45c9-b037-09a3affbd984\\\"}\"', 200, '18.433334ms', '2024-01-21 20:56:52', '2024-01-21 20:56:51');
INSERT INTO `operation_log` VALUES (20, 3, 'admin@qq.com', '127.0.0.1', '', 'Remark', '创建留言', '/api/v1/remark', 'POST', '', '{\"avatar\":\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\",\"nickname\":\"admin@qq.com\",\"messageContent\":\"hhh\",\"time\":7}', '\"{\\\"data\\\":{\\\"avatar\\\":\\\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\\\",\\\"message_content\\\":\\\"\\\",\\\"ip_address\\\":\\\"\\\",\\\"created_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"nickname\\\":\\\"admin@qq.com\\\",\\\"ip_source\\\":\\\"\\\",\\\"time\\\":7,\\\"is_review\\\":1,\\\"updated_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"id\\\":3918},\\\"trace_id\\\":\\\"926b9664-4fbb-4c0c-aa18-b6b57e7cfd31\\\",\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\"}\"', 200, '38.621916ms', '2024-01-23 15:47:56', '2024-01-23 15:47:56');
INSERT INTO `operation_log` VALUES (21, 3, 'admin@qq.com', '127.0.0.1', '', 'Remark', '创建留言', '/api/v1/remark', 'POST', '', '{\"avatar\":\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\",\"nickname\":\"admin@qq.com\",\"messageContent\":\"hhh\",\"time\":8}', '\"{\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"nickname\\\":\\\"admin@qq.com\\\",\\\"avatar\\\":\\\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\\\",\\\"time\\\":8,\\\"is_review\\\":1,\\\"created_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"updated_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"id\\\":3919,\\\"message_content\\\":\\\"\\\",\\\"ip_address\\\":\\\"\\\",\\\"ip_source\\\":\\\"\\\"},\\\"trace_id\\\":\\\"a489a288-1dc9-4020-8c92-9a87d51efab4\\\",\\\"code\\\":200}\"', 200, '33.223417ms', '2024-01-23 15:48:00', '2024-01-23 15:48:00');
INSERT INTO `operation_log` VALUES (22, 3, 'admin@qq.com', '127.0.0.1', '', 'Remark', '创建留言', '/api/v1/remark', 'POST', '', '{\"time\":8,\"avatar\":\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\",\"nickname\":\"admin@qq.com\",\"messageContent\":\"sss\"}', '\"{\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"ip_source\\\":\\\"\\\",\\\"is_review\\\":1,\\\"created_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"updated_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"id\\\":3920,\\\"nickname\\\":\\\"admin@qq.com\\\",\\\"avatar\\\":\\\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\\\",\\\"message_content\\\":\\\"\\\",\\\"ip_address\\\":\\\"\\\",\\\"time\\\":8},\\\"trace_id\\\":\\\"c4854f6e-e35d-42bf-942b-5e2e7a69c68c\\\",\\\"code\\\":200}\"', 200, '40.504792ms', '2024-01-23 16:04:24', '2024-01-23 16:04:24');
INSERT INTO `operation_log` VALUES (23, 3, 'admin@qq.com', '127.0.0.1', '', 'Remark', '创建留言', '/api/v1/remark', 'POST', '', '{\"time\":9,\"avatar\":\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\",\"nickname\":\"admin@qq.com\",\"messageContent\":\"sss\"}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"id\\\":3921,\\\"ip_address\\\":\\\"\\\",\\\"ip_source\\\":\\\"\\\",\\\"is_review\\\":1,\\\"nickname\\\":\\\"admin@qq.com\\\",\\\"avatar\\\":\\\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\\\",\\\"message_content\\\":\\\"\\\",\\\"time\\\":9,\\\"created_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"updated_at\\\":\\\"0001-01-01T00:00:00Z\\\"},\\\"trace_id\\\":\\\"8594b261-6ecb-47f3-be9f-10b5ae00f9cf\\\"}\"', 200, '36.300041ms', '2024-01-23 16:04:29', '2024-01-23 16:04:29');
INSERT INTO `operation_log` VALUES (24, 3, 'admin@qq.com', '127.0.0.1', '', 'Remark', '创建留言', '/api/v1/remark', 'POST', '', '{\"nickname\":\"admin@qq.com\",\"message_content\":\"ss\",\"time\":7,\"avatar\":\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\"}', '\"{\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"ip_address\\\":\\\"\\\",\\\"ip_source\\\":\\\"\\\",\\\"time\\\":7,\\\"is_review\\\":1,\\\"id\\\":3922,\\\"message_content\\\":\\\"ss\\\",\\\"created_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"updated_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"nickname\\\":\\\"admin@qq.com\\\",\\\"avatar\\\":\\\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\\\"},\\\"trace_id\\\":\\\"377f02ea-0e27-4ee8-909a-6fd8a4ebceff\\\",\\\"code\\\":200}\"', 200, '45.192334ms', '2024-01-23 16:05:31', '2024-01-23 16:05:31');
INSERT INTO `operation_log` VALUES (25, 3, 'admin@qq.com', '127.0.0.1', '', 'Remark', '创建留言', '/api/v1/remark', 'POST', '', '{\"avatar\":\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\",\"nickname\":\"admin@qq.com\",\"message_content\":\"hhh\",\"time\":8}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"avatar\\\":\\\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\\\",\\\"message_content\\\":\\\"hhh\\\",\\\"ip_address\\\":\\\"\\\",\\\"time\\\":8,\\\"created_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"updated_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"id\\\":3923,\\\"ip_source\\\":\\\"\\\",\\\"is_review\\\":1,\\\"nickname\\\":\\\"admin@qq.com\\\"},\\\"trace_id\\\":\\\"b3da15cf-72d3-4587-b96c-b80b5347a4b6\\\"}\"', 200, '32.549125ms', '2024-01-23 16:05:38', '2024-01-23 16:05:37');
INSERT INTO `operation_log` VALUES (26, 3, 'admin@qq.com', '127.0.0.1', '', 'Remark', '创建留言', '/api/v1/remark', 'POST', '', '{\"message_content\":\"ss\",\"time\":9,\"avatar\":\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\",\"nickname\":\"admin@qq.com\"}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"is_review\\\":1,\\\"updated_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"id\\\":3924,\\\"nickname\\\":\\\"admin@qq.com\\\",\\\"avatar\\\":\\\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\\\",\\\"message_content\\\":\\\"ss\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"time\\\":9,\\\"created_at\\\":\\\"0001-01-01T00:00:00Z\\\"},\\\"trace_id\\\":\\\"87c62978-57b0-411e-8f2f-0bc2fabab099\\\"}\"', 200, '58.3ms', '2024-01-23 16:36:10', '2024-01-23 16:36:10');
INSERT INTO `operation_log` VALUES (27, 3, 'admin@qq.com', '127.0.0.1', '', 'Remark', '创建留言', '/api/v1/remark', 'POST', '', '{\"avatar\":\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\",\"nickname\":\"admin@qq.com\",\"message_content\":\"ss\",\"time\":9}', '\"{\\\"data\\\":{\\\"time\\\":9,\\\"is_review\\\":1,\\\"updated_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"id\\\":3925,\\\"nickname\\\":\\\"admin@qq.com\\\",\\\"message_content\\\":\\\"ss\\\",\\\"created_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"avatar\\\":\\\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"ip_source\\\":\\\"本机地址\\\"},\\\"trace_id\\\":\\\"80f89ec7-e9b2-4480-a21d-3ec333977ec4\\\",\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\"}\"', 200, '32.976708ms', '2024-01-23 16:43:11', '2024-01-23 16:43:11');
INSERT INTO `operation_log` VALUES (28, 3, 'admin@qq.com', '127.0.0.1', '', 'Remark', '创建留言', '/api/v1/remark', 'POST', '', '{\"message_content\":\"哈哈哈哈\",\"time\":9,\"avatar\":\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\",\"nickname\":\"admin@qq.com\"}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"id\\\":3926,\\\"message_content\\\":\\\"哈哈哈哈\\\",\\\"ip_address\\\":\\\"127.0.0.1\\\",\\\"time\\\":9,\\\"nickname\\\":\\\"admin@qq.com\\\",\\\"avatar\\\":\\\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\\\",\\\"ip_source\\\":\\\"本机地址\\\",\\\"is_review\\\":1,\\\"created_at\\\":\\\"0001-01-01T00:00:00Z\\\",\\\"updated_at\\\":\\\"0001-01-01T00:00:00Z\\\"},\\\"trace_id\\\":\\\"6babf2ea-1f12-4b8d-bf0d-17caef75c2cd\\\"}\"', 200, '35.599958ms', '2024-01-23 17:23:59', '2024-01-23 17:23:58');
INSERT INTO `operation_log` VALUES (29, 3, 'admin@qq.com', '127.0.0.1', '', 'User', '修改用户信息', '/api/v1/user/info', 'POST', '', '{\"nickname\":\"admin@qq.com\",\"intro\":\"测试\",\"website\":\"admin@qq.com\",\"avatar\":\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\"}', '\"{\\\"code\\\":200,\\\"message\\\":\\\"操作成功\\\",\\\"data\\\":{\\\"phone\\\":\\\"\\\",\\\"website\\\":\\\"admin@qq.com\\\",\\\"intro\\\":\\\"测试\\\",\\\"created_at\\\":\\\"2023-05-17T15:34:25+08:00\\\",\\\"updated_at\\\":\\\"2024-01-23T17:30:02.716+08:00\\\",\\\"id\\\":3,\\\"user_id\\\":3,\\\"email\\\":\\\"\\\",\\\"nickname\\\":\\\"admin@qq.com\\\",\\\"avatar\\\":\\\"https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png\\\"},\\\"trace_id\\\":\\\"9e819cc1-ba78-466a-bf98-2a3ef6a1ba81\\\"}\"', 200, '57.072333ms', '2024-01-23 17:30:03', '2024-01-23 17:30:02');
COMMIT;

-- ----------------------------
-- Table structure for page
-- ----------------------------
DROP TABLE IF EXISTS `page`;
CREATE TABLE `page` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '页面id',
  `page_name` varchar(32) NOT NULL DEFAULT '' COMMENT '页面名',
  `page_label` varchar(32) NOT NULL DEFAULT '' COMMENT '页面标签',
  `page_cover` varchar(255) NOT NULL DEFAULT '' COMMENT '页面封面',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='页面';

-- ----------------------------
-- Records of page
-- ----------------------------
BEGIN;
INSERT INTO `page` VALUES (1, '首页', 'home', 'https://veport.oss-cn-beijing.aliyuncs.com/config/f9fa18da262910eb13f802b003147915.jpg', '2021-08-07 10:32:36', '2022-01-20 21:36:15');
INSERT INTO `page` VALUES (2, '归档', 'archive', 'https://veport.oss-cn-beijing.aliyuncs.com/config/82fc9c41de3c511ca1532d978b36fec7.jpg', '2021-08-07 10:32:36', '2022-01-19 22:14:50');
INSERT INTO `page` VALUES (3, '分类', 'category', 'https://veport.oss-cn-beijing.aliyuncs.com/config/f9fa18da262910eb13f802b003147915.jpg', '2021-08-07 10:32:36', '2022-01-19 23:11:11');
INSERT INTO `page` VALUES (4, '标签', 'tag', 'https://veport.oss-cn-beijing.aliyuncs.com/config/dd3678e409cab21ff2e5f875976058a6.jpg', '2021-08-07 10:32:36', '2021-10-04 15:43:38');
INSERT INTO `page` VALUES (5, '相册', 'album', 'https://veport.oss-cn-beijing.aliyuncs.com/config/dd3678e409cab21ff2e5f875976058a6.jpg', '2021-08-07 10:32:36', '2022-01-19 22:21:47');
INSERT INTO `page` VALUES (6, '友链', 'link', 'https://veport.oss-cn-beijing.aliyuncs.com/config/8b03884995623eab1a76772f23b58875.jpg', '2021-08-07 10:32:36', '2022-01-20 21:36:35');
INSERT INTO `page` VALUES (7, '关于', 'about', 'https://veport.oss-cn-beijing.aliyuncs.com/config/3a4b4e40fb8aa5fcc016f0228938d321.jpg', '2021-08-07 10:32:36', '2022-01-19 23:10:52');
INSERT INTO `page` VALUES (8, '留言', 'message', 'https://veport.oss-cn-beijing.aliyuncs.com/config/75e976f3364ba013d62e99ff3ab65d19.jpg', '2021-08-07 10:32:36', '2022-01-19 22:13:52');
INSERT INTO `page` VALUES (9, '个人中心', 'user', 'https://veport.oss-cn-beijing.aliyuncs.com/config/4e319068b295ca52080979d5653c334d.jpg', '2021-08-07 10:32:36', '2022-01-19 22:21:03');
INSERT INTO `page` VALUES (10, '文章列表', 'articleList', 'https://veport.oss-cn-beijing.aliyuncs.com/config/3a4b4e40fb8aa5fcc016f0228938d321.jpg', '2021-08-10 15:36:19', '2022-01-19 22:21:12');
INSERT INTO `page` VALUES (11, '说说', 'talk', 'https://veport.oss-cn-beijing.aliyuncs.com/config/f9fa18da262910eb13f802b003147915.jpg', '2022-01-23 00:51:24', '2022-02-11 12:23:15');
COMMIT;

-- ----------------------------
-- Table structure for photo
-- ----------------------------
DROP TABLE IF EXISTS `photo`;
CREATE TABLE `photo` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `album_id` int NOT NULL DEFAULT '0' COMMENT '相册id',
  `photo_name` varchar(32) NOT NULL DEFAULT '' COMMENT '照片名',
  `photo_desc` varchar(64) NOT NULL DEFAULT '' COMMENT '照片描述',
  `photo_src` varchar(255) NOT NULL DEFAULT '' COMMENT '照片地址',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='照片';

-- ----------------------------
-- Records of photo
-- ----------------------------
BEGIN;
INSERT INTO `photo` VALUES (31, 10, '1484534479992193025', '', 'https://veport.oss-cn-beijing.aliyuncs.com/photos/cfeb11ab6be04ca78f24a0d8974a296d.png', 0, '2022-01-21 22:32:56', '2024-01-15 19:43:08');
INSERT INTO `photo` VALUES (32, 10, '1484534479992193026', '', 'https://veport.oss-cn-beijing.aliyuncs.com/photos/dd3678e409cab21ff2e5f875976058a6.jpg', 0, '2022-01-21 22:32:56', '2024-01-15 19:43:08');
COMMIT;

-- ----------------------------
-- Table structure for photo_album
-- ----------------------------
DROP TABLE IF EXISTS `photo_album`;
CREATE TABLE `photo_album` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `album_name` varchar(32) NOT NULL DEFAULT '' COMMENT '相册名',
  `album_desc` varchar(64) NOT NULL DEFAULT '' COMMENT '相册描述',
  `album_cover` varchar(255) NOT NULL DEFAULT '' COMMENT '相册封面',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态值 1公开 2私密',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='相册';

-- ----------------------------
-- Records of photo_album
-- ----------------------------
BEGIN;
INSERT INTO `photo_album` VALUES (10, '图片', '图片', 'https://veport.oss-cn-beijing.aliyuncs.com/photos/4e319068b295ca52080979d5653c334d.jpg', 0, 1, '2022-01-19 22:20:25', '2024-01-15 19:43:08');
INSERT INTO `photo_album` VALUES (11, '', '', '', 0, 1, '2023-10-30 20:05:40', '2023-10-30 20:05:40');
INSERT INTO `photo_album` VALUES (12, '', '', '', 0, 2, '2023-10-30 20:10:26', '2023-10-30 20:10:26');
COMMIT;

-- ----------------------------
-- Table structure for remark
-- ----------------------------
DROP TABLE IF EXISTS `remark`;
CREATE TABLE `remark` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `nickname` varchar(64) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `message_content` varchar(255) NOT NULL DEFAULT '' COMMENT '留言内容',
  `ip_address` varchar(64) NOT NULL DEFAULT '' COMMENT '用户ip',
  `ip_source` varchar(255) NOT NULL DEFAULT '' COMMENT '用户地址',
  `time` int NOT NULL DEFAULT '0' COMMENT '弹幕速度',
  `is_review` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否审核',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3927 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='留言';

-- ----------------------------
-- Records of remark
-- ----------------------------
BEGIN;
INSERT INTO `remark` VALUES (3906, '管理员', 'https://veport.oss-cn-beijing.aliyuncs.com/icon/qqface/1.png', '测试留言', '127.0.0.1', '', 7, 1, '2022-01-18 00:31:12', '2024-01-15 19:43:08');
INSERT INTO `remark` VALUES (3910, '游客', 'https://veport.oss-cn-beijing.aliyuncs.com/config/59bf7acbb23bc6697ce334b978218d63.png', '77到此一游', '110.42.180.40', '北京市北京市', 8, 1, '2022-01-20 21:49:12', '2024-01-15 19:43:08');
INSERT INTO `remark` VALUES (3911, '梦梦', 'https://veport.oss-cn-beijing.aliyuncs.com/avatar/8b9ef1e87e98892795eed286c8778d75.png', '测试留言', '113.15.88.31', '广西壮族自治区来宾市 电信', 8, 1, '2022-01-21 11:47:59', '2024-01-15 19:43:08');
INSERT INTO `remark` VALUES (3914, '游客', 'https://veport.oss-cn-beijing.aliyuncs.com/config/6287609448bab8b7b627aa93fe270cf1.jpg', '图片挺好看', '113.15.88.31', '广西壮族自治区来宾市 电信', 7, 1, '2022-01-23 13:07:30', '2024-01-15 19:43:08');
INSERT INTO `remark` VALUES (3915, '游客', 'https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif', '嘿嘿', '117.182.36.251', '广西壮族自治区河池市 移动', 9, 1, '2022-01-25 04:37:24', '2024-01-15 19:43:08');
INSERT INTO `remark` VALUES (3916, '游客', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', 'ss', '', '', 9, 1, '2023-07-16 01:52:16', '2023-07-16 01:52:16');
INSERT INTO `remark` VALUES (3917, '游客', 'https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg', 'hh', '', '', 7, 1, '2023-07-16 01:53:52', '2023-07-16 01:53:52');
INSERT INTO `remark` VALUES (3918, 'admin@qq.com', 'https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png', '', '', '', 7, 1, '2024-01-23 15:47:56', '2024-01-23 15:47:56');
INSERT INTO `remark` VALUES (3919, 'admin@qq.com', 'https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png', '', '', '', 8, 1, '2024-01-23 15:48:00', '2024-01-23 15:48:00');
INSERT INTO `remark` VALUES (3920, 'admin@qq.com', 'https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png', '', '', '', 8, 1, '2024-01-23 16:04:24', '2024-01-23 16:04:24');
INSERT INTO `remark` VALUES (3921, 'admin@qq.com', 'https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png', '', '', '', 9, 1, '2024-01-23 16:04:29', '2024-01-23 16:04:29');
INSERT INTO `remark` VALUES (3922, 'admin@qq.com', 'https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png', 'ss', '', '', 7, 1, '2024-01-23 16:05:31', '2024-01-23 16:05:31');
INSERT INTO `remark` VALUES (3923, 'admin@qq.com', 'https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png', 'hhh', '', '', 8, 1, '2024-01-23 16:05:37', '2024-01-23 16:05:37');
INSERT INTO `remark` VALUES (3924, 'admin@qq.com', 'https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png', 'ss', '127.0.0.1', '本机地址', 9, 1, '2024-01-23 16:36:10', '2024-01-23 16:36:10');
INSERT INTO `remark` VALUES (3925, 'admin@qq.com', 'https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png', 'ss', '127.0.0.1', '本机地址', 9, 1, '2024-01-23 16:43:11', '2024-01-23 16:43:11');
INSERT INTO `remark` VALUES (3926, 'admin@qq.com', 'https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png', '哈哈哈哈', '127.0.0.1', '本机地址', 9, 1, '2024-01-23 17:23:58', '2024-01-23 17:23:58');
COMMIT;

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `role_pid` int NOT NULL DEFAULT '0' COMMENT '父角色id',
  `role_domain` varchar(64) NOT NULL DEFAULT '0' COMMENT '角色域',
  `role_name` varchar(64) NOT NULL DEFAULT '' COMMENT '角色名',
  `role_comment` varchar(64) NOT NULL DEFAULT '' COMMENT '角色备注',
  `is_disable` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否禁用  0否 1是',
  `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认角色 0否 1是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色';

-- ----------------------------
-- Records of role
-- ----------------------------
BEGIN;
INSERT INTO `role` VALUES (1, 0, 'blog', 'admin', '管理员', 1, 0, '2021-03-22 14:10:21', '2024-01-14 01:13:34');
INSERT INTO `role` VALUES (2, 0, 'blog', 'user', '用户', 1, 0, '2021-03-22 14:25:25', '2024-01-15 20:09:15');
INSERT INTO `role` VALUES (3, 0, 'system', 'test', '测试', 1, 0, '2021-03-22 14:42:23', '2024-01-15 20:09:14');
INSERT INTO `role` VALUES (4, 0, 'blog', 'super-admin', '超级管理员', 1, 0, '2023-05-30 20:53:04', '2024-01-18 19:52:19');
COMMIT;

-- ----------------------------
-- Table structure for role_api
-- ----------------------------
DROP TABLE IF EXISTS `role_api`;
CREATE TABLE `role_api` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `role_id` int NOT NULL DEFAULT '0' COMMENT '角色id',
  `api_id` int NOT NULL DEFAULT '0' COMMENT '接口id',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色-api关联';

-- ----------------------------
-- Records of role_api
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for role_menu
-- ----------------------------
DROP TABLE IF EXISTS `role_menu`;
CREATE TABLE `role_menu` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `role_id` int NOT NULL DEFAULT '0' COMMENT '角色id',
  `menu_id` int NOT NULL DEFAULT '0' COMMENT '菜单id',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色-菜单关联';

-- ----------------------------
-- Records of role_menu
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tag
-- ----------------------------
DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `tag_name` varchar(32) NOT NULL DEFAULT '' COMMENT '标签名',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='标签';

-- ----------------------------
-- Records of tag
-- ----------------------------
BEGIN;
INSERT INTO `tag` VALUES (29, '测试标签', '2022-01-18 00:29:02', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (30, '学习', '2022-01-19 22:35:50', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (31, 'spring', '2022-01-19 22:35:50', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (32, 'vue', '2022-01-19 22:35:50', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (33, '网站', '2022-01-21 12:04:02', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (34, '技术', '2022-01-21 12:21:31', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (35, '算法', '2022-01-22 12:59:49', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (36, 'tomcat', '2022-01-22 13:27:50', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (37, 'nginx', '2022-01-22 23:55:55', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (38, 'websocket', '2022-02-10 10:57:51', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (39, 'springboot', '2022-02-11 23:31:32', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (40, 'https', '2022-02-11 23:31:32', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (41, '算法日记', '2022-02-12 16:28:13', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (42, '动态规划', '2022-02-12 16:28:13', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (43, '背包问题', '2022-02-12 16:28:13', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (44, 'session', '2022-02-15 17:25:15', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (45, 'token', '2022-02-15 17:25:15', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (46, 'cookie', '2022-02-15 17:25:15', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (47, '冒泡排序', '2022-02-19 10:58:00', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (48, '面试笔记', '2022-02-19 11:31:18', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (49, '线程池', '2022-03-11 02:59:54', '2024-01-15 19:43:08');
INSERT INTO `tag` VALUES (50, 'hashmap', '2022-03-11 03:24:24', '2024-01-15 19:43:08');
COMMIT;

-- ----------------------------
-- Table structure for talk
-- ----------------------------
DROP TABLE IF EXISTS `talk`;
CREATE TABLE `talk` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '说说id',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户id',
  `content` varchar(2048) NOT NULL DEFAULT '' COMMENT '说说内容',
  `images` varchar(2048) NOT NULL DEFAULT '' COMMENT '图片',
  `is_top` tinyint NOT NULL DEFAULT '0' COMMENT '是否置顶',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1.公开 2.私密',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=61 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='说说';

-- ----------------------------
-- Records of talk
-- ----------------------------
BEGIN;
INSERT INTO `talk` VALUES (49, 2, '用户需要查看、发表文章、修改其他信息请登录后台管理系统。网站后台管理系统-&gt;https://ve77.cn/admin。 \n网站搭建问题请联系站长QQ791422171。\n', '', 1, 1, '2022-01-24 23:34:59', '2022-02-11 23:19:59');
INSERT INTO `talk` VALUES (50, 2, '如何搭建个人网站：必须熟练掌握Java语言 和 具备 html+css+JavaScript 基础。熟悉前端vue框架和后端springboot框架。网站发布过程中还需要了解如何在云服务器上部署 dns域名解析、https安全传输、websocket即时通讯、mysql数据库、rabbitmq消息队列、redis非关系数据库、tomcat网站部署、nginx反向代理、oss数据存储。', '', 1, 1, '2022-02-10 19:39:08', '2022-02-11 12:24:42');
INSERT INTO `talk` VALUES (52, 0, 'ssss', '', 1, 1, '2023-12-20 17:45:15', '2023-12-20 17:45:15');
INSERT INTO `talk` VALUES (53, 0, 'ssss', '', 1, 1, '2023-12-20 17:45:18', '2023-12-20 17:45:18');
INSERT INTO `talk` VALUES (58, 3, '这是一条私密说说', '', 0, 2, '2023-12-20 18:09:19', '2023-12-20 18:09:19');
INSERT INTO `talk` VALUES (59, 3, '这是一条公开说说', '', 0, 1, '2023-12-20 18:09:26', '2023-12-20 18:09:26');
INSERT INTO `talk` VALUES (60, 3, '这是一条置顶说说', '', 1, 1, '2023-12-20 18:09:37', '2023-12-20 18:58:08');
COMMIT;

-- ----------------------------
-- Table structure for unique_view
-- ----------------------------
DROP TABLE IF EXISTS `unique_view`;
CREATE TABLE `unique_view` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `views_count` int NOT NULL DEFAULT '0' COMMENT '访问量',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=703 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='页面访问数量';

-- ----------------------------
-- Records of unique_view
-- ----------------------------
BEGIN;
INSERT INTO `unique_view` VALUES (528, 1, '2022-01-18 11:51:50', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (529, 1, '2022-01-19 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (530, 9, '2022-01-20 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (531, 14, '2022-01-21 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (532, 21, '2022-01-22 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (533, 15, '2022-01-23 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (534, 12, '2022-01-24 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (535, 11, '2022-01-25 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (536, 2, '2022-01-26 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (537, 3, '2022-01-27 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (538, 1, '2022-01-28 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (539, 2, '2022-01-29 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (540, 0, '2022-01-30 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (541, 0, '2022-01-31 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (542, 0, '2022-02-01 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (543, 0, '2022-02-02 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (544, 0, '2022-02-03 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (545, 0, '2022-02-04 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (546, 0, '2022-02-05 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (547, 4, '2022-02-06 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (548, 0, '2022-02-07 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (549, 0, '2022-02-08 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (550, 4, '2022-02-09 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (551, 5, '2022-02-10 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (552, 20, '2022-02-11 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (553, 3, '2022-02-12 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (554, 9, '2022-02-13 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (555, 2, '2022-02-14 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (556, 2, '2022-02-15 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (557, 5, '2022-02-16 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (558, 1, '2022-02-17 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (559, 7, '2022-02-18 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (560, 3, '2022-02-19 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (561, 4, '2022-02-20 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (562, 14, '2022-02-21 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (563, 3, '2022-02-22 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (564, 0, '2022-02-23 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (565, 0, '2022-02-24 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (566, 0, '2022-02-25 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (567, 1, '2022-02-26 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (568, 5, '2022-02-27 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (569, 2, '2022-02-28 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (570, 1, '2022-03-01 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (571, 8, '2022-03-02 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (572, 3, '2022-03-03 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (573, 3, '2022-03-04 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (574, 1, '2022-03-05 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (575, 0, '2022-03-06 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (576, 0, '2022-03-07 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (577, 0, '2022-03-08 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (578, 1, '2022-03-09 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (579, 4, '2022-03-10 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (580, 3, '2022-03-11 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (581, 2, '2022-03-12 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (582, 0, '2022-03-13 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (583, 0, '2022-03-14 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (584, 1, '2022-03-15 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (585, 3, '2022-03-16 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (586, 5, '2022-03-17 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (587, 1, '2022-03-18 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (588, 1, '2022-03-19 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (589, 0, '2022-03-20 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (590, 4, '2022-03-21 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (591, 1, '2022-03-22 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (592, 1, '2022-03-23 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (593, 4, '2022-03-24 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (594, 0, '2022-03-25 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (595, 0, '2022-03-26 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (596, 0, '2022-03-27 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (597, 1, '2022-03-28 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (598, 1, '2022-03-29 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (599, 3, '2022-03-30 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (600, 0, '2022-03-31 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (601, 0, '2022-04-01 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (602, 1, '2022-04-02 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (603, 1, '2022-04-03 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (604, 1, '2022-04-04 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (605, 0, '2022-04-05 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (606, 1, '2022-04-06 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (607, 1, '2022-04-07 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (608, 0, '2022-04-08 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (609, 3, '2022-04-09 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (610, 4, '2022-04-10 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (611, 1, '2022-04-11 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (612, 5, '2022-04-12 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (613, 3, '2022-04-13 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (614, 6, '2022-04-14 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (615, 1, '2022-04-15 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (616, 0, '2022-04-16 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (617, 0, '2022-04-17 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (618, 1, '2022-04-18 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (619, 1, '2022-04-19 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (620, 3, '2022-04-20 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (621, 2, '2022-04-21 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (622, 0, '2022-04-22 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (623, 2, '2022-04-23 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (624, 6, '2022-04-24 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (625, 3, '2022-04-25 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (626, 0, '2022-04-26 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (627, 0, '2022-04-27 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (628, 0, '2022-04-28 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (629, 0, '2022-04-29 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (630, 3, '2022-04-30 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (631, 1, '2022-05-01 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (632, 0, '2022-05-02 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (633, 0, '2022-05-03 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (634, 0, '2022-05-04 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (635, 1, '2022-05-05 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (636, 5, '2022-05-06 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (637, 3, '2022-05-07 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (638, 0, '2022-05-08 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (639, 0, '2022-05-09 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (640, 2, '2022-05-10 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (641, 5, '2022-05-11 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (642, 1, '2022-05-12 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (643, 3, '2022-05-13 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (644, 1, '2022-05-14 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (645, 5, '2022-05-15 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (646, 1, '2022-05-16 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (647, 0, '2022-05-17 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (648, 0, '2022-05-18 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (649, 0, '2022-05-19 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (650, 1, '2022-05-20 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (651, 0, '2022-05-21 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (652, 0, '2022-05-22 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (653, 0, '2022-05-23 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (654, 0, '2022-05-24 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (655, 0, '2022-05-25 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (656, 0, '2022-05-26 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (657, 0, '2022-05-27 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (658, 0, '2022-05-28 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (659, 0, '2022-05-29 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (660, 0, '2022-05-30 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (661, 1, '2022-05-31 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (662, 2, '2022-06-01 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (663, 2, '2022-06-02 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (664, 1, '2022-06-03 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (665, 0, '2022-06-04 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (666, 0, '2022-06-05 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (667, 0, '2022-06-06 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (668, 1, '2022-06-07 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (669, 0, '2022-06-08 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (670, 0, '2022-06-09 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (671, 0, '2022-06-10 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (672, 0, '2022-06-11 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (673, 0, '2022-06-12 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (674, 0, '2022-06-13 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (675, 1, '2022-06-14 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (676, 0, '2022-06-15 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (677, 0, '2022-06-16 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (678, 0, '2022-06-17 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (679, 0, '2022-06-18 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (680, 1, '2022-06-19 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (681, 0, '2022-06-20 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (682, 0, '2022-06-21 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (683, 0, '2022-06-22 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (684, 0, '2022-06-23 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (685, 0, '2022-06-24 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (686, 1, '2022-06-25 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (687, 1, '2022-06-26 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (688, 1, '2022-06-27 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (689, 1, '2022-06-28 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (690, 1, '2022-06-29 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (691, 0, '2022-06-30 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (692, 0, '2022-07-01 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (693, 0, '2022-07-02 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (694, 0, '2022-07-03 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (695, 0, '2022-07-04 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (696, 0, '2022-07-05 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (697, 0, '2022-07-06 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (698, 0, '2022-07-07 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (699, 0, '2022-07-08 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (700, 1, '2023-06-15 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (701, 1, '2023-06-17 00:00:00', '2024-01-15 19:43:08');
INSERT INTO `unique_view` VALUES (702, 1, '2023-07-04 00:00:00', '2024-01-15 19:43:08');
COMMIT;

-- ----------------------------
-- Table structure for upload_record
-- ----------------------------
DROP TABLE IF EXISTS `upload_record`;
CREATE TABLE `upload_record` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户id',
  `label` varchar(128) NOT NULL DEFAULT '' COMMENT '标签',
  `file_name` varchar(64) NOT NULL DEFAULT '' COMMENT '文件名称',
  `file_size` int NOT NULL DEFAULT '0' COMMENT '文件大小',
  `file_md5` varchar(128) NOT NULL DEFAULT '' COMMENT '文件md5值',
  `file_url` varchar(256) NOT NULL DEFAULT '' COMMENT '上传路径',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_uid` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='上传记录';

-- ----------------------------
-- Records of upload_record
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for user_account
-- ----------------------------
DROP TABLE IF EXISTS `user_account`;
CREATE TABLE `user_account` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(128) NOT NULL DEFAULT '' COMMENT '密码',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态: -1删除 0正常 1禁用',
  `register_type` varchar(64) NOT NULL DEFAULT '' COMMENT '注册方式',
  `ip_address` varchar(255) NOT NULL DEFAULT '' COMMENT '注册ip',
  `ip_source` varchar(255) NOT NULL DEFAULT '' COMMENT '注册ip 源',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_username` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户登录信息';

-- ----------------------------
-- Records of user_account
-- ----------------------------
BEGIN;
INSERT INTO `user_account` VALUES (1, 've', '$2a$10$SRdVcOdTaly6wYDzZR4BIODH2QqEUTtW05VRrviHtG2jHriibW2NO', 0, '0', '127.0.0.1', '广西壮族自治区河池市 电信', '2023-05-12 15:16:48', '2023-07-07 18:39:40');
INSERT INTO `user_account` VALUES (2, 'admin', '$2a$10$Q.D5YRdVQ7ORZ/ui36PT5.NGd/eMzq/I3Malhyw76z3eB2VxD.2Rq', 0, '0', '127.0.0.1', '广西壮族自治区来宾市 电信', '2023-05-15 19:02:54', '2023-07-07 18:39:16');
INSERT INTO `user_account` VALUES (3, 'admin@qq.com', '$2a$10$ZINovpDg.FxFQRj6nhKDLOH55k19RDViybnVVn5EGuKQAcqChRs1e', 0, '0', '127.0.0.1', '广西壮族自治区梧州市 移动', '2023-05-17 15:34:25', '2023-08-31 20:36:04');
INSERT INTO `user_account` VALUES (4, 've77@qq.com', '$2a$10$123eQHGD/wa4neT1xxADo.9F.sCF4jJ2H33Rq7njMGhNT8tjxLYda', 0, '0', '127.0.0.1', '广西壮族自治区来宾市 电信', '2023-06-28 20:03:26', '2023-07-31 17:33:07');
INSERT INTO `user_account` VALUES (5, 've777@qq.com', '$2a$10$vvre.NmECAcm7IPrnD1qpOHhvI1AhsDfVgOR34l31zE0lFeaZcY1u', 0, '0', '127.0.0.1', '湖北省武汉市 广电网', '2023-06-28 20:06:40', '2023-07-31 17:33:09');
INSERT INTO `user_account` VALUES (6, '791422171@qq.com', '$2a$10$6fJR32Zj2KA0MlVJj7L6M.G0X4I/ydqTS.QTsk.J2SkL0bxNHBfLK', 0, 'email', '127.0.0.1', '广西壮族自治区梧州市 移动', '2023-06-29 18:35:40', '2024-01-14 01:20:43');
INSERT INTO `user_account` VALUES (11, 'ou_ef9abf85b8905510f005587d9717a596', '6', 0, 'feishu', '127.0.0.1', '广西壮族自治区梧州市 移动', '2023-07-04 11:50:22', '2024-01-21 20:46:02');
COMMIT;

-- ----------------------------
-- Table structure for user_information
-- ----------------------------
DROP TABLE IF EXISTS `user_information`;
CREATE TABLE `user_information` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户id',
  `email` varchar(128) NOT NULL DEFAULT '' COMMENT '用户邮箱',
  `nickname` varchar(128) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `avatar` varchar(1024) NOT NULL DEFAULT '' COMMENT '用户头像',
  `phone` varchar(32) NOT NULL DEFAULT '' COMMENT '用户手机号',
  `intro` varchar(255) NOT NULL DEFAULT '' COMMENT '个人简介',
  `website` varchar(255) NOT NULL DEFAULT '' COMMENT '个人网站',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_uuid` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户信息';

-- ----------------------------
-- Records of user_information
-- ----------------------------
BEGIN;
INSERT INTO `user_information` VALUES (1, 1, '', 've', 'https://mms1.baidu.com/it/u=2815887849,1501151317&fm=253&app=138&f=JPEG', '', '测试', '', '2023-05-12 15:16:48', '2023-07-02 19:18:36');
INSERT INTO `user_information` VALUES (2, 2, '', 'admin', 'https://mms1.baidu.com/it/u=2815887849,1501151317&fm=253&app=138&f=JPEG', '', '测试', '', '2023-05-15 19:02:54', '2023-07-02 19:18:40');
INSERT INTO `user_information` VALUES (3, 3, '', 'admin@qq.com', 'https://static.veweiyi.cn/blog/3/upload/article/cover/avatar_20231205202844.png', '', '测试', 'admin@qq.com', '2023-05-17 15:34:25', '2024-01-23 17:30:03');
INSERT INTO `user_information` VALUES (4, 4, '', 've77@qq.com', 'https://mms1.baidu.com/it/u=2815887849,1501151317&fm=253&app=138&f=JPEG', '', '测试', '', '2023-06-28 20:03:26', '2023-07-02 19:18:47');
INSERT INTO `user_information` VALUES (5, 5, '', 've777@qq.com', 'https://mms1.baidu.com/it/u=2815887849,1501151317&fm=253&app=138&f=JPEG', '', '测试', '', '2023-06-28 20:06:40', '2023-07-02 19:18:50');
INSERT INTO `user_information` VALUES (6, 6, '', '791422171@qq.com', 'https://mms1.baidu.com/it/u=2815887849,1501151317&fm=253&app=138&f=JPEG', '', '测试', '', '2023-06-29 18:35:40', '2023-07-02 19:18:52');
INSERT INTO `user_information` VALUES (11, 11, '', '游客d21e88ea', 'https://mms1.baidu.com/it/u=2815887849,1501151317&fm=253&app=138&f=JPEG', '', '这个人很神秘，什么都没有写！', '', '2023-07-04 11:50:22', '2023-10-08 18:55:33');
COMMIT;

-- ----------------------------
-- Table structure for user_login_history
-- ----------------------------
DROP TABLE IF EXISTS `user_login_history`;
CREATE TABLE `user_login_history` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户id',
  `login_type` varchar(64) NOT NULL DEFAULT '' COMMENT '登录类型',
  `agent` varchar(255) NOT NULL DEFAULT '' COMMENT '代理',
  `ip_address` varchar(255) NOT NULL DEFAULT '' COMMENT 'ip host',
  `ip_source` varchar(255) NOT NULL DEFAULT '' COMMENT 'ip 源',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `uk_uuid` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户登录历史';

-- ----------------------------
-- Records of user_login_history
-- ----------------------------
BEGIN;
INSERT INTO `user_login_history` VALUES (1, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-18 20:43:46', '2024-01-18 20:43:46');
INSERT INTO `user_login_history` VALUES (2, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-18 20:44:30', '2024-01-18 20:44:29');
INSERT INTO `user_login_history` VALUES (3, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-18 20:51:52', '2024-01-18 20:51:51');
INSERT INTO `user_login_history` VALUES (4, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-19 10:16:10', '2024-01-19 10:16:10');
INSERT INTO `user_login_history` VALUES (5, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-19 10:17:16', '2024-01-19 10:17:16');
INSERT INTO `user_login_history` VALUES (6, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-19 10:27:05', '2024-01-19 10:27:04');
INSERT INTO `user_login_history` VALUES (7, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-19 11:33:58', '2024-01-19 11:33:58');
INSERT INTO `user_login_history` VALUES (8, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-19 15:38:21', '2024-01-19 15:38:21');
INSERT INTO `user_login_history` VALUES (9, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-19 17:28:25', '2024-01-19 17:28:25');
INSERT INTO `user_login_history` VALUES (10, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-19 17:29:13', '2024-01-19 17:29:12');
INSERT INTO `user_login_history` VALUES (11, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-19 17:29:40', '2024-01-19 17:29:40');
INSERT INTO `user_login_history` VALUES (12, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-19 17:35:00', '2024-01-19 17:35:00');
INSERT INTO `user_login_history` VALUES (13, 3, 'email', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36', '127.0.0.1', '本机地址', '2024-01-19 17:47:51', '2024-01-19 17:47:51');
INSERT INTO `user_login_history` VALUES (14, 3, 'email', 'PostmanRuntime/7.36.1', '127.0.0.1', '本机地址', '2024-02-23 14:51:09', '2024-02-23 14:51:08');
COMMIT;

-- ----------------------------
-- Table structure for user_oauth
-- ----------------------------
DROP TABLE IF EXISTS `user_oauth`;
CREATE TABLE `user_oauth` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户id',
  `open_id` varchar(128) NOT NULL DEFAULT '' COMMENT '开发平台id，标识唯一用户',
  `platform` varchar(64) NOT NULL DEFAULT '' COMMENT '平台:手机号、邮箱、微信、飞书',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_oid_plat` (`open_id`,`platform`) USING BTREE,
  KEY `idx_uuid` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='第三方登录信息';

-- ----------------------------
-- Records of user_oauth
-- ----------------------------
BEGIN;
INSERT INTO `user_oauth` VALUES (6, 11, 'ou_ef9abf85b8905510f005587d9717a596', 'feishu', '2023-07-04 11:50:22', '2023-07-04 11:50:55');
COMMIT;

-- ----------------------------
-- Table structure for user_role
-- ----------------------------
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户id',
  `role_id` int NOT NULL DEFAULT '0' COMMENT '角色id',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户-角色关联';

-- ----------------------------
-- Records of user_role
-- ----------------------------
BEGIN;
INSERT INTO `user_role` VALUES (1, 3, 1);
INSERT INTO `user_role` VALUES (2, 3, 2);
INSERT INTO `user_role` VALUES (3, 3, 3);
INSERT INTO `user_role` VALUES (35, 11, 2);
INSERT INTO `user_role` VALUES (36, 6, 1);
INSERT INTO `user_role` VALUES (40, 5, 3);
INSERT INTO `user_role` VALUES (41, 5, 56);
INSERT INTO `user_role` VALUES (42, 4, 2);
INSERT INTO `user_role` VALUES (43, 4, 3);
INSERT INTO `user_role` VALUES (44, 2, 1);
INSERT INTO `user_role` VALUES (45, 1, 1);
COMMIT;

-- ----------------------------
-- Table structure for website_config
-- ----------------------------
DROP TABLE IF EXISTS `website_config`;
CREATE TABLE `website_config` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `key` varchar(32) NOT NULL DEFAULT '' COMMENT '关键词',
  `config` varchar(2048) NOT NULL DEFAULT '' COMMENT '配置信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='网站配置表';

-- ----------------------------
-- Records of website_config
-- ----------------------------
BEGIN;
INSERT INTO `website_config` VALUES (1, 'website_config', '{\"alipay_qr_code\":\"https://veport.oss-cn-beijing.aliyuncs.com/config/17f234dc487c1bb5bbb732869be0eb53.jpg\",\"is_chat_room\":1,\"is_reward\":1,\"tourist_avatar\":\"https://veport.oss-cn-beijing.aliyuncs.com/config/5bfb96809bee5ba80a36811f0bf1d1ea.gif\",\"website_create_time\":\"2022-01-19\",\"gitee\":\"https://gitee.com/wy791422171\",\"is_music_player\":0,\"website_author\":\"与梦\",\"website_avatar\":\"https://veport.oss-cn-beijing.aliyuncs.com/config/84aa08357bf6e74fc1d4f33552475f91.gif\",\"website_intro\":\"分享美好生活。\",\"wei_xin_qr_code\":\"https://veport.oss-cn-beijing.aliyuncs.com/config/6bed8a1130b170546341ece729e8819f.jpg\",\"github\":\"https://github.com/7914-ve\",\"is_email_notice\":1,\"is_message_review\":0,\"social_url_list\":[\"qq\",\"github\",\"gitee\"],\"user_avatar\":\"https://veport.oss-cn-beijing.aliyuncs.com/config/041a0d1c7fdfb5a610c307e7e44d4f39.jpg\",\"website_notice\":\"用户需要查看、发表文章、修改其他信息请登录后台管理系统。网站后台管理系统->https://veweiyi.cn/admin。     \\n网站搭建问题请联系站长QQ791422171。\",\"is_comment_review\":0,\"qq\":\"791422171\",\"social_login_list\":[\"weibo\",\"qq\",\"feishu\"],\"website_name\":\"静闻弦语\",\"website_record_no\":\"桂ICP备2023013735号-1\",\"websocket_url\":\"wss://ve77.cn:8088/api/websocket\"}', '2021-08-09 19:37:30', '2023-12-20 12:25:10');
INSERT INTO `website_config` VALUES (2, 'about', 'hhh**hhh**', '2023-07-16 02:47:23', '2023-07-16 02:53:02');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
