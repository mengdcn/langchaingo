package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/tmc/langchaingo/tgis/clipdropapi"
	clipdropapiParams "github.com/tmc/langchaingo/tgis/clipdropapi/params"
	"os"
	"strconv"
	"time"
)

const (
	textToImage       = "TextToImage"       // 文生图
	unCrop            = "UnCrop"            // 取消剪裁
	textInpainting    = "TextInpainting"    // 文本修复
	sketchToImage     = "SketchToImage"     // 草图到图
	replaceBackground = "ReplaceBackground" // 替换背景
	removeText        = "RemoveText"        // 删除文本
	removeBackground  = "RemoveBackground"  // 删除背景
	reimagine         = "Reimagine"         // 重新想象
	portraitSurface   = "PortraitSurface"   // 肖像表面法线
	portraitDepth     = "PortraitDepth"     // 人像深度估计
	imageUpscale      = "ImageUpscale"      // 图像放大
	cleanup           = "Cleanup"           // 清理
	saveFileDir       = "./images/"         // 保存目录 时间-图片类型
	imagesResource    = "./resource/"       // 原图片
)

var (
	clipdropApi *clipdropapi.ClipDropApi
	ctx         context.Context
)

func main() {
	apiKey := os.Getenv("CLIPDROP_KEY")
	fmt.Println(apiKey)

	clipdropLogic, err := clipdropapi.New(clipdropapi.WithAuthToken(apiKey))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	clipdropApi = clipdropLogic
	fmt.Println("clipdropApi", clipdropApi)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover ", r)
			msg := fmt.Sprintf("%v", r)
			err = errors.New(msg)
			return
		}
	}()

	ctx = context.Background()
	dType := unCrop // cleanup
	imagineResponse := &clipdropapiParams.ImagesResponse{}

	switch dType {
	case textToImage:
		imagineResponse, err = testImage()
		break
	case unCrop:
		imagineResponse, err = testUnCrop()

		break
	case textInpainting:
		imagineResponse, err = testTextInpainting()

		break
	case sketchToImage:
		imagineResponse, err = testSketchToImage()

		break
	case replaceBackground:
		imagineResponse, err = testReplaceBackground()

		break
	case removeText:
		imagineResponse, err = testRemoveText()

		break
	case removeBackground:
		imagineResponse, err = testRemoveBackground()

		break
	case reimagine:
		imagineResponse, err = testReimagine()

		break
	case portraitSurface:
		imagineResponse, err = testPortraitSurface()

		break
	case portraitDepth:
		imagineResponse, err = testPortraitDepth()

		break
	case imageUpscale:
		imagineResponse, err = testImageUpscale()

		break
	case cleanup:
		imagineResponse, err = testCleanup()

		break
	default:
		fmt.Println("please select the corresponding painting function")
		return
	}

	if err != nil {
		fmt.Printf("imagineResponse type %v exec , imagineResponse %v ,err %v \n ", dType, imagineResponse, err)
		return
	}

	err = saveFile(imagineResponse.ImgFile, imagineResponse.ImgExt, dType)
	if err != nil {
		fmt.Println("testExec saveFile err", err)
		return
	}

	fmt.Println("imagineResponse==imagineResponse.Success ", imagineResponse.Success)
	fmt.Println("imagineResponse==imagineResponse.XRemainingCredits ", imagineResponse.XRemainingCredits)
	fmt.Println("imagineResponse==imagineResponse.XReditsConsumed ", imagineResponse.XReditsConsumed)
	fmt.Println("imagineResponse==imagineResponse.ImgExt ", imagineResponse.ImgExt)
	fmt.Println("imagineResponse==imagineResponse.Error ", imagineResponse.Error)
}

