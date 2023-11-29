package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Игрок
type Player struct {
	inventory []string //Инвентарь
}

type Room struct {
	Name       string              //Имя комнаты
	Items      map[string][]string //Предметы в комнате
	Exits      []string            //Выходы из комнаты
	Order      []string            //Порядок элементов в комнате
	DoorIsOpen bool
}

type World struct { //Мир
	Player      Player          //Добавляем Игрока
	Rooms       map[string]Room //Добавляем комнаты
	CurrentRoom string          //Текущая комната
	Street      Room
}

var world World
var startLocationName string
var goals []string
var itemsToCollect = []string{"конспекты"}

func initGame() {
	world.Player = Player{}
	world.Rooms = make(map[string]Room)
	startLocationName = "кухня"
	goals = []string{"собрать рюкзак", "идти в универ"}
	kitchen := Room{ //Инициализация кухни
		Name:       "на кухне",
		Items:      map[string][]string{"на столе": {"чай"}},
		Exits:      []string{"коридор"},
		DoorIsOpen: true,
		Order:      []string{"на столе"},
	}
	world.Rooms["кухня"] = kitchen
	mainroom := Room{
		Name:       "в своей комнате",
		Items:      map[string][]string{"на столе": {"ключи", "конспекты"}, "на стуле": {"рюкзак"}},
		Exits:      []string{"коридор"},
		DoorIsOpen: true,
		Order:      []string{"на столе", "на стуле"},
	}

	world.Rooms["комната"] = mainroom

	corridor := Room{ //Инициализация коридора
		Name:       "в коридоре",
		Items:      map[string][]string{},
		Exits:      []string{"кухня", "комната", "улица"},
		DoorIsOpen: true,
	}
	world.Rooms["коридор"] = corridor
	street := Room{ //Инициализация улицы
		Name:       "на улице",
		Items:      map[string][]string{},
		Exits:      []string{"домой"},
		DoorIsOpen: false,
	}
	world.Street = street
	world.CurrentRoom = "кухня"
}

func handleCommand(command string) string {
	parts := strings.Split(command, " ") //Разбили команду на части
	if len(parts) == 0 {
		return "пустая команда"
	}

	cur := parts[0]
	var param1, param2 string
	if len(parts) > 1 {
		param1 = parts[1]
	}
	if len(parts) > 2 {
		param2 = parts[2]
	}

	switch cur {
	case "осмотреться":
		return lookAround()
	case "идти":
		return goRoom(param1)
	case "надеть":
		return wearItem(param1)
	case "взять":
		return takeItem(param1)
	case "применить":
		return useItem(param1, param2)
	default:
		return "неизвестная команда"
	}

}
func formatGoals() string {
	var parts []string
	for _, goal := range goals {
		switch goal {
		case "собрать рюкзак":
			if !isBackpackCollected() {
				parts = append(parts, goal)
			}
		default:
			parts = append(parts, goal)
		}
	}
	if len(parts) > 0 {
		return "надо " + strings.Join(parts, " и ")
	}
	return ""
}

