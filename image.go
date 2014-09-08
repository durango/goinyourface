package main

/*
#cgo LDFLAGS: -I/usr/local/lib -I/usr/local/include -lccv -lpng -ljpeg -lcblas -lfftw3 -lfftw3f
#include "ccv.h"
*/
import "C"

import (
  "fmt"
  "unsafe"
  "os"
  "log"
  "path/filepath"
  "bytes"
  "strings"
  "os/exec"
  "github.com/lazywei/go-opencv/opencv"
)

type Image struct {
  image *C.ccv_dense_matrix_t
  Detections []C.ccv_comp_t
  Type string
  Source string
  Color string
}

type DPMMixture struct {
  mixture *C.ccv_dpm_mixture_model_t
}

type ClassifierICF struct {
  classifier *C.ccv_icf_classifier_cascade_t
}

type ClassifierBBF struct {
  classifier *C.ccv_bbf_classifier_cascade_t
}

type CCVArray struct {
  array *C.ccv_array_t
}

type Comp struct {
  comp *C.ccv_comp_t
}

func main() {
  fmt.Println("Start")

  // compile source into CString
  imageSrc, _ := filepath.Abs(os.Args[1])

  cmd := exec.Command("jhead", "-autorot", fmt.Sprintf("%s", imageSrc))
  if err := cmd.Run(); err != nil {
    log.Fatal(err)
  }

  src := C.CString(imageSrc)
  defer C.free(unsafe.Pointer(src))

  imageICF := Image{Type: "ICF", Source: imageSrc, Color: "green"}
  imageBBF := Image{Type: "BBF", Source: imageSrc, Color: "green"}
  imageDPM := Image{Type: "DPM", Source: imageSrc, Color: "green"}

  // quickly draw opencv values...
  OpenCVClassifier("./haarcascade_frontalface_alt.xml", "OpenCV-frontal")
  OpenCVClassifier("./haarcascade_profileface.xml", "OpenCV-profile")

  // read for icf
  C.ccv_read_impl(unsafe.Pointer(src), &imageICF.image, C.CCV_IO_RGB_COLOR | C.CCV_IO_ANY_FILE, 0, 0, 0)

  // read for bbf
  C.ccv_read_impl(unsafe.Pointer(src), &imageBBF.image, C.CCV_IO_GRAY | C.CCV_IO_ANY_FILE, 0, 0, 0)

  // for dpm
  C.ccv_read_impl(unsafe.Pointer(src), &imageDPM.image, C.CCV_IO_ANY_FILE, 0, 0, 0)

  g1 := imageICF.icf()
  g2 := imageBBF.bbf()
  g3 := imageDPM.dpm()

  imageICF.Detections = <-g1
  imageBBF.Detections = <-g2
  imageDPM.Detections = <-g3

  if err := imageICF.drawBoxes(); err != nil {
    log.Fatal(err)
  }

  if err := imageDPM.drawBoxes(); err != nil {
    log.Fatal(err)
  }

  if err := imageBBF.drawBoxes(); err != nil {
    log.Fatal(err)
  }

  fmt.Println("Done!")
}

func OpenCVClassifier(cascadePath string, imageType string) error {
  imageSrc, _ := filepath.Abs(os.Args[1])
  imageOCV := opencv.LoadImage(imageSrc)
  classifier, _ := filepath.Abs(cascadePath)
  cascade := opencv.LoadHaarClassifierCascade(classifier)
  faces := cascade.DetectObjects(imageOCV)

  var ocvBuffer bytes.Buffer

  ocvBuffer.WriteString("")

  log.Printf("Faces: %+v", faces)
  for _, value := range faces {
    log.Printf("Face: %+v", value)
    ocvBuffer.WriteString(fmt.Sprintf("-draw \"rectangle %d,%d, %d,%d\" ", value.X(), value.Y(), value.X() + value.Width(), value.Y() + value.Height()))
  }

  if ocvBuffer.String() != "" {
    log.Println("Running opencv")
    cmd := fmt.Sprintf("`convert %s -fill none -stroke %s -strokewidth 2 %s%s`", imageSrc, "yellow", ocvBuffer.String(), fileName(imageSrc, imageType))

    if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
      log.Println(err)
    }
  }

  return nil
}

func (image *Image) fileName() string {
  return fileName(image.Source, image.Type)
}

