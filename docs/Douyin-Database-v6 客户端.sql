CREATE DATABASE `douyin` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

use douyin;

CREATE TABLE `tb_user` (
  `id` varchar(64) NOT NULL,
  `name` varchar(30) NOT NULL COMMENT '用户名',
  `password` varchar(60) NOT NULL,
  `follow_count` int DEFAULT '0' COMMENT '关注数量',
  `follower_count` int DEFAULT '0' COMMENT '粉丝数量',
  `register_time` datetime DEFAULT NULL COMMENT '注册时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ID_UNIQUE` (`id`),
  UNIQUE KEY `User_name_UNIQUE` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `tb_video` (
  `id` varchar(64) NOT NULL,
  `user_id` varchar(64) NOT NULL COMMENT '视频作者id',
  `title` varchar(50) DEFAULT NULL COMMENT '视频标题',
  `play_url` varchar(255) NOT NULL COMMENT '视频播放地址',
  `cover_url` varchar(255) DEFAULT NULL COMMENT '视频封面地址',
  `favorite_count` int DEFAULT '0' COMMENT '点赞数',
  `comment_count` int DEFAULT '0' COMMENT '评论数',
  `status` tinyint DEFAULT NULL COMMENT '状态位',
  `publish_time` datetime DEFAULT NULL COMMENT '视频发布时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ID_UNIQUE` (`id`),
  KEY `User_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `tb_favorite` (
  `id` varchar(64) NOT NULL,
  `user_id` varchar(64) NOT NULL COMMENT '点赞用户id',
  `video_id` varchar(64) NOT NULL COMMENT '视频id',
  `create_time` datetime NOT NULL COMMENT '点赞时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ID_UNIQUE` (`id`) USING BTREE,
  UNIQUE KEY `User_Video` (`user_id`,`video_id`) USING BTREE,
  KEY `User_id` (`user_id`) /*!80000 INVISIBLE */,
  KEY `Video_id` (`video_id`) /*!80000 INVISIBLE */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `tb_comment` (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `user_id` varchar(64) NOT NULL COMMENT '评论用户id',
  `video_id` varchar(64) NOT NULL COMMENT '视频id',
  `content` text COMMENT '评论内容',
  `create_date` datetime NOT NULL COMMENT '评论时间（yy-dd）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ID_UNIQUE` (`id`) USING BTREE,
  KEY `User_id` (`user_id`) /*!80000 INVISIBLE */,
  KEY `Video_id` (`video_id`) /*!80000 INVISIBLE */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `tb_follow` (
  `id` varchar(64) NOT NULL,
  `user_id` varchar(64) NOT NULL COMMENT '用户id（被关注用户）',
  `follow_id` varchar(64) NOT NULL COMMENT '粉丝用户id',
  `create_time` datetime DEFAULT NULL COMMENT '关注时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ID_UNIQUE` (`id`),
  UNIQUE KEY `User_Follow_UNIQUE` (`user_id`,`follow_id`) /*!80000 INVISIBLE */,
  KEY `User_id` (`user_id`) /*!80000 INVISIBLE */,
  KEY `Follower_id` (`follow_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
