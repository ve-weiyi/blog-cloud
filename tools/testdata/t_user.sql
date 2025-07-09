DROP TABLE IF EXISTS `api`;
CREATE TABLE `api` (
                       `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                       `parent_id` int NOT NULL DEFAULT '0' COMMENT '分组id',
                       `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'api名称',
                       `path` varchar(128) NOT NULL DEFAULT '' COMMENT 'api路径',
                       `method` varchar(16) NOT NULL DEFAULT '' COMMENT 'api请求方法',
                       `traceable` tinyint NOT NULL DEFAULT '0' COMMENT '是否追溯操作记录 0需要，1是',
                       `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1开，2关',
                       `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                       `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                       PRIMARY KEY (`id`) USING BTREE,
                       UNIQUE KEY `idx_path_method` (`path`,`method`,`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=166 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='接口';
