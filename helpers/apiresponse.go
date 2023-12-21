package helpers

import (
	"encoding/json"
	"net/http"
)

type ResponseSuccess struct {
	ResStatus  int         `json:"status"`
	Message    string      `json:"message"`
	Pagination interface{} `json:"pagination"`
	Result     interface{} `json:"data"`
}

type ResponseFailed struct {
	ResStatus int         `json:"status"`
	Message   interface{} `json:"message"`
}

type PaginationRes struct {
	PreviousPage    int `json:"previous_page"`
	CurrentPage     int `json:"current_page"`
	NextPage        int `json:"next_page"`
	LastPage        int `json:"last_page"`
	PerpageRecords  int `json:"perpage_records"`
	TotalPages      int `json:"total_pages"`
	TotalRecords    int `json:"total_records"`
	FilterDataMatch int `json:"total_match_records"`
}

func ApiFailedResponse(w http.ResponseWriter, message interface{}) {
	response := ResponseFailed{
		Message:   message,
		ResStatus: 0,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func ApiSuccessResponse(w http.ResponseWriter, result interface{}, message string, pagination_data interface{}) {
	if result == "" {
		result = map[int]int{}
	}
	var pagination PaginationRes

	if pagination_data != nil && pagination_data != "" {
		paginationMap, ok := pagination_data.(map[string]interface{})
		if !ok {
			http.Error(w, "Invalid pagination data", http.StatusInternalServerError)
			return
		}

		currentPage, currentPageExists := paginationMap["CurrentPage"].(int)
		totalPages, totalPagesExists := paginationMap["TotalPages"].(int)
		totalRecords, totalRecordsExists := paginationMap["TotalRows"].(int)
		perPageRecords, perpageRecordsExists := paginationMap["PerPageRecord"].(int)
		nextPage, nextPageExists := paginationMap["NextPage"].(int)
		prevPage, prevPageExists := paginationMap["PreviousPage"].(int)
		lastPage, lastPageExists := paginationMap["LastPage"].(int)
		matchRecords := paginationMap["matchCount"].(int)
		matchData := 0
		if matchRecords > 0 {
			matchData = matchRecords
		}

		if !currentPageExists || !totalPagesExists || !perpageRecordsExists || !totalRecordsExists || !nextPageExists || !prevPageExists || !lastPageExists {
			http.Error(w, "Missing required keys in pagination data", http.StatusInternalServerError)
			return
		}

		pagination = PaginationRes{
			PreviousPage:    prevPage,
			CurrentPage:     currentPage,
			NextPage:        nextPage,
			LastPage:        lastPage,
			PerpageRecords:  perPageRecords,
			TotalPages:      totalPages,
			TotalRecords:    totalRecords,
			FilterDataMatch: matchData,
		}

	}

	responseMap := map[string]interface{}{
		"message":   message,
		"status": 1,
		"result":    result,
	}

	if pagination_data != nil && pagination_data != "" {
		responseMap["Pagination"] = pagination
	}

	jsonResponse, err := json.Marshal(responseMap)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
