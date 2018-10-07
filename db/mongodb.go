package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"strconv"
)

type JobDto struct {
	ChatId      int64  `bson:"_id"`
	RemindTime  string `bson:"remindTime"`
	SkipReminds int    `bson:"skipReminds"`
}

const (
	dbName             = "reminder"
	jobsCollectionName = "jobs"
)

var (
	db *mgo.Database
)

func OpenDbConnection(dbUrl string) mgo.Session {
	log.Println("connect to mongodb using url: " + dbUrl)
	session, error := mgo.Dial(dbUrl)
	if error != nil {
		panic(error)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session.DB(dbName)
	return *session
}

func RegisterJob(chatId int64, remindTime string) JobDto {
	jobs := db.C(jobsCollectionName)
	jobDto := JobDto{ChatId: chatId, RemindTime: remindTime, SkipReminds: 0}

	err := jobs.Insert(&jobDto)
	if err != nil {
		log.Println(err)
	}
	return jobDto
}

func GetSkipReminds(chatId int64) int {
	jobs := db.C(jobsCollectionName)
	job := JobDto{}
	err := jobs.FindId(chatId).One(&job)
	if err != nil {
		return job.SkipReminds
	} else {
		return 0
	}
}

func UpdateSkipReminds(chatId int64, newSkipReminds int) {
	jobs := db.C(jobsCollectionName)
	jobToUpdate := JobDto{}
	err := jobs.FindId(chatId).One(&jobToUpdate)
	if err != nil {
		jobToUpdate.SkipReminds = newSkipReminds
		jobs.UpdateId(chatId, jobToUpdate)
	} else {
		log.Println("no job to update for chatId: " + strconv.FormatInt(chatId, 10))
	}
}
