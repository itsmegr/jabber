package handlers

import (
	"encoding/json"
	"jabber/server/service"
	"log"
	"net/http"
)

/*
	stats params
	1. Active groups
	2. Total active members
*/

type StatRes struct {
	ActiveGroups int `json:"total_active_groups"`
	ActiveMembers int `json:"total_active_members"`
}

func StatsHandler(w http.ResponseWriter, r *http.Request){
	//getting total number of active groups
	activeGroups := len(service.GlobalHub.Groups)
	activeMembers := 0
	for group, _ := range service.GlobalHub.Groups {
		activeMembers = activeMembers + len(group.Clients)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	log.Println( activeGroups,activeMembers)
	res := StatRes{
		ActiveGroups: activeGroups,
		ActiveMembers: activeMembers,
	}
	json.NewEncoder(w).Encode(res)
}