func fileName(name string, imageType string) string {
  var buffer bytes.Buffer

  ext := filepath.Ext(name)
  buffer.WriteString(strings.TrimSuffix(name, ext))
  buffer.WriteString("-")
  buffer.WriteString(imageType)
  buffer.WriteString(ext)

  return buffer.String()
}

func (image *Image) drawBoxes() error {
  var buffer bytes.Buffer

  if len(image.Detections) < 1 {
    return nil
  }

  log.Printf("%+v", image.Detections)

  for _, comp := range image.Detections {
    box_x := (comp.rect.x + comp.rect.width)
    box_y := (comp.rect.y + comp.rect.height)
    final_box_x := box_x
    final_box_y := box_y
    final_x := comp.rect.x
    final_y := comp.rect.y

    // if the box is drawn out of scope.. no point in drawing it
    if final_y > image.image.rows || final_x > image.image.cols {
      continue
    }

    buffer.WriteString(fmt.Sprintf("-draw \"rectangle %d,%d, %d,%d\" ", final_x, final_y, final_box_x, final_box_y))
  }

  cmd := fmt.Sprintf("`convert %s -fill none -stroke %s -strokewidth 2 %s%s`", image.Source, image.Color, buffer.String(), image.fileName())

  if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
    log.Println(err)
    return err
  }

  return nil
}

func (image *Image) dpm() <- chan []C.ccv_comp_t {
  c := make(chan []C.ccv_comp_t)

  go func() {
    cascadeSrc, _ := filepath.Abs(os.Args[4])
    cascade := new(DPMMixture)
    cascadeName := C.CString(cascadeSrc)
    defer C.free(unsafe.Pointer(cascadeName))

    cascade.mixture = C.ccv_dpm_read_mixture_model(cascadeName)
    defer C.free(unsafe.Pointer(cascade.mixture))

    faces := new(CCVArray)
    faces.array = C.ccv_dpm_detect_objects(image.image, &cascade.mixture, 1, C.ccv_dpm_default_params)
    if faces.array == nil {
      c <- make([]C.ccv_comp_t, 0)
    } else {
      defer C.free(unsafe.Pointer(faces.array))
      slice := (*[1 << 30]C.ccv_comp_t)(unsafe.Pointer(faces.array.data))[:faces.array.rnum:faces.array.rnum]
      c <- slice
    }
  }()

  return c
}

func (image *Image) bbf() <- chan []C.ccv_comp_t {
  c := make(chan []C.ccv_comp_t)

  go func() {
    cascadeSrc, _ := filepath.Abs(os.Args[3])

    cascade := new(ClassifierBBF)
    cascadeName := C.CString(cascadeSrc)
    defer C.free(unsafe.Pointer(cascadeName))
    cascade.classifier = C.ccv_bbf_read_classifier_cascade(cascadeName)
    defer C.free(unsafe.Pointer(cascade.classifier))

    faces := new(CCVArray)
    faces.array = C.ccv_bbf_detect_objects(image.image, &cascade.classifier, 1, C.ccv_bbf_default_params)
    if faces.array == nil {
      c <- make([]C.ccv_comp_t, 0)
    } else {
      defer C.free(unsafe.Pointer(faces.array))
      slice := (*[1 << 30]C.ccv_comp_t)(unsafe.Pointer(faces.array.data))[:faces.array.rnum:faces.array.rnum]
      c <- slice
    }
  }()

  return c
}

func (image *Image) icf() <- chan []C.ccv_comp_t {
  c := make(chan []C.ccv_comp_t)

  go func() {
    cascadeSrc, _ := filepath.Abs(os.Args[2])

    cascade := new(ClassifierICF)
    cascadeName := C.CString(cascadeSrc)
    defer C.free(unsafe.Pointer(cascadeName))
    cascade.classifier = C.ccv_icf_read_classifier_cascade(cascadeName)
    defer C.free(unsafe.Pointer(cascade.classifier))

    // fmt.Println("doe")

    faces := new(CCVArray)
    faces.array = C.ccv_icf_detect_objects(image.image, unsafe.Pointer(&cascade.classifier), 1, C.ccv_icf_default_params)
    if faces.array == nil {
      c <- make([]C.ccv_comp_t, 0)
    } else {
      defer C.free(unsafe.Pointer(faces.array))
      slice := (*[1 << 30]C.ccv_comp_t)(unsafe.Pointer(faces.array.data))[:faces.array.rnum:faces.array.rnum]
      c <- slice
    }
  }()

  return c
}

