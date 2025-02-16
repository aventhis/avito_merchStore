// Package models содержит определения структур данных, используемых в проекте.
package models

// Merch представляет мерч с именем и ценой.
type Merch struct {
	Name  string // Название мерча
	Price int    // Цена мерча в монетах
}

// MerchList содержит список доступного мерча и их цены.
var MerchList = []Merch{
	{"t-shirt", 80},
	{"cup", 20},
	{"book", 50},
	{"pen", 10},
	{"powerbank", 200},
	{"hoody", 300},
	{"umbrella", 200},
	{"socks", 10},
	{"wallet", 50},
	{"pink-hoody", 500},
}
