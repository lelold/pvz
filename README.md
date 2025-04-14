## **Сервис для ПВЗ**
Для запуска:
1) Клонируем: 

> git clone https://github.com/lelold/pvz.git

2) Меняем директорию:

> cd pvz
3) Собираем и запускаем контейнер:
> docker-compose up --build

	Сервис доступен по адресу: http://localhost:8080
4) Для повторного запуска без миграций используем:

> docker-compose up avito-pvz-service --build

**

## Endpoints

![image](https://github.com/user-attachments/assets/074781ac-1600-4596-8b36-923816bbb9f3)
![image](https://github.com/user-attachments/assets/4aeddc54-7d70-440a-9137-0ed1f04835b8)
![image](https://github.com/user-attachments/assets/dc585939-b891-40e2-a1ca-3bd756fe7268)
![image](https://github.com/user-attachments/assets/2e6280bc-3163-4748-862e-69148e62d97b)
![image](https://github.com/user-attachments/assets/1d842d9c-8866-4389-8afb-cbe9e8b2f424)
![image](https://github.com/user-attachments/assets/511b4776-5190-40b8-a540-041c005d0deb) 
![image](https://github.com/user-attachments/assets/df6e615d-5c11-4538-89bd-a0550b05b5ff)
![image](https://github.com/user-attachments/assets/18b3a691-c04c-4ef8-a1fa-23b562871cab)
![image](https://github.com/user-attachments/assets/0c6006b2-cd97-4de7-8b53-be1d98c2af9c)


## Стек
**Go** как главный ЯП
**PostgreSQL**  как СУБД
**Docker** для деплоя
**k6** для нагрузочных тестов

## Тесты
Для юнит-тестов использовался пакет sqlmock, а так же утилита mockery для генерации моков, k6 для нагрузочных тестов

![image](https://github.com/user-attachments/assets/b152ab62-fe0a-4a52-b9f7-d058119ebe3d)
![image](https://github.com/user-attachments/assets/a8ea6bdf-78bf-43fc-baaf-aea3642dda31)

Общее покрытие юзкейсов, хендлеров и репозиториев более 75%
