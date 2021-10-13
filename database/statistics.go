package database

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func updateStatistics(stats *Statistics) error {
	statsColl := &Statistics{}
	return mgm.Coll(statsColl).Update(stats)
}

func createStatistics(statistics *Statistics) error {
	stats := &Statistics{}
	return mgm.Coll(stats).Create(statistics)
}

func NodeOnline() error {
	stats, err := GetStatistics()
	if err != nil {
		log.Println(err)
		return err
	}
	stats.ActiveNodes = stats.ActiveNodes + 1
	return updateStatistics(stats)
}

func NodeOffline() error {
	stats, err := GetStatistics()
	if err != nil {
		log.Println(err)
		return err
	}

	// number of active nodes cannot exceed the number of total nodes
	if stats.ActiveNodes > 0 && stats.TotalNodes > stats.ActiveNodes {
		stats.ActiveNodes = stats.ActiveNodes - 1
	}
	return updateStatistics(stats)
}

func RemoveNode() error {
	stats, err := GetStatistics()
	if err != nil {
		log.Println(err)
		return err
	}
	if stats.TotalNodes > 0 {
		stats.TotalNodes = stats.TotalNodes - 1
	}
	return updateStatistics(stats)
}

func AddNode() error {
	stats, err := GetStatistics()
	if err != nil {
		log.Println(err)
		return err
	}
	stats.TotalNodes = stats.TotalNodes + 1
	return updateStatistics(stats)
}

func UpdateLastMessageId(lastMessageId int) error {
	stats, err := GetStatistics()
	if err != nil {
		log.Println(err)
		return err
	}
	stats.UpdateId = lastMessageId
	return updateStatistics(stats)
}

func GetStatistics() (*Statistics, error) {
	stats := &Statistics{}

	err := mgm.Coll(stats).First(bson.M{}, stats)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return stats, nil
}

func CreateStatisticsIfNotExists() {
	statistics := &Statistics{
		ActiveNodes: 0,
		TotalNodes:  0,
		UpdateId:    0,
	}
	if stats, err := GetStatistics(); stats == nil {
		err := createStatistics(statistics)
		if err != nil {
			log.Println(err)
			return
		}
	} else if err != nil {
		log.Println(err)
	}
}
