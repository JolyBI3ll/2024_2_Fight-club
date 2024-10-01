package main

type Credentials struct {
	ID         int     `json:"id"`
	Username   string  `json:"username"`
	Password   string  `json:"password"`
	Email      string  `json:"email"`
	Name       string  `json:"name,omitempty"`
	Score      float64 `json:"score,omitempty"`
	Avatar     string  `json:"avatar,omitempty"`
	Sex        rune    `json:"sex,omitempty"`
	GuestCount int     `json:"guestCount,omitempty"`
	Age        int     `json:"age,omitempty"`
	Address    string  `json:"address,omitempty"`
	IsHost     bool    `json:"isHost,omitempty"`
}

type Author struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Score      float64 `json:"score"`
	Avatar     string  `json:"avatar"`
	Sex        rune    `json:"sex"`
	GuestCount int     `json:"guestCount"`
}

type Place struct {
	ID              int       `json:"id"`
	LocationMain    string    `json:"locationMain"`
	LocationStreet  string    `json:"locationStreet"`
	Position        []float64 `json:"position"`
	Images          []string  `json:"pictures"`
	Author          Author    `json:"author"`
	PublicationDate string    `json:"publicationDate"`
	AvailableDates  []string  `json:"avaibleDates"`
	Distance        float64   `json:"distance"`
}

var Users = []Credentials{
	{ID: 1, Username: "johndoe", Password: "password123", Email: "johndoe@example.com", Name: "Leo D.", Score: 4.98, Avatar: "images/avatar1.jpg", Sex: 'M', GuestCount: 50},
	{ID: 2, Username: "oleg", Password: "oleg123", Email: "oleg228@example.com", Name: "Oleg S.", Score: 4.50, Avatar: "images/avatar2.jpg", Sex: 'M', GuestCount: 20},
}

