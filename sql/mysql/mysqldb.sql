/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `categories` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `id_s` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `category_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `category_desc` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `num_views` int(10) unsigned DEFAULT 0,
  `num_topics` int(10) unsigned DEFAULT 0,
  `levelc` tinyint(4) DEFAULT NULL,
  `parent_id` int(10) unsigned DEFAULT NULL,
  `num_chd` smallint(6) DEFAULT NULL,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_categories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `category_chds` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `category_id` int(10) unsigned DEFAULT NULL,
  `category_chd_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category_chds_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mdrafts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `id_s` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mtext` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mattach1` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mattach2` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mattach3` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mattach4` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mattach5` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `category_id` int(10) unsigned DEFAULT NULL,
  `topic_id` int(10) unsigned DEFAULT NULL,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_mdrafts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `message_attachments` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `mattach` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `category_id` int(10) unsigned DEFAULT NULL,
  `topic_id` int(10) unsigned DEFAULT NULL,
  `message_id` int(10) unsigned DEFAULT NULL,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_message_attachments_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `message_texts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `mtext` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `category_id` int(10) unsigned DEFAULT NULL,
  `topic_id` int(10) unsigned DEFAULT NULL,
  `message_id` int(10) unsigned DEFAULT NULL,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_message_texts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `messages` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `id_s` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `num_likes` int(10) unsigned DEFAULT NULL,
  `num_upvotes` int(10) unsigned DEFAULT NULL,
  `num_downvotes` int(10) unsigned DEFAULT NULL,
  `category_id` int(10) unsigned DEFAULT NULL,
  `topic_id` int(10) unsigned DEFAULT NULL,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_messages_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `topics` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `id_s` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `topic_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `topic_desc` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `num_tags` int(10) unsigned DEFAULT NULL,
  `tag1` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tag2` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tag3` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tag4` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tag5` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tag6` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tag7` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tag8` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tag9` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tag10` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `num_views` int(10) unsigned DEFAULT 0,
  `num_messages` int(10) unsigned DEFAULT 0,
  `category_id` int(10) unsigned DEFAULT NULL,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_topics_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `topics_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `id_s` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `topic_id` int(10) unsigned DEFAULT NULL,
  `num_messages` int(10) unsigned DEFAULT 0,
  `num_views` int(10) unsigned DEFAULT 0,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_topics_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ubadges` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `id_s` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ubadge_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ubadge_desc` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_ubadges_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ubadges_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `id_s` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ubadge_id` int(10) unsigned DEFAULT NULL,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_ubadges_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ugroup_chds` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `ugroup_id` int(10) unsigned DEFAULT NULL,
  `ugroup_chd_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_ugroup_chds_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ugroups` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `id_s` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ugroup_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ugroup_desc` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `levelc` tinyint(4) DEFAULT NULL,
  `parent_id` int(10) unsigned DEFAULT NULL,
  `num_chd` smallint(6) DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_ugroups_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ugroups_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `id_s` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ugroup_id` int(10) unsigned DEFAULT NULL,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_ugroups_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_bookmarks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `topic_id` int(10) unsigned DEFAULT NULL,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_bookmarks_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_likes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `topic_id` int(10) unsigned DEFAULT NULL,
  `message_id` int(10) unsigned DEFAULT 0,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_likes_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_replies` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `topic_id` int(10) unsigned DEFAULT NULL,
  `message_id` int(10) unsigned DEFAULT 0,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_replies_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_topics` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `topic_id` int(10) unsigned DEFAULT NULL,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_topics_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_votes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `topic_id` int(10) unsigned DEFAULT NULL,
  `message_id` int(10) unsigned DEFAULT 0,
  `vote` int(10) unsigned DEFAULT 0,
  `ugroup_id` int(10) unsigned DEFAULT 0,
  `user_id` int(10) unsigned DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_votes_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `id_s` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `auth_token` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `username` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `first_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `last_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `role` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password` varbinary(255) DEFAULT NULL,
  `active` tinyint(1) DEFAULT 0,
  `email_confirmation_token` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email_selector` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email_verifier` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email_token_sent_at` timestamp NULL DEFAULT NULL,
  `email_token_expiry` timestamp NULL DEFAULT NULL,
  `email_confirmed_at` timestamp NULL DEFAULT NULL,
  `new_email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `new_email_reset_token` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `new_email_selector` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `new_email_verifier` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `new_email_token_sent_at` timestamp NULL DEFAULT NULL,
  `new_email_token_expiry` timestamp NULL DEFAULT NULL,
  `new_email_confirmed_at` timestamp NULL DEFAULT NULL,
  `password_reset_token` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password_selector` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password_verifier` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password_token_sent_at` timestamp NULL DEFAULT NULL,
  `password_token_expiry` timestamp NULL DEFAULT NULL,
  `password_confirmed_at` timestamp NULL DEFAULT NULL,
  `timezone` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT 'Asia/Kolkata',
  `sign_in_count` int(10) unsigned DEFAULT NULL,
  `current_sign_in_at` timestamp NULL DEFAULT NULL,
  `last_sign_in_at` timestamp NULL DEFAULT NULL,
  `statusc` tinyint(3) unsigned DEFAULT NULL,
  `created_day` smallint(5) unsigned DEFAULT NULL,
  `created_week` tinyint(3) unsigned DEFAULT NULL,
  `created_month` tinyint(3) unsigned DEFAULT NULL,
  `created_year` smallint(5) unsigned DEFAULT NULL,
  `updated_day` smallint(5) unsigned DEFAULT NULL,
  `updated_week` tinyint(3) unsigned DEFAULT NULL,
  `updated_month` tinyint(3) unsigned DEFAULT NULL,
  `updated_year` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
