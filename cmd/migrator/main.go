package main

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/module/dsn"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func migrate() (err error) {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&domain.User{}, &domain.City{}, &domain.Ad{}, &domain.AdPosition{}, &domain.AdAvailableDate{}, &domain.Image{}, &domain.Request{}, &domain.Review{})
	if err != nil {
		return err
	}
	if err := seedCities(db); err != nil {
		return err
	}
	fmt.Println("Database migrated")
	return nil
}

func main() {
	err := migrate()
	if err != nil {
		log.Fatal(err)
	}
}

func seedCities(db *gorm.DB) error {
	cities := []domain.City{
		{Title: "Москва", EnTitle: "Moscow", Description: "Столица России", Image: ""},
		{Title: "Санкт-Петербург", EnTitle: "Saint-Petersburg", Description: "Культурная столица России", Image: ""},
		{Title: "Новосибирск", EnTitle: "Novosibirsk", Description: "Третий по численности город России", Image: ""},
		{Title: "Екатеринбург", EnTitle: "Yekaterinburg", Description: "Административный центр Урала", Image: ""},
		{Title: "Казань", EnTitle: "Kazan", Description: "Столица Республики Татарстан", Image: ""},
		{Title: "Нижний Новгород", EnTitle: "Nizhny Novgorod", Description: "Важный культурный и экономический центр", Image: ""},
		{Title: "Челябинск", EnTitle: "Chelyabinsk", Description: "Крупный промышленный город на Урале", Image: ""},
		{Title: "Самара", EnTitle: "Samara", Description: "Крупный город на Волге", Image: ""},
		{Title: "Омск", EnTitle: "Omsk", Description: "Один из крупнейших городов Сибири", Image: ""},
		{Title: "Ростов-на-Дону", EnTitle: "Rostov-on-Don", Description: "Город на юге России", Image: ""},
		{Title: "Уфа", EnTitle: "Ufa", Description: "Столица Башкортостана", Image: ""},
		{Title: "Красноярск", EnTitle: "Krasnoyarsk", Description: "Крупный центр Восточной Сибири", Image: ""},
		{Title: "Воронеж", EnTitle: "Voronezh", Description: "Культурный и промышленный центр", Image: ""},
		{Title: "Пермь", EnTitle: "Perm", Description: "Город на Урале", Image: ""},
		{Title: "Волгоград", EnTitle: "Volgograd", Description: "Город-герой на Волге", Image: ""},
		{Title: "Краснодар", EnTitle: "Krasnodar", Description: "Центр Краснодарского края", Image: ""},
		{Title: "Тюмень", EnTitle: "Tyumen", Description: "Один из старейших сибирских городов", Image: ""},
		{Title: "Ижевск", EnTitle: "Izhevsk", Description: "Столица Удмуртии", Image: ""},
		{Title: "Барнаул", EnTitle: "Barnaul", Description: "Крупный город в Алтайском крае", Image: ""},
		{Title: "Ульяновск", EnTitle: "Ulyanovsk", Description: "Родина В.И. Ленина", Image: ""},
		{Title: "Иркутск", EnTitle: "Irkutsk", Description: "Крупный город на Байкале", Image: ""},
		{Title: "Хабаровск", EnTitle: "Khabarovsk", Description: "Один из крупнейших городов Дальнего Востока", Image: ""},
		{Title: "Ярославль", EnTitle: "Yaroslavl", Description: "Один из старейших городов России", Image: ""},
		{Title: "Махачкала", EnTitle: "Makhachkala", Description: "Столица Дагестана", Image: ""},
		{Title: "Томск", EnTitle: "Tomsk", Description: "Крупный университетский город", Image: ""},
		{Title: "Оренбург", EnTitle: "Orenburg", Description: "Город на границе Европы и Азии", Image: ""},
		{Title: "Кемерово", EnTitle: "Kemerovo", Description: "Центр Кузбасса", Image: ""},
		{Title: "Рязань", EnTitle: "Ryazan", Description: "Один из древних городов России", Image: ""},
		{Title: "Астрахань", EnTitle: "Astrakhan", Description: "Крупный порт на Каспии", Image: ""},
	}

	var count int64
	if err := db.Model(&domain.City{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		if err := db.Create(&cities).Error; err != nil {
			return err
		}
	}

	return nil
}