var Places = []Place{
	{
		ID: 1, LocationMain: "Москва", LocationStreet: "Тверская", Position: []float64{55.7558, 37.6173},
		Images: []string{"images/pic1.jpg", "images/pic2.jpg", "images/pic3.jpg"}, Author: Author{ID: 1, Name: "Лев Д.", Score: 4.98, Avatar: "images/avatar1.jpg", Sex: 'M', GuestCount: 50},
		PublicationDate: "2024-09-01", AvailableDates: []string{"2024-09-15", "2024-09-20"}, Distance: 5.0,
	},
	{
		ID: 2, LocationMain: "Санкт-Петербург", LocationStreet: "Невский проспект", Position: []float64{59.9311, 30.3609},
		Images: []string{"images/pic4.jpg", "images/pic5.jpg", "images/pic6.jpg", "images/pic7.jpg"}, Author: Author{ID: 2, Name: "Анна К.", Score: 4.85, Avatar: "images/avatar2.jpg", Sex: 'F', GuestCount: 45},
		PublicationDate: "2024-08-21", AvailableDates: []string{"2024-09-18", "2024-09-25"}, Distance: 2.3,
	},
	{
		ID: 3, LocationMain: "Казань", LocationStreet: "Улица Баумана", Position: []float64{55.7894, 49.1221},
		Images: []string{"images/pic8.jpg", "images/pic9.jpg", "images/pic10.jpg", "images/pic11.jpg", "images/pic12.jpg"}, Author: Author{ID: 3, Name: "Сергей И.", Score: 4.75, Avatar: "images/avatar3.jpg", Sex: 'M', GuestCount: 30},
		PublicationDate: "2024-07-15", AvailableDates: []string{"2024-09-10", "2024-09-28"}, Distance: 3.5,
	},
	{
		ID: 4, LocationMain: "Новосибирск", LocationStreet: "Красный проспект", Position: []float64{55.0415, 82.9346},
		Images: []string{"images/pic13.jpg", "images/pic14.jpg", "images/pic15.jpg"}, Author: Author{ID: 4, Name: "Олег С.", Score: 4.90, Avatar: "images/avatar4.jpg", Sex: 'M', GuestCount: 40},
		PublicationDate: "2024-10-05", AvailableDates: []string{"2024-10-12", "2024-10-20"}, Distance: 4.0,
	},
	{
		ID: 5, LocationMain: "Екатеринбург", LocationStreet: "Проспект Ленина", Position: []float64{56.8389, 60.6057},
		Images: []string{"images/pic16.jpg", "images/pic17.jpg", "images/pic18.jpg", "images/pic19.jpg"}, Author: Author{ID: 5, Name: "Наталья П.", Score: 4.88, Avatar: "images/avatar5.jpg", Sex: 'F', GuestCount: 35},
		PublicationDate: "2024-11-01", AvailableDates: []string{"2024-11-15", "2024-11-25"}, Distance: 3.7,
	},
	{
		ID: 6, LocationMain: "Нижний Новгород", LocationStreet: "Большая Покровская", Position: []float64{56.3287, 44.0020},
		Images: []string{"images/pic20.jpg", "images/pic21.jpg", "images/pic22.jpg", "images/pic23.jpg", "images/pic24.jpg", "images/pic25.jpg"}, Author: Author{ID: 6, Name: "Иван К.", Score: 4.92, Avatar: "images/avatar6.jpg", Sex: 'M', GuestCount: 42},
		PublicationDate: "2024-06-12", AvailableDates: []string{"2024-06-20", "2024-06-28"}, Distance: 6.1,
	},
	{
		ID: 7, LocationMain: "Красноярск", LocationStreet: "Проспект Мира", Position: []float64{56.0153, 92.8932},
		Images: []string{"images/pic26.jpg", "images/pic27.jpg", "images/pic28.jpg"}, Author: Author{ID: 7, Name: "Светлана З.", Score: 4.70, Avatar: "images/avatar7.jpg", Sex: 'F', GuestCount: 38},
		PublicationDate: "2024-12-02", AvailableDates: []string{"2024-12-10", "2024-12-15"}, Distance: 5.5,
	},
	{
		ID: 8, LocationMain: "Владивосток", LocationStreet: "Светланская", Position: []float64{43.1155, 131.8855},
		Images: []string{"images/pic29.jpg", "images/pic30.jpg", "images/pic31.jpg", "images/pic32.jpg"}, Author: Author{ID: 8, Name: "Михаил Р.", Score: 4.80, Avatar: "images/avatar8.jpg", Sex: 'M', GuestCount: 29},
		PublicationDate: "2024-05-05", AvailableDates: []string{"2024-05-12", "2024-05-20"}, Distance: 7.2,
	},
	{
		ID: 9, LocationMain: "Ростов-на-Дону", LocationStreet: "Большая Садовая", Position: []float64{47.2357, 39.7015},
		Images: []string{"images/pic33.jpg", "images/pic34.jpg", "images/pic35.jpg", "images/pic36.jpg", "images/pic37.jpg"}, Author: Author{ID: 9, Name: "Елена Г.", Score: 4.95, Avatar: "images/avatar9.jpg", Sex: 'F', GuestCount: 31},
		PublicationDate: "2024-04-18", AvailableDates: []string{"2024-04-25", "2024-05-02"}, Distance: 2.8,
	},
	{
		ID: 10, LocationMain: "Самара", LocationStreet: "Улица Куйбышева", Position: []float64{53.1959, 50.1002},
		Images: []string{"images/pic38.jpg", "images/pic39.jpg", "images/pic40.jpg", "images/pic41.jpg", "images/pic42.jpg"}, Author: Author{ID: 10, Name: "Олег И.", Score: 4.87, Avatar: "images/avatar10.jpg", Sex: 'M', GuestCount: 34},
		PublicationDate: "2024-03-21", AvailableDates: []string{"2024-04-01", "2024-04-08"}, Distance: 6.5,
	},
	{
		ID: 11, LocationMain: "Омск", LocationStreet: "Улица Ленина", Position: []float64{54.9885, 73.3242},
		Images: []string{"images/pic43.jpg", "images/pic44.jpg", "images/pic45.jpg", "images/pic46.jpg"}, Author: Author{ID: 11, Name: "Мария С.", Score: 4.93, Avatar: "images/avatar11.jpg", Sex: 'F', GuestCount: 29},
		PublicationDate: "2024-02-10", AvailableDates: []string{"2024-02-20", "2024-02-28"}, Distance: 4.4,
	},
	{
		ID: 12, LocationMain: "Волгоград", LocationStreet: "Проспект Ленина", Position: []float64{48.7080, 44.5133},
		Images: []string{"images/pic47.jpg", "images/pic48.jpg", "images/pic49.jpg", "images/pic50.jpg", "images/pic51.jpg"}, Author: Author{ID: 12, Name: "Алексей Н.", Score: 4.90, Avatar: "images/avatar12.jpg", Sex: 'M', GuestCount: 40},
		PublicationDate: "2024-01-15", AvailableDates: []string{"2024-01-25", "2024-02-05"}, Distance: 3.2,
	},
	{
		ID: 13, LocationMain: "Воронеж", LocationStreet: "Проспект Революции", Position: []float64{51.6720, 39.1843},
		Images: []string{"images/pic52.jpg", "images/pic53.jpg", "images/pic54.jpg"}, Author: Author{ID: 13, Name: "Евгений Ф.", Score: 4.85, Avatar: "images/avatar13.jpg", Sex: 'M', GuestCount: 45},
		PublicationDate: "2024-05-01", AvailableDates: []string{"2024-05-10", "2024-05-20"}, Distance: 5.6,
	},
	{
		ID: 14, LocationMain: "Краснодар", LocationStreet: "Улица Красная", Position: []float64{45.0355, 38.9753},
		Images: []string{"images/pic56.jpg", "images/pic57.jpg", "images/pic58.jpg", "images/pic59.jpg", "images/pic60.jpg"}, Author: Author{ID: 14, Name: "Игорь Т.", Score: 4.77, Avatar: "images/avatar14.jpg", Sex: 'M', GuestCount: 37},
		PublicationDate: "2024-05-25", AvailableDates: []string{"2024-06-05", "2024-06-15"}, Distance: 5.8,
	},
	{
		ID: 15, LocationMain: "Уфа", LocationStreet: "Улица Ленина", Position: []float64{54.7388, 55.9721},
		Images: []string{"images/pic61.jpg", "images/pic62.jpg", "images/pic63.jpg", "images/pic64.jpg"}, Author: Author{ID: 15, Name: "Ольга Е.", Score: 4.94, Avatar: "images/avatar15.jpg", Sex: 'F', GuestCount: 41},
		PublicationDate: "2024-09-01", AvailableDates: []string{"2024-09-10", "2024-09-20"}, Distance: 7.0,
	},
	{
		ID: 16, LocationMain: "Ярославль", LocationStreet: "Проспект Ленина", Position: []float64{57.6260, 39.8845},
		Images: []string{"images/pic65.jpg", "images/pic66.jpg", "images/pic67.jpg", "images/pic68.jpg", "images/pic69.jpg"}, Author: Author{ID: 16, Name: "Александр Ф.", Score: 4.92, Avatar: "images/avatar16.jpg", Sex: 'M', GuestCount: 33},
		PublicationDate: "2024-07-20", AvailableDates: []string{"2024-07-30", "2024-08-05"}, Distance: 4.3,
	},
	{
		ID: 17, LocationMain: "Тюмень", LocationStreet: "Площадь Ленина", Position: []float64{57.1530, 65.5343},
		Images: []string{"images/pic70.jpg", "images/pic71.jpg", "images/pic72.jpg", "images/pic73.jpg"}, Author: Author{ID: 17, Name: "Дмитрий Л.", Score: 4.78, Avatar: "images/avatar17.jpg", Sex: 'M', GuestCount: 36},
		PublicationDate: "2024-06-01", AvailableDates: []string{"2024-06-10", "2024-06-20"}, Distance: 3.9,
	},
	{
		ID: 18, LocationMain: "Тула", LocationStreet: "Проспект Ленина", Position: []float64{54.1920, 37.6156},
		Images: []string{"images/pic74.jpg", "images/pic75.jpg", "images/pic76.jpg", "images/pic77.jpg", "images/pic78.jpg"}, Author: Author{ID: 18, Name: "Юлия З.", Score: 4.99, Avatar: "images/avatar18.jpg", Sex: 'F', GuestCount: 50},
		PublicationDate: "2024-04-10", AvailableDates: []string{"2024-04-20", "2024-04-30"}, Distance: 6.2,
	},
	{
		ID: 19, LocationMain: "Ижевск", LocationStreet: "Улица Ленина", Position: []float64{56.8526, 53.2048},
		Images: []string{"images/pic79.jpg", "images/pic80.jpg"}, Author: Author{ID: 19, Name: "Павел К.", Score: 4.91, Avatar: "images/avatar19.jpg", Sex: 'M', GuestCount: 27},
		PublicationDate: "2024-02-25", AvailableDates: []string{"2024-03-05", "2024-03-15"}, Distance: 5.4,
	},
	{
		ID: 20, LocationMain: "Челябинск", LocationStreet: "Проспект Ленина", Position: []float64{55.1599, 61.4026},
		Images: []string{"images/pic81.jpg", "images/pic82.jpg", "images/pic83.jpg", "images/pic84.jpg", "images/pic85.jpg"}, Author: Author{ID: 20, Name: "Виктория М.", Score: 4.76, Avatar: "images/avatar20.jpg", Sex: 'F', GuestCount: 32},
		PublicationDate: "2024-03-05", AvailableDates: []string{"2024-03-15", "2024-03-25"}, Distance: 7.1,
	},
}
