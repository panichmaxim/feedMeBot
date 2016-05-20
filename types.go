package main

type Cities struct {
  Status string `json:"status"`
  Items int `json:"items"`
  PageID int `json:"page_id"`
  PageLimit int `json:"page_limit"`
  Response[] City `json:"response"`
}

type City struct {
  ID string `json:"id"`
  Title string `json:"title"`
  Icon string `json:"icon"`
  Locale string `json:"locale"`
  Lat string `json:"lat"`
  Lon string `json:"lon"`
}

type Restaurants struct {
  Status      string `json:"status"`
  Items      int `json:"items"`
  PageID   int `json:"page_id"`
  PageLimit  int `json:"page_limit"`
  Response[]   Restaurant `json:"response"`
}
type Restaurant struct{
  ID         string    `json:"id"`
  LogoSquare     string `json:"logo_square"`
  CoverInner   string `json:"cover_inner"`
  Title       string `json:"city_id"`
  CityID     string  `json:"city_id"`
  Logo       string  `json:"logo"`
  Cover       string `json:"cover"`
  Description   string `json:"description"`
  Images[]     Images  `json:"images"`
  News       News `json:"news"`
  MenuFiles     MenuFiles `json:"menu_files"`
}

type MenuFiles struct {
  ID       int  `json:"id"`
  File     string  `json:"file"`
  Title     string  `json:"title"`
  Description string  `json:"description"`
  Size    string  `json:"size"`
}

type News struct {
  ID int `json:"id"`
  PlaceID int `json:"place_id"`
  Title string `json:"title"`
  Text string `json:"text"`
  Cover string `json:"cover"`
  CoverInner string `json:"cover_inner"`
}

type Images struct {
  ID string `json:"id"`
  Path string `json:"path"`
}