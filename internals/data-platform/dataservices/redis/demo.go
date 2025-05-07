package redis

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"strconv"

// 	"github.com/redis/go-redis/v9"
// )

// func main() {
// 	ctx := context.Background()

// 	// Connect to Redis
// 	client := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379", // Change to your Redis server
// 	})

// 	// Get number of databases
// 	config, err := client.ConfigGet(ctx, "databases").Result()
// 	if err != nil {
// 		log.Fatalf("Error getting databases: %v", err)
// 	}

// 	// Extract the number of databases
// 	numDatabases, err := strconv.Atoi(config[1].(string))
// 	if err != nil {
// 		log.Fatalf("Error converting database count: %v", err)
// 	}

// 	fmt.Printf("Total Redis Databases: %d\n", numDatabases)

// 	// Iterate through each database and check keys
// 	for i := 0; i < numDatabases; i++ {
// 		client.Do(ctx, "SELECT", i) // Switch to DB i
// 		keys, err := client.Keys(ctx, "*").Result()
// 		if err != nil {
// 			log.Printf("Error fetching keys from DB %d: %v\n", i, err)
// 			continue
// 		}

// 		// Print database number and its keys
// 		if len(keys) > 0 {
// 			fmt.Printf("Database %d: %d keys\n", i, len(keys))
// 		}
// 	}

// 	// Close Redis connection
// 	defer client.Close()
// }
