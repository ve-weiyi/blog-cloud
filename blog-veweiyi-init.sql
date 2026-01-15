/*
 Navicat Premium Dump SQL

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80034 (8.0.34)
 Source Host           : localhost:3306
 Source Schema         : blog-veweiyi

 Target Server Type    : MySQL
 Target Server Version : 80034 (8.0.34)
 File Encoding         : 65001

 Date: 18/01/2026 00:49:06
*/

SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for t_album
-- ----------------------------
DROP TABLE IF EXISTS `t_album`;
CREATE TABLE `t_album`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `album_name`  varchar(32)  NOT NULL DEFAULT '' COMMENT '相册名',
    `album_desc`  varchar(64)  NOT NULL DEFAULT '' COMMENT '相册描述',
    `album_cover` varchar(255) NOT NULL DEFAULT '' COMMENT '相册封面',
    `is_delete`   tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
    `status`      tinyint      NOT NULL DEFAULT '1' COMMENT '状态值 1公开 2私密',
    `created_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='相册';

-- ----------------------------
-- Table structure for t_api
-- ----------------------------
DROP TABLE IF EXISTS `t_api`;
CREATE TABLE `t_api`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `parent_id`  int          NOT NULL DEFAULT '0' COMMENT '分组id',
    `name`       varchar(128) NOT NULL DEFAULT '' COMMENT 'api名称',
    `path`       varchar(128) NOT NULL DEFAULT '' COMMENT 'api路径',
    `method`     varchar(16)  NOT NULL DEFAULT '' COMMENT 'api请求方法',
    `traceable`  tinyint      NOT NULL DEFAULT '0' COMMENT '是否追溯操作记录 0需要，1是',
    `status`     tinyint      NOT NULL DEFAULT '0' COMMENT '状态 0正常 1禁用',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `idx_path_method` (`path`,`method`,`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=985 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='接口';

