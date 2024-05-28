package function

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/hsmtkk/ur-available-room/function/mail"
	"github.com/hsmtkk/ur-available-room/function/room"
)

const BASE_URL = "https://www.ur-net.go.jp"

func init() {
	functions.HTTP("EntryPoint", entryPoint)
}

func entryPoint(w http.ResponseWriter, r *http.Request) {
	req := room.Requirement{RentHigh: 150000, FloorSpaceLow: 50, Room: []string{"3K", "3DK", "3LDK"}, Prefecture: room.Tokyo, Area: room.East23}
	rooms, err := room.New().Get(req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(rooms)

	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")
	recipient := os.Getenv("MAIL_RECIPIENT")
	body := formatMailBody(rooms)
	status, id, err := mail.New(domain, apiKey).Send(r.Context(), "UR available room", body, recipient)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("status=%s, id=%s\n", status, id)
}

func formatMailBody(rooms []room.Room) string {
	lines := []string{}
	for _, room := range rooms {
		if room.RoomCount == 0 {
			continue
		}
		lines = append(lines, room.Name)
		lines = append(lines, room.Location)
		lines = append(lines, BASE_URL+room.URL)
		lines = append(lines, "---")
	}
	if len(lines) == 0 {
		return ""
	} else {
		return strings.Join(lines, "\n")
	}
}
