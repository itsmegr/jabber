package handlers

import (
	"fmt"
	"jabber/server/service"
	"net/http"
)

func findGroup(name string) (*service.Group, bool) {
	allGroups := service.GlobalHub.Groups
	for group, ok := range allGroups {
		if(group.Name==name&&ok){
			return group, true
		}
	}
	return nil, false
}

func handleGroup(){
	
}

func JoinHandler(w http.ResponseWriter, r *http.Request){
	queryParams := r.URL.Query();
	groupName := queryParams.Get("group")
	clientName := queryParams.Get("name")


	// fmt.Println(groupName, clientName)
	w.Write([]byte(fmt.Sprintf("%v, %v", groupName, clientName)))
}





