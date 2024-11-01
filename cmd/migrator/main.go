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
		{Title: "Москва", EnTitle: "Moscow", Description: "Столица России"},
		{Title: "Санкт-Петербург", EnTitle: "Saint-Petersburg", Description: "Культурная столица России"},
		{Title: "Новосибирск", EnTitle: "Novosibirsk", Description: "Третий по численности город России"},
		{Title: "Екатеринбург", EnTitle: "Yekaterinburg", Description: "Административный центр Урала"},
		{Title: "Казань", EnTitle: "Kazan", Description: "Столица Республики Татарстан"},
		{Title: "Нижний Новгород", EnTitle: "Nizhny Novgorod", Description: "Важный культурный и экономический центр"},
		{Title: "Челябинск", EnTitle: "Chelyabinsk", Description: "Крупный промышленный город на Урале"},
		{Title: "Самара", EnTitle: "Samara", Description: "Крупный город на Волге"},
		{Title: "Омск", EnTitle: "Omsk", Description: "Один из крупнейших городов Сибири"},
		{Title: "Ростов-на-Дону", EnTitle: "Rostov-on-Don", Description: "Город на юге России"},
		{Title: "Уфа", EnTitle: "Ufa", Description: "Столица Башкортостана"},
		{Title: "Красноярск", EnTitle: "Krasnoyarsk", Description: "Крупный центр Восточной Сибири"},
		{Title: "Воронеж", EnTitle: "Voronezh", Description: "Культурный и промышленный центр"},
		{Title: "Пермь", EnTitle: "Perm", Description: "Город на Урале"},
		{Title: "Волгоград", EnTitle: "Volgograd", Description: "Город-герой на Волге"},
		{Title: "Краснодар", EnTitle: "Krasnodar", Description: "Центр Краснодарского края"},
		{Title: "Тюмень", EnTitle: "Tyumen", Description: "Один из старейших сибирских городов"},
		{Title: "Ижевск", EnTitle: "Izhevsk", Description: "Столица Удмуртии"},
		{Title: "Барнаул", EnTitle: "Barnaul", Description: "Крупный город в Алтайском крае"},
		{Title: "Ульяновск", EnTitle: "Ulyanovsk", Description: "Родина В.И. Ленина"},
		{Title: "Иркутск", EnTitle: "Irkutsk", Description: "Крупный город на Байкале"},
		{Title: "Хабаровск", EnTitle: "Khabarovsk", Description: "Один из крупнейших городов Дальнего Востока"},
		{Title: "Ярославль", EnTitle: "Yaroslavl", Description: "Один из старейших городов России"},
		{Title: "Махачкала", EnTitle: "Makhachkala", Description: "Столица Дагестана"},
		{Title: "Томск", EnTitle: "Tomsk", Description: "Крупный университетский город"},
		{Title: "Оренбург", EnTitle: "Orenburg", Description: "Город на границе Европы и Азии"},
		{Title: "Кемерово", EnTitle: "Kemerovo", Description: "Центр Кузбасса"},
		{Title: "Рязань", EnTitle: "Ryazan", Description: "Один из древних городов России"},
		{Title: "Астрахань", EnTitle: "Astrakhan", Description: "Крупный порт на Каспии"},
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
