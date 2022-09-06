package check

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sonyamoonglade/delivery-service/internal/delivery/transport/dto"
	"github.com/sonyamoonglade/delivery-service/pkg/helpers"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
)

var FileDoesNotExist = errors.New("file does not exist")
var ApiKeyHasExpired = errors.New("api key has expired")
var NoApiKeysLeft = errors.New("no api keys left")

const (
	CheckWriteTimeout = time.Millisecond * 5000
	keysFilename      = "keys.txt"
	checkFilename     = "check.docx"
)

var pathToKeys string
var pathToCheck string

//initPaths is made for testing purposes
func initPaths(path string) {
	pathToKeys = path + "/" + keysFilename
	pathToCheck = path + "/" + checkFilename
}

type Service interface {
	Format(doc *document.Document, dto dto.CheckDto)
	OpenTemplate(path string) (*document.Document, error)
	SetLicense(key string) error
	GetFirstKey() (string, error)
	RestoreKey() error
	Copy(w http.ResponseWriter) error
}

type checkService struct {
	mut sync.Mutex
}

func NewCheckService(path string) Service {
	initPaths(path)
	return &checkService{mut: sync.Mutex{}}
}

func (c *checkService) Format(doc *document.Document, dto dto.CheckDto) {
	c.mut.Lock()
	defer c.mut.Unlock()
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
					text := fmt.Sprintf("#%s", ordIdStr)
					r.AddText(text)
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
					strSum := strconv.Itoa(int(dto.Order.Amount))
					text := fmt.Sprintf("%s.0₽", strSum)
					r.AddText(text)
				case "punish":
					//If not delivered dont care
					if !dto.Order.IsDelivered {
						r.ClearContent()
						break
					}

					var punishmentv int64 = 0

					//Compare total price for isPunished or not
					actualSum := helpers.CalculateTotalAmount(dto.Order.Cart)

					//Order is punished for delivery
					if actualSum < dto.Order.Amount {
						punishmentv = dto.Order.Amount - actualSum
					}
					//Fill punishment text
					r.ClearContent()
					text1 := "Оплата доставки"
					r.AddText(text1)

					r.AddTab()

					text2 := fmt.Sprintf("%d.0 ₽", punishmentv)
					if punishmentv == 0 {
						text2 = "бесплатно"
					}
					r.AddText(text2)
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
func (c *checkService) OpenTemplate(path string) (*document.Document, error) {
	c.mut.Lock()
	defer c.mut.Unlock()

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
func (c *checkService) SetLicense(key string) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	if err := license.SetMeteredKey(key); err != nil {
		return err
	}
	return nil

}
func (c *checkService) GetFirstKey() (string, error) {

	c.mut.Lock()
	defer c.mut.Unlock()
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
	keys := strings.Split(content, "\r\n")
	if len(keys) == 1 && keys[0] == "" {
		return "", NoApiKeysLeft
	}

	return keys[0], nil
}
func (c *checkService) RestoreKey() error {
	c.mut.Lock()
	defer c.mut.Unlock()

	file, err := os.Open(pathToKeys)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	file.Close()

	keys := strings.Split(string(content), "\r\n")

	//Check if 0 keys left
	if len(keys) == 1 && keys[0] == "" {
		return NoApiKeysLeft
	}

	//Remove dead key
	keys = keys[1:]

	//Open file again with trunc and write permissions
	file, err = os.OpenFile(pathToKeys, os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	//Clear file content
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	buff := bytes.NewBufferString(strings.Join(keys, "\r\n"))

	//Write alive keys back
	_, err = file.Write(buff.Bytes())
	if err != nil {
		return err
	}

	return nil
}
func (c *checkService) Copy(w http.ResponseWriter) error {
	c.mut.Lock()
	defer c.mut.Unlock()

	file, err := os.Open(pathToCheck)
	defer file.Close()

	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	w.Header().Add("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.WriteHeader(http.StatusOK)

	if _, err = io.Copy(w, file); err != nil {
		return err
	}

	return nil
}
