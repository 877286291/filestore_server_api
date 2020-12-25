CREATE TABLE `tbl_file` (
  `id` int NOT NULL AUTO_INCREMENT,
  `file_sha1` char(40) NOT NULL COMMENT '文件hash',
  `file_name` varchar(255) NOT NULL COMMENT '文件名',
  `file_size` bigint(20) unsigned zerofill NOT NULL COMMENT '文件大小',
  `file_addr` varchar(1024) NOT NULL COMMENT '文件存储位置',
  `create_at` datetime DEFAULT NULL COMMENT '创建日期',
  `update_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  `status` int NOT NULL DEFAULT '0' COMMENT '可用/禁用/已删除',
  `ext1` int DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `file_sha1` (`file_sha1`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;