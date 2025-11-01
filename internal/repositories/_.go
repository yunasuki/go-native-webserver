package repositories

import "go-native-webserver/internal/dal"

type BaseRepository struct {
	db dal.DatabaseConnection
}
