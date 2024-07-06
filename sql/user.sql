-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: db:3306
-- Generation Time: Jul 06, 2024 at 01:18 AM
-- Server version: 8.0.33
-- PHP Version: 8.2.8

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `graphql_with_casbin_db`
--

-- --------------------------------------------------------

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `user_id` int UNSIGNED NOT NULL,
  `user_nm` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `pass_hash` varchar(61) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `token` varchar(132) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `is_admin` tinyint(1) NOT NULL DEFAULT '0',
  `status_cd` enum('active','non_active','nullified','blocked') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
  `created_dttm` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_user_id` int UNSIGNED NOT NULL DEFAULT '0',
  `updated_dttm` datetime DEFAULT NULL,
  `updated_user_id` int UNSIGNED NOT NULL DEFAULT '0',
  `nullified_dttm` datetime DEFAULT NULL,
  `nullified_user_id` int UNSIGNED NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `user`
--

INSERT INTO `user` (`user_id`, `user_nm`, `pass_hash`, `token`, `is_admin`, `status_cd`, `created_dttm`, `created_user_id`, `updated_dttm`, `updated_user_id`, `nullified_dttm`, `nullified_user_id`) VALUES
(1, 'superadmin', '$2a$09$x2Xa4OiY3FzIbOKNBAMhZOSXJsCKYtbkbphP3KiPJS0ZieidVnzbu', '', 1, 'active', '2024-06-28 10:46:53', 1, NULL, 0, NULL, 0),
(4, 'superadmin2', '$2a$09$mh3IQn3DjIyhoGr5Ix0wTeBeYQBbqIwlw/8Oh6EeUstkDHLOHMgbi', '', 1, 'active', '2024-07-05 00:56:28', 111, NULL, 0, NULL, 0);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`user_id`),
  ADD UNIQUE KEY `user_nm` (`user_nm`),
  ADD KEY `x_status` (`status_cd`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `user`
--
ALTER TABLE `user`
  MODIFY `user_id` int UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