func isBackpackCollected() bool {
	for _, item := range itemsToCollect {
		found := false
		for _, invItem := range world.Player.inventory {
			if invItem == item {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func getItemsInRoom(room Room) map[string][]string { //Получаем предметы, которые находятся в комнате
	return room.Items
}
func formatRoomDescription(room Room) string {
	items := getItemsInRoom(room)
	var parts []string
	for _, location := range room.Order {
		if itemList, ok := items[location]; ok && len(itemList) > 0 {
			parts = append(parts, location+": "+strings.Join(itemList, ", "))
		}
	}

	return strings.Join(parts, ", ")
}

func lookAround() string { //Реализация функции осмотреться
	currentRoom, exists := world.Rooms[world.CurrentRoom]
	if !exists {
		currentRoom = world.Street
	}
	var description string
	exits := strings.Join(currentRoom.Exits, ", ")
	if world.CurrentRoom == startLocationName {
		description = "ты находишься" + " " + currentRoom.Name + ","
		goalsDescription := formatGoals()
		if len(currentRoom.Items) > 0 {
			description += " " + formatRoomDescription(currentRoom) + ","
		}
		if goalsDescription != "" {
			description += " " + goalsDescription + ". "
		}
		if exits != "" {
			description += "можно пройти -" + " " + exits
		}
		return description
	} else {
		emptyRoom := true
		for _, itemList := range currentRoom.Items {
			if len(itemList) > 0 {
				emptyRoom = false
				break
			}
		}
		if emptyRoom {
			description = "пустая комната. " + "можно пройти -" + " " + exits
			return description
		}
		if len(world.Rooms[world.CurrentRoom].Items) > 0 {
			description = formatRoomDescription(currentRoom) + ". можно пройти -" + " " + exits
		}
		return description
	}
}

func goRoom(targetRoom string) string {
	currentRoom, exists := world.Rooms[world.CurrentRoom] //Проверяем существует ли комната в списке комнат
	if !exists {
		currentRoom = world.Street //Если нет => улица
	}

	var description string
	hasExit := false //Проверяем есть ли выход из текущей комнаты в данную
	for _, exit := range currentRoom.Exits {
		if exit == targetRoom {
			hasExit = true
			break
		}
	}
	if hasExit {
		if world.Street.DoorIsOpen && targetRoom == "улица" {
			streetRoom := world.Street
			exits := strings.Join(streetRoom.Exits, ", ")
			world.CurrentRoom = targetRoom
			return "на улице весна. можно пройти - " + exits
		} else if !world.Rooms[targetRoom].DoorIsOpen { //Если дверь закрыта в данную комнату
			return "дверь закрыта"
		} else {
			if targetRoom == startLocationName {
				exits := strings.Join(world.Rooms[targetRoom].Exits, ", ")
				world.CurrentRoom = targetRoom
				return targetRoom + ", ничего интересного. можно пройти - " + exits
			} else if targetRoom == "коридор" {
				exits := strings.Join(world.Rooms[targetRoom].Exits, ", ")
				description = "ничего интересного. "
				description += "можно пройти - " + exits
				world.CurrentRoom = targetRoom
				return description
			} else {
				exits := strings.Join(world.Rooms[targetRoom].Exits, ", ")
				description = "ты " + world.Rooms[targetRoom].Name + ". "
				description += "можно пройти - " + exits
				world.CurrentRoom = targetRoom
				return description
			}
		}
	} else {
		return "нет пути в " + targetRoom
	}
}

func canwearItem(item string) bool { //Проверка можно ли надеть предмет
	return item == "рюкзак"
}

func findItemsLocation(items map[string][]string, itemName string) (string, int, bool) { //Поиск элемента в комнате
	for location, itemList := range items {
		for i, item := range itemList {
			if item == itemName {
				return location, i, true
			}
		}
	}
	return "", -1, false
}

func wearItem(itemName string) string { //Реализация функции надеть предмет
	currentRoom, exists := world.Rooms[world.CurrentRoom]
	if !exists {
		currentRoom = world.Street
	}
	location, index, itemExists := findItemsLocation(currentRoom.Items, itemName)
	if !itemExists {
		return "нет такого"
	}
	if canwearItem(itemName) {
		world.Player.inventory = append(world.Player.inventory, itemName)
		currentRoom.Items[location] = append(currentRoom.Items[location][:index], currentRoom.Items[location][index+1:]...)
		return "вы надели" + ": " + itemName
	}
	return "нельзя надеть"

}

func takeItem(itemName string) string {
	currentRoom, exists := world.Rooms[world.CurrentRoom]
	if !exists {
		currentRoom = world.Street
	}
	location, index, itemExists := findItemsLocation(currentRoom.Items, itemName)
	if !itemExists {
		return "нет такого"
	}

	hasBackpack := false
	for _, item := range world.Player.inventory {
		if item == "рюкзак" {
			hasBackpack = true
			break
		}
	}
	if !hasBackpack {
		return "некуда класть"
	}

	world.Player.inventory = append(world.Player.inventory, itemName)
	currentRoom.Items[location] = append(currentRoom.Items[location][:index], currentRoom.Items[location][index+1:]...) //Удаляем предмет из комнаты по локации и индексу

	return "предмет добавлен в инвентарь:" + " " + itemName
}
func useKey(target string, currentRoom Room) string { //Применяет ключ
	if target != "дверь" {
		return "не к чему применить"
	} else {
		world.Street.DoorIsOpen = true
		return "дверь открыта"
	}
}
func useItem(itemName string, target string) string {
	currentRoom, exists := world.Rooms[world.CurrentRoom] //Проверяем существует ли текущая комната
	if !exists {
		currentRoom = world.Street
	}

	var foundItem string
	for _, item := range world.Player.inventory {
		if item == itemName {
			foundItem = item
			break
		}
	}
	if foundItem == "" {
		return "нет предмета в инвентаре - " + itemName
	}
	switch itemName {
	case "ключи":
		return useKey(target, currentRoom)
	default:
		return "не к чему применить"
	}

}

func main() {
	initGame()

	fmt.Println("Добро пожаловать в игру!")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			return
		}

		command = strings.TrimSpace(command)
		if strings.ToLower(command) == "выход" {
			fmt.Println("We will be glad to see you again!")
			break
		}
		result := handleCommand(command)

		fmt.Println(result)

	}
}
