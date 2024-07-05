CREATE TABLE `user` (
   `user_id` int unsigned NOT NULL AUTO_INCREMENT,
   `user_nm` varchar(30) NOT NULL DEFAULT '',
   `pass_hash` varchar(61) NOT NULL DEFAULT '',
   `token` varchar(132) NOT NULL DEFAULT '',
   `is_admin` tinyint(1) NOT NULL DEFAULT 0,
   `status_cd` enum('active','non_active','nullified','blocked') NOT NULL DEFAULT 'active',
   `created_dttm` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `created_user_id` int unsigned NOT NULL DEFAULT 0,
   `updated_dttm` datetime DEFAULT NULL,
   `updated_user_id` int unsigned NOT NULL DEFAULT 0,
   `nullified_dttm` datetime DEFAULT NULL,
   `nullified_user_id` int unsigned NOT NULL DEFAULT 0,
   PRIMARY KEY (`user_id`),
   UNIQUE KEY `user_nm` (`user_nm`),
   KEY `x_status` (`status_cd`)
 ) ENGINE=InnoDB