package http_handlers

type JoinResponse struct {
	Token        string `json:"token"`
	PlayerID     int    `json:"player_id"`
	MapSizeX     int    `json:"map_size_x"`
	MapSizeY     int    `json:"map_size_y"`
	InterestSize int    `json:"interest_size"`
}
