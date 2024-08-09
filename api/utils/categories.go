// categories.go
package utils

type CategoryInfo struct {
	Category string
	MaxPage  int
}

var Categories = []CategoryInfo{
	{Category: "bebidas-alcoholicas", MaxPage: 10},
	{Category: "bebidas-jugos-y-aguas", MaxPage: 8},
	{Category: "carniceria", MaxPage: 2},
	{Category: "cuidado-personal", MaxPage: 8},
	{Category: "desayuno", MaxPage: 5},
	{Category: "despensa", MaxPage: 13},
	{Category: "dulces-y-snacks", MaxPage: 5},
	{Category: "ferreteria", MaxPage: 1},
	{Category: "la-gran-feria-cugat", MaxPage: 3},
	{Category: "del-mundo-a-tu-despensa", MaxPage: 7},
	{Category: "lacteos", MaxPage: 7},
	{Category: "limpieza-y-aseo", MaxPage: 15},
	{Category: "mascotas", MaxPage: 1},
	{Category: "mundo-bebe", MaxPage: 3},
	{Category: "mundo-congelados", MaxPage: 9},
	{Category: "navidad", MaxPage: 1},
	{Category: "panaderia-y-pasteleria", MaxPage: 2},
	{Category: "preparados", MaxPage: 1},
	{Category: "quesos-y-fiambreria", MaxPage: 7},
}
