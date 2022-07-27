package check

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/helpers"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var FileDoesNotExist = errors.New("file does not exist")
var ApiKeyHasExpired = errors.New("api key has expired")
var NoApiKeysLeft = errors.New("no api keys left")

const pathToKeys = "check/keys.txt"
const pathToCheck = "check/check.docx"

func Format(doc *document.Document, dto dto.CheckDto) {
	var paragraphs []document.Paragraph
	for _, p := range doc.Paragraphs() {
		paragraphs = append(paragraphs, p)
	}
	//Filling header placeholders
	hh := doc.Headers()
	for _, h := range hh {
		for _, p := range h.Paragraphs() {
			for _, r := range p.Runs() {
				switch r.Text() {
				case "ord":
					r.ClearContent()
					ordIdStr := helpers.SixifyOrderId(dto.Order.OrderID)
					r.AddText(fmt.Sprintf("#%s", ordIdStr))
				}
			}
		}
	}
	//Filling footer placeholders
	ff := doc.Footers()
	for _, f := range ff {
		for _, p := range f.Paragraphs() {
			for _, r := range p.Runs() {
				switch r.Text() {
				case "sum":
					r.ClearContent()
					strSum := strconv.Itoa(int(dto.Order.TotalCartPrice))
					r.AddText(fmt.Sprintf("%s.0₽", strSum))
				}

			}
		}
	}
	//Filling body placeholders
	for _, p := range paragraphs {
		for _, r := range p.Runs() {
			switch r.Text() {
			case "username":
				r.ClearContent()
				r.AddText(dto.User.Username)
			case "phoneNumber":
				r.ClearContent()
				r.AddText(dto.User.PhoneNumber)
			case "address":
				r.ClearContent()
				if dto.Order.IsDelivered {
					address := "Адрес - ул. %s"
					r.AddText(fmt.Sprintf(address, dto.Order.DeliveryDetails.Address))
				}
			case "delivery":
				r.ClearContent()
				r.AddText(helpers.IsDeliveredTranslate(dto.Order.IsDelivered))
			case "PAY":
				r.ClearContent()
				r.AddText(helpers.PayTranslate(dto.Order.Pay))
			case "ent":
				r.ClearContent()
				if dto.Order.IsDelivered {
					r.AddText(fmt.Sprintf("Подъезд %d,", dto.Order.DeliveryDetails.EntranceNumber))
					continue
				}
			case "fl":
				r.ClearContent()
				if dto.Order.IsDelivered {
					r.AddText(fmt.Sprintf("Квартира %d.", dto.Order.DeliveryDetails.FlatCall))
					continue
				}
			case "gr":
				r.ClearContent()
				if dto.Order.IsDelivered {
					r.AddText(fmt.Sprintf("Этаж %d,", dto.Order.DeliveryDetails.Floor))
					continue
				}
			case "Содержимое":
				r.ClearContent()
				r.AddText("Содержимое:")
				r.AddBreak()
				pp := dto.Order.Cart
				for _, p := range pp {
					r.AddText(" - ")
					name := p.Name
					words := strings.Split(name, " ")
					if len(words) > 1 {
						name = words[0]
					}
					r.AddText(name)
					r.AddTab()
					r.AddTab()
					r.AddText(fmt.Sprintf("%d * %d.0₽", p.Quantity, p.Price))
					r.AddBreak()
				}
			}

		}
	}
}
func OpenTemplate(path string) (*document.Document, error) {

	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, FileDoesNotExist
	}

	doc, err := document.Open(path)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "credits") {
			return nil, ApiKeyHasExpired
		}
		return nil, err
	}
	return doc, nil
}
func SetLicense(key string) error {

	if err := license.SetMeteredKey(key); err != nil {
		return err
	}
	return nil

}
func GetFirstKey() (string, error) {
	//Open keys file
	file, err := os.Open(pathToKeys)
	defer file.Close()
	if err != nil {
		return "", err
	}

	//Read
	byt, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	//Get content
	content := string(byt)
	//Split by new-line, get first key
	keys := strings.Split(content, "\r")
	if len(keys) == 1 && keys[0] == "" {
		return "", NoApiKeysLeft
	}

	return keys[0], nil
}
func RestoreKey() error {
	//Open keys file to read content
	file, err := os.Open(pathToKeys)
	if err != nil {
		return err
	}
	//Read keys
	byt, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	keys := strings.Split(string(byt), "\n")

	//Check if 0 keys left
	if len(keys) == 1 && keys[0] == "" {
		return NoApiKeysLeft
	}

	//Remove dead key
	keys = keys[1:]

	//Close the file
	file.Close()

	//Open file again with trunc and write permissions
	file, err = os.OpenFile(pathToKeys, os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}

	//Remove file content
	err = file.Truncate(0)
	if err != nil {
		return err
	}
	//Make a buff
	buff := bytes.NewBufferString(strings.Join(keys, "\n"))

	//Write alive keys back
	_, err = file.Write(buff.Bytes())

	defer file.Close()
	return err
}

func Copy(w http.ResponseWriter) error {
	//mutex here
	file, err := os.Open(pathToCheck)
	stat, _ := file.Stat()

	w.Header().Add("Content-Length", fmt.Sprintf("%d", stat.Size()))

	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(w, file); err != nil {
		return err
	}

	return nil
}
