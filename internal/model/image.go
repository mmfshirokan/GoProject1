package model

type Image struct {
	LocalImgPath  string `json:"path" validate:"filepath"`
	ServerImgPath string `env:"IMAGE_STORE_PATH" envDefault:"/home/andreishyrakanau/projects/project1/GoProject1/images/" validate:"dir"`
	Name          string `json:"imgname"`
}
