DROP TABLE IF EXISTS `t_device_model_info`;

CREATE TABLE `t_device_model_info`
(
    `id`           int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `brand`        varchar(32)                                                    NOT NULL DEFAULT '' COMMENT '品牌',
    `slug`         varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci   NOT NULL DEFAULT '' COMMENT 'slug',
    `device_model` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '' COMMENT '设备型号',
    `device_name`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '' COMMENT '设备名称',
    `device_type`  varchar(64)                                                    NOT NULL DEFAULT '' COMMENT '设备类型',
    `device_id`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '' COMMENT '设备id',
    `description`  varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '' COMMENT '设备描述',
    `image_url`    varchar(255)                                                   NOT NULL DEFAULT '' COMMENT '图片',
    `capacity`     varchar(255)                                                   NOT NULL DEFAULT '' COMMENT '电池容量',
    `data`         varchar(8196) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'json',
    `created_at`   datetime                                                       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`   datetime                                                       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_slug` (`slug`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='设备型号';


CREATE TABLE `t_device_brand`
(
    `id`                int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name`              varchar(32)  NOT NULL DEFAULT '' COMMENT '名称',
    `slug`              varchar(64)  NOT NULL DEFAULT '' COMMENT 'slug',
    `brand_id`          varchar(255) NOT NULL DEFAULT '' COMMENT '品牌id',
    `number_of_devices` int          NOT NULL DEFAULT '0' COMMENT '产品数量',
    `created_at`        datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`        datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_slug` (`slug`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='设备品牌';