func testCleanup() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imageRequest := &clipdropapiParams.CleanupRequest{}
	imageRequest.ImageFile = imagesResource + "clean.jpeg"
	imageRequest.MaskFile = imagesResource + "clean-mask.png"
	imageRequest.Mode = "quality" //  可选字段，可以设置为fast或quality控制速度和质量之间的权衡, fast是默认模式，速度更快，但可能会在结果图像中产生伪影; quality速度较慢，但会产生更好的结果

	imagesResponse, err = clipdropApi.Cleanup(ctx, imageRequest)
	if err != nil {
		fmt.Println("clipdropApi.Cleanup err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil
}

func testImageUpscale() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imageRequest := &clipdropapiParams.ImageUpscaleRequest{}
	imageRequest.ImageFile = imagesResource + "image-upscaling.png"
	imageRequest.TargetWidth = 4096
	imageRequest.TargetHeight = 4096

	imagesResponse, err = clipdropApi.ImageUpscale(ctx, imageRequest)
	if err != nil {
		fmt.Println("clipdropApi.ImageUpscale err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil
}

func testPortraitDepth() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imageRequest := &clipdropapiParams.PortraitDepthEstimationRequest{}
	imageRequest.ImageFile = imagesResource + "reimagine_1024x1024.jpg"

	imagesResponse, err = clipdropApi.PortraitDepth(ctx, imageRequest)
	if err != nil {
		fmt.Println("clipdropApi.PortraitDepth err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil
}

func testPortraitSurface() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imageRequest := &clipdropapiParams.PortraitSurfaceNormalsRequest{}
	imageRequest.ImageFile = imagesResource + "reimagine_1024x1024.jpg"

	imagesResponse, err = clipdropApi.PortraitSurface(ctx, imageRequest)
	if err != nil {
		fmt.Println("clipdropApi.PortraitSurface err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil
}

func testReimagine() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imageRequest := &clipdropapiParams.ReimagineRequest{}
	imageRequest.ImageFile = imagesResource + "reimagine_1024x1024.jpg"

	imagesResponse, err = clipdropApi.Reimagine(ctx, imageRequest)
	if err != nil {
		fmt.Println("clipdropApi.Reimagine err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil
}

func testRemoveBackground() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imageRequest := &clipdropapiParams.RemoveBackgroundRequest{}
	imageRequest.ImageFile = imagesResource + "remove-background.jpeg"

	imagesResponse, err = clipdropApi.RemoveBackground(ctx, imageRequest)
	if err != nil {
		fmt.Println("clipdropApi.RemoveBackground err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil
}

func testRemoveText() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imageRequest := &clipdropapiParams.RemoveTextRequest{}
	imageRequest.ImageFile = imagesResource + "remove-text-2_923x693.png"

	imagesResponse, err = clipdropApi.RemoveText(ctx, imageRequest)
	if err != nil {
		fmt.Println("clipdropApi.RemoveText err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil
}

func testReplaceBackground() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imageRequest := &clipdropapiParams.ReplaceBackgroundRequest{}
	imageRequest.ImageFile = imagesResource + "replace-background.jpg"
	imageRequest.Prompt = "a cozy marble kitchen with wine glasses"

	imagesResponse, err = clipdropApi.ReplaceBackground(ctx, imageRequest)
	if err != nil {
		fmt.Println("clipdropApi.ReplaceBackground err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil
}

func testSketchToImage() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imageRequest := &clipdropapiParams.SketchToImageRequest{}
	imageRequest.SketchFile = imagesResource + "Sketch-to-image_1024x1024.png"
	imageRequest.Prompt = "an owl on a branch, cinematic"

	imagesResponse, err = clipdropApi.SketchToImage(ctx, imageRequest)
	if err != nil {
		fmt.Println("clipdropApi.SketchToImage err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil
}

func testTextInpainting() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imageRequest := &clipdropapiParams.TextInpaintingRequest{}
	imageRequest.ImageFile = imagesResource + "text-inpainting.jpeg"
	imageRequest.MaskFile = imagesResource + "TextInpainting-mask_file.png"
	imageRequest.TextPrompt = "A woman with a red scarf"

	imagesResponse, err = clipdropApi.TextInpainting(ctx, imageRequest)
	if err != nil {
		fmt.Println("clipdropApi.TextInpainting err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil
}

func testUnCrop() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imageRequest := &clipdropapiParams.UnCropRequest{}
	imageRequest.ImageFile = imagesResource + "reimagine_1024x1024.jpg" //"image-upscaling.png"
	imageRequest.ExtendLeft = -700                                      // 可选 最大为 2k，默认为 0 【正负2k】
	imageRequest.ExtendRight = 0
	imageRequest.ExtendUp = 0
	imageRequest.ExtendDown = 0

	imagesResponse, err = clipdropApi.UnCrop(ctx, imageRequest)
	if err != nil {
		fmt.Println("clipdropApi.UnCrop err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil
}

func testImage() (imagesResponse *clipdropapiParams.ImagesResponse, err error) {
	imagineRequest := &clipdropapiParams.ImagesRequest{}
	imagineRequest.Prompt = "shot of vaporwave fashion dog in miami"
	//imagineRequest.Prompt = "Long, long ago, there was an ancient small village where the people lived a simple and simple life. This village has a special tradition, which is that every year a young person is selected as the moistening soil, responsible for nourishing the earth and ensuring a bountiful harvest. This year, the youngest boy in the village - Xiaoming, was selected as the Runtu. He is a kind and hardworking young man, full of awe and gratitude towards nature. Xiaoming was well aware of his heavy responsibility, so he began working hard day and night to ensure that every inch of land was fully nourished. However, this year there was particularly little rainfall in the village, the land dried up and cracked, and the crops turned yellow. Xiaoming saw it in his eyes and was anxious in his heart, so he decided to seek help. So he embarked on a journey to search for rainwater. Xiaoming climbed mountains and crossed mountains, traversed forests and deserts, and went through countless hardships. Finally, in a remote cave, he discovered a mysterious spring. This spring water is clear and sweet, and has magical power to nourish the earth. Xiaoming was overjoyed and immediately decided to bring this magical spring back to the village. On the way back to the village, Xiaoming encountered a difficult problem. That is a turbulent river, with a wide surface and turbulent currents. Xiao Ming was unable to cross the river directly, and he fell into a predicament. Just when he was at a loss, a kind turtle appeared. The turtle told Xiaoming that it could help him cross the river, but Xiaoming needed to agree to one condition. Xiaoming eagerly agreed to the turtle's conditions, so the turtle carried him across the river. Xiaoming was extremely grateful and thanked the turtle before continuing on his journey. Finally, Xiaoming returned to the village with the magical spring water. He immediately poured the spring water into the village canal, nourishing every inch of the land.urtle suddenly appeared. It told Xiaoming that the condition it had promised before was to ask Xiaoming to become its friend. " //"shot of vaporwave fashion dog in miami"
	// 长度计算 官方文档要求在1000个字符以内，测试1999个 （430个单词）也可以。2130 个字符就不可以（450个单词）
	//sl := strings.Split(imagineRequest.Prompt, " ")
	//slen := len(sl)
	//var solen, stlen int
	//for _, i2 := range sl {
	//	if strings.Contains(i2, ",") {
	//		slo := strings.Split(i2, ",")
	//		solen += len(slo)
	//	}
	//	if strings.Contains(i2, ".") {
	//		slt := strings.Split(i2, ".")
	//		stlen += len(slt)
	//	}
	//}
	//fmt.Println("len==", utf8.RuneCountInString(imagineRequest.Prompt), "slen ", slen, "slen, ", solen, "slen. ", stlen, "total len", slen+solen+stlen)

	// len== 2101 slen  355 slen,  50 slen.  46 total len 451 。报错500

	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("recover ", r)
			msg := fmt.Sprintf("recover %v", r)
			err = errors.New(msg)
			return
		}
	}()

	imagesResponse, err = clipdropApi.Images(ctx, imagineRequest)
	if err != nil {
		fmt.Println("clipdropApi.Images err", err.Error())
		return imagesResponse, err
	}

	return imagesResponse, nil

}

func saveFile(b []byte, imgExt, dType string) (err error) {
	if len(imgExt) < 1 {
		err = errors.New("img ext empty")
		return
	}
	fileName := saveFileDir + strconv.Itoa(int(time.Now().Unix())) + "-" + dType + "." + imgExt

	err = os.WriteFile(fileName, b, os.ModePerm)
	if err != nil {
		fmt.Println("error2:", err.Error())
		return err
	}
	return nil
}
