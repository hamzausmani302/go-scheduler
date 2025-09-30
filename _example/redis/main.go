package main

import (
	"log"
	"time"

	"github.com/algorythma/go-scheduler/storage"

	"github.com/algorythma/go-scheduler"
)

func TaskWithoutArgs() {
	log.Println("TaskWithoutArgs is executed")
}

func TaskWithArgs(message string) {
	log.Println("TaskWithArgs is executed. message:", message)
}

func main() {
	redisStore, err := storage.NewRedisStorage(storage.RedisConfig{
		Host: "localhost",
		Port: 6379,
		Password: "TiyNzDdjvY",
		Db: 0,
	})
	if err != nil {
		log.Fatalf("Couldn't create scheduler storage : %v", err)
	}
	s := scheduler.New(redisStore)
	// Start a task without arguments
	if _, err := s.RunAfter(30*time.Second, TaskWithoutArgs); err != nil {
		log.Fatal(err)
	}

	// Start a task with arguments
	if _, err := s.RunEvery(5*time.Second, TaskWithArgs, "Hello from recurring task 1"); err != nil {
		log.Fatal(err)
	}

	// Start the same task as above with a different argument
	if _, err := s.RunEvery(10*time.Second, TaskWithArgs, "Hello from recurring task 2"); err != nil {
		log.Fatal(err)
	}
	// Start the same task as above with a different argument
	s.Start()
	s.Wait()
}
