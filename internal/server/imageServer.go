package server

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"strconv"

	"github.com/mmfshirokan/GoProject1/proto/pb"
	"google.golang.org/grpc/metadata"
)

type ImageServer struct {
	pb.UnimplementedImageServer
}

func NewImageServer() pb.ImageServer {
	return &ImageServer{}
}

func (serv *ImageServer) DownloadImage(req *pb.RequestDownloadImage, stream pb.Image_DownloadImageServer) error {
	err := stream.SetHeader(metadata.Pairs(
		"authorization", req.GetAuthToken(),
	))
	if err != nil {
		logError(err)
		return err
	}

	imgFull, err := os.ReadFile(ImgNameWrap(req.GetUserID(), req.GetImageName()))
	if err != nil {
		logError(err)
		return err
	}

	bytesLeftToRead := len(imgFull)
	imgPiece := make([]byte, 128)
	imgReader := bytes.NewReader(imgFull)

	for bytesLeftToRead > 0 {
		n, err := imgReader.Read(imgPiece)
		if err != nil {
			logError(err)
			return err
		}

		stream.Send(&pb.ResponseDownloadImage{
			ImagePiece:       imgPiece,
			StreamIsFinished: false,
		})

		bytesLeftToRead -= n
	}

	stream.Send(&pb.ResponseDownloadImage{
		ImagePiece:       nil, // rm?
		StreamIsFinished: true,
	})

	return nil
}

func (serv *ImageServer) UploadImage(stream pb.Image_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		logError(err)
		return err
	}

	err = stream.SetHeader(metadata.Pairs(
		"authorization", req.GetAuthToken(),
	))
	if err != nil {
		logError(err)
		return err
	}

	imgFull := make([]byte, 2048)
	streamIsFinished := false

	for !streamIsFinished {
		req, err = stream.Recv()
		if err != nil {
			logError(err)
			return err
		}

		imgFull = append(imgFull, req.GetImagePiece()...)
		streamIsFinished = req.GetStreamIsFinished()
	}

	img, _, err := image.Decode(bytes.NewReader(imgFull))
	if err != nil {
		logError(err)
		return err
	}

	destFile, err := os.Create(ImgNameWrap(req.GetUserID(), req.GetImageName()))
	if err != nil {
		logError(err)
		return err
	}

	err = png.Encode(destFile, img)
	if err != nil {
		logError(err)
		return err
	}

	return nil
}

func ImgNameWrap(id int64, name string) string {
	return "../../images/" + strconv.FormatInt(id, 10) + "-" + name + ".png"
}

// type ImageServer interface {
//     DownloadImage(*RequestDownloadImage, Image_DownloadImageServer) error
//     UploadImage(Image_UploadImageServer) error
//     mustEmbedUnimplementedImageServer()
// }

// imgFile, err := os.Open(req.GetImageLocation())
// if err != nil {
// 	logError(err)
// 	return err
// }
// defer imgFile.Close()

// imgInfo, err := os.Stat(req.GetImageLocation())
// if err != nil {
// 	logError(err)
// 	return err
// }

// imgBytesToRead := int(imgInfo.Size())

// imgFull := make([]byte, imgBytesToRead)

// for imgBytesToRead > 0 {
// 	n, err := imgFile.Read(imgBuffer)
// 	if err != nil {
// 		logError(err)
// 		return err
// 	}
//os.ReadFile("../../images/" + strconv.FormatInt(req.GetUserID(), 10) + "-" + req.GetImageName() + ".png")
// }

// imgFile, err := os.Open()
// if err != nil {
// 	logError(err)
// 	return err
// }
