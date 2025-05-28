package main

// import (
// 	"context"
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/user"
// 	"path/filepath"
// 	"time"

// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/google"
// 	"google.golang.org/api/calendar/v3"
// 	"google.golang.org/api/option"
// )

// func getClient(config *oauth2.Config) *http.Client {
//     usr, _ := user.Current()
//     tokenFile := filepath.Join(usr.HomeDir, ".credentials", "calendar_token.json")

//     tok, err := tokenFromFile(tokenFile)
//     if err != nil {
//         tok = getTokenFromWeb(config)
//         saveToken(tokenFile, tok)
//     }
//     return config.Client(context.Background(), tok)
// }

// func tokenFromFile(file string) (*oauth2.Token, error) {
//     f, err := os.Open(file)
//     if err != nil {
//         return nil, err
//     }
//     defer f.Close()
//     tok := &oauth2.Token{}
//     err = json.NewDecoder(f).Decode(tok)
//     return tok, err
// }

// func saveToken(path string, token *oauth2.Token) {
//     os.MkdirAll(filepath.Dir(path), 0700)
//     f, _ := os.Create(path)
//     defer f.Close()
//     json.NewEncoder(f).Encode(token)
//     fmt.Println("Saved token to", path)
// }

// func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
//     authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
//     fmt.Printf("Go to the following link in your browser then type the authorization code:\n%v\n", authURL)

//     var code string
//     fmt.Print("Enter the code: ")
//     fmt.Scan(&code)

//     tok, err := config.Exchange(context.TODO(), code)
//     if err != nil {
//         log.Fatalf("Unable to retrieve token from web: %v", err)
//     }
//     return tok
// }

// func CreateEvent(summary, location, description string, startTime, endTime time.Time) error {

//     b, err := os.ReadFile("/home/vapusdata/Downloads/credentials.json")
//     if err != nil {
//         return fmt.Errorf("unable to read client secret file: %v", err)
//     }

//     config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
//     if err != nil {
//         return fmt.Errorf("unable to parse client secret file to config: %v", err)
//     }

//     client := getClient(config)

//     srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
//     if err != nil {
//         return fmt.Errorf("unable to create calendar service: %v", err)
//     }

//     // here the calendar is creating the event
//     event := &calendar.Event{
//         Summary:     summary,
//         Location:    location,
//         Description: description,
//         Start: &calendar.EventDateTime{
//             DateTime: startTime.Format(time.RFC3339),
//             TimeZone: "Asia/Kolkata",
//         },
//         End: &calendar.EventDateTime{
//             DateTime: endTime.Format(time.RFC3339),
//             TimeZone: "Asia/Kolkata",
//         },
//     }

//     // here creating the meeting or event in the calendar
//     createdEvent, err := srv.Events.Insert("primary", event).Do()
//     if err != nil {
//         return fmt.Errorf("unable to create event: %v", err)
//     }

//     fmt.Printf(" Event created: %s\n", createdEvent.HtmlLink)
//     return nil
// }

// // this is for delete the meeting or event
// func DeleteEvent(eventID string) error {

//     b, err := os.ReadFile("/home/vapusdata/Downloads/credentials.json")
//     if err != nil {
//         return fmt.Errorf("unable to read client secret file: %v", err)
//     }

//     config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
//     if err != nil {
//         return fmt.Errorf("unable to parse config: %v", err)
//     }

//     client := getClient(config)

//     srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
//     if err != nil {
//         return fmt.Errorf("unable to create calendar service: %v", err)
//     }

//     err = srv.Events.Delete("primary", eventID).Do()
//     if err != nil {
//         return fmt.Errorf("unable to delete event: %v", err)
//     }

//     fmt.Println(" Event deleted successfully")
//     return nil
// }

// func UpdateEvent(eventID, summary, location, description string) error {
//     b, err := os.ReadFile("/home/vapusdata/Downloads/credentials.json")
//     if err != nil {
//         return fmt.Errorf("unable to read client secret file: %v", err)
//     }

