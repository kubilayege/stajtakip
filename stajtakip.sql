-- --------------------------------------------------------
-- Host:                         localhost
-- Server version:               5.7.18 - MySQL Community Server (GPL)
-- Server OS:                    Win64
-- HeidiSQL Version:             9.4.0.5125
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- Dumping database structure for stajtest
CREATE DATABASE IF NOT EXISTS `stajtest` /*!40100 DEFAULT CHARACTER SET latin1 */;
USE `stajtest`;

-- Dumping structure for table stajtest.denkstaj
CREATE TABLE IF NOT EXISTS `denkstaj` (
  `OgrenciNo` int(11) NOT NULL,
  `KurumAdi` varchar(50) NOT NULL,
  `OncekiOkul` varchar(50) NOT NULL,
  `KabulGun` tinyint(3) unsigned NOT NULL,
  `ToplamGun` tinyint(3) unsigned NOT NULL,
  PRIMARY KEY (`OgrenciNo`),
  CONSTRAINT `FK_DenkStaj_ogrenci` FOREIGN KEY (`OgrenciNo`) REFERENCES `ogrenci` (`No`)
) ENGINE=InnoDB DEFAULT CHARSET=latin5;

-- Data exporting was unselected.
-- Dumping structure for table stajtest.komisyon
CREATE TABLE IF NOT EXISTS `komisyon` (
  `AdSoyad` varchar(50) NOT NULL,
  `Dahil` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`AdSoyad`)
) ENGINE=InnoDB DEFAULT CHARSET=latin5;

-- Data exporting was unselected.
-- Dumping structure for table stajtest.konu
CREATE TABLE IF NOT EXISTS `konu` (
  `Baslik` varchar(50) NOT NULL,
  `Aktif` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`Baslik`)
) ENGINE=InnoDB DEFAULT CHARSET=latin5;

-- Data exporting was unselected.
-- Dumping structure for view stajtest.kurum
-- Creating temporary table to overcome VIEW dependency errors
CREATE TABLE `kurum` (
	`kurumadi` VARCHAR(50) NOT NULL COLLATE 'latin5_turkish_ci'
) ENGINE=MyISAM;

-- Dumping structure for table stajtest.mulakat
CREATE TABLE IF NOT EXISTS `mulakat` (
  `OgrenciNo` int(11) NOT NULL,
  `StajBaslangic` date NOT NULL,
  `TarihSaat` datetime DEFAULT NULL,
  `KomisyonUye1` varchar(50) DEFAULT NULL,
  `KomisyonUye2` varchar(50) DEFAULT NULL,
  `PuanDevam` tinyint(2) unsigned DEFAULT NULL,
  `PuanCaba` tinyint(2) unsigned DEFAULT NULL,
  `PuanVakit` tinyint(2) unsigned DEFAULT NULL,
  `PuanAmireDavranis` tinyint(2) unsigned DEFAULT NULL,
  `PuanIsArkadasaDavranis` tinyint(2) unsigned DEFAULT NULL,
  `PuanProje` tinyint(3) unsigned DEFAULT NULL,
  `PuanDuzen` tinyint(3) unsigned DEFAULT NULL,
  `PuanSunum` tinyint(3) unsigned DEFAULT NULL,
  `PuanIcerik` tinyint(3) unsigned DEFAULT NULL,
  `PuanMulakat` tinyint(3) unsigned DEFAULT NULL,
  PRIMARY KEY (`OgrenciNo`,`StajBaslangic`),
  KEY `FK_mulakat_komisyon` (`KomisyonUye1`),
  KEY `FK_mulakat_komisyon_2` (`KomisyonUye2`),
  CONSTRAINT `FK_mulakat_komisyon` FOREIGN KEY (`KomisyonUye1`) REFERENCES `komisyon` (`AdSoyad`),
  CONSTRAINT `FK_mulakat_komisyon_2` FOREIGN KEY (`KomisyonUye2`) REFERENCES `komisyon` (`AdSoyad`),
  CONSTRAINT `FK_mulakat_staj` FOREIGN KEY (`OgrenciNo`, `StajBaslangic`) REFERENCES `staj` (`OgrenciNo`, `Baslangic`)
) ENGINE=InnoDB DEFAULT CHARSET=latin5;

-- Data exporting was unselected.
-- Dumping structure for table stajtest.ogrenci
CREATE TABLE IF NOT EXISTS `ogrenci` (
  `No` int(11) NOT NULL,
  `Ad` varchar(50) NOT NULL,
  `Soyad` varchar(50) NOT NULL,
  `Ogretim` tinyint(1) unsigned NOT NULL,
  PRIMARY KEY (`No`)
) ENGINE=InnoDB DEFAULT CHARSET=latin5;

-- Data exporting was unselected.
-- Dumping structure for table stajtest.ogrenciek
CREATE TABLE IF NOT EXISTS `ogrenciek` (
  `OgrenciNo` int(11) NOT NULL,
  `Dosya` varchar(100) NOT NULL,
  PRIMARY KEY (`OgrenciNo`),
  CONSTRAINT `FK_OgrenciEk_ogrenci` FOREIGN KEY (`OgrenciNo`) REFERENCES `ogrenci` (`No`)
) ENGINE=InnoDB DEFAULT CHARSET=latin5;

-- Data exporting was unselected.
-- Dumping structure for table stajtest.staj
CREATE TABLE IF NOT EXISTS `staj` (
  `OgrenciNo` int(11) NOT NULL,
  `KurumAdi` varchar(50) NOT NULL,
  `Sehir` varchar(50) NOT NULL,
  `KonuBaslik` varchar(50) NOT NULL,
  `Baslangic` date NOT NULL,
  `Bitis` date NOT NULL,
  `Sinif` tinyint(3) unsigned NOT NULL,
  `ToplamGun` tinyint(3) unsigned NOT NULL,
  `KabulGun` tinyint(3) unsigned DEFAULT NULL,
  `Degerlendirildi` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`OgrenciNo`,`Baslangic`),
  KEY `FK_staj_konu` (`KonuBaslik`),
  CONSTRAINT `FK__ogrenci` FOREIGN KEY (`OgrenciNo`) REFERENCES `ogrenci` (`No`),
  CONSTRAINT `FK_staj_konu` FOREIGN KEY (`KonuBaslik`) REFERENCES `konu` (`Baslik`)
) ENGINE=InnoDB DEFAULT CHARSET=latin5;

-- Data exporting was unselected.
-- Dumping structure for view stajtest.kurum
-- Removing temporary table and create final VIEW structure
DROP TABLE IF EXISTS `kurum`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` VIEW `kurum` AS SELECT DISTINCT kurumadi FROM staj UNION SELECT DISTINCT kurumadi FROM denkstaj ;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