-- ----------------------------
-- Table structure for t_article
-- ----------------------------
DROP TABLE IF EXISTS `t_article`;
CREATE TABLE `t_article`
(
    `id`              int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `user_id`         varchar(64)   NOT NULL DEFAULT '' COMMENT '作者',
    `category_id`     int           NOT NULL DEFAULT '0' COMMENT '文章分类',
    `article_cover`   varchar(1024) NOT NULL DEFAULT '' COMMENT '文章缩略图',
    `article_title`   varchar(64)   NOT NULL DEFAULT '' COMMENT '标题',
    `article_content` longtext      NOT NULL COMMENT '内容',
    `article_type`    tinyint       NOT NULL DEFAULT '0' COMMENT '文章类型 1原创 2转载 3翻译',
    `original_url`    varchar(255)  NOT NULL DEFAULT '' COMMENT '原文链接',
    `is_top`          tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否置顶 0否 1是',
    `is_delete`       tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除  0否 1是',
    `status`          tinyint       NOT NULL DEFAULT '1' COMMENT '状态值 1公开 2私密 3草稿 4评论可见',
    `like_count`      int           NOT NULL DEFAULT '0' COMMENT '点赞数',
    `view_count`      int           NOT NULL DEFAULT '0' COMMENT '查看数',
    `created_at`      datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
    `updated_at`      datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='文章';

-- ----------------------------
-- Table structure for t_article_tag
-- ----------------------------
DROP TABLE IF EXISTS `t_article_tag`;
CREATE TABLE `t_article_tag`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `article_id` int NOT NULL DEFAULT '0' COMMENT '文章id',
    `tag_id`     int NOT NULL DEFAULT '0' COMMENT '标签id',
    PRIMARY KEY (`id`) USING BTREE,
    KEY          `fk_article_tag_1` (`article_id`) USING BTREE,
    KEY          `fk_article_tag_2` (`tag_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='文章-标签关联';

-- ----------------------------
-- Table structure for t_category
-- ----------------------------
DROP TABLE IF EXISTS `t_category`;
CREATE TABLE `t_category`
(
    `id`            int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `category_name` varchar(32) NOT NULL DEFAULT '' COMMENT '分类名',
    `created_at`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_name` (`category_name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='文章分类';

-- ----------------------------
-- Table structure for t_chat
-- ----------------------------
DROP TABLE IF EXISTS `t_chat`;
CREATE TABLE `t_chat`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id`     varchar(64)   NOT NULL DEFAULT '' COMMENT '用户id',
    `terminal_id` varchar(64)   NOT NULL DEFAULT '' COMMENT '设备id',
    `nickname`    varchar(64)   NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar`      varchar(255)  NOT NULL DEFAULT '' COMMENT '头像',
    `ip_address`  varchar(64)   NOT NULL DEFAULT '' COMMENT '用户ip 127.0.0.1',
    `ip_source`   varchar(128)  NOT NULL DEFAULT '' COMMENT '用户地址 广东省深圳市',
    `type`        varchar(64)   NOT NULL DEFAULT '' COMMENT '类型:chatgpt chatroom',
    `content`     varchar(4096) NOT NULL DEFAULT '' COMMENT '聊天内容',
    `status`      int           NOT NULL DEFAULT '0' COMMENT '状态:0正常 1编辑 2撤回 3删除',
    `created_at`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY           `idx_uid` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='聊天消息';

-- ----------------------------
-- Table structure for t_comment
-- ----------------------------
DROP TABLE IF EXISTS `t_comment`;
CREATE TABLE `t_comment`
(
    `id`              int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id`         varchar(64)  NOT NULL DEFAULT '' COMMENT '评论用户id',
    `terminal_id`     varchar(64)  NOT NULL DEFAULT '' COMMENT '设备id',
    `topic_id`        int          NOT NULL DEFAULT '0' COMMENT '主题id',
    `parent_id`       int          NOT NULL DEFAULT '0' COMMENT '父评论id',
    `reply_id`        int          NOT NULL DEFAULT '0' COMMENT '回复评论id',
    `reply_user_id`   varchar(255) NOT NULL DEFAULT '' COMMENT '评论回复用户id',
    `comment_content` text         NOT NULL COMMENT '评论内容',
    `type`            tinyint      NOT NULL DEFAULT '0' COMMENT '评论类型 1.文章 2.友链 3.说说',
    `status`          tinyint      NOT NULL DEFAULT '0' COMMENT '状态 0.正常 1.已编辑 2.已删除',
    `like_count`      int          NOT NULL DEFAULT '0' COMMENT '评论点赞数量',
    `created_at`      datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`      datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY               `fk_comment_user` (`user_id`) USING BTREE,
    KEY               `fk_comment_parent` (`parent_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='评论';

-- ----------------------------
-- Table structure for t_file_log
-- ----------------------------
DROP TABLE IF EXISTS `t_file_log`;
CREATE TABLE `t_file_log`
(
    `id`          int          NOT NULL AUTO_INCREMENT COMMENT 'id',
    `user_id`     varchar(64)  NOT NULL DEFAULT '' COMMENT '用户id',
    `terminal_id` varchar(64)  NOT NULL DEFAULT '' COMMENT '设备id',
    `file_path`   varchar(128) NOT NULL DEFAULT '' COMMENT '文件路径',
    `file_name`   varchar(128) NOT NULL DEFAULT '' COMMENT '文件名称',
    `file_type`   varchar(128) NOT NULL DEFAULT '' COMMENT '文件类型',
    `file_size`   int          NOT NULL DEFAULT '0' COMMENT '文件大小',
    `file_md5`    varchar(128) NOT NULL DEFAULT '' COMMENT '文件md5值',
    `file_url`    varchar(256) NOT NULL DEFAULT '' COMMENT '上传路径',
    `created_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY           `idx_uid` (`user_id`) USING BTREE,
    KEY           `idx_path` (`file_path`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='文件上传记录';

-- ----------------------------
-- Table structure for t_friend
-- ----------------------------
DROP TABLE IF EXISTS `t_friend`;
CREATE TABLE `t_friend`
(
    `id`           int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `link_name`    varchar(32)  NOT NULL DEFAULT '' COMMENT '链接名',
    `link_avatar`  varchar(255) NOT NULL DEFAULT '' COMMENT '链接头像',
    `link_address` varchar(64)  NOT NULL DEFAULT '' COMMENT '链接地址',
    `link_intro`   varchar(100) NOT NULL DEFAULT '' COMMENT '链接介绍',
    `created_at`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_name` (`link_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='友链';

-- ----------------------------
-- Table structure for t_login_log
-- ----------------------------
DROP TABLE IF EXISTS `t_login_log`;
CREATE TABLE `t_login_log`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `user_id`     varchar(64) NOT NULL DEFAULT '' COMMENT '用户id',
    `terminal_id` varchar(64) NOT NULL DEFAULT '' COMMENT '设备id',
    `login_type`  varchar(64) NOT NULL DEFAULT '' COMMENT '登录类型',
    `app_name`    varchar(64) NOT NULL DEFAULT '' COMMENT 'app名称',
    `login_at`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    `logout_at`   datetime             DEFAULT NULL COMMENT '登出时间',
    `created_at`  datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY           `idx_uid` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=207 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户登录历史';

-- ----------------------------
-- Table structure for t_menu
-- ----------------------------
DROP TABLE IF EXISTS `t_menu`;
CREATE TABLE `t_menu`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `parent_id`   int           NOT NULL DEFAULT '0' COMMENT '父id',
    `path`        varchar(64)   NOT NULL DEFAULT '' COMMENT '路由路径',
    `name`        varchar(64)   NOT NULL DEFAULT '' COMMENT '路由名称',
    `component`   varchar(256)  NOT NULL DEFAULT '' COMMENT '路由组件',
    `redirect`    varchar(256)  NOT NULL DEFAULT '' COMMENT '路由重定向',
    `type`        varchar(64)   NOT NULL DEFAULT '0' COMMENT '菜单类型',
    `title`       varchar(64)   NOT NULL DEFAULT '' COMMENT '菜单标题',
    `icon`        varchar(64)   NOT NULL DEFAULT '' COMMENT '菜单图标',
    `rank`        int           NOT NULL DEFAULT '0' COMMENT '排序',
    `perm`        varchar(64)   NOT NULL DEFAULT '' COMMENT '权限标识',
    `params`      varchar(256)  NOT NULL DEFAULT '' COMMENT '路由参数',
    `keep_alive`  tinyint       NOT NULL DEFAULT '0' COMMENT '是否缓存',
    `always_show` tinyint       NOT NULL DEFAULT '0' COMMENT '是否一直显示菜单',
    `is_hidden`   tinyint       NOT NULL DEFAULT '0' COMMENT '是否隐藏',
    `status`      tinyint       NOT NULL DEFAULT '0' COMMENT '是否禁用',
    `extra`       varchar(1024) NOT NULL DEFAULT '' COMMENT '菜单元数据',
    `created_at`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_path` (`path`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=813 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='菜单';

-- ----------------------------
-- Table structure for t_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `t_operation_log`;
CREATE TABLE `t_operation_log`
(
    `id`              int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `user_id`         varchar(64)   NOT NULL DEFAULT '' COMMENT '用户id',
    `terminal_id`     varchar(64)   NOT NULL DEFAULT '' COMMENT '设备id',
    `opt_module`      varchar(32)   NOT NULL DEFAULT '' COMMENT '操作模块',
    `opt_desc`        varchar(255)  NOT NULL DEFAULT '' COMMENT '操作描述',
    `request_uri`     varchar(255)  NOT NULL DEFAULT '' COMMENT '请求地址',
    `request_method`  varchar(32)   NOT NULL DEFAULT '' COMMENT '请求方式',
    `request_data`    varchar(4096) NOT NULL DEFAULT '' COMMENT '请求参数',
    `response_data`   varchar(4096) NOT NULL DEFAULT '' COMMENT '返回数据',
    `response_status` int           NOT NULL DEFAULT '0' COMMENT '响应状态码',
    `cost`            varchar(32)   NOT NULL DEFAULT '' COMMENT '耗时（ms）',
    `created_at`      datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`      datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=560 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='操作记录';

-- ----------------------------
-- Table structure for t_page
-- ----------------------------
DROP TABLE IF EXISTS `t_page`;
CREATE TABLE `t_page`
(
    `id`              int unsigned NOT NULL AUTO_INCREMENT COMMENT '页面id',
    `page_name`       varchar(32)   NOT NULL DEFAULT '' COMMENT '页面名',
    `page_label`      varchar(32)   NOT NULL DEFAULT '' COMMENT '页面标签',
    `page_cover`      varchar(255)  NOT NULL DEFAULT '' COMMENT '页面封面',
    `is_carousel`     tinyint       NOT NULL DEFAULT '0' COMMENT '是否轮播',
    `carousel_covers` varchar(1024) NOT NULL DEFAULT '' COMMENT '轮播图片列表',
    `created_at`      datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`      datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='页面';

-- ----------------------------
-- Table structure for t_photo
-- ----------------------------
DROP TABLE IF EXISTS `t_photo`;
CREATE TABLE `t_photo`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `album_id`   int          NOT NULL DEFAULT '0' COMMENT '相册id',
    `photo_name` varchar(32)  NOT NULL DEFAULT '' COMMENT '照片名',
    `photo_desc` varchar(64)  NOT NULL DEFAULT '' COMMENT '照片描述',
    `photo_src`  varchar(255) NOT NULL DEFAULT '' COMMENT '照片地址',
    `is_delete`  tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='照片';

-- ----------------------------
-- Table structure for t_remark
-- ----------------------------
DROP TABLE IF EXISTS `t_remark`;
CREATE TABLE `t_remark`
(
    `id`              int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `user_id`         varchar(255) NOT NULL DEFAULT '0' COMMENT '用户id',
    `terminal_id`     varchar(255) NOT NULL DEFAULT '' COMMENT '终端id',
    `message_content` varchar(255) NOT NULL DEFAULT '' COMMENT '留言内容',
    `status`          int          NOT NULL DEFAULT '0' COMMENT '状态:0正常 1编辑 2撤回 3删除',
    `created_at`      datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
    `updated_at`      datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='留言';

-- ----------------------------
-- Table structure for t_role
-- ----------------------------
DROP TABLE IF EXISTS `t_role`;
CREATE TABLE `t_role`
(
    `id`           int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `parent_id`    int         NOT NULL DEFAULT '0' COMMENT '父角色id',
    `role_key`     varchar(64) NOT NULL DEFAULT '' COMMENT '角色标识',
    `role_label`   varchar(64) NOT NULL DEFAULT '' COMMENT '角色标签',
    `role_comment` varchar(64) NOT NULL DEFAULT '' COMMENT '角色备注',
    `is_default`   tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认角色 0否 1是',
    `status`       tinyint     NOT NULL DEFAULT '0' COMMENT '状态  0正常 1禁用',
    `created_at`   datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`   datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色';

-- ----------------------------
-- Table structure for t_role_api
-- ----------------------------
DROP TABLE IF EXISTS `t_role_api`;
CREATE TABLE `t_role_api`
(
    `id`      int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `role_id` int NOT NULL DEFAULT '0' COMMENT '角色id',
    `api_id`  int NOT NULL DEFAULT '0' COMMENT '接口id',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色-api关联';

-- ----------------------------
-- Table structure for t_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `t_role_menu`;
CREATE TABLE `t_role_menu`
(
    `id`      int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `role_id` int NOT NULL DEFAULT '0' COMMENT '角色id',
    `menu_id` int NOT NULL DEFAULT '0' COMMENT '菜单id',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=189 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色-菜单关联';

-- ----------------------------
-- Table structure for t_tag
-- ----------------------------
DROP TABLE IF EXISTS `t_tag`;
CREATE TABLE `t_tag`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `tag_name`   varchar(32) NOT NULL DEFAULT '' COMMENT '标签名',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_name` (`tag_name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='标签';

-- ----------------------------
-- Table structure for t_talk
-- ----------------------------
DROP TABLE IF EXISTS `t_talk`;
CREATE TABLE `t_talk`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT COMMENT '说说id',
    `user_id`    varchar(64)   NOT NULL DEFAULT '' COMMENT '用户id',
    `content`    varchar(2048) NOT NULL DEFAULT '' COMMENT '说说内容',
    `images`     varchar(2048) NOT NULL DEFAULT '' COMMENT '图片',
    `is_top`     tinyint       NOT NULL DEFAULT '0' COMMENT '是否置顶',
    `status`     tinyint       NOT NULL DEFAULT '1' COMMENT '状态 1.公开 2.私密',
    `like_count` int           NOT NULL DEFAULT '0' COMMENT '点赞数',
    `created_at` datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='说说';

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user`
(
    `id`            int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `user_id`       varchar(64)   NOT NULL DEFAULT '' COMMENT '用户id',
    `username`      varchar(64)   NOT NULL DEFAULT '' COMMENT '用户名',
    `password`      varchar(128)  NOT NULL DEFAULT '' COMMENT '用户密码',
    `nickname`      varchar(64)   NOT NULL DEFAULT '' COMMENT '用户昵称',
    `avatar`        varchar(255)  NOT NULL DEFAULT '' COMMENT '用户头像',
    `email`         varchar(64)   NOT NULL DEFAULT '' COMMENT '邮箱',
    `phone`         varchar(64)   NOT NULL DEFAULT '' COMMENT '手机号',
    `info`          varchar(1024) NOT NULL DEFAULT '' COMMENT '用户信息',
    `status`        tinyint       NOT NULL DEFAULT '0' COMMENT '状态: -1删除 0正常 1禁用',
    `register_type` varchar(64)   NOT NULL DEFAULT '' COMMENT '注册方式',
    `ip_address`    varchar(255)  NOT NULL DEFAULT '' COMMENT '注册ip',
    `ip_source`     varchar(255)  NOT NULL DEFAULT '' COMMENT '注册ip 源',
    `created_at`    datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_uid` (`user_id`) USING BTREE,
    UNIQUE KEY `uk_username` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户登录信息';

-- ----------------------------
-- Table structure for t_user_oauth
-- ----------------------------
DROP TABLE IF EXISTS `t_user_oauth`;
CREATE TABLE `t_user_oauth`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `user_id`    varchar(64)  NOT NULL DEFAULT '' COMMENT '用户id',
    `platform`   varchar(64)  NOT NULL DEFAULT '' COMMENT '平台:手机号、邮箱、微信、飞书',
    `open_id`    varchar(128) NOT NULL DEFAULT '' COMMENT '第三方平台id，标识唯一用户',
    `nickname`   varchar(128) NOT NULL DEFAULT '' COMMENT '第三方平台昵称',
    `avatar`     varchar(256) NOT NULL DEFAULT '' COMMENT '第三方平台头像',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_uid_plat` (`user_id`,`platform`) USING BTREE,
    UNIQUE KEY `uk_oid_plat` (`open_id`,`platform`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='第三方登录信息';

-- ----------------------------
-- Table structure for t_user_role
-- ----------------------------
DROP TABLE IF EXISTS `t_user_role`;
CREATE TABLE `t_user_role`
(
    `id`      int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `user_id` varchar(64) NOT NULL DEFAULT '' COMMENT '用户id',
    `role_id` int         NOT NULL DEFAULT '0' COMMENT '角色id',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户-角色关联';

-- ----------------------------
-- Table structure for t_visit_daily_stats
-- ----------------------------
DROP TABLE IF EXISTS `t_visit_daily_stats`;
CREATE TABLE `t_visit_daily_stats`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `date`       varchar(10) NOT NULL DEFAULT '' COMMENT '日期',
    `view_count` int         NOT NULL DEFAULT '0' COMMENT '访问量',
    `visit_type` tinyint     NOT NULL DEFAULT '1' COMMENT '1 访客数 2 浏览数',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_date_type` (`date`,`visit_type`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=548 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='页面访问数量';

-- ----------------------------
-- Table structure for t_visit_log
-- ----------------------------
DROP TABLE IF EXISTS `t_visit_log`;
CREATE TABLE `t_visit_log`
(
    `id`          int         NOT NULL AUTO_INCREMENT COMMENT 'id',
    `user_id`     varchar(64) NOT NULL DEFAULT '' COMMENT '用户id',
    `terminal_id` varchar(64) NOT NULL DEFAULT '' COMMENT '设备id',
    `page_name`   varchar(64) NOT NULL DEFAULT '' COMMENT '访问页面',
    `created_at`  datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY           `idx_uid` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=868 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_visitor
-- ----------------------------
DROP TABLE IF EXISTS `t_visitor`;
CREATE TABLE `t_visitor`
(
    `id`          int          NOT NULL AUTO_INCREMENT COMMENT 'id',
    `terminal_id` varchar(64)  NOT NULL DEFAULT '' COMMENT '设备id',
    `os`          varchar(50)  NOT NULL DEFAULT '' COMMENT '操作系统',
    `browser`     varchar(50)  NOT NULL DEFAULT '' COMMENT '浏览器',
    `ip_address`  varchar(255) NOT NULL DEFAULT '' COMMENT '操作ip',
    `ip_source`   varchar(255) NOT NULL DEFAULT '' COMMENT '操作地址',
    `created_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_tid` (`terminal_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1468 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_website_config
-- ----------------------------
DROP TABLE IF EXISTS `t_website_config`;
CREATE TABLE `t_website_config`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `key`        varchar(32)   NOT NULL DEFAULT '' COMMENT '关键词',
    `config`     varchar(4096) NOT NULL DEFAULT '' COMMENT '配置信息',
    `created_at` datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='网站配置表';

SET
FOREIGN_KEY_CHECKS = 1;
