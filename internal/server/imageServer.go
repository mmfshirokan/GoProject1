package server

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"os"
	"strconv"

	"github.com/mmfshirokan/GoProject1/proto/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type ImageServer struct {
	pb.UnimplementedImageServer
}

func NewImageServer() pb.ImageServer {
	return &ImageServer{}
}

func (serv *ImageServer) UploadImage(req *pb.RequestUploadImage, stream pb.Image_UploadImageServer) error {

	id := req.GetUserID()
	imgName := req.GetImageName()

	err := stream.SetHeader(metadata.Pairs(
		"authorization", req.GetAuthToken(),
	))
	if err != nil {
		logError(err)
		return err
	}

	imgFull := make([]byte, 4200)
	imgPiece := &pb.ResponseUploadImage{}

	for {
		err = stream.RecvMsg(imgPiece)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("recv error att uploadImage:", err)
		}

		imgFull = append(imgFull, imgPiece.ImagePiece...)
	}

	img, _, err := image.Decode(bytes.NewBuffer(imgFull))
	if err != nil {
		logError(err)
		return err
	}

	destFile, err := os.Create(ImgNameWrap(id, imgName))
	if err != nil {
		logError(err)
		return err
	}

	err = png.Encode(destFile, img)
	if err != nil {
		logError(err)
		return err
	}

	log.Info("Upload finished")
	return nil
}

func (serv *ImageServer) DownloadImage(req *pb.RequestDownloadImage, stream pb.Image_DownloadImageServer) error {

	err := stream.SetHeader(metadata.Pairs(
		"authorization", req.GetAuthToken(),
	))
	if err != nil {
		log.Fatal(err)
	}

	imgFull, err := os.ReadFile(ImgNameWrap(req.GetUserID(), req.GetImageName()))
	if err != nil {
		log.Fatal(err)
	}

	imgPiece := make([]byte, 128)
	imgReader := bytes.NewReader(imgFull)

	for {
		_, err := imgReader.Read(imgPiece)
		if err == io.EOF {
			log.Info("Download finished")
			return err
		}
		if err != nil {
			log.Fatal(err)
		}

		stream.Send(&pb.ResponseDownloadImage{
			ImagePiece: imgPiece,
		})
	}
}

func ImgNameWrap(id int64, name string) string {
	return "/home/andreishyrakanau/projects/project1/GoProject1/images/" + strconv.FormatInt(id, 10) + "-" + name + ".png"
}