//     config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
//     if err != nil {
//         return fmt.Errorf("unable to parse config: %v", err)
//     }

//     client := getClient(config)
//     srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
//     if err != nil {
//         return fmt.Errorf("unable to create calendar client: %v", err)
//     }

//     // Get the existing event first
//     event, err := srv.Events.Get("primary", eventID).Do()
//     if err != nil {
//         return fmt.Errorf("unable to fetch event: %v", err)
//     }

//     // Update fields
//     if summary != "" {
//         event.Summary = summary
//     }
//     if location != "" {
//         event.Location = location
//     }
//     if description != "" {
//         event.Description = description
//     }

//     updatedEvent, err := srv.Events.Update("primary", event.Id, event).Do()
//     if err != nil {
//         return fmt.Errorf("unable to update event: %v", err)
//     }

//     fmt.Printf("Event updated: %s\n", updatedEvent.HtmlLink)
//     return nil
// }

// // this is for fetching all the meetings or event

// func ListUpcomingEvents() error {

//     b, err := os.ReadFile("/home/vapusdata/Downloads/credentials.json")
//     if err != nil {
//         return fmt.Errorf("unable to read client secret file: %v", err)
//     }

//     config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
//     if err != nil {
//         return fmt.Errorf("unable to parse client secret file to config: %v", err)
//     }

//     client := getClient(config)

//     srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
//     if err != nil {
//         return fmt.Errorf("unable to retrieve Calendar client: %v", err)
//     }

//     t := time.Now().Format(time.RFC3339)
//     events, err := srv.Events.List("primary").
//         ShowDeleted(false).
//         SingleEvents(true).
//         TimeMin(t).
//         MaxResults(10).
//         OrderBy("startTime").
//         Do()
//     if err != nil {
//         return fmt.Errorf("unable to retrieve events: %v", err)
//     }

//     fmt.Println("Upcoming events:")
//     if len(events.Items) == 0 {
//         fmt.Println("No upcoming events found.")
//     } else {
//         for _, item := range events.Items {
//             date := item.Start.DateTime
//             if date == "" {
//                 date = item.Start.Date
//             }
//             fmt.Printf("%v - %v (ID: %v)\n", date, item.Summary, item.Id)
//         }
//     }

//     return nil
// }

// func main() {

//     action := flag.String("action", "list", "Action to perform: create | list | update | delete")
//     eventID := flag.String("id", "", "Event ID for update/delete")
//     summary := flag.String("summary", "", "Evnt Title")
//     location := flag.String("location", "", "Event location")
//     description := flag.String("description", "", "Event description")

//     flag.Parse()

//     switch *action {
//     case "create":
//         start := time.Now().Add(2 * time.Hour)
//         end := start.Add(30 * time.Minute)
//         if *summary == "" || *location == "" || *description == "" {
//             log.Fatalf("Missing required fields for create: summary, location, description")
//         }
//         err := CreateEvent(*summary, *location, *description, start, end)
//         if err != nil {
//             log.Fatalf("Error creating event: %v", err)
//         }

//     case "list":
//         err := ListUpcomingEvents()
//         if err != nil {
//             log.Fatalf("List error: %v", err)
//         }

//     case "update":

//         if *eventID == "" {
//         log.Fatal("Please provide event ID with -id")
//         }
//         err := UpdateEvent(*eventID, *summary, *location, *description)
//         if err != nil {
//             log.Fatalf("Error updating event: %v", err)
//         }

//     case "delete":
//         if *eventID == "" {
//             log.Fatal("Provide -id for delete")
//         }
//         err := DeleteEvent(*eventID)
//         if err != nil {
//             log.Fatalf("Delete error: %v", err)
//         }

//     default:
//         log.Fatalf("Unknown action: %s", *action)
//     }
// }
